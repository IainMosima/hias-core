package tasks

import (
	"context"
	"log"

	preauthService "github.com/bitbiz/hias-core/domains/preauth/service"
)

// PreAuthExpiryTask expires pre-authorizations past their validity end date.
type PreAuthExpiryTask struct {
	schedule       string
	preauthService preauthService.PreAuthService
}

func NewPreAuthExpiryTask(schedule string, preauthService preauthService.PreAuthService) *PreAuthExpiryTask {
	return &PreAuthExpiryTask{schedule: schedule, preauthService: preauthService}
}

func (t *PreAuthExpiryTask) Name() string     { return "preauth-expiry" }
func (t *PreAuthExpiryTask) Schedule() string { return t.schedule }

func (t *PreAuthExpiryTask) Execute(ctx context.Context) error {
	log.Println("Running pre-auth expiry task")

	resp := t.preauthService.ExpirePreAuths(ctx)
	if resp.Error != nil {
		return resp.Error
	}
	log.Printf("Pre-auth expiry complete: %d pre-auths expired", resp.Data)

	return nil
}
