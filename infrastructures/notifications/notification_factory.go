package notifications

import (
	"github.com/aws/aws-sdk-go-v2/service/ses"
)

type NotificationFactory struct {
	smsAdapter   *SMSAdapter
	emailAdapter *EmailAdapter
}

func NewNotificationFactory(
	smsAPIKey, smsUsername, smsSenderID string,
	sesClient *ses.Client, fromEmail string,
) *NotificationFactory {
	return &NotificationFactory{
		smsAdapter:   NewSMSAdapter(smsAPIKey, smsUsername, smsSenderID),
		emailAdapter: NewEmailAdapter(sesClient, fromEmail),
	}
}

func (f *NotificationFactory) GetSMSAdapter() *SMSAdapter {
	return f.smsAdapter
}

func (f *NotificationFactory) GetEmailAdapter() *EmailAdapter {
	return f.emailAdapter
}
