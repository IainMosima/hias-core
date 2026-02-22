package notifications

import "context"

type NotificationManager interface {
	SendSMS(ctx context.Context, phone, message string) error
	SendEmail(ctx context.Context, to, subject, body string) error
	SendPush(ctx context.Context, userID, title, body string) error
}
