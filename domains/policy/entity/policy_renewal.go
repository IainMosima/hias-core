package entity

import (
	"time"

	"github.com/google/uuid"
)

type PolicyRenewal struct {
	ID                  uuid.UUID  `json:"id"`
	PolicyID            uuid.UUID  `json:"policy_id"`
	RenewedPolicyID     uuid.UUID  `json:"renewed_policy_id"`
	Status              string     `json:"status"`
	RenewalDate         time.Time  `json:"renewal_date"`
	NewPremium          int64      `json:"new_premium"`
	PremiumChangeReason string     `json:"premium_change_reason"`
	NewPlanID           uuid.UUID  `json:"new_plan_id"`
	ApprovedBy          uuid.UUID  `json:"approved_by"`
	ApprovedAt          *time.Time `json:"approved_at,omitempty"`
	CompletedAt         *time.Time `json:"completed_at,omitempty"`
	ExpiresAt           *time.Time `json:"expires_at,omitempty"`
	CreatedBy           uuid.UUID  `json:"created_by"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}
