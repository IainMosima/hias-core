package entity

import (
	"time"

	"github.com/google/uuid"
)

type Lead struct {
	ID                 uuid.UUID  `json:"id"`
	LeadNumber         string     `json:"lead_number"`
	ContactName        string     `json:"contact_name"`
	ContactEmail       string     `json:"contact_email"`
	ContactPhone       string     `json:"contact_phone"`
	CompanyName        string     `json:"company_name"`
	Source             string     `json:"source"`
	Segment            string     `json:"segment"`
	PlanType           string     `json:"plan_type"`
	EstimatedMembers   int        `json:"estimated_members"`
	ExpectedPremium    int64      `json:"expected_premium"`
	ClosureProbability int        `json:"closure_probability"`
	Currency           string     `json:"currency"`
	Status             string     `json:"status"`
	AssignedTo         uuid.UUID  `json:"assigned_to"`
	NextFollowUpDate   *time.Time `json:"next_follow_up_date,omitempty"`
	Notes              string     `json:"notes"`
	CreatedBy          uuid.UUID  `json:"created_by"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}
