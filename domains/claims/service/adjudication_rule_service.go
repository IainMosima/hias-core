package service

import (
	"context"

	claimsSchema "github.com/bitbiz/hias-core/domains/claims/schema"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/google/uuid"
)

type AdjudicationRuleService interface {
	CreateRule(ctx context.Context, req claimsSchema.CreateAdjudicationRuleRequest) *schema.ServiceResponse[claimsSchema.AdjudicationRuleResponse]
	GetRule(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[claimsSchema.AdjudicationRuleResponse]
	ListRules(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]claimsSchema.AdjudicationRuleResponse]
	UpdateRule(ctx context.Context, id uuid.UUID, req claimsSchema.UpdateAdjudicationRuleRequest) *schema.ServiceResponse[claimsSchema.AdjudicationRuleResponse]
	DeleteRule(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[string]
}
