package entity

import (
	"github.com/google/uuid"
	"time"
)

type Plan struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Type             string    `json:"type"`    // individual, group
	Segment          string    `json:"segment"` // retail, corporate, sme
	BasePremium      int64     `json:"base_premium"`
	PremiumFrequency string    `json:"premium_frequency"` // monthly, quarterly, semi_annual, annual
	Currency         string    `json:"currency"`
	Status           string    `json:"status"`
	Description      string    `json:"description"`
	CreatedBy        uuid.UUID `json:"created_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
