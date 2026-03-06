package service

import (
	"context"

	claimsSchema "github.com/bitbiz/hias-core/domains/claims/schema"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/google/uuid"
)

type EscalationRuleService interface {
	CreateRule(ctx context.Context, req claimsSchema.CreateEscalationRuleRequest) *schema.ServiceResponse[claimsSchema.EscalationRuleResponse]
	GetRule(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[claimsSchema.EscalationRuleResponse]
	ListRules(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]claimsSchema.EscalationRuleResponse]
	UpdateRule(ctx context.Context, id uuid.UUID, req claimsSchema.UpdateEscalationRuleRequest) *schema.ServiceResponse[claimsSchema.EscalationRuleResponse]
	DeleteRule(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[string]
}
