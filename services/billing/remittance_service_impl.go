package billing

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bitbiz/hias-core/domains/billing/entity"
	billingRepo "github.com/bitbiz/hias-core/domains/billing/repository"
	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	"github.com/bitbiz/hias-core/domains/billing/service"
	claimRepo "github.com/bitbiz/hias-core/domains/claims/repository"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	providerRepo "github.com/bitbiz/hias-core/domains/provider/repository"
	"github.com/bitbiz/hias-core/shared"
	"github.com/google/uuid"
)

type remittanceServiceImpl struct {
	remittanceRepo billingRepo.RemittanceRepository
	claimRepo      claimRepo.ClaimRepository
	providerRepo   providerRepo.ProviderRepository
}

func NewRemittanceService(
	remittanceRepo billingRepo.RemittanceRepository,
	claimRepo claimRepo.ClaimRepository,
	providerRepo providerRepo.ProviderRepository,
) service.RemittanceService {
	return &remittanceServiceImpl{
		remittanceRepo: remittanceRepo,
		claimRepo:      claimRepo,
		providerRepo:   providerRepo,
	}
}

func (s *remittanceServiceImpl) CreateRemittance(ctx context.Context, providerID uuid.UUID) *schema.ServiceResponse[billingSchema.RemittanceResponse] {
	claims, err := s.claimRepo.GetApprovedForRemittance(ctx, providerID)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.RemittanceResponse](http.StatusInternalServerError, "Failed to get approved claims", err)
	}

	if len(claims) == 0 {
		return schema.NewServiceErrorResponse[billingSchema.RemittanceResponse](http.StatusBadRequest, "No approved claims for remittance", nil)
	}

	var totalAmount int64
	claimIDs := make([]string, len(claims))
	for i, c := range claims {
		totalAmount += c.ApprovedAmount
		claimIDs[i] = c.ID.String()
	}

	claimIDsJSON, _ := json.Marshal(claimIDs)
	now := time.Now()

	remittance := &entity.Remittance{
		ProviderID:  providerID,
		ClaimIDs:    claimIDsJSON,
		TotalAmount: totalAmount,
		Currency:    string(shared.CurrencyKES),
		Status:      string(shared.RemittanceStatusPending),
		PeriodStart: now.AddDate(0, -1, 0),
		PeriodEnd:   now,
	}

	created, err := s.remittanceRepo.Create(ctx, remittance)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.RemittanceResponse](http.StatusInternalServerError, "Failed to create remittance", err)
	}

	return schema.NewServiceResponse(billingSchema.ToRemittanceResponse(created), http.StatusCreated, "Remittance created")
}

func (s *remittanceServiceImpl) RunRemittanceCycle(ctx context.Context) *schema.ServiceResponse[int] {
	providers, err := s.providerRepo.List(ctx, 1000, 0)
	if err != nil {
		return schema.NewServiceErrorResponse[int](http.StatusInternalServerError, "Failed to list providers", err)
	}

	created := 0
	for _, p := range providers {
		if p.Status == string(shared.ProviderStatusActive) {
			resp := s.CreateRemittance(ctx, p.ID)
			if resp.Error == nil {
				created++
			}
		}
	}

	return schema.NewServiceResponse(created, http.StatusOK, fmt.Sprintf("%d remittances created", created))
}

func (s *remittanceServiceImpl) SendRemittanceAdvice(ctx context.Context, remittanceID uuid.UUID) *schema.ServiceResponse[string] {
	_, err := s.remittanceRepo.MarkAdviceSent(ctx, remittanceID)
	if err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusInternalServerError, "Failed to send advice", err)
	}
	return schema.NewServiceResponse("sent", http.StatusOK, "Remittance advice sent")
}

func (s *remittanceServiceImpl) GetRemittance(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[billingSchema.RemittanceResponse] {
	remittance, err := s.remittanceRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.RemittanceResponse](http.StatusNotFound, "Remittance not found", err)
	}
	return schema.NewServiceResponse(billingSchema.ToRemittanceResponse(remittance), http.StatusOK, "Remittance retrieved")
}

func (s *remittanceServiceImpl) ListRemittances(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]billingSchema.RemittanceResponse] {
	offset := (page - 1) * pageSize
	remittances, err := s.remittanceRepo.List(ctx, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]billingSchema.RemittanceResponse](http.StatusInternalServerError, "Failed to list remittances", err)
	}

	responses := make([]billingSchema.RemittanceResponse, len(remittances))
	for i, r := range remittances {
		responses[i] = billingSchema.ToRemittanceResponse(r)
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Remittances retrieved")
}
