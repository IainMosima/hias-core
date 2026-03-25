package service

import (
	"context"
	"time"

	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/google/uuid"
)

type InvoiceService interface {
	GetInvoice(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[billingSchema.InvoiceResponse]
	ListInvoices(ctx context.Context, dateFrom, dateTo *time.Time, page, pageSize int) *schema.ServiceResponse[billingSchema.InvoiceListResponse]
	CreateInvoice(ctx context.Context, req billingSchema.CreateInvoiceRequest) *schema.ServiceResponse[billingSchema.InvoiceResponse]
}
