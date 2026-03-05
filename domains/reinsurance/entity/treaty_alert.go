package entity

import (
	"time"

	"github.com/google/uuid"
)

type TreatyAlert struct {
	ID             uuid.UUID  `json:"id"`
	TreatyID       uuid.UUID  `json:"treaty_id"`
	TreatyLayerID  uuid.UUID  `json:"treaty_layer_id,omitempty"`
	AlertType      string     `json:"alert_type"`
	Severity       string     `json:"severity"`
	Message        string     `json:"message"`
	ThresholdValue int64      `json:"threshold_value"`
	CurrentValue   int64      `json:"current_value"`
	IsAcknowledged bool       `json:"is_acknowledged"`
	AcknowledgedBy uuid.UUID  `json:"acknowledged_by,omitempty"`
	AcknowledgedAt *time.Time `json:"acknowledged_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}
