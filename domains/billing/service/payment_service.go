package service

import (
	"context"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	"github.com/google/uuid"
)

type PaymentService interface {
	InitiatePayment(ctx context.Context, req billingSchema.InitiatePaymentRequest, createdBy uuid.UUID) *schema.ServiceResponse[billingSchema.PaymentResponse]
	ProcessWebhook(ctx context.Context, data interface{}) *schema.ServiceResponse[string]
	RetryPayment(ctx context.Context, paymentID uuid.UUID) *schema.ServiceResponse[billingSchema.PaymentResponse]
	ReconcilePayment(ctx context.Context, paymentID uuid.UUID) *schema.ServiceResponse[billingSchema.PaymentResponse]
	GetPayment(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[billingSchema.PaymentResponse]
	ListPayments(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]billingSchema.PaymentResponse]
}
