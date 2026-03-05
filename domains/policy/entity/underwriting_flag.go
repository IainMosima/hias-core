package entity

import (
	"time"

	"github.com/google/uuid"
)

type UnderwritingFlag struct {
	ID           uuid.UUID  `json:"id"`
	AssessmentID uuid.UUID  `json:"assessment_id"`
	PolicyID     uuid.UUID  `json:"policy_id"`
	MemberID     uuid.UUID  `json:"member_id"`
	FlagType     string     `json:"flag_type"`
	Severity     string     `json:"severity"`
	Details      string     `json:"details"`
	Status       string     `json:"status"`
	ResolvedBy   uuid.UUID  `json:"resolved_by"`
	ResolvedAt   *time.Time `json:"resolved_at,omitempty"`
	Resolution   string     `json:"resolution"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
