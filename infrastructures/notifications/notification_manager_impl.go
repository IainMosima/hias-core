package notifications

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/shared/utils"
)

type notificationManagerImpl struct {
	smsAdapter   *SMSAdapter
	emailAdapter *EmailAdapter
}

func NewNotificationManager(factory *NotificationFactory) NotificationManager {
	return &notificationManagerImpl{
		smsAdapter:   factory.GetSMSAdapter(),
		emailAdapter: factory.GetEmailAdapter(),
	}
}

func (m *notificationManagerImpl) SendSMS(ctx context.Context, phone, message string) error {
	utils.LogInfo("Sending SMS to %s", phone)
	if err := m.smsAdapter.Send(ctx, phone, message); err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	utils.LogInfo("SMS sent successfully to %s", phone)
	return nil
}

func (m *notificationManagerImpl) SendEmail(ctx context.Context, to, subject, body string) error {
	utils.LogInfo("Sending email to %s", to)
	if err := m.emailAdapter.Send(ctx, to, subject, body); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	utils.LogInfo("Email sent successfully to %s", to)
	return nil
}

func (m *notificationManagerImpl) SendPush(_ context.Context, userID, title, body string) error {
	utils.LogInfo("Push notification to user %s: %s - %s (not implemented)", userID, title, body)
	return nil
}
