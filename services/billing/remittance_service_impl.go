package billing

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/csv"
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

	// Calculate Withholding Tax
	whtRate := shared.DefaultWHTRate
	whtAmount := int64(float64(totalAmount) * whtRate)
	netAmount := totalAmount - whtAmount

	remittance := &entity.Remittance{
		ProviderID:  providerID,
		ClaimIDs:    claimIDsJSON,
		TotalAmount: totalAmount,
		Currency:    string(shared.CurrencyKES),
		Status:      string(shared.RemittanceStatusPending),
		PeriodStart: now.AddDate(0, -1, 0),
		PeriodEnd:   now,
		WHTRate:     whtRate,
		WHTAmount:   whtAmount,
		NetAmount:   netAmount,
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

func (s *remittanceServiceImpl) ExportPaymentFile(ctx context.Context, remittanceID uuid.UUID) *schema.ServiceResponse[billingSchema.PaymentExportResponse] {
	remittance, err := s.remittanceRepo.GetByID(ctx, remittanceID)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.PaymentExportResponse](http.StatusNotFound, "Remittance not found", err)
	}

	// Parse claim IDs from remittance
	var claimIDStrs []string
	if err := json.Unmarshal(remittance.ClaimIDs, &claimIDStrs); err != nil {
		return schema.NewServiceErrorResponse[billingSchema.PaymentExportResponse](http.StatusInternalServerError, "Failed to parse claim IDs", err)
	}

	// Fetch provider
	provider, err := s.providerRepo.GetByID(ctx, remittance.ProviderID)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.PaymentExportResponse](http.StatusInternalServerError, "Failed to get provider", err)
	}

	// Build claim details
	exportClaims := make([]billingSchema.PaymentExportClaim, 0, len(claimIDStrs))
	for _, idStr := range claimIDStrs {
		claimID, parseErr := uuid.Parse(idStr)
		if parseErr != nil {
			continue
		}
		claim, claimErr := s.claimRepo.GetByID(ctx, claimID)
		if claimErr != nil {
			continue
		}
		exportClaims = append(exportClaims, billingSchema.PaymentExportClaim{
			ClaimNumber: claim.ClaimNumber,
			Amount:      claim.ApprovedAmount,
			ServiceDate: claim.ServiceDate,
		})
	}

	// Generate CSV payment file
	var csvBuf bytes.Buffer
	csvWriter := csv.NewWriter(&csvBuf)

	// Header row
	csvWriter.Write([]string{"Provider", "Amount (KES)", "Currency", "Reference", "WHT Amount", "Net Amount"})
	// Summary row
	csvWriter.Write([]string{
		provider.Name,
		fmt.Sprintf("%.2f", float64(remittance.TotalAmount)/100),
		remittance.Currency,
		remittance.ID.String(),
		fmt.Sprintf("%.2f", float64(remittance.WHTAmount)/100),
		fmt.Sprintf("%.2f", float64(remittance.NetAmount)/100),
	})
	// Blank separator
	csvWriter.Write([]string{})
	// Claim detail header
	csvWriter.Write([]string{"Claim Number", "Amount (KES)", "Service Date"})
	for _, c := range exportClaims {
		csvWriter.Write([]string{
			c.ClaimNumber,
			fmt.Sprintf("%.2f", float64(c.Amount)/100),
			c.ServiceDate.Format("2006-01-02"),
		})
	}
	csvWriter.Flush()

	paymentFileCSV := base64.StdEncoding.EncodeToString(csvBuf.Bytes())

	export := billingSchema.PaymentExportResponse{
		RemittanceID:   remittance.ID,
		ProviderID:     remittance.ProviderID,
		ProviderName:   provider.Name,
		TotalAmount:    remittance.TotalAmount,
		Currency:       remittance.Currency,
		PeriodStart:    remittance.PeriodStart,
		PeriodEnd:      remittance.PeriodEnd,
		Claims:         exportClaims,
		PaymentFileCSV: paymentFileCSV,
	}

	return schema.NewServiceResponse(export, http.StatusOK, "Payment file exported")
}
