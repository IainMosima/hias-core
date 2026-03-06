package service

import (
	"context"
	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/google/uuid"
)

type PremiumLedgerService interface {
	RecordEntry(ctx context.Context, req billingSchema.CreatePremiumLedgerRequest, createdBy uuid.UUID) *schema.ServiceResponse[billingSchema.PremiumLedgerResponse]
	GetRegister(ctx context.Context, policyID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]billingSchema.PremiumLedgerResponse]
	GetBalance(ctx context.Context, policyID uuid.UUID) *schema.ServiceResponse[billingSchema.PremiumBalanceResponse]
}
