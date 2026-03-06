package entity

import (
	"time"

	"github.com/google/uuid"
)

type EscalationRule struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	ConditionType   string    `json:"condition_type"`
	ThresholdAmount int64     `json:"threshold_amount"`
	EscalationRole  string    `json:"escalation_role"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
