package tasks

import (
	"context"
	"log"
)

// PaymentRetryTask retries failed payments up to max retries.
type PaymentRetryTask struct {
	schedule string
	// paymentService would be injected
}

func NewPaymentRetryTask(schedule string) *PaymentRetryTask {
	return &PaymentRetryTask{schedule: schedule}
}

func (t *PaymentRetryTask) Name() string     { return "payment-retry" }
func (t *PaymentRetryTask) Schedule() string  { return t.schedule }

func (t *PaymentRetryTask) Execute(ctx context.Context) error {
	log.Println("Running payment retry task")

	// Get failed payments where retry_count < max_retries
	// For each: t.paymentService.RetryPayment(ctx, paymentID)

	return nil
}
