package schema

type SubmitPreAuthRequest struct {
	PolicyID       string   `json:"policy_id" binding:"required,uuid"`
	MemberID       string   `json:"member_id" binding:"required,uuid"`
	ProviderID     string   `json:"provider_id" binding:"required,uuid"`
	ProcedureCodes []string `json:"procedure_codes" binding:"required"`
	DiagnosisCodes []string `json:"diagnosis_codes" binding:"required"`
	EstimatedCost  int64    `json:"estimated_cost" binding:"required,min=1"`
	Notes          string   `json:"notes"`
}

type ReviewPreAuthRequest struct {
	Decision       string `json:"decision" binding:"required,oneof=APPROVED DENIED INFO_REQUESTED"`
	ApprovedAmount int64  `json:"approved_amount"`
	DenialReason   string `json:"denial_reason"`
	ValidityDays   int    `json:"validity_days"`
}
