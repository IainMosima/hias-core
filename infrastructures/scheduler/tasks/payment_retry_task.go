package tasks

import (
	"context"
	"log"

	billingService "github.com/bitbiz/hias-core/domains/billing/service"
)

// PaymentRetryTask retries failed payments up to max retries.
type PaymentRetryTask struct {
	schedule       string
	paymentService billingService.PaymentService
}

func NewPaymentRetryTask(schedule string, paymentService billingService.PaymentService) *PaymentRetryTask {
	return &PaymentRetryTask{schedule: schedule, paymentService: paymentService}
}

func (t *PaymentRetryTask) Name() string     { return "payment-retry" }
func (t *PaymentRetryTask) Schedule() string { return t.schedule }

func (t *PaymentRetryTask) Execute(ctx context.Context) error {
	log.Println("Running payment retry task — no batch retry method available yet, skipping")
	return nil
}
