package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/notification/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/notification/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type notificationRepository struct {
	store db.Store
}

func NewNotificationRepository(store db.Store) domainRepo.NotificationRepository {
	return &notificationRepository{store: store}
}

func (r *notificationRepository) Create(ctx context.Context, notification *entity.Notification) (*entity.Notification, error) {
	dbNotification, err := r.store.CreateNotification(ctx, db.CreateNotificationParams{
		UserID:   notification.UserID,
		Channel:  notification.Channel,
		Type:     notification.Type,
		Subject:  notification.Subject,
		Body:     notification.Body,
		Metadata: notification.Metadata,
		Status:   notification.Status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}
	return sqlcNotificationToDomain(dbNotification), nil
}

func (r *notificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Notification, error) {
	dbNotification, err := r.store.GetNotificationByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification by ID: %w", err)
	}
	return sqlcNotificationToDomain(dbNotification), nil
}

func (r *notificationRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Notification, error) {
	dbNotifications, err := r.store.ListNotificationsByUser(ctx, db.ListNotificationsByUserParams{
		UserID: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list notifications by user: %w", err)
	}
	notifications := make([]*entity.Notification, len(dbNotifications))
	for i, n := range dbNotifications {
		notifications[i] = sqlcNotificationToDomain(n)
	}
	return notifications, nil
}

func (r *notificationRepository) ListUnreadByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Notification, error) {
	dbNotifications, err := r.store.ListUnreadNotificationsByUser(ctx, db.ListUnreadNotificationsByUserParams{
		UserID: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list unread notifications by user: %w", err)
	}
	notifications := make([]*entity.Notification, len(dbNotifications))
	for i, n := range dbNotifications {
		notifications[i] = sqlcNotificationToDomain(n)
	}
	return notifications, nil
}

func (r *notificationRepository) CountUnreadByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	count, err := r.store.CountUnreadNotificationsByUser(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to count unread notifications by user: %w", err)
	}
	return count, nil
}

func (r *notificationRepository) MarkRead(ctx context.Context, id uuid.UUID) (*entity.Notification, error) {
	dbNotification, err := r.store.MarkNotificationRead(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to mark notification as read: %w", err)
	}
	return sqlcNotificationToDomain(dbNotification), nil
}

func (r *notificationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Notification, error) {
	dbNotification, err := r.store.UpdateNotificationStatus(ctx, db.UpdateNotificationStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update notification status: %w", err)
	}
	return sqlcNotificationToDomain(dbNotification), nil
}

func (r *notificationRepository) MarkSent(ctx context.Context, id uuid.UUID) (*entity.Notification, error) {
	dbNotification, err := r.store.MarkNotificationSent(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to mark notification as sent: %w", err)
	}
	return sqlcNotificationToDomain(dbNotification), nil
}

func (r *notificationRepository) IncrementRetry(ctx context.Context, id uuid.UUID) (*entity.Notification, error) {
	dbNotification, err := r.store.IncrementNotificationRetry(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to increment notification retry: %w", err)
	}
	return sqlcNotificationToDomain(dbNotification), nil
}

func (r *notificationRepository) GetFailedForRetry(ctx context.Context, limit int) ([]*entity.Notification, error) {
	dbNotifications, err := r.store.GetFailedNotificationsForRetry(ctx, int32(limit))
	if err != nil {
		return nil, fmt.Errorf("failed to get failed notifications for retry: %w", err)
	}
	notifications := make([]*entity.Notification, len(dbNotifications))
	for i, n := range dbNotifications {
		notifications[i] = sqlcNotificationToDomain(n)
	}
	return notifications, nil
}

func sqlcNotificationToDomain(n db.Notification) *entity.Notification {
	return &entity.Notification{
		ID:         n.ID,
		UserID:     n.UserID,
		Channel:    n.Channel,
		Type:       n.Type,
		Subject:    n.Subject,
		Body:       n.Body,
		Metadata:   n.Metadata,
		Status:     n.Status,
		RetryCount: int(n.RetryCount),
		MaxRetries: int(n.MaxRetries),
		SentAt:     pgtypeTimestamptzToTimePtr(n.SentAt),
		ReadAt:     pgtypeTimestamptzToTimePtr(n.ReadAt),
		CreatedAt:  n.CreatedAt,
		UpdatedAt:  n.UpdatedAt,
	}
}
