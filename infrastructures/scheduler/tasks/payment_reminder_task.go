package tasks

import (
	"context"
	"log"
)

// PaymentReminderTask sends escalating reminders for unpaid invoices.
// Day 7: gentle reminder, Day 14: second notice, Day 21: final warning.
type PaymentReminderTask struct {
	schedule string
	// billingService, notificationService would be injected
}

func NewPaymentReminderTask(schedule string) *PaymentReminderTask {
	return &PaymentReminderTask{schedule: schedule}
}

func (t *PaymentReminderTask) Name() string     { return "payment-reminder" }
func (t *PaymentReminderTask) Schedule() string  { return t.schedule }

func (t *PaymentReminderTask) Execute(ctx context.Context) error {
	log.Println("Running payment reminder task")

	// Get overdue invoices at day 7, 14, 21 thresholds
	// For each, send appropriate escalation notification
	// t.billingService.SendReminder(ctx)

	return nil
}
