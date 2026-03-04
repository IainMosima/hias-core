package entity

import (
	"time"

	"github.com/google/uuid"
)

type LeadActivity struct {
	ID           uuid.UUID  `json:"id"`
	LeadID       uuid.UUID  `json:"lead_id"`
	ActivityType string     `json:"activity_type"`
	Description  string     `json:"description"`
	ScheduledAt  *time.Time `json:"scheduled_at,omitempty"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	CreatedBy    uuid.UUID  `json:"created_by"`
	CreatedAt    time.Time  `json:"created_at"`
}
