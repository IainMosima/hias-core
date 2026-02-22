package notification

import (
	"context"
	"log"
	"net/http"

	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/domains/notification/entity"
	"github.com/bitbiz/hias-core/domains/notification/repository"
	notifSchema "github.com/bitbiz/hias-core/domains/notification/schema"
	"github.com/bitbiz/hias-core/domains/notification/service"
	"github.com/bitbiz/hias-core/infrastructures/notifications"
	"github.com/google/uuid"
)

type notificationServiceImpl struct {
	notificationRepo    repository.NotificationRepository
	notificationManager notifications.NotificationManager
}

func NewNotificationService(
	notificationRepo repository.NotificationRepository,
	notificationManager notifications.NotificationManager,
) service.NotificationService {
	return &notificationServiceImpl{
		notificationRepo:    notificationRepo,
		notificationManager: notificationManager,
	}
}

func (s *notificationServiceImpl) Send(ctx context.Context, userID uuid.UUID, channel, notifType, subject, body string) *schema.ServiceResponse[string] {
	notification := &entity.Notification{
		UserID:     userID,
		Channel:    channel,
		Type:       notifType,
		Subject:    subject,
		Body:       body,
		Status:     "PENDING",
		RetryCount: 0,
		MaxRetries: 3,
	}

	created, err := s.notificationRepo.Create(ctx, notification)
	if err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusInternalServerError, "Failed to create notification record", err)
	}

	// Dispatch via notification manager
	dispatchErr := s.dispatch(ctx, channel, subject, body, userID.String())
	if dispatchErr != nil {
		// Mark as failed but still return success for record creation
		_, _ = s.notificationRepo.UpdateStatus(ctx, created.ID, "FAILED")
		return schema.NewServiceErrorResponse[string](http.StatusInternalServerError, "Notification created but dispatch failed", dispatchErr)
	}

	// Mark as sent
	_, _ = s.notificationRepo.MarkSent(ctx, created.ID)

	return schema.NewServiceResponse(created.ID.String(), http.StatusCreated, "Notification sent successfully")
}

func (s *notificationServiceImpl) SendBulk(ctx context.Context, userIDs []uuid.UUID, channel, notifType, subject, body string) *schema.ServiceResponse[int] {
	sentCount := 0

	for _, userID := range userIDs {
		notification := &entity.Notification{
			UserID:     userID,
			Channel:    channel,
			Type:       notifType,
			Subject:    subject,
			Body:       body,
			Status:     "PENDING",
			RetryCount: 0,
			MaxRetries: 3,
		}

		created, err := s.notificationRepo.Create(ctx, notification)
		if err != nil {
			log.Printf("Failed to create notification for user %s: %v", userID, err)
			continue
		}

		dispatchErr := s.dispatch(ctx, channel, subject, body, userID.String())
		if dispatchErr != nil {
			_, _ = s.notificationRepo.UpdateStatus(ctx, created.ID, "FAILED")
			log.Printf("Failed to dispatch notification for user %s: %v", userID, dispatchErr)
			continue
		}

		_, _ = s.notificationRepo.MarkSent(ctx, created.ID)
		sentCount++
	}

	return schema.NewServiceResponse(sentCount, http.StatusOK, "Bulk notifications processed")
}

func (s *notificationServiceImpl) MarkRead(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[string] {
	_, err := s.notificationRepo.MarkRead(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusInternalServerError, "Failed to mark notification as read", err)
	}

	return schema.NewServiceResponse(id.String(), http.StatusOK, "Notification marked as read")
}

func (s *notificationServiceImpl) GetUnreadCount(ctx context.Context, userID uuid.UUID) *schema.ServiceResponse[int64] {
	count, err := s.notificationRepo.CountUnreadByUser(ctx, userID)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to get unread count", err)
	}

	return schema.NewServiceResponse(count, http.StatusOK, "Unread count retrieved")
}

func (s *notificationServiceImpl) ListByUser(ctx context.Context, userID uuid.UUID, page, pageSize int) *schema.ServiceResponse[interface{}] {
	offset := (page - 1) * pageSize
	notifs, err := s.notificationRepo.ListByUser(ctx, userID, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[interface{}](http.StatusInternalServerError, "Failed to list notifications", err)
	}

	responses := notifSchema.ToNotificationResponseList(notifs)
	return schema.NewServiceResponse[interface{}](responses, http.StatusOK, "Notifications retrieved")
}

func (s *notificationServiceImpl) RetryFailed(ctx context.Context) *schema.ServiceResponse[int] {
	failedNotifs, err := s.notificationRepo.GetFailedForRetry(ctx, 50)
	if err != nil {
		return schema.NewServiceErrorResponse[int](http.StatusInternalServerError, "Failed to get failed notifications", err)
	}

	retriedCount := 0
	for _, notif := range failedNotifs {
		dispatchErr := s.dispatch(ctx, notif.Channel, notif.Subject, notif.Body, notif.UserID.String())
		if dispatchErr != nil {
			_, _ = s.notificationRepo.IncrementRetry(ctx, notif.ID)
			log.Printf("Retry failed for notification %s: %v", notif.ID, dispatchErr)
			continue
		}

		_, _ = s.notificationRepo.MarkSent(ctx, notif.ID)
		retriedCount++
	}

	return schema.NewServiceResponse(retriedCount, http.StatusOK, "Failed notifications retried")
}

// dispatch sends a notification through the appropriate channel via the notification manager.
func (s *notificationServiceImpl) dispatch(ctx context.Context, channel, subject, body, userID string) error {
	switch channel {
	case "SMS":
		return s.notificationManager.SendSMS(ctx, userID, body)
	case "EMAIL":
		return s.notificationManager.SendEmail(ctx, userID, subject, body)
	case "PUSH":
		return s.notificationManager.SendPush(ctx, userID, subject, body)
	case "IN_APP":
		// IN_APP notifications are stored in DB and delivered via SSE/WebSocket.
		// No external dispatch needed; the record itself is the notification.
		return nil
	default:
		return nil
	}
}
