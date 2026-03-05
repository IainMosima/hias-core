package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/claims/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type claimDocumentRepository struct {
	store db.Store
}

func NewClaimDocumentRepository(store db.Store) domainRepo.ClaimDocumentRepository {
	return &claimDocumentRepository{store: store}
}

func (r *claimDocumentRepository) Create(ctx context.Context, doc *entity.ClaimDocument) (*entity.ClaimDocument, error) {
	dbDoc, err := r.store.CreateClaimDocument(ctx, db.CreateClaimDocumentParams{
		ClaimID:    doc.ClaimID,
		FileName:   doc.FileName,
		FileType:   doc.FileType,
		FileSize:   doc.FileSize,
		S3Key:      doc.S3Key,
		UploadedBy: doc.UploadedBy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create claim document: %w", err)
	}
	return sqlcClaimDocumentToDomain(dbDoc), nil
}

func (r *claimDocumentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ClaimDocument, error) {
	dbDoc, err := r.store.GetClaimDocumentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get claim document by ID: %w", err)
	}
	return sqlcClaimDocumentToDomain(dbDoc), nil
}

func (r *claimDocumentRepository) ListByClaim(ctx context.Context, claimID uuid.UUID) ([]*entity.ClaimDocument, error) {
	dbDocs, err := r.store.ListClaimDocumentsByClaim(ctx, claimID)
	if err != nil {
		return nil, fmt.Errorf("failed to list claim documents by claim: %w", err)
	}
	docs := make([]*entity.ClaimDocument, len(dbDocs))
	for i, d := range dbDocs {
		docs[i] = sqlcClaimDocumentToDomain(d)
	}
	return docs, nil
}

func (r *claimDocumentRepository) SoftDelete(ctx context.Context, id uuid.UUID) (*entity.ClaimDocument, error) {
	dbDoc, err := r.store.SoftDeleteClaimDocument(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to soft delete claim document: %w", err)
	}
	return sqlcClaimDocumentToDomain(dbDoc), nil
}

func sqlcClaimDocumentToDomain(d db.ClaimDocument) *entity.ClaimDocument {
	return &entity.ClaimDocument{
		ID:         d.ID,
		ClaimID:    d.ClaimID,
		FileName:   d.FileName,
		FileType:   d.FileType,
		FileSize:   d.FileSize,
		S3Key:      d.S3Key,
		UploadedBy: d.UploadedBy,
		IsDeleted:  d.IsDeleted,
		CreatedAt:  d.CreatedAt,
		UpdatedAt:  d.UpdatedAt,
	}
}
