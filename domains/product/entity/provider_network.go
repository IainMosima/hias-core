package entity

import (
	"time"

	"github.com/google/uuid"
)

type ProviderNetwork struct {
	ID              uuid.UUID `json:"id"`
	PlanID          uuid.UUID `json:"plan_id"`
	ProviderID      uuid.UUID `json:"provider_id"`
	BenefitCategory string    `json:"benefit_category,omitempty"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
