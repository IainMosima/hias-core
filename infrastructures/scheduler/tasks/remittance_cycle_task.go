package tasks

import (
	"context"
	"log"

	billingService "github.com/bitbiz/hias-core/domains/billing/service"
)

// RemittanceCycleTask runs weekly provider remittance (pay approved claims to providers).
type RemittanceCycleTask struct {
	schedule          string
	remittanceService billingService.RemittanceService
}

func NewRemittanceCycleTask(schedule string, remittanceService billingService.RemittanceService) *RemittanceCycleTask {
	return &RemittanceCycleTask{schedule: schedule, remittanceService: remittanceService}
}

func (t *RemittanceCycleTask) Name() string     { return "remittance-cycle" }
func (t *RemittanceCycleTask) Schedule() string { return t.schedule }

func (t *RemittanceCycleTask) Execute(ctx context.Context) error {
	log.Println("Running weekly remittance cycle")

	resp := t.remittanceService.RunRemittanceCycle(ctx)
	if resp.Error != nil {
		return resp.Error
	}
	log.Printf("Remittance cycle complete: %d remittances created", resp.Data)

	return nil
}
