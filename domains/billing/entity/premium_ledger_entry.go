package entity

import (
	"github.com/google/uuid"
	"time"
)

type PremiumLedgerEntry struct {
	ID              uuid.UUID `json:"id"`
	PolicyID        uuid.UUID `json:"policy_id"`
	EntryType       string    `json:"entry_type"`
	Amount          int64     `json:"amount"`
	Description     string    `json:"description"`
	ReferenceNumber string    `json:"reference_number"`
	EffectiveDate   time.Time `json:"effective_date"`
	BalanceAfter    int64     `json:"balance_after"`
	CreatedBy       uuid.UUID `json:"created_by,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
