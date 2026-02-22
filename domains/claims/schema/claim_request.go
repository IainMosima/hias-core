package schema

import "time"

type SubmitClaimRequest struct {
	PolicyID       string            `json:"policy_id" binding:"required,uuid"`
	MemberID       string            `json:"member_id" binding:"required,uuid"`
	ProviderID     string            `json:"provider_id" binding:"required,uuid"`
	PreAuthID      string            `json:"preauth_id,omitempty"`
	DiagnosisCodes []string          `json:"diagnosis_codes" binding:"required"`
	ServiceDate    time.Time         `json:"service_date" binding:"required"`
	AdmissionDate  *time.Time        `json:"admission_date,omitempty"`
	DischargeDate  *time.Time        `json:"discharge_date,omitempty"`
	Notes          string            `json:"notes"`
	LineItems      []LineItemRequest `json:"line_items" binding:"required,min=1"`
}

type LineItemRequest struct {
	ProcedureCode string `json:"procedure_code" binding:"required"`
	ProcedureName string `json:"procedure_name" binding:"required"`
	DiagnosisCode string `json:"diagnosis_code"`
	Quantity      int    `json:"quantity" binding:"required,min=1"`
	UnitPrice     int64  `json:"unit_price" binding:"required,min=1"`
}

type ReviewClaimRequest struct {
	Decision string `json:"decision" binding:"required,oneof=APPROVED REJECTED"`
	Reason   string `json:"reason"`
}
