package entity

import (
	"encoding/json"
	"time"
	"github.com/google/uuid"
)

type Exclusion struct {
	ID          uuid.UUID       `json:"id"`
	PlanID      uuid.UUID       `json:"plan_id"`
	Description string          `json:"description"`
	Type        string          `json:"type"` // pre_existing, cosmetic, experimental
	ICDCodes    json.RawMessage `json:"icd_codes"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
