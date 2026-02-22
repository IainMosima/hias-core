package schema

import (
	"encoding/json"
	"time"
	"github.com/bitbiz/hias-core/domains/claims/entity"
	"github.com/google/uuid"
)

type ClaimResponse struct {
	ID                   uuid.UUID             `json:"id"`
	ClaimNumber          string                `json:"claim_number"`
	PolicyID             uuid.UUID             `json:"policy_id"`
	MemberID             uuid.UUID             `json:"member_id"`
	ProviderID           uuid.UUID             `json:"provider_id"`
	Status               string                `json:"status"`
	TotalAmount          int64                 `json:"total_amount"`
	ApprovedAmount       int64                 `json:"approved_amount"`
	CoPayAmount          int64                 `json:"co_pay_amount"`
	MemberResponsibility int64                 `json:"member_responsibility"`
	DiagnosisCodes       json.RawMessage       `json:"diagnosis_codes"`
	ServiceDate          time.Time             `json:"service_date"`
	Notes                string                `json:"notes"`
	RejectionReason      string                `json:"rejection_reason,omitempty"`
	LineItems            []LineItemResponse     `json:"line_items,omitempty"`
	Decision             *AdjudicationResponse `json:"decision,omitempty"`
	FraudFlags           []FraudFlagResponse   `json:"fraud_flags,omitempty"`
	CreatedAt            time.Time             `json:"created_at"`
	UpdatedAt            time.Time             `json:"updated_at"`
}

type LineItemResponse struct {
	ID             uuid.UUID `json:"id"`
	ProcedureCode  string    `json:"procedure_code"`
	ProcedureName  string    `json:"procedure_name"`
	DiagnosisCode  string    `json:"diagnosis_code"`
	Quantity       int       `json:"quantity"`
	UnitPrice      int64     `json:"unit_price"`
	TotalPrice     int64     `json:"total_price"`
	ApprovedAmount int64     `json:"approved_amount"`
}

type AdjudicationResponse struct {
	Decision             string          `json:"decision"`
	PayableAmount        int64           `json:"payable_amount"`
	MemberResponsibility int64           `json:"member_responsibility"`
	Reasons              json.RawMessage `json:"reasons"`
	RuleResults          json.RawMessage `json:"rule_results"`
	AdjudicatedAt        time.Time       `json:"adjudicated_at"`
}

type FraudFlagResponse struct {
	ID       uuid.UUID `json:"id"`
	FlagType string    `json:"flag_type"`
	Severity string    `json:"severity"`
	Details  string    `json:"details"`
	Resolved bool      `json:"resolved"`
}

func ToClaimResponse(c *entity.Claim) ClaimResponse {
	return ClaimResponse{
		ID: c.ID, ClaimNumber: c.ClaimNumber, PolicyID: c.PolicyID,
		MemberID: c.MemberID, ProviderID: c.ProviderID, Status: c.Status,
		TotalAmount: c.TotalAmount, ApprovedAmount: c.ApprovedAmount,
		CoPayAmount: c.CoPayAmount, MemberResponsibility: c.MemberResponsibility,
		DiagnosisCodes: c.DiagnosisCodes, ServiceDate: c.ServiceDate,
		Notes: c.Notes, RejectionReason: c.RejectionReason,
		CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt,
	}
}
