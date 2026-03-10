package entity

import (
	"github.com/google/uuid"
	"time"
)

type Policy struct {
	ID                uuid.UUID  `json:"id"`
	PlanID            uuid.UUID  `json:"plan_id"`
	PolicyholderName  string     `json:"policyholder_name"`
	PolicyholderEmail string     `json:"policyholder_email"`
	PolicyholderPhone string     `json:"policyholder_phone"`
	PolicyNumber      string     `json:"policy_number"`
	Status            string     `json:"status"` // DRAFT, ACTIVE, LAPSED, TERMINATED, SUSPENDED
	StartDate         time.Time  `json:"start_date"`
	EndDate           time.Time  `json:"end_date"`
	PremiumAmount     int64      `json:"premium_amount"`
	Currency          string     `json:"currency"`
	RenewedFromID     *uuid.UUID `json:"renewed_from_id,omitempty"`
	ActivatedAt       *time.Time `json:"activated_at,omitempty"`
	CreatedBy         uuid.UUID  `json:"created_by"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
