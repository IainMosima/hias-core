package smart

import "time"

type ClaimSubmission struct {
	ClaimNumber    string              `json:"claim_number"`
	FacilityCode   string              `json:"facility_code"`
	MemberNumber   string              `json:"member_number"`
	ServiceDate    time.Time           `json:"service_date"`
	DiagnosisCodes []string            `json:"diagnosis_codes"`
	LineItems      []ClaimLineItem     `json:"line_items"`
	TotalAmount    int64               `json:"total_amount"`
}

type ClaimLineItem struct {
	ProcedureCode string `json:"procedure_code"`
	ProcedureName string `json:"procedure_name"`
	Quantity      int    `json:"quantity"`
	UnitPrice     int64  `json:"unit_price"`
	TotalPrice    int64  `json:"total_price"`
}

type ClaimSubmissionResponse struct {
	ReferenceNumber string `json:"reference_number"`
	Status          string `json:"status"`
	Message         string `json:"message"`
}

type ClaimStatusResponse struct {
	ReferenceNumber string `json:"reference_number"`
	Status          string `json:"status"`
	ApprovedAmount  int64  `json:"approved_amount,omitempty"`
	RejectionReason string `json:"rejection_reason,omitempty"`
}
