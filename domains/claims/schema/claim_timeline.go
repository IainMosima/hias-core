package schema

import (
	"time"

	"github.com/google/uuid"
)

type ClaimTimelineEntry struct {
	ID              uuid.UUID `json:"id"`
	Action          string    `json:"action"`
	FromStatus      string    `json:"from_status"`
	ToStatus        string    `json:"to_status"`
	PerformedBy     uuid.UUID `json:"performed_by"`
	PerformedByName string    `json:"performed_by_name"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
}
