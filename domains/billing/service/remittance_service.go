package service

import (
	"context"
	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/google/uuid"
)

type RemittanceService interface {
	CreateRemittance(ctx context.Context, providerID uuid.UUID) *schema.ServiceResponse[billingSchema.RemittanceResponse]
	RunRemittanceCycle(ctx context.Context) *schema.ServiceResponse[int]
	SendRemittanceAdvice(ctx context.Context, remittanceID uuid.UUID) *schema.ServiceResponse[string]
	GetRemittance(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[billingSchema.RemittanceResponse]
	ListRemittances(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]billingSchema.RemittanceResponse]
	ExportPaymentFile(ctx context.Context, remittanceID uuid.UUID) *schema.ServiceResponse[billingSchema.PaymentExportResponse]
}
