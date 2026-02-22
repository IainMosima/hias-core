package entity

import (
	"encoding/json"
	"time"
	"github.com/google/uuid"
)

type Payment struct {
	ID              uuid.UUID       `json:"id"`
	InvoiceID       uuid.UUID       `json:"invoice_id,omitempty"`
	ClaimID         uuid.UUID       `json:"claim_id,omitempty"`
	Type            string          `json:"type"` // PREMIUM, REMITTANCE
	Amount          int64           `json:"amount"`
	Currency        string          `json:"currency"`
	Method          string          `json:"method"` // MPESA, BANK_TRANSFER
	ReferenceNumber string          `json:"reference_number"`
	Status          string          `json:"status"`
	RetryCount      int             `json:"retry_count"`
	MaxRetries      int             `json:"max_retries"`
	GatewayResponse json.RawMessage `json:"gateway_response,omitempty"`
	PaidAt          *time.Time      `json:"paid_at,omitempty"`
	ReconciledAt    *time.Time      `json:"reconciled_at,omitempty"`
	CreatedBy       uuid.UUID       `json:"created_by"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}
