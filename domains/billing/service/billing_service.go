package service

import (
	"context"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
)

type BillingService interface {
	GenerateInvoice(ctx context.Context, policyID string) *schema.ServiceResponse[billingSchema.InvoiceResponse]
	RunBillingCycle(ctx context.Context) *schema.ServiceResponse[int]
	SendReminder(ctx context.Context, invoiceID string) *schema.ServiceResponse[string]
	HandleOverdue(ctx context.Context) *schema.ServiceResponse[int]
}
