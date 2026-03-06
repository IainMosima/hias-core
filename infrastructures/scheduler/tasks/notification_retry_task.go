package tasks

import (
	"context"
	"log"

	notifService "github.com/bitbiz/hias-core/domains/notification/service"
)

// NotificationRetryTask retries failed SMS/email notifications every 30 minutes.
type NotificationRetryTask struct {
	schedule            string
	notificationService notifService.NotificationService
}

func NewNotificationRetryTask(schedule string, notificationService notifService.NotificationService) *NotificationRetryTask {
	return &NotificationRetryTask{schedule: schedule, notificationService: notificationService}
}

func (t *NotificationRetryTask) Name() string     { return "notification-retry" }
func (t *NotificationRetryTask) Schedule() string { return t.schedule }

func (t *NotificationRetryTask) Execute(ctx context.Context) error {
	log.Println("Running notification retry task")

	resp := t.notificationService.RetryFailed(ctx)
	if resp.Error != nil {
		return resp.Error
	}
	log.Printf("Retried %d failed notifications", resp.Data)

	return nil
}
