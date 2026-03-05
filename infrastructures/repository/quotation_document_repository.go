package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/sales/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/sales/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type quotationDocumentRepository struct {
	store db.Store
}

func NewQuotationDocumentRepository(store db.Store) domainRepo.QuotationDocumentRepository {
	return &quotationDocumentRepository{store: store}
}

func (r *quotationDocumentRepository) Create(ctx context.Context, doc *entity.QuotationDocument) (*entity.QuotationDocument, error) {
	dbDoc, err := r.store.CreateQuotationDocument(ctx, db.CreateQuotationDocumentParams{
		QuotationID:    doc.QuotationID,
		VersionNumber:  pgtype.Int4{Int32: int32(doc.VersionNumber), Valid: doc.VersionNumber > 0},
		FileName:       doc.FileName,
		FileType:       doc.FileType,
		FileSize:       doc.FileSize,
		S3Key:          doc.S3Key,
		UploadedBy:     uuidToPgtype(doc.UploadedBy),
		CanEditRoles:   doc.CanEditRoles,
		CanDeleteRoles: doc.CanDeleteRoles,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create quotation document: %w", err)
	}
	return sqlcQuotationDocumentToDomain(dbDoc), nil
}

func (r *quotationDocumentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.QuotationDocument, error) {
	dbDoc, err := r.store.GetQuotationDocumentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotation document by ID: %w", err)
	}
	return sqlcQuotationDocumentToDomain(dbDoc), nil
}

func (r *quotationDocumentRepository) ListByQuotation(ctx context.Context, quotationID uuid.UUID) ([]*entity.QuotationDocument, error) {
	dbDocs, err := r.store.ListQuotationDocumentsByQuotation(ctx, quotationID)
	if err != nil {
		return nil, fmt.Errorf("failed to list quotation documents: %w", err)
	}
	docs := make([]*entity.QuotationDocument, len(dbDocs))
	for i, d := range dbDocs {
		docs[i] = sqlcQuotationDocumentToDomain(d)
	}
	return docs, nil
}

func (r *quotationDocumentRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	err := r.store.SoftDeleteQuotationDocument(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to soft delete quotation document: %w", err)
	}
	return nil
}

func (r *quotationDocumentRepository) Update(ctx context.Context, doc *entity.QuotationDocument) (*entity.QuotationDocument, error) {
	dbDoc, err := r.store.UpdateQuotationDocument(ctx, db.UpdateQuotationDocumentParams{
		ID:             doc.ID,
		FileName:       doc.FileName,
		CanEditRoles:   doc.CanEditRoles,
		CanDeleteRoles: doc.CanDeleteRoles,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update quotation document: %w", err)
	}
	return sqlcQuotationDocumentToDomain(dbDoc), nil
}

func sqlcQuotationDocumentToDomain(d db.QuotationDocument) *entity.QuotationDocument {
	versionNumber := 0
	if d.VersionNumber.Valid {
		versionNumber = int(d.VersionNumber.Int32)
	}
	return &entity.QuotationDocument{
		ID:             d.ID,
		QuotationID:    d.QuotationID,
		VersionNumber:  versionNumber,
		FileName:       d.FileName,
		FileType:       d.FileType,
		FileSize:       d.FileSize,
		S3Key:          d.S3Key,
		UploadedBy:     pgtypeToUUID(d.UploadedBy),
		CanEditRoles:   d.CanEditRoles,
		CanDeleteRoles: d.CanDeleteRoles,
		IsDeleted:      d.IsDeleted,
		CreatedAt:      d.CreatedAt,
		UpdatedAt:      d.UpdatedAt,
	}
}
