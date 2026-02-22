package schema

import (
	"encoding/json"
	"time"
	"github.com/bitbiz/hias-core/domains/preauth/entity"
	"github.com/google/uuid"
)

type PreAuthResponse struct {
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
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

func ToPreAuthResponse(p *entity.PreAuthorization) PreAuthResponse {
	return PreAuthResponse{
		ID: p.ID, PolicyID: p.PolicyID, MemberID: p.MemberID, ProviderID: p.ProviderID,
		AuthCode: p.AuthCode, ProcedureCodes: p.ProcedureCodes, DiagnosisCodes: p.DiagnosisCodes,
		EstimatedCost: p.EstimatedCost, ApprovedAmount: p.ApprovedAmount, Status: p.Status,
		ValidityStart: p.ValidityStart, ValidityEnd: p.ValidityEnd, Notes: p.Notes,
		DenialReason: p.DenialReason, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}
