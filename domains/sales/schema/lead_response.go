package schema

import (
	"time"

	"github.com/bitbiz/hias-core/domains/sales/entity"
	"github.com/google/uuid"
)

type LeadResponse struct {
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

type LeadActivityResponse struct {
	ID           uuid.UUID  `json:"id"`
	LeadID       uuid.UUID  `json:"lead_id"`
	ActivityType string     `json:"activity_type"`
	Description  string     `json:"description"`
	ScheduledAt  *time.Time `json:"scheduled_at,omitempty"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	CreatedBy    uuid.UUID  `json:"created_by"`
	CreatedAt    time.Time  `json:"created_at"`
}

func ToLeadResponse(l *entity.Lead) LeadResponse {
	return LeadResponse{
		ID: l.ID, LeadNumber: l.LeadNumber, ContactName: l.ContactName,
		ContactEmail: l.ContactEmail, ContactPhone: l.ContactPhone,
		CompanyName: l.CompanyName, Source: l.Source, Segment: l.Segment,
		PlanType: l.PlanType, EstimatedMembers: l.EstimatedMembers,
		ExpectedPremium: l.ExpectedPremium, ClosureProbability: l.ClosureProbability,
		Currency: l.Currency, Status: l.Status, AssignedTo: l.AssignedTo,
		NextFollowUpDate: l.NextFollowUpDate, Notes: l.Notes,
		CreatedBy: l.CreatedBy, CreatedAt: l.CreatedAt, UpdatedAt: l.UpdatedAt,
	}
}

func ToLeadActivityResponse(a *entity.LeadActivity) LeadActivityResponse {
	return LeadActivityResponse{
		ID: a.ID, LeadID: a.LeadID, ActivityType: a.ActivityType,
		Description: a.Description, ScheduledAt: a.ScheduledAt,
		CompletedAt: a.CompletedAt, CreatedBy: a.CreatedBy, CreatedAt: a.CreatedAt,
	}
}
