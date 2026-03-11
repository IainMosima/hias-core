package document

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	"github.com/bitbiz/hias-core/domains/document/entity"
	docRepo "github.com/bitbiz/hias-core/domains/document/repository"
	"github.com/bitbiz/hias-core/domains/document/schema"
	docService "github.com/bitbiz/hias-core/domains/document/service"
	identitySchema "github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/shared"
	awsSvc "github.com/bitbiz/hias-core/shared/aws"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/google/uuid"
)

type DocumentServiceConfig struct {
	MaxFileSize  int
	AllowedTypes []string
}

type documentServiceImpl struct {
	docRepo  docRepo.DocumentRepository
	s3Svc    awsSvc.S3Service
	auditSvc auditService.AuditService
	config   DocumentServiceConfig
}

func NewDocumentService(
	repo docRepo.DocumentRepository,
	s3Svc awsSvc.S3Service,
	auditSvc auditService.AuditService,
	config DocumentServiceConfig,
) docService.DocumentService {
	return &documentServiceImpl{
		docRepo:  repo,
		s3Svc:    s3Svc,
		auditSvc: auditSvc,
		config:   config,
	}
}

func (s *documentServiceImpl) RequestUploadURL(ctx context.Context, req schema.UploadURLRequest, uploadedBy uuid.UUID) *identitySchema.ServiceResponse[schema.UploadURLResponse] {
	if err := utils.ValidateFileName(req.FileName); err != nil {
		return identitySchema.NewServiceErrorResponse[schema.UploadURLResponse](http.StatusBadRequest, err.Error(), err)
	}
	if err := utils.ValidateFileSize(req.FileSize, s.config.MaxFileSize); err != nil {
		return identitySchema.NewServiceErrorResponse[schema.UploadURLResponse](http.StatusBadRequest, err.Error(), err)
	}
	if err := utils.ValidateMimeType(req.MimeType, s.config.AllowedTypes); err != nil {
		return identitySchema.NewServiceErrorResponse[schema.UploadURLResponse](http.StatusBadRequest, err.Error(), err)
	}

	entityID, err := uuid.Parse(req.EntityID)
	if err != nil {
		return identitySchema.NewServiceErrorResponse[schema.UploadURLResponse](http.StatusBadRequest, "Invalid entity ID", err)
	}

	s3Key := utils.GenerateS3Key(req.EntityType, req.EntityID, req.DocumentType, req.FileName)

	doc, err := s.docRepo.Create(ctx, &entity.Document{
		EntityType:   req.EntityType,
		EntityID:     entityID,
		DocumentType: req.DocumentType,
		Status:       string(shared.DocumentStatusPendingUpload),
		FileName:     req.FileName,
		FileSize:     req.FileSize,
		MimeType:     req.MimeType,
		S3Key:        s3Key,
		UploadedBy:   uploadedBy,
	})
	if err != nil {
		return identitySchema.NewServiceErrorResponse[schema.UploadURLResponse](http.StatusInternalServerError, "Failed to create document record", err)
	}

	var expiresIn int64 = 900
	uploadURL, err := s.s3Svc.GetPresignedPutURL(ctx, s3Key, req.MimeType, expiresIn)
	if err != nil {
		return identitySchema.NewServiceErrorResponse[schema.UploadURLResponse](http.StatusInternalServerError, "Failed to generate upload URL", err)
	}

	return identitySchema.NewServiceResponse(schema.UploadURLResponse{
		DocumentID: doc.ID,
		UploadURL:  uploadURL,
		S3Key:      s3Key,
		ExpiresIn:  expiresIn,
	}, http.StatusOK, "Upload URL generated")
}

func (s *documentServiceImpl) ConfirmUpload(ctx context.Context, id uuid.UUID, uploadedBy uuid.UUID) *identitySchema.ServiceResponse[schema.DocumentResponse] {
	doc, err := s.docRepo.GetByID(ctx, id)
	if err != nil {
		return identitySchema.NewServiceErrorResponse[schema.DocumentResponse](http.StatusNotFound, "Document not found", err)
	}

	if doc.Status != string(shared.DocumentStatusPendingUpload) {
		return identitySchema.NewServiceErrorResponse[schema.DocumentResponse](http.StatusBadRequest, fmt.Sprintf("Document status is %s, expected PENDING_UPLOAD", doc.Status), fmt.Errorf("invalid status"))
	}

	_, err = s.s3Svc.HeadObject(ctx, doc.S3Key)
	if err != nil {
		return identitySchema.NewServiceErrorResponse[schema.DocumentResponse](http.StatusBadRequest, "File not found in S3 — upload may have failed", err)
	}

	now := time.Now()
	err = s.docRepo.UpdateStatus(ctx, id, string(shared.DocumentStatusActive), &now)
	if err != nil {
		return identitySchema.NewServiceErrorResponse[schema.DocumentResponse](http.StatusInternalServerError, "Failed to confirm upload", err)
	}

	doc.Status = string(shared.DocumentStatusActive)
	doc.ConfirmedAt = &now

	newVal, _ := json.Marshal(map[string]string{"status": "ACTIVE"})
	s.auditSvc.LogEvent(ctx, uploadedBy, string(shared.AuditEntityTypeDocument), id, "CONFIRM_UPLOAD", nil, newVal, "", "")

	return identitySchema.NewServiceResponse(schema.ToDocumentResponse(doc), http.StatusOK, "Upload confirmed")
}

func (s *documentServiceImpl) GetDownloadURL(ctx context.Context, id uuid.UUID) *identitySchema.ServiceResponse[schema.DownloadURLResponse] {
	doc, err := s.docRepo.GetByID(ctx, id)
	if err != nil {
		return identitySchema.NewServiceErrorResponse[schema.DownloadURLResponse](http.StatusNotFound, "Document not found", err)
	}

	if doc.Status != string(shared.DocumentStatusActive) {
		return identitySchema.NewServiceErrorResponse[schema.DownloadURLResponse](http.StatusBadRequest, "Document is not active", fmt.Errorf("invalid status: %s", doc.Status))
	}

	url, err := s.s3Svc.GetPresignedURL(ctx, doc.S3Key, 900)
	if err != nil {
		return identitySchema.NewServiceErrorResponse[schema.DownloadURLResponse](http.StatusInternalServerError, "Failed to generate download URL", err)
	}

	return identitySchema.NewServiceResponse(schema.DownloadURLResponse{
		DownloadURL: url,
	}, http.StatusOK, "Download URL generated")
}

func (s *documentServiceImpl) ListByEntity(ctx context.Context, entityType string, entityID uuid.UUID) *identitySchema.ServiceResponse[[]schema.DocumentResponse] {
	docs, err := s.docRepo.ListByEntity(ctx, entityType, entityID)
	if err != nil {
		return identitySchema.NewServiceErrorResponse[[]schema.DocumentResponse](http.StatusInternalServerError, "Failed to list documents", err)
	}

	responses := make([]schema.DocumentResponse, len(docs))
	for i, d := range docs {
		responses[i] = schema.ToDocumentResponse(d)
	}

	return identitySchema.NewServiceResponse(responses, http.StatusOK, "Documents retrieved")
}

func (s *documentServiceImpl) BulkRequestUploadURLs(ctx context.Context, req schema.BulkUploadURLRequest, uploadedBy uuid.UUID) *identitySchema.ServiceResponse[[]schema.UploadURLResponse] {
	results := make([]schema.UploadURLResponse, 0, len(req.Files))
	for _, file := range req.Files {
		resp := s.RequestUploadURL(ctx, file, uploadedBy)
		if resp.Error != nil {
			return identitySchema.NewServiceErrorResponse[[]schema.UploadURLResponse](resp.StatusCode, fmt.Sprintf("Failed for file %s: %s", file.FileName, resp.Message), resp.Error)
		}
		results = append(results, resp.Data)
	}
	return identitySchema.NewServiceResponse(results, http.StatusOK, "Upload URLs generated")
}

func (s *documentServiceImpl) DeleteDocument(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) *identitySchema.ServiceResponse[string] {
	doc, err := s.docRepo.GetByID(ctx, id)
	if err != nil {
		return identitySchema.NewServiceErrorResponse[string](http.StatusNotFound, "Document not found", err)
	}

	now := time.Now()
	err = s.docRepo.SoftDelete(ctx, id, now)
	if err != nil {
		return identitySchema.NewServiceErrorResponse[string](http.StatusInternalServerError, "Failed to delete document", err)
	}

	if err := s.s3Svc.Delete(ctx, doc.S3Key); err != nil {
		log.Printf("Warning: failed to delete S3 object %s: %v", doc.S3Key, err)
	}

	oldVal, _ := json.Marshal(map[string]string{"status": doc.Status})
	newVal, _ := json.Marshal(map[string]string{"status": "DELETED"})
	s.auditSvc.LogEvent(ctx, deletedBy, string(shared.AuditEntityTypeDocument), id, "DELETE", oldVal, newVal, "", "")

	return identitySchema.NewServiceResponse("Document deleted", http.StatusOK, "Document deleted successfully")
}
