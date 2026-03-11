package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/bitbiz/hias-core/domains/document/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/document/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type documentRepository struct {
	store db.Store
}

func NewDocumentRepository(store db.Store) domainRepo.DocumentRepository {
	return &documentRepository{store: store}
}

func (r *documentRepository) Create(ctx context.Context, doc *entity.Document) (*entity.Document, error) {
	dbD, err := r.store.CreateDocument(ctx, db.CreateDocumentParams{
		EntityType:   doc.EntityType,
		EntityID:     doc.EntityID,
		DocumentType: doc.DocumentType,
		Status:       doc.Status,
		FileName:     doc.FileName,
		FileSize:     doc.FileSize,
		MimeType:     doc.MimeType,
		S3Key:        doc.S3Key,
		UploadedBy:   doc.UploadedBy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}
	return sqlcDocumentToDomain(dbD), nil
}

func (r *documentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Document, error) {
	dbD, err := r.store.GetDocumentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	return sqlcDocumentToDomain(dbD), nil
}

func (r *documentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, confirmedAt *time.Time) error {
	err := r.store.UpdateDocumentStatus(ctx, db.UpdateDocumentStatusParams{
		ID:          id,
		Status:      status,
		ConfirmedAt: timePtrToPgtypeTimestamptz(confirmedAt),
	})
	if err != nil {
		return fmt.Errorf("failed to update document status: %w", err)
	}
	return nil
}

func (r *documentRepository) SoftDelete(ctx context.Context, id uuid.UUID, deletedAt time.Time) error {
	err := r.store.SoftDeleteDocument(ctx, db.SoftDeleteDocumentParams{
		ID:        id,
		DeletedAt: timeToPgtypeTimestamptz(deletedAt),
	})
	if err != nil {
		return fmt.Errorf("failed to soft delete document: %w", err)
	}
	return nil
}

func (r *documentRepository) ListByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]*entity.Document, error) {
	dbDs, err := r.store.ListDocumentsByEntity(ctx, db.ListDocumentsByEntityParams{
		EntityType: entityType,
		EntityID:   entityID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list documents by entity: %w", err)
	}
	result := make([]*entity.Document, len(dbDs))
	for i, d := range dbDs {
		result[i] = sqlcDocumentToDomain(d)
	}
	return result, nil
}

func sqlcDocumentToDomain(d db.Document) *entity.Document {
	return &entity.Document{
		ID:           d.ID,
		EntityType:   d.EntityType,
		EntityID:     d.EntityID,
		DocumentType: d.DocumentType,
		Status:       d.Status,
		FileName:     d.FileName,
		FileSize:     d.FileSize,
		MimeType:     d.MimeType,
		S3Key:        d.S3Key,
		UploadedBy:   d.UploadedBy,
		ConfirmedAt:  pgtypeTimestamptzToTimePtr(d.ConfirmedAt),
		DeletedAt:    pgtypeTimestamptzToTimePtr(d.DeletedAt),
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}
