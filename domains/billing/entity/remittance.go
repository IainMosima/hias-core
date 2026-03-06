package entity

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type Remittance struct {
	ID                   uuid.UUID       `json:"id"`
	ProviderID           uuid.UUID       `json:"provider_id"`
	ClaimIDs             json.RawMessage `json:"claim_ids"`
	TotalAmount          int64           `json:"total_amount"`
	Currency             string          `json:"currency"`
	Status               string          `json:"status"`
	RemittanceAdviceSent bool            `json:"remittance_advice_sent"`
	PeriodStart          time.Time       `json:"period_start"`
	PeriodEnd            time.Time       `json:"period_end"`
	PaymentID            uuid.UUID       `json:"payment_id,omitempty"`
	WHTRate              float64         `json:"wht_rate"`
	WHTAmount            int64           `json:"wht_amount"`
	NetAmount            int64           `json:"net_amount"`
	CreatedBy            uuid.UUID       `json:"created_by"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}
