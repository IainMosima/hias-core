package entity

import (
	"github.com/google/uuid"
	"time"
)

type Refund struct {
	ID           uuid.UUID  `json:"id"`
	PolicyID     uuid.UUID  `json:"policy_id"`
	CreditNoteID uuid.UUID  `json:"credit_note_id,omitempty"`
	Amount       int64      `json:"amount"`
	Currency     string     `json:"currency"`
	Status       string     `json:"status"`
	Reason       string     `json:"reason"`
	ApprovedBy   uuid.UUID  `json:"approved_by,omitempty"`
	ApprovedAt   *time.Time `json:"approved_at,omitempty"`
	ProcessedAt  *time.Time `json:"processed_at,omitempty"`
	CreatedBy    uuid.UUID  `json:"created_by,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
