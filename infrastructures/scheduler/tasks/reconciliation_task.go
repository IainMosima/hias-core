package tasks

import (
	"context"
	"log"
)

// ReconciliationTask matches confirmed payments to bank statements daily.
type ReconciliationTask struct {
	schedule string
	// paymentService would be injected
}

func NewReconciliationTask(schedule string) *ReconciliationTask {
	return &ReconciliationTask{schedule: schedule}
}

func (t *ReconciliationTask) Name() string     { return "reconciliation" }
func (t *ReconciliationTask) Schedule() string  { return t.schedule }

func (t *ReconciliationTask) Execute(ctx context.Context) error {
	log.Println("Running daily reconciliation task")

	// Get confirmed payments not yet reconciled
	// Match against bank statement entries
	// For matched: t.paymentService.ReconcilePayment(ctx, paymentID)

	return nil
}
