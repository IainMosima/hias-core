package tasks

import (
	"context"
	"log"

	billingService "github.com/bitbiz/hias-core/domains/billing/service"
)

// PaymentReminderTask sends escalating reminders for unpaid invoices.
// Day 7: gentle reminder, Day 14: second notice, Day 21: final warning.
type PaymentReminderTask struct {
	schedule       string
	billingService billingService.BillingService
}

func NewPaymentReminderTask(schedule string, billingService billingService.BillingService) *PaymentReminderTask {
	return &PaymentReminderTask{schedule: schedule, billingService: billingService}
}

func (t *PaymentReminderTask) Name() string     { return "payment-reminder" }
func (t *PaymentReminderTask) Schedule() string { return t.schedule }

func (t *PaymentReminderTask) Execute(ctx context.Context) error {
	log.Println("Running payment reminder task")

	resp := t.billingService.HandleOverdue(ctx)
	if resp.Error != nil {
		return resp.Error
	}
	log.Printf("Payment reminder complete: %d overdue invoices handled", resp.Data)

	return nil
}
