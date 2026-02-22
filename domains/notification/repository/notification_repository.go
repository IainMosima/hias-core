package repository

import (
	"context"
	"github.com/bitbiz/hias-core/domains/notification/entity"
	"github.com/google/uuid"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *entity.Notification) (*entity.Notification, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Notification, error)
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Notification, error)
	ListUnreadByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Notification, error)
	CountUnreadByUser(ctx context.Context, userID uuid.UUID) (int64, error)
	MarkRead(ctx context.Context, id uuid.UUID) (*entity.Notification, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Notification, error)
	MarkSent(ctx context.Context, id uuid.UUID) (*entity.Notification, error)
	IncrementRetry(ctx context.Context, id uuid.UUID) (*entity.Notification, error)
	GetFailedForRetry(ctx context.Context, limit int) ([]*entity.Notification, error)
}
