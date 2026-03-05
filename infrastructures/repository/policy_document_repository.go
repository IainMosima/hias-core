package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/policy/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/policy/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type policyDocumentRepository struct {
	store db.Store
}

func NewPolicyDocumentRepository(store db.Store) domainRepo.PolicyDocumentRepository {
	return &policyDocumentRepository{store: store}
}

func (r *policyDocumentRepository) Create(ctx context.Context, doc *entity.PolicyDocument) (*entity.PolicyDocument, error) {
	dbD, err := r.store.CreatePolicyDocument(ctx, db.CreatePolicyDocumentParams{
		PolicyID:     doc.PolicyID,
		MemberID:     uuidToPgtype(doc.MemberID),
		DocumentType: doc.DocumentType,
		FileName:     doc.FileName,
		FileSize:     int64ToPgtypeInt8(doc.FileSize),
		S3Key:        doc.S3Key,
		GeneratedBy:  doc.GeneratedBy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create policy document: %w", err)
	}
	return sqlcPolicyDocumentToDomain(dbD), nil
}

func (r *policyDocumentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.PolicyDocument, error) {
	dbD, err := r.store.GetPolicyDocumentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get policy document: %w", err)
	}
	return sqlcPolicyDocumentToDomain(dbD), nil
}

func (r *policyDocumentRepository) ListByPolicy(ctx context.Context, policyID uuid.UUID) ([]*entity.PolicyDocument, error) {
	dbDs, err := r.store.ListPolicyDocumentsByPolicy(ctx, policyID)
	if err != nil {
		return nil, fmt.Errorf("failed to list policy documents: %w", err)
	}
	result := make([]*entity.PolicyDocument, len(dbDs))
	for i, d := range dbDs {
		result[i] = sqlcPolicyDocumentToDomain(d)
	}
	return result, nil
}

func (r *policyDocumentRepository) ListByMember(ctx context.Context, memberID uuid.UUID) ([]*entity.PolicyDocument, error) {
	dbDs, err := r.store.ListPolicyDocumentsByMember(ctx, uuidToPgtype(memberID))
	if err != nil {
		return nil, fmt.Errorf("failed to list member documents: %w", err)
	}
	result := make([]*entity.PolicyDocument, len(dbDs))
	for i, d := range dbDs {
		result[i] = sqlcPolicyDocumentToDomain(d)
	}
	return result, nil
}

func (r *policyDocumentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.store.DeletePolicyDocument(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete policy document: %w", err)
	}
	return nil
}

func sqlcPolicyDocumentToDomain(d db.PolicyDocument) *entity.PolicyDocument {
	var fileSize int64
	if d.FileSize.Valid {
		fileSize = d.FileSize.Int64
	}
	return &entity.PolicyDocument{
		ID:           d.ID,
		PolicyID:     d.PolicyID,
		MemberID:     pgtypeToUUID(d.MemberID),
		DocumentType: d.DocumentType,
		FileName:     d.FileName,
		FileSize:     fileSize,
		S3Key:        d.S3Key,
		GeneratedBy:  d.GeneratedBy,
		CreatedAt:    d.CreatedAt,
	}
}
