package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/policy/entity"
	"github.com/google/uuid"
)

type PolicyDocumentRepository interface {
	Create(ctx context.Context, doc *entity.PolicyDocument) (*entity.PolicyDocument, error)
	CreateV2(ctx context.Context, doc *entity.PolicyDocument) (*entity.PolicyDocument, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.PolicyDocument, error)
	ListByPolicy(ctx context.Context, policyID uuid.UUID) ([]*entity.PolicyDocument, error)
	ListByMember(ctx context.Context, memberID uuid.UUID) ([]*entity.PolicyDocument, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetNextVersion(ctx context.Context, entityType string, entityID uuid.UUID, docType string) (int, error)
	GetLatestByEntity(ctx context.Context, entityType string, entityID uuid.UUID, docType string) (*entity.PolicyDocument, error)
	UpdateGenerated(ctx context.Context, id uuid.UUID, fileName string, fileSize int64, s3Key string, status string) (*entity.PolicyDocument, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, errorMessage string) (*entity.PolicyDocument, error)
	Supersede(ctx context.Context, id uuid.UUID, supersededBy uuid.UUID) error
	ListByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]*entity.PolicyDocument, error)
	ConfirmUpload(ctx context.Context, id uuid.UUID, fileSize int64) (*entity.PolicyDocument, error)
}
