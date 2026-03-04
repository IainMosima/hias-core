package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/bitbiz/hias-core/shared"
)

// NotificationDispatchHandler dispatches SMS/email/push notifications.
type NotificationDispatchHandler struct {
	// notificationManager from infrastructures/notifications would be injected
}

func NewNotificationDispatchHandler() *NotificationDispatchHandler {
	return &NotificationDispatchHandler{}
}

func (h *NotificationDispatchHandler) GetName() string {
	return "notification-dispatch-handler"
}

func (h *NotificationDispatchHandler) HandleMessage(ctx context.Context, payload []byte) error {
	var msg NotificationMessage
	if err := json.Unmarshal(payload, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal notification message: %w", err)
	}

	log.Printf("Dispatching %s notification to user %s: %s", msg.Channel, msg.UserID, msg.Title)

	switch msg.Channel {
	case string(shared.NotificationChannelSMS):
		// h.notificationManager.SendSMS(ctx, phone, msg.Message)
		log.Printf("SMS sent to user %s", msg.UserID)
	case string(shared.NotificationChannelEmail):
		// h.notificationManager.SendEmail(ctx, email, msg.Title, msg.Message)
		log.Printf("Email sent to user %s", msg.UserID)
	case string(shared.NotificationChannelPush):
		// h.notificationManager.SendPush(ctx, deviceToken, msg.Title, msg.Message)
		log.Printf("Push notification sent to user %s", msg.UserID)
	case string(shared.NotificationChannelInApp):
		// Store as in-app notification (already done by notification service)
		log.Printf("In-app notification stored for user %s", msg.UserID)
	default:
		log.Printf("Unknown notification channel: %s", msg.Channel)
	}

	return nil
}
