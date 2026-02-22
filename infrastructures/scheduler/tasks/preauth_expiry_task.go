package tasks

import (
	"context"
	"log"
)

// PreAuthExpiryTask expires pre-authorizations past their validity end date.
type PreAuthExpiryTask struct {
	schedule string
	// preauthService would be injected
}

func NewPreAuthExpiryTask(schedule string) *PreAuthExpiryTask {
	return &PreAuthExpiryTask{schedule: schedule}
}

func (t *PreAuthExpiryTask) Name() string     { return "preauth-expiry" }
func (t *PreAuthExpiryTask) Schedule() string  { return t.schedule }

func (t *PreAuthExpiryTask) Execute(ctx context.Context) error {
	log.Println("Running pre-auth expiry task")

	// Find approved pre-auths where validity_end < now
	// For each: update status to EXPIRED

	return nil
}
