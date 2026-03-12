package service

import (
	"context"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/google/uuid"
)

type PolicyDocumentService interface {
	GenerateWelcomeLetter(ctx context.Context, policyID uuid.UUID, generatedBy uuid.UUID) *schema.ServiceResponse[policySchema.PolicyDocumentResponse]
	GenerateMemberCard(ctx context.Context, memberID uuid.UUID, generatedBy uuid.UUID) *schema.ServiceResponse[policySchema.PolicyDocumentResponse]
	GeneratePolicySchedule(ctx context.Context, policyID uuid.UUID, generatedBy uuid.UUID) *schema.ServiceResponse[policySchema.PolicyDocumentResponse]
	GenerateRenewalNotice(ctx context.Context, renewalID uuid.UUID, generatedBy uuid.UUID) *schema.ServiceResponse[policySchema.PolicyDocumentResponse]
	ListByPolicy(ctx context.Context, policyID uuid.UUID) *schema.ServiceResponse[[]policySchema.PolicyDocumentResponse]
	GetDocument(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[policySchema.PolicyDocumentResponse]
	DeleteDocument(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[string]
	BulkGenerateMemberCards(ctx context.Context, policyID uuid.UUID, generatedBy uuid.UUID) *schema.ServiceResponse[[]policySchema.PolicyDocumentResponse]
	GenerateLOU(ctx context.Context, preauthID uuid.UUID, generatedBy uuid.UUID) *schema.ServiceResponse[policySchema.PolicyDocumentResponse]
	GenerateDeclineLetter(ctx context.Context, policyID uuid.UUID, memberName, claimNumber, rejectionReason string, generatedBy uuid.UUID) *schema.ServiceResponse[policySchema.PolicyDocumentResponse]

	// V1 unified document generation
	GenerateDocument(ctx context.Context, req policySchema.GenerateDocumentRequest) *schema.ServiceResponse[policySchema.PolicyDocumentResponse]
	CanGenerateDocument(ctx context.Context, entityType string, entityID uuid.UUID, docType string) *schema.ServiceResponse[policySchema.DocumentReadinessResponse]
	GetDocumentAvailability(ctx context.Context, entityType string, entityID uuid.UUID) *schema.ServiceResponse[[]policySchema.DocumentAvailabilityItem]

	// Upload flow
	RequestUploadURL(ctx context.Context, policyID uuid.UUID, req policySchema.UploadPolicyDocumentURLRequest, uploadedBy uuid.UUID) *schema.ServiceResponse[policySchema.PolicyDocumentUploadURLResponse]
	ConfirmUpload(ctx context.Context, documentID uuid.UUID, uploadedBy uuid.UUID) *schema.ServiceResponse[policySchema.PolicyDocumentResponse]
}
