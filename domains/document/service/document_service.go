package service

import (
	"context"

	"github.com/bitbiz/hias-core/domains/document/schema"
	identitySchema "github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/google/uuid"
)

type DocumentService interface {
	RequestUploadURL(ctx context.Context, req schema.UploadURLRequest, uploadedBy uuid.UUID) *identitySchema.ServiceResponse[schema.UploadURLResponse]
	ConfirmUpload(ctx context.Context, id uuid.UUID, uploadedBy uuid.UUID) *identitySchema.ServiceResponse[schema.DocumentResponse]
	GetDownloadURL(ctx context.Context, id uuid.UUID) *identitySchema.ServiceResponse[schema.DownloadURLResponse]
	ListByEntity(ctx context.Context, entityType string, entityID uuid.UUID) *identitySchema.ServiceResponse[[]schema.DocumentResponse]
	BulkRequestUploadURLs(ctx context.Context, req schema.BulkUploadURLRequest, uploadedBy uuid.UUID) *identitySchema.ServiceResponse[[]schema.UploadURLResponse]
	DeleteDocument(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) *identitySchema.ServiceResponse[string]
}
