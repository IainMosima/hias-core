package tasks

import (
	"context"
	"fmt"
	"log"

	claimRepo "github.com/bitbiz/hias-core/domains/claims/repository"
	notificationService "github.com/bitbiz/hias-core/domains/notification/service"
	"github.com/bitbiz/hias-core/shared"
)

// ClaimSLATask checks for SLA-breached and approaching-SLA claims and sends notifications.
type ClaimSLATask struct {
	schedule        string
	claimRepo       claimRepo.ClaimRepository
	notificationSvc notificationService.NotificationService
}

func NewClaimSLATask(
	schedule string,
	claimRepo claimRepo.ClaimRepository,
	notificationSvc notificationService.NotificationService,
) *ClaimSLATask {
	return &ClaimSLATask{
		schedule:        schedule,
		claimRepo:       claimRepo,
		notificationSvc: notificationSvc,
	}
}

func (t *ClaimSLATask) Name() string     { return "claim-sla-enforcement" }
func (t *ClaimSLATask) Schedule() string { return t.schedule }

func (t *ClaimSLATask) Execute(ctx context.Context) error {
	// 1. Process SLA-breached claims
	breachedClaims, err := t.claimRepo.ListSLABreached(ctx, 100, 0)
	if err != nil {
		return fmt.Errorf("failed to list SLA breached claims: %w", err)
	}

	for _, claim := range breachedClaims {
		log.Printf("SLA breached: claim %s (status: %s, breach at: %v)", claim.ClaimNumber, claim.Status, claim.SLABreachAt)

		// Send notification to claim creator about SLA breach
		if t.notificationSvc != nil {
			t.notificationSvc.Send(ctx,
				claim.CreatedBy,
				string(shared.NotificationChannelInApp),
				string(shared.NotificationTypeClaim),
				"SLA Breach Alert",
				fmt.Sprintf("Claim %s has breached its SLA deadline. Current status: %s. Please take immediate action.", claim.ClaimNumber, claim.Status),
			)
		}
	}

	// 2. Send warnings for claims approaching SLA (within 24 hours)
	approachingClaims, err := t.claimRepo.ListApproachingSLA(ctx, 100, 0)
	if err != nil {
		return fmt.Errorf("failed to list approaching SLA claims: %w", err)
	}

	for _, claim := range approachingClaims {
		log.Printf("SLA approaching: claim %s (breach at: %v)", claim.ClaimNumber, claim.SLABreachAt)

		if t.notificationSvc != nil {
			t.notificationSvc.Send(ctx,
				claim.CreatedBy,
				string(shared.NotificationChannelInApp),
				string(shared.NotificationTypeClaim),
				"SLA Warning",
				fmt.Sprintf("Claim %s is approaching its SLA deadline. Please process before %v.", claim.ClaimNumber, claim.SLABreachAt),
			)
		}
	}

	log.Printf("SLA enforcement: %d breached, %d approaching", len(breachedClaims), len(approachingClaims))

	return nil
}
