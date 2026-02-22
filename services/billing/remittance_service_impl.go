package billing

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	billingEntity "github.com/bitbiz/hias-core/domains/billing/entity"
	billingRepo "github.com/bitbiz/hias-core/domains/billing/repository"
	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	billingService "github.com/bitbiz/hias-core/domains/billing/service"
	claimRepo "github.com/bitbiz/hias-core/domains/claims/repository"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	providerRepo "github.com/bitbiz/hias-core/domains/provider/repository"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/google/uuid"
)

type remittanceServiceImpl struct {
	remittanceRepo billingRepo.RemittanceRepository
	claimRepo      claimRepo.ClaimRepository
	providerRepo   providerRepo.ProviderRepository
	paymentRepo    billingRepo.PaymentRepository
}

func NewRemittanceService(
	remittanceRepo billingRepo.RemittanceRepository,
	claimRepo claimRepo.ClaimRepository,
	providerRepo providerRepo.ProviderRepository,
	paymentRepo billingRepo.PaymentRepository,
) billingService.RemittanceService {
	return &remittanceServiceImpl{
		remittanceRepo: remittanceRepo,
		claimRepo:      claimRepo,
		providerRepo:   providerRepo,
		paymentRepo:    paymentRepo,
	}
}

// CreateRemittance aggregates all approved claims for a given provider into
// a single remittance record. The remittance represents the total amount
// owed to the provider for services rendered during the period.
func (s *remittanceServiceImpl) CreateRemittance(ctx context.Context, providerID uuid.UUID) *schema.ServiceResponse[billingSchema.RemittanceResponse] {
	// Verify provider exists and is active
	provider, err := s.providerRepo.GetByID(ctx, providerID)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.RemittanceResponse](http.StatusNotFound, "Provider not found", err)
	}
	if provider.Status != string(shared.ProviderStatusActive) {
		return schema.NewServiceErrorResponse[billingSchema.RemittanceResponse](
			http.StatusConflict,
			fmt.Sprintf("Provider is not active (status: %s)", provider.Status),
			fmt.Errorf("provider %s is not active", providerID),
		)
	}

	// Get approved claims for this provider that haven't been remitted yet
	claims, err := s.claimRepo.GetApprovedForRemittance(ctx, providerID)
	if err != nil {
		utils.LogError("Failed to get approved claims for provider %s: %v", providerID, err)
		return schema.NewServiceErrorResponse[billingSchema.RemittanceResponse](http.StatusInternalServerError, "Failed to get approved claims", err)
	}

	if len(claims) == 0 {
		return schema.NewServiceErrorResponse[billingSchema.RemittanceResponse](
			http.StatusNotFound,
			"No approved claims available for remittance",
			fmt.Errorf("no approved claims for provider %s", providerID),
		)
	}

	// Calculate total and collect claim IDs
	var totalAmount int64
	claimIDs := make([]string, 0, len(claims))
	for _, c := range claims {
		totalAmount += c.ApprovedAmount
		claimIDs = append(claimIDs, c.ID.String())
	}

	claimIDsJSON, err := json.Marshal(claimIDs)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.RemittanceResponse](http.StatusInternalServerError, "Failed to encode claim IDs", err)
	}

	// Set the remittance period (covers the past 7 days by default)
	now := time.Now()
	periodStart := now.AddDate(0, 0, -7)

	remittance := &billingEntity.Remittance{
		ProviderID:  providerID,
		ClaimIDs:    claimIDsJSON,
		TotalAmount: totalAmount,
		Currency:    "KES",
		Status:      string(shared.RemittanceStatusPending),
		PeriodStart: periodStart,
		PeriodEnd:   now,
	}

	remittance, err = s.remittanceRepo.Create(ctx, remittance)
	if err != nil {
		utils.LogError("Failed to create remittance for provider %s: %v", providerID, err)
		return schema.NewServiceErrorResponse[billingSchema.RemittanceResponse](http.StatusInternalServerError, "Failed to create remittance", err)
	}

	utils.LogInfo("Remittance %s created for provider %s (%s): %d claims, total %d KES",
		remittance.ID, provider.Name, providerID, len(claims), totalAmount)

	return schema.NewServiceResponse(billingSchema.ToRemittanceResponse(remittance), http.StatusCreated, "Remittance created")
}

// RunRemittanceCycle iterates through all active providers and creates
// remittances for those with approved claims. This is intended to be run
// by a scheduled job (e.g., weekly remittance cycle).
func (s *remittanceServiceImpl) RunRemittanceCycle(ctx context.Context) *schema.ServiceResponse[int] {
	providers, err := s.providerRepo.ListByStatus(ctx, string(shared.ProviderStatusActive), 1000, 0)
	if err != nil {
		utils.LogError("Failed to list active providers for remittance cycle: %v", err)
		return schema.NewServiceErrorResponse[int](http.StatusInternalServerError, "Failed to list active providers", err)
	}

	createdCount := 0
	skippedCount := 0
	for _, p := range providers {
		resp := s.CreateRemittance(ctx, p.ID)
		if resp.Error != nil {
			// No approved claims or other issue -- not an error for the cycle
			skippedCount++
			continue
		}
		createdCount++
	}

	utils.LogInfo("Remittance cycle completed: %d remittances created, %d providers skipped (out of %d active)",
		createdCount, skippedCount, len(providers))

	return schema.NewServiceResponse(createdCount, http.StatusOK, fmt.Sprintf("Remittance cycle completed: %d remittances created", createdCount))
}

// SendRemittanceAdvice marks the remittance advice as sent. In a full
// implementation, this would generate a PDF remittance advice document
// and send it to the provider via email.
func (s *remittanceServiceImpl) SendRemittanceAdvice(ctx context.Context, remittanceID uuid.UUID) *schema.ServiceResponse[string] {
	remittance, err := s.remittanceRepo.GetByID(ctx, remittanceID)
	if err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusNotFound, "Remittance not found", err)
	}

	if remittance.RemittanceAdviceSent {
		return schema.NewServiceResponse(
			"Remittance advice already sent",
			http.StatusOK,
			"No action needed",
		)
	}

	// Mark as sent
	remittance, err = s.remittanceRepo.MarkAdviceSent(ctx, remittanceID)
	if err != nil {
		utils.LogError("Failed to mark remittance advice as sent for %s: %v", remittanceID, err)
		return schema.NewServiceErrorResponse[string](http.StatusInternalServerError, "Failed to send remittance advice", err)
	}

	// In a full implementation, this would:
	// 1. Generate a PDF remittance advice with claim details
	// 2. Look up the provider's email from the provider profile
	// 3. Send via the notification service (email channel)
	// 4. Log in the audit trail

	utils.LogInfo("Remittance advice sent for %s (provider: %s, amount: %d %s)",
		remittanceID, remittance.ProviderID, remittance.TotalAmount, remittance.Currency)

	return schema.NewServiceResponse(
		fmt.Sprintf("Remittance advice sent for remittance %s", remittanceID),
		http.StatusOK,
		"Remittance advice sent",
	)
}

// GetRemittance retrieves a single remittance by ID.
func (s *remittanceServiceImpl) GetRemittance(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[billingSchema.RemittanceResponse] {
	remittance, err := s.remittanceRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.RemittanceResponse](http.StatusNotFound, "Remittance not found", err)
	}
	return schema.NewServiceResponse(billingSchema.ToRemittanceResponse(remittance), http.StatusOK, "Remittance retrieved")
}

// ListRemittances returns a paginated list of remittances.
func (s *remittanceServiceImpl) ListRemittances(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]billingSchema.RemittanceResponse] {
	offset := (page - 1) * pageSize
	remittances, err := s.remittanceRepo.List(ctx, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]billingSchema.RemittanceResponse](http.StatusInternalServerError, "Failed to list remittances", err)
	}

	responses := make([]billingSchema.RemittanceResponse, 0, len(remittances))
	for _, r := range remittances {
		responses = append(responses, billingSchema.ToRemittanceResponse(r))
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Remittances retrieved")
}
