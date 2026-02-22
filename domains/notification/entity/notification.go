package entity

import (
	"encoding/json"
	"time"
	"github.com/google/uuid"
)

type Notification struct {
	ID         uuid.UUID       `json:"id"`
	UserID     uuid.UUID       `json:"user_id"`
	Channel    string          `json:"channel"` // SMS, EMAIL, IN_APP, PUSH
	Type       string          `json:"type"`
	Subject    string          `json:"subject"`
	Body       string          `json:"body"`
	Metadata   json.RawMessage `json:"metadata,omitempty"`
	Status     string          `json:"status"`
	RetryCount int             `json:"retry_count"`
	MaxRetries int             `json:"max_retries"`
	SentAt     *time.Time      `json:"sent_at,omitempty"`
	ReadAt     *time.Time      `json:"read_at,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}
