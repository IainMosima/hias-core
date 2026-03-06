package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AdjudicationRule struct {
	ID         uuid.UUID       `json:"id"`
	Name       string          `json:"name"`
	RuleType   string          `json:"rule_type"`
	Parameters json.RawMessage `json:"parameters"`
	Priority   int             `json:"priority"`
	IsActive   bool            `json:"is_active"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}
