package entity

import (
	"time"

	"github.com/google/uuid"
)

type Treaty struct {
	ID             uuid.UUID `json:"id"`
	TreatyNumber   string    `json:"treaty_number"`
	Name           string    `json:"name"`
	TreatyType     string    `json:"treaty_type"`
	Status         string    `json:"status"`
	EffectiveDate  time.Time `json:"effective_date"`
	ExpiryDate     time.Time `json:"expiry_date"`
	RetentionLimit int64     `json:"retention_limit"`
	Currency       string    `json:"currency"`
	Notes          string    `json:"notes,omitempty"`
	CreatedBy      uuid.UUID `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
