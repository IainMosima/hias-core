package schema

import (
	"encoding/json"
	"time"
)

type CreatePolicyRequest struct {
	PlanID            string    `json:"plan_id" binding:"required,uuid"`
	PolicyholderName  string    `json:"policyholder_name" binding:"required"`
	PolicyholderEmail string    `json:"policyholder_email" binding:"required,email"`
	PolicyholderPhone string    `json:"policyholder_phone" binding:"required"`
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`
}

type EnrollMemberRequest struct {
	NationalID   string `json:"national_id"`
	Name         string `json:"name" binding:"required"`
	DateOfBirth  string `json:"date_of_birth" binding:"required"`
	Gender       string `json:"gender" binding:"required"`
	Relationship string `json:"relationship" binding:"required"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	KRAPin       string `json:"kra_pin"`
	County       string `json:"county"`
	Address      string `json:"address"`
}

type ActivatePolicyRequest struct {
	PaymentReference string `json:"payment_reference" binding:"required"`
}

type UpdatePolicyRequest struct {
	PolicyholderName  *string    `json:"policyholder_name"`
	PolicyholderEmail *string    `json:"policyholder_email"`
	PolicyholderPhone *string    `json:"policyholder_phone"`
	StartDate         *time.Time `json:"start_date"`
	EndDate           *time.Time `json:"end_date"`
}

type ChangePlanRequest struct {
	NewPlanID string `json:"new_plan_id" binding:"required,uuid"`
	Reason    string `json:"reason"`
}

type UpdateMemberRequest struct {
	Name    *string `json:"name"`
	Phone   *string `json:"phone"`
	Email   *string `json:"email"`
	KRAPin  *string `json:"kra_pin"`
	County  *string `json:"county"`
	Address *string `json:"address"`
}

type RemoveMemberRequest struct {
	Reason string `json:"reason"`
}

type BulkIDsRequest struct {
	IDs []string `json:"ids" binding:"required"`
}

type BulkEnrollRequest struct {
	Members []EnrollMemberRequest `json:"members" binding:"required"`
}

type BulkRemoveRequest struct {
	MemberIDs []string `json:"member_ids" binding:"required"`
	Reason    string   `json:"reason"`
}

type CreateEndorsementRequest struct {
	PolicyID        string          `json:"policy_id" binding:"required,uuid"`
	EndorsementType string          `json:"endorsement_type" binding:"required"`
	EffectiveDate   string          `json:"effective_date" binding:"required"`
	Changes         json.RawMessage `json:"changes" binding:"required"`
	Reason          string          `json:"reason"`
}

type RejectEndorsementRequest struct {
	Reason string `json:"reason" binding:"required"`
}

type InitiateRenewalRequest struct {
	PolicyID    string `json:"policy_id" binding:"required,uuid"`
	NewPlanID   string `json:"new_plan_id"`
	RenewalDate string `json:"renewal_date" binding:"required"`
	ExpiresAt   string `json:"expires_at"`
}

type RejectRenewalRequest struct {
	Reason string `json:"reason" binding:"required"`
}

type BulkRenewalRequest struct {
	PolicyIDs []string `json:"policy_ids" binding:"required"`
}

type SubmitAssessmentRequest struct {
	PolicyID            string          `json:"policy_id" binding:"required,uuid"`
	MemberID            string          `json:"member_id"`
	Questionnaire       json.RawMessage `json:"questionnaire" binding:"required"`
	MedicalDeclarations json.RawMessage `json:"medical_declarations"`
}

type ReviewAssessmentRequest struct {
	Status         string `json:"status" binding:"required"`
	RiskScore      int    `json:"risk_score"`
	DecisionReason string `json:"decision_reason"`
}

type ResolveFlagRequest struct {
	Resolution string `json:"resolution" binding:"required"`
}

type OverrideFlagRequest struct {
	Reason string `json:"reason" binding:"required"`
}

type CreateUnderwritingRuleRequest struct {
	PlanID          string `json:"plan_id" binding:"required,uuid"`
	RuleType        string `json:"rule_type" binding:"required"`
	Relationship    string `json:"relationship"`
	ParameterKey    string `json:"parameter_key" binding:"required"`
	ParameterValue  string `json:"parameter_value" binding:"required"`
	Severity        string `json:"severity"`
	RiskScoreWeight int    `json:"risk_score_weight"`
	IsBlocking      bool   `json:"is_blocking"`
	IsActive        *bool  `json:"is_active"`
	Description     string `json:"description"`
}

type UpdateUnderwritingRuleRequest struct {
	RuleType        *string `json:"rule_type"`
	Relationship    *string `json:"relationship"`
	ParameterKey    *string `json:"parameter_key"`
	ParameterValue  *string `json:"parameter_value"`
	Severity        *string `json:"severity"`
	RiskScoreWeight *int    `json:"risk_score_weight"`
	IsBlocking      *bool   `json:"is_blocking"`
	IsActive        *bool   `json:"is_active"`
	Description     *string `json:"description"`
}

type ApplyCreditNoteRequest struct {
	InvoiceID string `json:"invoice_id" binding:"required,uuid"`
}
