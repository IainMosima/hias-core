package tasks

import (
	"context"
	"log"

	billingService "github.com/bitbiz/hias-core/domains/billing/service"
)

// ReconciliationTask matches confirmed payments to bank statements daily.
type ReconciliationTask struct {
	schedule       string
	paymentService billingService.PaymentService
}

func NewReconciliationTask(schedule string, paymentService billingService.PaymentService) *ReconciliationTask {
	return &ReconciliationTask{schedule: schedule, paymentService: paymentService}
}

func (t *ReconciliationTask) Name() string     { return "reconciliation" }
func (t *ReconciliationTask) Schedule() string { return t.schedule }

func (t *ReconciliationTask) Execute(ctx context.Context) error {
	log.Println("Running daily reconciliation task — no batch reconcile method available yet, skipping")
	return nil
}
