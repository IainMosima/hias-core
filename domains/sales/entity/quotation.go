package entity

import (
	"time"

	"github.com/google/uuid"
)

type Quotation struct {
	ID              uuid.UUID  `json:"id"`
	QuotationNumber string     `json:"quotation_number"`
	LeadID          uuid.UUID  `json:"lead_id"`
	PlanID          uuid.UUID  `json:"plan_id"`
	QuotationType   string     `json:"quotation_type"`
	Status          string     `json:"status"`
	CurrentVersion  int        `json:"current_version"`
	PolicyID        uuid.UUID  `json:"policy_id,omitempty"`
	ValidFrom       *time.Time `json:"valid_from,omitempty"`
	ValidUntil      *time.Time `json:"valid_until,omitempty"`
	ClientName      string     `json:"client_name"`
	ClientEmail     string     `json:"client_email"`
	ClientPhone     string     `json:"client_phone"`
	Currency        string     `json:"currency"`
	CreatedBy       uuid.UUID  `json:"created_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
