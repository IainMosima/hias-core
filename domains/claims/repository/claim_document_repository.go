package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	"github.com/google/uuid"
)

type ClaimDocumentRepository interface {
	Create(ctx context.Context, doc *entity.ClaimDocument) (*entity.ClaimDocument, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ClaimDocument, error)
	ListByClaim(ctx context.Context, claimID uuid.UUID) ([]*entity.ClaimDocument, error)
	SoftDelete(ctx context.Context, id uuid.UUID) (*entity.ClaimDocument, error)
}
