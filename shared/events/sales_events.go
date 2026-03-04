package events

import (
	"time"

	"github.com/google/uuid"
)

const (
	EventLeadCreated        = "lead.created"
	EventLeadStatusChanged  = "lead.status_changed"
	EventQuotationCreated   = "quotation.created"
	EventQuotationIssued    = "quotation.issued"
	EventQuotationAccepted  = "quotation.accepted"
	EventQuotationConverted = "quotation.converted"
	EventApprovalRequested  = "approval.requested"
	EventApprovalGranted    = "approval.granted"
	EventApprovalRejected   = "approval.rejected"
)

type LeadCreatedEvent struct {
	LeadID     uuid.UUID `json:"lead_id"`
	LeadNumber string    `json:"lead_number"`
	Source     string    `json:"source"`
	AssignedTo uuid.UUID `json:"assigned_to"`
	Timestamp  time.Time `json:"timestamp"`
}

type LeadStatusChangedEvent struct {
	LeadID    uuid.UUID `json:"lead_id"`
	OldStatus string    `json:"old_status"`
	NewStatus string    `json:"new_status"`
	ChangedBy uuid.UUID `json:"changed_by"`
	Timestamp time.Time `json:"timestamp"`
}

type QuotationCreatedEvent struct {
	QuotationID     uuid.UUID `json:"quotation_id"`
	QuotationNumber string    `json:"quotation_number"`
	LeadID          uuid.UUID `json:"lead_id"`
	PlanID          uuid.UUID `json:"plan_id"`
	FinalPremium    int64     `json:"final_premium"`
	Timestamp       time.Time `json:"timestamp"`
}

type QuotationConvertedEvent struct {
	QuotationID     uuid.UUID `json:"quotation_id"`
	QuotationNumber string    `json:"quotation_number"`
	PolicyID        uuid.UUID `json:"policy_id"`
	PolicyNumber    string    `json:"policy_number"`
	FinalPremium    int64     `json:"final_premium"`
	MemberCount     int       `json:"member_count"`
	Timestamp       time.Time `json:"timestamp"`
}

type ApprovalRequestedEvent struct {
	QuotationID    uuid.UUID `json:"quotation_id"`
	VersionNumber  int       `json:"version_number"`
	RequestedBy    uuid.UUID `json:"requested_by"`
	EscalationRole string    `json:"escalation_role"`
	DiscountValue  int64     `json:"discount_value"`
	LoadingValue   int64     `json:"loading_value"`
	Timestamp      time.Time `json:"timestamp"`
}

type ApprovalGrantedEvent struct {
	QuotationID   uuid.UUID `json:"quotation_id"`
	VersionNumber int       `json:"version_number"`
	ApprovedBy    uuid.UUID `json:"approved_by"`
	Timestamp     time.Time `json:"timestamp"`
}

type ApprovalRejectedEvent struct {
	QuotationID   uuid.UUID `json:"quotation_id"`
	VersionNumber int       `json:"version_number"`
	RejectedBy    uuid.UUID `json:"rejected_by"`
	Reason        string    `json:"reason"`
	Timestamp     time.Time `json:"timestamp"`
}
