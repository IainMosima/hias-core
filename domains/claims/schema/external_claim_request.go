package schema

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ExternalClaimRequest struct {
	IdempotencyKey  string            `json:"idempotency_key" binding:"required"`
	ExternalClaimID string            `json:"external_claim_id"`
	MemberNumber    string            `json:"member_number" binding:"required"`
	PolicyNumber    string            `json:"policy_number"`
	ProviderCode    string            `json:"provider_code"`
	DiagnosisCodes  []string          `json:"diagnosis_codes" binding:"required"`
	ServiceDate     time.Time         `json:"service_date" binding:"required"`
	AdmissionDate   *time.Time        `json:"admission_date,omitempty"`
	DischargeDate   *time.Time        `json:"discharge_date,omitempty"`
	ClaimType       string            `json:"claim_type"`
	Notes           string            `json:"notes"`
	LineItems       []LineItemRequest `json:"line_items" binding:"required,min=1"`
	Metadata        json.RawMessage   `json:"metadata,omitempty"`
}

type ExternalClaimResponse struct {
	ClaimID     uuid.UUID `json:"claim_id"`
	ClaimNumber string    `json:"claim_number"`
	Status      string    `json:"status"`
	ReceivedAt  time.Time `json:"received_at"`
}

type ExternalClaimStatusResponse struct {
	ClaimID        uuid.UUID `json:"claim_id"`
	ClaimNumber    string    `json:"claim_number"`
	Status         string    `json:"status"`
	TotalAmount    int64     `json:"total_amount"`
	ApprovedAmount int64     `json:"approved_amount"`
	UpdatedAt      time.Time `json:"updated_at"`
}
