package service

import (
	"context"
	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/google/uuid"
)

type CommissionService interface {
	CreateRule(ctx context.Context, req billingSchema.CreateCommissionRuleRequest, createdBy uuid.UUID) *schema.ServiceResponse[billingSchema.CommissionRuleResponse]
	ListRulesByPlan(ctx context.Context, planID uuid.UUID) *schema.ServiceResponse[[]billingSchema.CommissionRuleResponse]
	CalculateCommission(ctx context.Context, req billingSchema.CalculateCommissionRequest) *schema.ServiceResponse[billingSchema.CalculateCommissionResponse]
	ListPayments(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]billingSchema.CommissionPaymentResponse]
	ProcessPayments(ctx context.Context) *schema.ServiceResponse[int]
}
