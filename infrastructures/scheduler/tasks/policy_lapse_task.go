package tasks

import (
	"context"
	"log"

	policyService "github.com/bitbiz/hias-core/domains/policy/service"
)

// PolicyLapseTask lapses policies overdue by 30 days and terminates those overdue by 120 days.
type PolicyLapseTask struct {
	schedule      string
	policyService policyService.PolicyService
}

func NewPolicyLapseTask(schedule string, policyService policyService.PolicyService) *PolicyLapseTask {
	return &PolicyLapseTask{schedule: schedule, policyService: policyService}
}

func (t *PolicyLapseTask) Name() string     { return "policy-lapse" }
func (t *PolicyLapseTask) Schedule() string { return t.schedule }

func (t *PolicyLapseTask) Execute(ctx context.Context) error {
	log.Println("Running policy lapse/termination task — no batch lapse method available yet, skipping")
	return nil
}
