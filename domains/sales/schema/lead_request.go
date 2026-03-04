package schema

import "time"

type CreateLeadRequest struct {
	ContactName        string     `json:"contact_name" binding:"required"`
	ContactEmail       string     `json:"contact_email"`
	ContactPhone       string     `json:"contact_phone"`
	CompanyName        string     `json:"company_name"`
	Source             string     `json:"source" binding:"required,oneof=direct referral web agent broker"`
	Segment            string     `json:"segment" binding:"required,oneof=retail corporate sme"`
	PlanType           string     `json:"plan_type" binding:"required,oneof=individual group"`
	EstimatedMembers   int        `json:"estimated_members"`
	ExpectedPremium    int64      `json:"expected_premium"`
	ClosureProbability int        `json:"closure_probability"`
	NextFollowUpDate   *time.Time `json:"next_follow_up_date"`
	Notes              string     `json:"notes"`
}

type UpdateLeadRequest struct {
	ContactName        string     `json:"contact_name"`
	ContactEmail       string     `json:"contact_email"`
	ContactPhone       string     `json:"contact_phone"`
	CompanyName        string     `json:"company_name"`
	Source             string     `json:"source"`
	Segment            string     `json:"segment"`
	PlanType           string     `json:"plan_type"`
	EstimatedMembers   *int       `json:"estimated_members"`
	ExpectedPremium    *int64     `json:"expected_premium"`
	ClosureProbability *int       `json:"closure_probability"`
	AssignedTo         string     `json:"assigned_to"`
	NextFollowUpDate   *time.Time `json:"next_follow_up_date"`
	Notes              string     `json:"notes"`
}

type UpdateLeadStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=NEW CONTACTED QUALIFIED PROPOSAL_SENT NEGOTIATION WON LOST DORMANT"`
}

type CreateLeadActivityRequest struct {
	ActivityType string     `json:"activity_type" binding:"required,oneof=call email meeting note follow_up"`
	Description  string     `json:"description"`
	ScheduledAt  *time.Time `json:"scheduled_at"`
	CompletedAt  *time.Time `json:"completed_at"`
}
