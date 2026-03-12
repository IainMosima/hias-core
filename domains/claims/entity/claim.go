package entity

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type Claim struct {
	ID                   uuid.UUID       `json:"id"`
	ClaimNumber          string          `json:"claim_number"`
	PolicyID             uuid.UUID       `json:"policy_id"`
	MemberID             uuid.UUID       `json:"member_id"`
	ProviderID           uuid.UUID       `json:"provider_id"`
	PreAuthID            uuid.UUID       `json:"preauth_id,omitempty"`
	Status               string          `json:"status"`
	TotalAmount          int64           `json:"total_amount"`
	ApprovedAmount       int64           `json:"approved_amount"`
	CoPayAmount          int64           `json:"co_pay_amount"`
	MemberResponsibility int64           `json:"member_responsibility"`
	DiagnosisCodes       json.RawMessage `json:"diagnosis_codes"`
	ServiceDate          time.Time       `json:"service_date"`
	AdmissionDate        *time.Time      `json:"admission_date,omitempty"`
	DischargeDate        *time.Time      `json:"discharge_date,omitempty"`
	Notes                string          `json:"notes"`
	ClaimType            string          `json:"claim_type"`
	VettedAmount         *int64          `json:"vetted_amount,omitempty"`
	VettedBy             uuid.UUID       `json:"vetted_by,omitempty"`
	VettedAt             *time.Time      `json:"vetted_at,omitempty"`
	SLABreachAt          *time.Time      `json:"sla_breach_at,omitempty"`
	RejectionReason      string          `json:"rejection_reason,omitempty"`
	EscalatedTo          string          `json:"escalated_to,omitempty"`
	ClaimSource          string          `json:"claim_source"`
	IdempotencyKey       string          `json:"idempotency_key,omitempty"`
	ExternalClaimID      string          `json:"external_claim_id,omitempty"`
	SourceMetadata       json.RawMessage `json:"source_metadata,omitempty"`
	IsDraft              bool            `json:"is_draft"`
	DraftCompletedAt     *time.Time      `json:"draft_completed_at,omitempty"`
	CreatedBy            uuid.UUID       `json:"created_by"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}
