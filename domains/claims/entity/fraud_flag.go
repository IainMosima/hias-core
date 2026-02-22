package entity

import (
	"time"
	"github.com/google/uuid"
)

type FraudFlag struct {
	ID               uuid.UUID  `json:"id"`
	ClaimID          uuid.UUID  `json:"claim_id"`
	FlagType         string     `json:"flag_type"` // DUPLICATE, FREQUENCY, AMOUNT_THRESHOLD
	Severity         string     `json:"severity"`
	Details          string     `json:"details"`
	Resolved         bool       `json:"resolved"`
	ResolvedBy       uuid.UUID  `json:"resolved_by,omitempty"`
	ResolvedAt       *time.Time `json:"resolved_at,omitempty"`
	ReferenceClaimID uuid.UUID  `json:"reference_claim_id,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}
