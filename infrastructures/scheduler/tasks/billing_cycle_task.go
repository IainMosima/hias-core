package tasks

import (
	"context"
	"log"
)

// BillingCycleTask generates invoices for all active policies on the 1st of each month.
type BillingCycleTask struct {
	schedule string
	// billingService would be injected
}

func NewBillingCycleTask(schedule string) *BillingCycleTask {
	return &BillingCycleTask{schedule: schedule}
}

func (t *BillingCycleTask) Name() string     { return "billing-cycle" }
func (t *BillingCycleTask) Schedule() string  { return t.schedule }

func (t *BillingCycleTask) Execute(ctx context.Context) error {
	log.Println("Running billing cycle — generating invoices for active policies")

	// resp := t.billingService.RunBillingCycle(ctx)
	// if resp.Error != nil {
	//     return resp.Error
	// }
	// log.Printf("Billing cycle complete: %d invoices generated", resp.Data)

	return nil
}
