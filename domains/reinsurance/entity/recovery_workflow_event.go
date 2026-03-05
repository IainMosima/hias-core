package entity

import (
	"time"

	"github.com/google/uuid"
)

type RecoveryWorkflowEvent struct {
	ID          uuid.UUID `json:"id"`
	RecoveryID  uuid.UUID `json:"recovery_id"`
	FromStatus  string    `json:"from_status"`
	ToStatus    string    `json:"to_status"`
	EventType   string    `json:"event_type"`
	Notes       string    `json:"notes,omitempty"`
	PerformedBy uuid.UUID `json:"performed_by"`
	CreatedAt   time.Time `json:"created_at"`
}
