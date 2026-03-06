package service

import (
	"context"
	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/google/uuid"
)

type RefundService interface {
	RequestRefund(ctx context.Context, req billingSchema.CreateRefundRequest, createdBy uuid.UUID) *schema.ServiceResponse[billingSchema.RefundResponse]
	ApproveRefund(ctx context.Context, id uuid.UUID, approvedBy uuid.UUID) *schema.ServiceResponse[billingSchema.RefundResponse]
	ProcessRefund(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[billingSchema.RefundResponse]
	ListRefunds(ctx context.Context, policyID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]billingSchema.RefundResponse]
}
