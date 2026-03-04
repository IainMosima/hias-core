package notification

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/domains/notification/entity"
	"github.com/bitbiz/hias-core/domains/notification/repository"
	notifSchema "github.com/bitbiz/hias-core/domains/notification/schema"
	"github.com/bitbiz/hias-core/infrastructures/queue"
	"github.com/bitbiz/hias-core/shared"
	"github.com/google/uuid"
)

type notificationServiceImpl struct {
	notifRepo    repository.NotificationRepository
	queueManager queue.QueueManager
}

func NewNotificationService(
	notifRepo repository.NotificationRepository,
	queueManager queue.QueueManager,
) *notificationServiceImpl {
	return &notificationServiceImpl{
		notifRepo:    notifRepo,
		queueManager: queueManager,
	}
}

func (s *notificationServiceImpl) Send(ctx context.Context, userID uuid.UUID, channel, notifType, subject, body string) *schema.ServiceResponse[string] {
	notif := &entity.Notification{
		UserID:     userID,
		Channel:    channel,
		Type:       notifType,
		Subject:    subject,
		Body:       body,
		Status:     string(shared.NotificationStatusPending),
		RetryCount: 0,
		MaxRetries: shared.MaxNotificationRetries,
	}

	created, err := s.notifRepo.Create(ctx, notif)
	if err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusInternalServerError, "Failed to create notification", err)
	}

	// Publish to queue for async delivery
	go func() {
		if s.queueManager != nil {
			payload, _ := json.Marshal(map[string]interface{}{
				"notification_id": created.ID.String(),
				"channel":         channel,
				"user_id":         userID.String(),
			})
			pubErr := s.queueManager.Publish(context.Background(), "notification-events", payload)
			if pubErr != nil {
				log.Printf("Failed to publish notification event: %v", pubErr)
			}
		}
	}()

	return schema.NewServiceResponse(created.ID.String(), http.StatusCreated, "Notification queued")
}

func (s *notificationServiceImpl) SendBulk(ctx context.Context, userIDs []uuid.UUID, channel, notifType, subject, body string) *schema.ServiceResponse[int] {
	sent := 0
	for _, uid := range userIDs {
		resp := s.Send(ctx, uid, channel, notifType, subject, body)
		if resp.Error == nil {
			sent++
		}
	}
	return schema.NewServiceResponse(sent, http.StatusOK, "Bulk notifications sent")
}

func (s *notificationServiceImpl) MarkRead(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[string] {
	_, err := s.notifRepo.MarkRead(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusInternalServerError, "Failed to mark notification as read", err)
	}
	return schema.NewServiceResponse("marked", http.StatusOK, "Notification marked as read")
}

func (s *notificationServiceImpl) GetUnreadCount(ctx context.Context, userID uuid.UUID) *schema.ServiceResponse[int64] {
	count, err := s.notifRepo.CountUnreadByUser(ctx, userID)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to get unread count", err)
	}
	return schema.NewServiceResponse(count, http.StatusOK, "Unread count retrieved")
}

func (s *notificationServiceImpl) ListByUser(ctx context.Context, userID uuid.UUID, page, pageSize int) *schema.ServiceResponse[interface{}] {
	offset := (page - 1) * pageSize
	notifications, err := s.notifRepo.ListByUser(ctx, userID, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[interface{}](http.StatusInternalServerError, "Failed to list notifications", err)
	}

	responses := make([]notifSchema.NotificationResponse, len(notifications))
	for i, n := range notifications {
		responses[i] = notifSchema.ToNotificationResponse(n)
	}

	result := map[string]interface{}{
		"notifications": responses,
		"page":          page,
		"page_size":     pageSize,
	}

	return schema.NewServiceResponse[interface{}](result, http.StatusOK, "Notifications retrieved")
}

func (s *notificationServiceImpl) RetryFailed(ctx context.Context) *schema.ServiceResponse[int] {
	failed, err := s.notifRepo.GetFailedForRetry(ctx, 50)
	if err != nil {
		return schema.NewServiceErrorResponse[int](http.StatusInternalServerError, "Failed to get failed notifications", err)
	}

	retried := 0
	for _, n := range failed {
		if n.RetryCount < n.MaxRetries {
			_, incrErr := s.notifRepo.IncrementRetry(ctx, n.ID)
			if incrErr == nil {
				retried++
			}
		}
	}

	return schema.NewServiceResponse(retried, http.StatusOK, "Failed notifications retried")
}
