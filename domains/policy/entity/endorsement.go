package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Endorsement struct {
	ID                uuid.UUID       `json:"id"`
	PolicyID          uuid.UUID       `json:"policy_id"`
	EndorsementType   string          `json:"endorsement_type"`
	Status            string          `json:"status"`
	EffectiveDate     time.Time       `json:"effective_date"`
	Changes           json.RawMessage `json:"changes"`
	Reason            string          `json:"reason"`
	PremiumAdjustment int64           `json:"premium_adjustment"`
	RequestedBy       uuid.UUID       `json:"requested_by"`
	ApprovedBy        uuid.UUID       `json:"approved_by"`
	ApprovedAt        *time.Time      `json:"approved_at,omitempty"`
	AppliedAt         *time.Time      `json:"applied_at,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}
