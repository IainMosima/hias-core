package service

import (
	"context"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/google/uuid"
)

type PolicyService interface {
	CreatePolicy(ctx context.Context, req policySchema.CreatePolicyRequest, createdBy uuid.UUID) *schema.ServiceResponse[policySchema.PolicyResponse]
	GetPolicy(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[policySchema.PolicyResponse]
	ListPolicies(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]policySchema.PolicyResponse]
	ActivatePolicy(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[policySchema.PolicyResponse]
	LapsePolicy(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[policySchema.PolicyResponse]
	TerminatePolicy(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[policySchema.PolicyResponse]
	ReinstatePolicy(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[policySchema.PolicyResponse]
	CalculateProratedPremium(ctx context.Context, policyID uuid.UUID) *schema.ServiceResponse[int64]
	GetTotalCount(ctx context.Context) *schema.ServiceResponse[int64]
}
