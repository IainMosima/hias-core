package service

import (
	"context"

	"github.com/bitbiz/hias-core/domains/identity/schema"
	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/google/uuid"
)

type UnderwritingRuleService interface {
	CreateRule(ctx context.Context, req policySchema.CreateUnderwritingRuleRequest, createdBy uuid.UUID) *schema.ServiceResponse[policySchema.UnderwritingRuleResponse]
	ListByPlan(ctx context.Context, planID uuid.UUID) *schema.ServiceResponse[[]policySchema.UnderwritingRuleResponse]
	UpdateRule(ctx context.Context, id uuid.UUID, req policySchema.UpdateUnderwritingRuleRequest) *schema.ServiceResponse[policySchema.UnderwritingRuleResponse]
	DeleteRule(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[string]
}
