package repository

import (
	"context"
	"github.com/bitbiz/hias-core/domains/audit/entity"
	"github.com/google/uuid"
)

// AuditRepository is APPEND ONLY - no Update or Delete methods
type AuditRepository interface {
	Create(ctx context.Context, event *entity.AuditEvent) (*entity.AuditEvent, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.AuditEvent, error)
	List(ctx context.Context, limit, offset int) ([]*entity.AuditEvent, error)
	ListByEntity(ctx context.Context, entityType string, entityID uuid.UUID, limit, offset int) ([]*entity.AuditEvent, error)
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.AuditEvent, error)
	Count(ctx context.Context) (int64, error)
	CountByEntity(ctx context.Context, entityType string, entityID uuid.UUID) (int64, error)
	CountByUser(ctx context.Context, userID uuid.UUID) (int64, error)
}
