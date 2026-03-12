package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/bitbiz/hias-core/domains/policy/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/policy/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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

func (r *policyDocumentRepository) CreateV2(ctx context.Context, doc *entity.PolicyDocument) (*entity.PolicyDocument, error) {
	dbD, err := r.store.CreatePolicyDocumentV2(ctx, db.CreatePolicyDocumentV2Params{
		PolicyID:       doc.PolicyID,
		MemberID:       uuidToPgtype(doc.MemberID),
		DocumentType:   doc.DocumentType,
		FileName:       doc.FileName,
		FileSize:       int64ToPgtypeInt8(doc.FileSize),
		S3Key:          doc.S3Key,
		GeneratedBy:    doc.GeneratedBy,
		Version:        int32(doc.Version),
		Status:         doc.Status,
		GenerationMode: doc.GenerationMode,
		EntityType:     doc.EntityType,
		EntityID:       doc.EntityID,
		MimeType:       doc.MimeType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create policy document v2: %w", err)
	}
	return sqlcPolicyDocumentToDomain(dbD), nil
}

func (r *policyDocumentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.PolicyDocument, error) {
	dbD, err := r.store.GetPolicyDocumentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get policy document: %w", err)
	}
	return rowToDomain(policyDocRow{
		ID: dbD.ID, PolicyID: dbD.PolicyID, MemberID: dbD.MemberID,
		DocumentType: dbD.DocumentType, FileName: dbD.FileName, FileSize: dbD.FileSize,
		S3Key: dbD.S3Key, GeneratedBy: dbD.GeneratedBy, CreatedAt: dbD.CreatedAt,
		Version: dbD.Version, Status: dbD.Status, GenerationMode: dbD.GenerationMode,
		EntityType: dbD.EntityType, EntityID: dbD.EntityID, SupersededBy: dbD.SupersededBy,
		ErrorMessage: dbD.ErrorMessage, UpdatedAt: dbD.UpdatedAt, MimeType: dbD.MimeType,
		GeneratedByName: dbD.GeneratedByName,
	}), nil
}

func (r *policyDocumentRepository) ListByPolicy(ctx context.Context, policyID uuid.UUID) ([]*entity.PolicyDocument, error) {
	dbDs, err := r.store.ListPolicyDocumentsByPolicy(ctx, policyID)
	if err != nil {
		return nil, fmt.Errorf("failed to list policy documents: %w", err)
	}
	result := make([]*entity.PolicyDocument, len(dbDs))
	for i, d := range dbDs {
		result[i] = rowToDomain(policyDocRow{
			ID: d.ID, PolicyID: d.PolicyID, MemberID: d.MemberID,
			DocumentType: d.DocumentType, FileName: d.FileName, FileSize: d.FileSize,
			S3Key: d.S3Key, GeneratedBy: d.GeneratedBy, CreatedAt: d.CreatedAt,
			Version: d.Version, Status: d.Status, GenerationMode: d.GenerationMode,
			EntityType: d.EntityType, EntityID: d.EntityID, SupersededBy: d.SupersededBy,
			ErrorMessage: d.ErrorMessage, UpdatedAt: d.UpdatedAt, MimeType: d.MimeType,
			GeneratedByName: d.GeneratedByName,
		})
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
		result[i] = rowToDomain(policyDocRow{
			ID: d.ID, PolicyID: d.PolicyID, MemberID: d.MemberID,
			DocumentType: d.DocumentType, FileName: d.FileName, FileSize: d.FileSize,
			S3Key: d.S3Key, GeneratedBy: d.GeneratedBy, CreatedAt: d.CreatedAt,
			Version: d.Version, Status: d.Status, GenerationMode: d.GenerationMode,
			EntityType: d.EntityType, EntityID: d.EntityID, SupersededBy: d.SupersededBy,
			ErrorMessage: d.ErrorMessage, UpdatedAt: d.UpdatedAt, MimeType: d.MimeType,
			GeneratedByName: d.GeneratedByName,
		})
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

func (r *policyDocumentRepository) GetNextVersion(ctx context.Context, entityType string, entityID uuid.UUID, docType string) (int, error) {
	v, err := r.store.GetNextDocumentVersion(ctx, db.GetNextDocumentVersionParams{
		EntityType:   entityType,
		EntityID:     entityID,
		DocumentType: docType,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get next document version: %w", err)
	}
	return int(v), nil
}

func (r *policyDocumentRepository) GetLatestByEntity(ctx context.Context, entityType string, entityID uuid.UUID, docType string) (*entity.PolicyDocument, error) {
	dbD, err := r.store.GetLatestDocumentByEntity(ctx, db.GetLatestDocumentByEntityParams{
		EntityType:   entityType,
		EntityID:     entityID,
		DocumentType: docType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get latest document: %w", err)
	}
	return rowToDomain(policyDocRow{
		ID: dbD.ID, PolicyID: dbD.PolicyID, MemberID: dbD.MemberID,
		DocumentType: dbD.DocumentType, FileName: dbD.FileName, FileSize: dbD.FileSize,
		S3Key: dbD.S3Key, GeneratedBy: dbD.GeneratedBy, CreatedAt: dbD.CreatedAt,
		Version: dbD.Version, Status: dbD.Status, GenerationMode: dbD.GenerationMode,
		EntityType: dbD.EntityType, EntityID: dbD.EntityID, SupersededBy: dbD.SupersededBy,
		ErrorMessage: dbD.ErrorMessage, UpdatedAt: dbD.UpdatedAt, MimeType: dbD.MimeType,
		GeneratedByName: dbD.GeneratedByName,
	}), nil
}

func (r *policyDocumentRepository) UpdateGenerated(ctx context.Context, id uuid.UUID, fileName string, fileSize int64, s3Key string, status string) (*entity.PolicyDocument, error) {
	dbD, err := r.store.UpdatePolicyDocumentGenerated(ctx, db.UpdatePolicyDocumentGeneratedParams{
		ID:       id,
		FileName: fileName,
		FileSize: int64ToPgtypeInt8(fileSize),
		S3Key:    s3Key,
		Status:   status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update generated document: %w", err)
	}
	return sqlcPolicyDocumentToDomain(dbD), nil
}

func (r *policyDocumentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, errorMessage string) (*entity.PolicyDocument, error) {
	dbD, err := r.store.UpdatePolicyDocumentStatus(ctx, db.UpdatePolicyDocumentStatusParams{
		ID:           id,
		Status:       status,
		ErrorMessage: stringToPgtypeText(errorMessage),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update document status: %w", err)
	}
	return sqlcPolicyDocumentToDomain(dbD), nil
}

func (r *policyDocumentRepository) Supersede(ctx context.Context, id uuid.UUID, supersededBy uuid.UUID) error {
	err := r.store.SupersedePolicyDocument(ctx, db.SupersedePolicyDocumentParams{
		ID:           id,
		SupersededBy: uuidToPgtype(supersededBy),
	})
	if err != nil {
		return fmt.Errorf("failed to supersede document: %w", err)
	}
	return nil
}

func (r *policyDocumentRepository) ConfirmUpload(ctx context.Context, id uuid.UUID, fileSize int64) (*entity.PolicyDocument, error) {
	dbD, err := r.store.ConfirmPolicyDocumentUpload(ctx, db.ConfirmPolicyDocumentUploadParams{
		ID:       id,
		FileSize: int64ToPgtypeInt8(fileSize),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to confirm document upload: %w", err)
	}
	return sqlcPolicyDocumentToDomain(dbD), nil
}

func (r *policyDocumentRepository) ListByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]*entity.PolicyDocument, error) {
	dbDs, err := r.store.ListPolicyDocumentsByEntity(ctx, db.ListPolicyDocumentsByEntityParams{
		EntityType: entityType,
		EntityID:   entityID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list documents by entity: %w", err)
	}
	result := make([]*entity.PolicyDocument, len(dbDs))
	for i, d := range dbDs {
		result[i] = rowToDomain(policyDocRow{
			ID: d.ID, PolicyID: d.PolicyID, MemberID: d.MemberID,
			DocumentType: d.DocumentType, FileName: d.FileName, FileSize: d.FileSize,
			S3Key: d.S3Key, GeneratedBy: d.GeneratedBy, CreatedAt: d.CreatedAt,
			Version: d.Version, Status: d.Status, GenerationMode: d.GenerationMode,
			EntityType: d.EntityType, EntityID: d.EntityID, SupersededBy: d.SupersededBy,
			ErrorMessage: d.ErrorMessage, UpdatedAt: d.UpdatedAt, MimeType: d.MimeType,
			GeneratedByName: d.GeneratedByName,
		})
	}
	return result, nil
}

// sqlcPolicyDocumentToDomain converts base PolicyDocument (from create/update RETURNING *).
func sqlcPolicyDocumentToDomain(d db.PolicyDocument) *entity.PolicyDocument {
	var fileSize int64
	if d.FileSize.Valid {
		fileSize = d.FileSize.Int64
	}
	var supersededBy *uuid.UUID
	if d.SupersededBy.Valid {
		id := uuid.UUID(d.SupersededBy.Bytes)
		supersededBy = &id
	}
	var errorMessage string
	if d.ErrorMessage.Valid {
		errorMessage = d.ErrorMessage.String
	}
	return &entity.PolicyDocument{
		ID:             d.ID,
		PolicyID:       d.PolicyID,
		MemberID:       pgtypeToUUID(d.MemberID),
		DocumentType:   d.DocumentType,
		FileName:       d.FileName,
		FileSize:       fileSize,
		MimeType:       d.MimeType,
		S3Key:          d.S3Key,
		GeneratedBy:    d.GeneratedBy,
		Version:        int(d.Version),
		Status:         d.Status,
		GenerationMode: d.GenerationMode,
		EntityType:     d.EntityType,
		EntityID:       d.EntityID,
		SupersededBy:   supersededBy,
		ErrorMessage:   errorMessage,
		CreatedAt:      d.CreatedAt,
		UpdatedAt:      d.UpdatedAt,
	}
}

// policyDocRow is a common intermediate for SQLC Row types that include generated_by_name.
type policyDocRow struct {
	ID              uuid.UUID
	PolicyID        uuid.UUID
	MemberID        pgtype.UUID
	DocumentType    string
	FileName        string
	FileSize        pgtype.Int8
	S3Key           string
	GeneratedBy     uuid.UUID
	CreatedAt       time.Time
	Version         int32
	Status          string
	GenerationMode  string
	EntityType      string
	EntityID        uuid.UUID
	SupersededBy    pgtype.UUID
	ErrorMessage    pgtype.Text
	UpdatedAt       time.Time
	MimeType        string
	GeneratedByName string
}

func rowToDomain(d policyDocRow) *entity.PolicyDocument {
	var fileSize int64
	if d.FileSize.Valid {
		fileSize = d.FileSize.Int64
	}
	var supersededBy *uuid.UUID
	if d.SupersededBy.Valid {
		id := uuid.UUID(d.SupersededBy.Bytes)
		supersededBy = &id
	}
	var errorMessage string
	if d.ErrorMessage.Valid {
		errorMessage = d.ErrorMessage.String
	}
	return &entity.PolicyDocument{
		ID:              d.ID,
		PolicyID:        d.PolicyID,
		MemberID:        pgtypeToUUID(d.MemberID),
		DocumentType:    d.DocumentType,
		FileName:        d.FileName,
		FileSize:        fileSize,
		MimeType:        d.MimeType,
		S3Key:           d.S3Key,
		GeneratedBy:     d.GeneratedBy,
		Version:         int(d.Version),
		Status:          d.Status,
		GenerationMode:  d.GenerationMode,
		EntityType:      d.EntityType,
		EntityID:        d.EntityID,
		SupersededBy:    supersededBy,
		ErrorMessage:    errorMessage,
		GeneratedByName: d.GeneratedByName,
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}
}
