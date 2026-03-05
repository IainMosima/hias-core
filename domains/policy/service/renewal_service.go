package service

import (
	"context"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/google/uuid"
)

type RenewalService interface {
	InitiateRenewal(ctx context.Context, req policySchema.InitiateRenewalRequest, createdBy uuid.UUID) *schema.ServiceResponse[policySchema.RenewalResponse]
	GetRenewal(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[policySchema.RenewalResponse]
	ListByPolicy(ctx context.Context, policyID uuid.UUID) *schema.ServiceResponse[[]policySchema.RenewalResponse]
	ApproveRenewal(ctx context.Context, id uuid.UUID, approvedBy uuid.UUID) *schema.ServiceResponse[policySchema.RenewalResponse]
	RejectRenewal(ctx context.Context, id uuid.UUID, reason string) *schema.ServiceResponse[policySchema.RenewalResponse]
	CompleteRenewal(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[policySchema.RenewalResponse]
	ExpirePendingRenewals(ctx context.Context) *schema.ServiceResponse[int]
	BulkInitiateRenewals(ctx context.Context, policyIDs []uuid.UUID, createdBy uuid.UUID) *schema.ServiceResponse[policySchema.BulkResultResponse]
}
