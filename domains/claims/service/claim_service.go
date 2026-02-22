package service

import (
	"context"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	claimsSchema "github.com/bitbiz/hias-core/domains/claims/schema"
	"github.com/google/uuid"
)

type ClaimService interface {
	SubmitClaim(ctx context.Context, req claimsSchema.SubmitClaimRequest, createdBy uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse]
	GetClaim(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse]
	ListClaims(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]claimsSchema.ClaimResponse]
	ListClaimsByStatus(ctx context.Context, status string, page, pageSize int) *schema.ServiceResponse[[]claimsSchema.ClaimResponse]
	ListClaimsByProvider(ctx context.Context, providerID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]claimsSchema.ClaimResponse]
	ApproveClaim(ctx context.Context, id uuid.UUID, approvedBy uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse]
	RejectClaim(ctx context.Context, id uuid.UUID, reason string, rejectedBy uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse]
	GetTotalCount(ctx context.Context) *schema.ServiceResponse[int64]
}
