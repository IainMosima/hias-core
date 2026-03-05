package service

import (
	"context"
	claimsSchema "github.com/bitbiz/hias-core/domains/claims/schema"
	"github.com/bitbiz/hias-core/domains/identity/schema"
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
	VetClaim(ctx context.Context, id uuid.UUID, req claimsSchema.VetClaimRequest, vettedBy uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse]
	MarkReadyForPayment(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse]
	MarkPaid(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse]
	MarkPartPaid(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse]
	BulkSubmitClaims(ctx context.Context, req claimsSchema.BulkSubmitClaimsRequest, createdBy uuid.UUID) *schema.ServiceResponse[[]claimsSchema.ClaimResponse]
	ListSLABreached(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]claimsSchema.ClaimResponse]
	UploadClaimDocument(ctx context.Context, claimID uuid.UUID, fileName, fileType string, fileSize int64, s3Key string, uploadedBy uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimDocumentResponse]
	ListClaimDocuments(ctx context.Context, claimID uuid.UUID) *schema.ServiceResponse[[]claimsSchema.ClaimDocumentResponse]
	DeleteClaimDocument(ctx context.Context, docID uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimDocumentResponse]
	GetTotalCount(ctx context.Context) *schema.ServiceResponse[int64]
}
