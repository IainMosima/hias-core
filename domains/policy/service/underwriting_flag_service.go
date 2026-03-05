package service

import (
	"context"

	"github.com/bitbiz/hias-core/domains/identity/schema"
	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/google/uuid"
)

type UnderwritingFlagService interface {
	ListByPolicy(ctx context.Context, policyID uuid.UUID) *schema.ServiceResponse[[]policySchema.UnderwritingFlagResponse]
	ListByMember(ctx context.Context, memberID uuid.UUID) *schema.ServiceResponse[[]policySchema.UnderwritingFlagResponse]
	GetFlag(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[policySchema.UnderwritingFlagResponse]
	ListOpen(ctx context.Context, limit, offset int32) *schema.ServiceResponse[[]policySchema.UnderwritingFlagResponse]
	CountOpen(ctx context.Context) *schema.ServiceResponse[int64]
	ResolveFlag(ctx context.Context, id uuid.UUID, req policySchema.ResolveFlagRequest, resolvedBy uuid.UUID) *schema.ServiceResponse[policySchema.UnderwritingFlagResponse]
	OverrideFlag(ctx context.Context, id uuid.UUID, req policySchema.OverrideFlagRequest, overriddenBy uuid.UUID) *schema.ServiceResponse[policySchema.UnderwritingFlagResponse]
}
