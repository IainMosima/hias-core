package entity

import (
	"time"

	"github.com/google/uuid"
)

type Cession struct {
	ID               uuid.UUID `json:"id"`
	CessionNumber    string    `json:"cession_number"`
	TreatyID         uuid.UUID `json:"treaty_id"`
	PolicyID         uuid.UUID `json:"policy_id"`
	TreatyLayerID    uuid.UUID `json:"treaty_layer_id,omitempty"`
	CessionType      string    `json:"cession_type"`
	GrossAmount      int64     `json:"gross_amount"`
	CededAmount      int64     `json:"ceded_amount"`
	RetainedAmount   int64     `json:"retained_amount"`
	CommissionAmount int64     `json:"commission_amount"`
	SharePercentage  float64   `json:"share_percentage"`
	Status           string    `json:"status"`
	CreatedBy        uuid.UUID `json:"created_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
