package service

import (
	"context"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/google/uuid"
)

type EndorsementService interface {
	CreateEndorsement(ctx context.Context, req policySchema.CreateEndorsementRequest, requestedBy uuid.UUID) *schema.ServiceResponse[policySchema.EndorsementResponse]
	GetEndorsement(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[policySchema.EndorsementResponse]
	ListByPolicy(ctx context.Context, policyID uuid.UUID) *schema.ServiceResponse[[]policySchema.EndorsementResponse]
	ApproveEndorsement(ctx context.Context, id uuid.UUID, approvedBy uuid.UUID) *schema.ServiceResponse[policySchema.EndorsementResponse]
	RejectEndorsement(ctx context.Context, id uuid.UUID, reason string, rejectedBy uuid.UUID) *schema.ServiceResponse[policySchema.EndorsementResponse]
	ApplyEndorsement(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[policySchema.EndorsementResponse]
	CancelEndorsement(ctx context.Context, id uuid.UUID, reason string, cancelledBy uuid.UUID) *schema.ServiceResponse[policySchema.EndorsementResponse]
}
