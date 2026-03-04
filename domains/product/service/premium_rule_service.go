package service

import (
	"context"
	"encoding/json"

	"github.com/bitbiz/hias-core/domains/identity/schema"
	productSchema "github.com/bitbiz/hias-core/domains/product/schema"
	"github.com/google/uuid"
)

type PremiumRuleService interface {
	CreatePremiumRule(ctx context.Context, planID uuid.UUID, req productSchema.CreatePremiumRuleRequest) *schema.ServiceResponse[productSchema.PremiumRuleResponse]
	ListPremiumRulesByPlan(ctx context.Context, planID uuid.UUID) *schema.ServiceResponse[[]productSchema.PremiumRuleResponse]
	DeletePremiumRule(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[string]
	CalculatePremium(ctx context.Context, planID uuid.UUID, memberCount int, relationships []string) *schema.ServiceResponse[int64]
	CalculatePremiumWithMembers(ctx context.Context, planID uuid.UUID, memberCount int, proposedMembers json.RawMessage) *schema.ServiceResponse[int64]
}
