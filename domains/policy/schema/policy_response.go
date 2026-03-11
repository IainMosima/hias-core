package schema

import (
	"encoding/json"
	"time"

	"github.com/bitbiz/hias-core/domains/policy/entity"
	"github.com/google/uuid"
)

type PolicyResponse struct {
	ID                uuid.UUID  `json:"id"`
	PlanID            uuid.UUID  `json:"plan_id"`
	PlanName          string     `json:"plan_name"`
	PolicyholderName  string     `json:"policyholder_name"`
	PolicyholderEmail string     `json:"policyholder_email"`
	PolicyholderPhone string     `json:"policyholder_phone"`
	PolicyNumber      string     `json:"policy_number"`
	Status            string     `json:"status"`
	StartDate         time.Time  `json:"start_date"`
	EndDate           time.Time  `json:"end_date"`
	PremiumAmount     int64      `json:"premium_amount"`
	Currency          string     `json:"currency"`
	RenewedFromID     *uuid.UUID `json:"renewed_from_id,omitempty"`
	ActivatedAt       *time.Time `json:"activated_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type MemberResponse struct {
	ID                uuid.UUID  `json:"id"`
	PolicyID          uuid.UUID  `json:"policy_id"`
	NationalID        string     `json:"national_id"`
	Name              string     `json:"name"`
	DateOfBirth       time.Time  `json:"date_of_birth"`
	Gender            string     `json:"gender"`
	Relationship      string     `json:"relationship"`
	MemberNumber      string     `json:"member_number"`
	Phone             string     `json:"phone"`
	Email             string     `json:"email"`
	KRAPin            string     `json:"kra_pin"`
	County            string     `json:"county"`
	Address           string     `json:"address"`
	Status            string     `json:"status"`
	Verified          bool       `json:"verified"`
	VerifiedAt        *time.Time `json:"verified_at,omitempty"`
	CoverageStartDate *time.Time `json:"coverage_start_date,omitempty"`
	CoverageEndDate   *time.Time `json:"coverage_end_date,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
}

type BulkResultResponse struct {
	Succeeded int      `json:"succeeded"`
	Failed    int      `json:"failed"`
	Errors    []string `json:"errors,omitempty"`
}

type BulkMemberResultResponse struct {
	Succeeded int              `json:"succeeded"`
	Failed    int              `json:"failed"`
	Members   []MemberResponse `json:"members,omitempty"`
	Errors    []string         `json:"errors,omitempty"`
}

type EndorsementResponse struct {
	ID                uuid.UUID       `json:"id"`
	PolicyID          uuid.UUID       `json:"policy_id"`
	EndorsementType   string          `json:"endorsement_type"`
	Status            string          `json:"status"`
	EffectiveDate     time.Time       `json:"effective_date"`
	Changes           json.RawMessage `json:"changes"`
	Reason            string          `json:"reason"`
	PremiumAdjustment int64           `json:"premium_adjustment"`
	RequestedBy       uuid.UUID       `json:"requested_by"`
	ApprovedBy        uuid.UUID       `json:"approved_by,omitempty"`
	ApprovedAt        *time.Time      `json:"approved_at,omitempty"`
	AppliedAt         *time.Time      `json:"applied_at,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

type RenewalResponse struct {
	ID                  uuid.UUID  `json:"id"`
	PolicyID            uuid.UUID  `json:"policy_id"`
	RenewedPolicyID     uuid.UUID  `json:"renewed_policy_id,omitempty"`
	Status              string     `json:"status"`
	RenewalDate         time.Time  `json:"renewal_date"`
	NewPremium          int64      `json:"new_premium"`
	PremiumChangeReason string     `json:"premium_change_reason"`
	NewPlanID           uuid.UUID  `json:"new_plan_id,omitempty"`
	ApprovedBy          uuid.UUID  `json:"approved_by,omitempty"`
	ApprovedAt          *time.Time `json:"approved_at,omitempty"`
	CompletedAt         *time.Time `json:"completed_at,omitempty"`
	ExpiresAt           *time.Time `json:"expires_at,omitempty"`
	CreatedBy           uuid.UUID  `json:"created_by"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type UnderwritingResponse struct {
	ID                  uuid.UUID       `json:"id"`
	PolicyID            uuid.UUID       `json:"policy_id"`
	MemberID            uuid.UUID       `json:"member_id,omitempty"`
	Status              string          `json:"status"`
	Questionnaire       json.RawMessage `json:"questionnaire"`
	MedicalDeclarations json.RawMessage `json:"medical_declarations"`
	RiskScore           int             `json:"risk_score"`
	RiskFlags           json.RawMessage `json:"risk_flags"`
	DecisionReason      string          `json:"decision_reason"`
	AssessedBy          uuid.UUID       `json:"assessed_by,omitempty"`
	AssessedAt          *time.Time      `json:"assessed_at,omitempty"`
	CreatedBy           uuid.UUID       `json:"created_by"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
}

type PolicyDocumentResponse struct {
	ID           uuid.UUID `json:"id"`
	PolicyID     uuid.UUID `json:"policy_id"`
	MemberID     uuid.UUID `json:"member_id,omitempty"`
	DocumentType string    `json:"document_type"`
	FileName     string    `json:"file_name"`
	FileSize     int64     `json:"file_size"`
	S3Key        string    `json:"s3_key"`
	GeneratedBy  uuid.UUID `json:"generated_by"`
	CreatedAt    time.Time `json:"created_at"`
}

func ToPolicyResponse(p *entity.Policy) PolicyResponse {
	return PolicyResponse{
		ID: p.ID, PlanID: p.PlanID, PlanName: p.PlanName, PolicyholderName: p.PolicyholderName,
		PolicyholderEmail: p.PolicyholderEmail, PolicyholderPhone: p.PolicyholderPhone,
		PolicyNumber: p.PolicyNumber, Status: p.Status, StartDate: p.StartDate,
		EndDate: p.EndDate, PremiumAmount: p.PremiumAmount, Currency: p.Currency,
		RenewedFromID: p.RenewedFromID, ActivatedAt: p.ActivatedAt,
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}

func ToMemberResponse(m *entity.Member) MemberResponse {
	return MemberResponse{
		ID: m.ID, PolicyID: m.PolicyID, NationalID: m.NationalID,
		Name: m.Name, DateOfBirth: m.DateOfBirth, Gender: m.Gender,
		Relationship: m.Relationship, MemberNumber: m.MemberNumber,
		Phone: m.Phone, Email: m.Email, KRAPin: m.KRAPin,
		County: m.County, Address: m.Address, Status: m.Status,
		Verified: m.Verified, VerifiedAt: m.VerifiedAt,
		CoverageStartDate: m.CoverageStartDate, CoverageEndDate: m.CoverageEndDate,
		CreatedAt: m.CreatedAt,
	}
}

func ToEndorsementResponse(e *entity.Endorsement) EndorsementResponse {
	return EndorsementResponse{
		ID: e.ID, PolicyID: e.PolicyID, EndorsementType: e.EndorsementType,
		Status: e.Status, EffectiveDate: e.EffectiveDate, Changes: e.Changes,
		Reason: e.Reason, PremiumAdjustment: e.PremiumAdjustment,
		RequestedBy: e.RequestedBy, ApprovedBy: e.ApprovedBy,
		ApprovedAt: e.ApprovedAt, AppliedAt: e.AppliedAt,
		CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt,
	}
}

func ToRenewalResponse(r *entity.PolicyRenewal) RenewalResponse {
	return RenewalResponse{
		ID: r.ID, PolicyID: r.PolicyID, RenewedPolicyID: r.RenewedPolicyID,
		Status: r.Status, RenewalDate: r.RenewalDate,
		NewPremium: r.NewPremium, PremiumChangeReason: r.PremiumChangeReason,
		NewPlanID: r.NewPlanID, ApprovedBy: r.ApprovedBy,
		ApprovedAt: r.ApprovedAt, CompletedAt: r.CompletedAt,
		ExpiresAt: r.ExpiresAt, CreatedBy: r.CreatedBy,
		CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
}

func ToUnderwritingResponse(u *entity.UnderwritingAssessment) UnderwritingResponse {
	return UnderwritingResponse{
		ID: u.ID, PolicyID: u.PolicyID, MemberID: u.MemberID,
		Status: u.Status, Questionnaire: u.Questionnaire,
		MedicalDeclarations: u.MedicalDeclarations, RiskScore: u.RiskScore,
		RiskFlags: u.RiskFlags, DecisionReason: u.DecisionReason,
		AssessedBy: u.AssessedBy, AssessedAt: u.AssessedAt,
		CreatedBy: u.CreatedBy, CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt,
	}
}

func ToPolicyDocumentResponse(d *entity.PolicyDocument) PolicyDocumentResponse {
	return PolicyDocumentResponse{
		ID: d.ID, PolicyID: d.PolicyID, MemberID: d.MemberID,
		DocumentType: d.DocumentType, FileName: d.FileName,
		FileSize: d.FileSize, S3Key: d.S3Key,
		GeneratedBy: d.GeneratedBy, CreatedAt: d.CreatedAt,
	}
}

type UnderwritingFlagResponse struct {
	ID           uuid.UUID  `json:"id"`
	AssessmentID uuid.UUID  `json:"assessment_id,omitempty"`
	PolicyID     uuid.UUID  `json:"policy_id"`
	MemberID     uuid.UUID  `json:"member_id,omitempty"`
	FlagType     string     `json:"flag_type"`
	Severity     string     `json:"severity"`
	Details      string     `json:"details"`
	Status       string     `json:"status"`
	ResolvedBy   uuid.UUID  `json:"resolved_by,omitempty"`
	ResolvedAt   *time.Time `json:"resolved_at,omitempty"`
	Resolution   string     `json:"resolution,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func ToUnderwritingFlagResponse(f *entity.UnderwritingFlag) UnderwritingFlagResponse {
	return UnderwritingFlagResponse{
		ID: f.ID, AssessmentID: f.AssessmentID, PolicyID: f.PolicyID,
		MemberID: f.MemberID, FlagType: f.FlagType, Severity: f.Severity,
		Details: f.Details, Status: f.Status, ResolvedBy: f.ResolvedBy,
		ResolvedAt: f.ResolvedAt, Resolution: f.Resolution,
		CreatedAt: f.CreatedAt, UpdatedAt: f.UpdatedAt,
	}
}

type UnderwritingRuleResponse struct {
	ID              uuid.UUID `json:"id"`
	PlanID          uuid.UUID `json:"plan_id"`
	RuleType        string    `json:"rule_type"`
	Relationship    string    `json:"relationship,omitempty"`
	ParameterKey    string    `json:"parameter_key"`
	ParameterValue  string    `json:"parameter_value"`
	Severity        string    `json:"severity"`
	RiskScoreWeight int       `json:"risk_score_weight"`
	IsBlocking      bool      `json:"is_blocking"`
	IsActive        bool      `json:"is_active"`
	Description     string    `json:"description,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func ToUnderwritingRuleResponse(r *entity.UnderwritingRule) UnderwritingRuleResponse {
	return UnderwritingRuleResponse{
		ID: r.ID, PlanID: r.PlanID, RuleType: r.RuleType,
		Relationship: r.Relationship, ParameterKey: r.ParameterKey,
		ParameterValue: r.ParameterValue, Severity: r.Severity,
		RiskScoreWeight: r.RiskScoreWeight, IsBlocking: r.IsBlocking,
		IsActive: r.IsActive, Description: r.Description,
		CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
}
