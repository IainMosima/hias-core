package entity

import (
	"time"

	"github.com/google/uuid"
)

type TreatyParticipant struct {
	ID              uuid.UUID `json:"id"`
	TreatyID        uuid.UUID `json:"treaty_id"`
	ReinsurerName   string    `json:"reinsurer_name"`
	SharePercentage float64   `json:"share_percentage"`
	CommissionRate  float64   `json:"commission_rate"`
	IsLead          bool      `json:"is_lead"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
