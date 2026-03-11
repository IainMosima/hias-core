package repository

import (
	"context"
	"time"

	"github.com/bitbiz/hias-core/domains/document/entity"
	"github.com/google/uuid"
)

type DocumentRepository interface {
	Create(ctx context.Context, doc *entity.Document) (*entity.Document, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Document, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, confirmedAt *time.Time) error
	SoftDelete(ctx context.Context, id uuid.UUID, deletedAt time.Time) error
	ListByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]*entity.Document, error)
}
