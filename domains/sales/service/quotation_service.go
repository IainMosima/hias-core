package service

import (
	"context"

	"github.com/bitbiz/hias-core/domains/identity/schema"
	salesSchema "github.com/bitbiz/hias-core/domains/sales/schema"
	"github.com/google/uuid"
)

type QuotationService interface {
	// CRUD
	CreateQuotation(ctx context.Context, req salesSchema.CreateQuotationRequest, createdBy uuid.UUID) *schema.ServiceResponse[salesSchema.QuotationDetailResponse]
	GetQuotation(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[salesSchema.QuotationDetailResponse]
	ListQuotations(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]salesSchema.QuotationResponse]
	ListQuotationsByLead(ctx context.Context, leadID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]salesSchema.QuotationResponse]

	// Versioning
	CreateVersion(ctx context.Context, quotationID uuid.UUID, req salesSchema.CreateQuotationVersionRequest, createdBy uuid.UUID) *schema.ServiceResponse[salesSchema.QuotationVersionResponse]
	GetVersion(ctx context.Context, quotationID uuid.UUID, versionNumber int) *schema.ServiceResponse[salesSchema.QuotationVersionResponse]
	ListVersions(ctx context.Context, quotationID uuid.UUID) *schema.ServiceResponse[[]salesSchema.QuotationVersionResponse]
	CompareVersions(ctx context.Context, quotationID uuid.UUID, versionA, versionB int) *schema.ServiceResponse[salesSchema.VersionComparisonResponse]

	// Approval
	SubmitForApproval(ctx context.Context, quotationID uuid.UUID, versionNumber int, submittedBy uuid.UUID) *schema.ServiceResponse[salesSchema.QuotationVersionResponse]
	ApproveVersion(ctx context.Context, quotationID uuid.UUID, versionNumber int, req salesSchema.ApproveVersionRequest, approvedBy uuid.UUID, approverRole string) *schema.ServiceResponse[salesSchema.QuotationVersionResponse]
	RejectVersion(ctx context.Context, quotationID uuid.UUID, versionNumber int, req salesSchema.RejectVersionRequest, rejectedBy uuid.UUID) *schema.ServiceResponse[salesSchema.QuotationVersionResponse]

	// Lifecycle
	IssueQuotation(ctx context.Context, id uuid.UUID, issuedBy uuid.UUID) *schema.ServiceResponse[salesSchema.QuotationResponse]
	AcceptQuotation(ctx context.Context, id uuid.UUID, acceptedBy uuid.UUID) *schema.ServiceResponse[salesSchema.QuotationResponse]
	DeclineQuotation(ctx context.Context, id uuid.UUID, declinedBy uuid.UUID) *schema.ServiceResponse[salesSchema.QuotationResponse]
	ExpireQuotations(ctx context.Context) *schema.ServiceResponse[int]

	// Communication
	SendToClient(ctx context.Context, id uuid.UUID, req salesSchema.SendQuotationRequest, sentBy uuid.UUID) *schema.ServiceResponse[salesSchema.QuotationResponse]

	// Conversion
	ConvertToPolicy(ctx context.Context, id uuid.UUID, req salesSchema.ConvertToPolicyRequest, convertedBy uuid.UUID) *schema.ServiceResponse[salesSchema.ConversionResultResponse]

	// Documents
	UploadDocument(ctx context.Context, quotationID uuid.UUID, meta salesSchema.UploadDocumentMeta, s3Key string, uploadedBy uuid.UUID) *schema.ServiceResponse[salesSchema.QuotationDocumentResponse]
	ListDocuments(ctx context.Context, quotationID uuid.UUID) *schema.ServiceResponse[[]salesSchema.QuotationDocumentResponse]
	UpdateDocument(ctx context.Context, docID uuid.UUID, req salesSchema.UpdateDocumentMeta, updatedBy uuid.UUID, userRole string) *schema.ServiceResponse[salesSchema.QuotationDocumentResponse]
	DeleteDocument(ctx context.Context, docID uuid.UUID, deletedBy uuid.UUID, userRole string) *schema.ServiceResponse[bool]

	// Count
	GetTotalCount(ctx context.Context) *schema.ServiceResponse[int64]
}
