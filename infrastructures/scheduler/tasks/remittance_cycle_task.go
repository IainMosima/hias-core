package tasks

import (
	"context"
	"log"
)

// RemittanceCycleTask runs weekly provider remittance (pay approved claims to providers).
type RemittanceCycleTask struct {
	schedule string
	// remittanceService would be injected
}

func NewRemittanceCycleTask(schedule string) *RemittanceCycleTask {
	return &RemittanceCycleTask{schedule: schedule}
}

func (t *RemittanceCycleTask) Name() string     { return "remittance-cycle" }
func (t *RemittanceCycleTask) Schedule() string  { return t.schedule }

func (t *RemittanceCycleTask) Execute(ctx context.Context) error {
	log.Println("Running weekly remittance cycle")

	// resp := t.remittanceService.RunRemittanceCycle(ctx)
	// if resp.Error != nil {
	//     return resp.Error
	// }
	// log.Printf("Remittance cycle complete: %d remittances created", resp.Data)

	return nil
}
