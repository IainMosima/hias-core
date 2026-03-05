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
	ClaimType      string            `json:"claim_type,omitempty"`
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

type VetClaimRequest struct {
	VettedAmount int64  `json:"vetted_amount" binding:"required,min=0"`
	Notes        string `json:"notes"`
}

type BulkSubmitClaimsRequest struct {
	Claims []SubmitClaimRequest `json:"claims" binding:"required,min=1"`
}

type CreateCaseRequest struct {
	PreAuthID         string     `json:"preauth_id" binding:"required,uuid"`
	ExpectedDischarge *time.Time `json:"expected_discharge,omitempty"`
	Diagnosis         string     `json:"diagnosis"`
	TreatingDoctor    string     `json:"treating_doctor"`
	RoomType          string     `json:"room_type"`
	EstimatedCost     int64      `json:"estimated_cost"`
	Notes             string     `json:"notes"`
}

type UpdateCaseRequest struct {
	Diagnosis      string `json:"diagnosis"`
	TreatingDoctor string `json:"treating_doctor"`
	RoomType       string `json:"room_type"`
	EstimatedCost  *int64 `json:"estimated_cost,omitempty"`
	Notes          string `json:"notes"`
}

type AdmitCaseRequest struct {
	AdmissionDate time.Time `json:"admission_date" binding:"required"`
}

type DischargeCaseRequest struct {
	ActualDischarge time.Time `json:"actual_discharge" binding:"required"`
	ActualCost      int64     `json:"actual_cost" binding:"required,min=0"`
}
