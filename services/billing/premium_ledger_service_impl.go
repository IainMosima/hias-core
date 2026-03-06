package billing

import (
	"context"
	"net/http"

	"github.com/bitbiz/hias-core/domains/billing/entity"
	billingRepo "github.com/bitbiz/hias-core/domains/billing/repository"
	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	"github.com/bitbiz/hias-core/domains/billing/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/google/uuid"
)

type premiumLedgerServiceImpl struct {
	ledgerRepo billingRepo.PremiumLedgerRepository
}

func NewPremiumLedgerService(ledgerRepo billingRepo.PremiumLedgerRepository) service.PremiumLedgerService {
	return &premiumLedgerServiceImpl{ledgerRepo: ledgerRepo}
}

func (s *premiumLedgerServiceImpl) RecordEntry(ctx context.Context, req billingSchema.CreatePremiumLedgerRequest, createdBy uuid.UUID) *schema.ServiceResponse[billingSchema.PremiumLedgerResponse] {
	policyID, _ := uuid.Parse(req.PolicyID)

	// Get current balance to compute balance_after
	currentBalance, _ := s.ledgerRepo.GetBalanceByPolicy(ctx, policyID)

	var balanceAfter int64
	if req.EntryType == "CREDIT" {
		balanceAfter = currentBalance + req.Amount
	} else {
		balanceAfter = currentBalance - req.Amount
	}

	entry := &entity.PremiumLedgerEntry{
		PolicyID:        policyID,
		EntryType:       req.EntryType,
		Amount:          req.Amount,
		Description:     req.Description,
		ReferenceNumber: req.ReferenceNumber,
		EffectiveDate:   req.EffectiveDate,
		BalanceAfter:    balanceAfter,
		CreatedBy:       createdBy,
	}

	created, err := s.ledgerRepo.Create(ctx, entry)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.PremiumLedgerResponse](http.StatusInternalServerError, "Failed to record premium entry", err)
	}
	return schema.NewServiceResponse(billingSchema.ToPremiumLedgerResponse(created), http.StatusCreated, "Premium ledger entry created")
}

func (s *premiumLedgerServiceImpl) GetRegister(ctx context.Context, policyID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]billingSchema.PremiumLedgerResponse] {
	offset := (page - 1) * pageSize
	entries, err := s.ledgerRepo.ListByPolicy(ctx, policyID, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]billingSchema.PremiumLedgerResponse](http.StatusInternalServerError, "Failed to get premium register", err)
	}
	responses := make([]billingSchema.PremiumLedgerResponse, len(entries))
	for i, e := range entries {
		responses[i] = billingSchema.ToPremiumLedgerResponse(e)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Premium register retrieved")
}

func (s *premiumLedgerServiceImpl) GetBalance(ctx context.Context, policyID uuid.UUID) *schema.ServiceResponse[billingSchema.PremiumBalanceResponse] {
	balance, err := s.ledgerRepo.GetBalanceByPolicy(ctx, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.PremiumBalanceResponse](http.StatusInternalServerError, "Failed to get premium balance", err)
	}
	return schema.NewServiceResponse(billingSchema.PremiumBalanceResponse{
		PolicyID: policyID, Balance: balance,
	}, http.StatusOK, "Premium balance retrieved")
}
