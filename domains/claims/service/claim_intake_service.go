package service

import (
	"context"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	claimsSchema "github.com/bitbiz/hias-core/domains/claims/schema"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/google/uuid"
)

type ClaimIntakeService interface {
	SubmitExternal(ctx context.Context, req claimsSchema.ExternalClaimRequest, partner *entity.APIPartner) *schema.ServiceResponse[claimsSchema.ExternalClaimResponse]
	GetExternalStatus(ctx context.Context, claimID uuid.UUID, partnerID uuid.UUID) *schema.ServiceResponse[claimsSchema.ExternalClaimStatusResponse]
	CreateDraft(ctx context.Context, req claimsSchema.DraftClaimRequest, createdBy uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse]
	UpdateDraft(ctx context.Context, id uuid.UUID, req claimsSchema.DraftClaimRequest, updatedBy uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse]
	ListDrafts(ctx context.Context, createdBy uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]claimsSchema.ClaimResponse]
	SubmitDraft(ctx context.Context, id uuid.UUID, submittedBy uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse]
	DeleteDraft(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse]
}
