package tasks

import (
	"context"
	"log"
)

// PolicyLapseTask lapses policies overdue by 30 days and terminates those overdue by 120 days.
type PolicyLapseTask struct {
	schedule string
	// policyService would be injected
}

func NewPolicyLapseTask(schedule string) *PolicyLapseTask {
	return &PolicyLapseTask{schedule: schedule}
}

func (t *PolicyLapseTask) Name() string     { return "policy-lapse" }
func (t *PolicyLapseTask) Schedule() string  { return t.schedule }

func (t *PolicyLapseTask) Execute(ctx context.Context) error {
	log.Println("Running policy lapse/termination task")

	// Step 1: Find active policies with unpaid invoices > 30 days → lapse
	// for each: t.policyService.LapsePolicy(ctx, policyID)

	// Step 2: Find lapsed policies > 120 days since lapse → terminate
	// for each: t.policyService.TerminatePolicy(ctx, policyID)

	return nil
}
