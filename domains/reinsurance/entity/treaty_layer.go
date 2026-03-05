package entity

import (
	"time"

	"github.com/google/uuid"
)

type TreatyLayer struct {
	ID               uuid.UUID `json:"id"`
	TreatyID         uuid.UUID `json:"treaty_id"`
	LayerNumber      int       `json:"layer_number"`
	AttachmentPoint  int64     `json:"attachment_point"`
	LayerLimit       int64     `json:"layer_limit"`
	DeductibleAmount int64     `json:"deductible_amount"`
	PremiumRate      float64   `json:"premium_rate"`
	AggregateLimit   *int64    `json:"aggregate_limit,omitempty"`
	AggregateUsed    int64     `json:"aggregate_used"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
