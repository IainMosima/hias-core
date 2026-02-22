package entity

import (
	"encoding/json"
	"time"
	"github.com/google/uuid"
)

type PreAuthorization struct {
	ID             uuid.UUID       `json:"id"`
	PolicyID       uuid.UUID       `json:"policy_id"`
	MemberID       uuid.UUID       `json:"member_id"`
	ProviderID     uuid.UUID       `json:"provider_id"`
	AuthCode       string          `json:"auth_code,omitempty"`
	ProcedureCodes json.RawMessage `json:"procedure_codes"`
	DiagnosisCodes json.RawMessage `json:"diagnosis_codes"`
	EstimatedCost  int64           `json:"estimated_cost"`
	ApprovedAmount int64           `json:"approved_amount"`
	Status         string          `json:"status"`
	ValidityStart  *time.Time      `json:"validity_start,omitempty"`
	ValidityEnd    *time.Time      `json:"validity_end,omitempty"`
	Notes          string          `json:"notes"`
	DenialReason   string          `json:"denial_reason,omitempty"`
	ReviewedBy     uuid.UUID       `json:"reviewed_by,omitempty"`
	ReviewedAt     *time.Time      `json:"reviewed_at,omitempty"`
	CreatedBy      uuid.UUID       `json:"created_by"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}
