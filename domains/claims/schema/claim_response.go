package schema

import (
	"encoding/json"
	"time"

	claimEntity "github.com/bitbiz/hias-core/domains/claims/entity"
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
	ClaimType            string                `json:"claim_type"`
	VettedAmount         *int64                `json:"vetted_amount,omitempty"`
	VettedBy             uuid.UUID             `json:"vetted_by,omitempty"`
	VettedAt             *time.Time            `json:"vetted_at,omitempty"`
	SLABreachAt          *time.Time            `json:"sla_breach_at,omitempty"`
	RejectionReason      string                `json:"rejection_reason,omitempty"`
	LineItems            []LineItemResponse    `json:"line_items,omitempty"`
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
	DeductibleApplied    int64           `json:"deductible_applied"`
	CoPayApplied         int64           `json:"co_pay_applied"`
	SubLimitApplied      int64           `json:"sub_limit_applied,omitempty"`
	BenefitCategory      string          `json:"benefit_category,omitempty"`
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

func ToClaimResponse(c *claimEntity.Claim) ClaimResponse {
	return ClaimResponse{
		ID: c.ID, ClaimNumber: c.ClaimNumber, PolicyID: c.PolicyID,
		MemberID: c.MemberID, ProviderID: c.ProviderID, Status: c.Status,
		TotalAmount: c.TotalAmount, ApprovedAmount: c.ApprovedAmount,
		CoPayAmount: c.CoPayAmount, MemberResponsibility: c.MemberResponsibility,
		DiagnosisCodes: c.DiagnosisCodes, ServiceDate: c.ServiceDate,
		ClaimType: c.ClaimType, VettedAmount: c.VettedAmount,
		VettedBy: c.VettedBy, VettedAt: c.VettedAt,
		SLABreachAt: c.SLABreachAt,
		Notes:       c.Notes, RejectionReason: c.RejectionReason,
		CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt,
	}
}

type CaseRecordResponse struct {
	ID                 uuid.UUID  `json:"id"`
	CaseNumber         string     `json:"case_number"`
	PreAuthID          uuid.UUID  `json:"preauth_id"`
	PolicyID           uuid.UUID  `json:"policy_id"`
	MemberID           uuid.UUID  `json:"member_id"`
	ProviderID         uuid.UUID  `json:"provider_id"`
	Status             string     `json:"status"`
	AdmissionDate      *time.Time `json:"admission_date,omitempty"`
	ExpectedDischarge  *time.Time `json:"expected_discharge,omitempty"`
	ActualDischarge    *time.Time `json:"actual_discharge,omitempty"`
	Diagnosis          string     `json:"diagnosis,omitempty"`
	TreatingDoctor     string     `json:"treating_doctor,omitempty"`
	RoomType           string     `json:"room_type,omitempty"`
	TotalEstimatedCost int64      `json:"total_estimated_cost"`
	TotalActualCost    int64      `json:"total_actual_cost"`
	Notes              string     `json:"notes,omitempty"`
	ClosedAt           *time.Time `json:"closed_at,omitempty"`
	CreatedBy          uuid.UUID  `json:"created_by"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

func ToCaseRecordResponse(c *claimEntity.CaseRecord) CaseRecordResponse {
	return CaseRecordResponse{
		ID: c.ID, CaseNumber: c.CaseNumber, PreAuthID: c.PreAuthID,
		PolicyID: c.PolicyID, MemberID: c.MemberID, ProviderID: c.ProviderID,
		Status: c.Status, AdmissionDate: c.AdmissionDate,
		ExpectedDischarge: c.ExpectedDischarge, ActualDischarge: c.ActualDischarge,
		Diagnosis: c.Diagnosis, TreatingDoctor: c.TreatingDoctor, RoomType: c.RoomType,
		TotalEstimatedCost: c.TotalEstimatedCost, TotalActualCost: c.TotalActualCost,
		Notes: c.Notes, ClosedAt: c.ClosedAt, CreatedBy: c.CreatedBy,
		CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt,
	}
}

type ClaimDocumentResponse struct {
	ID         uuid.UUID `json:"id"`
	ClaimID    uuid.UUID `json:"claim_id"`
	FileName   string    `json:"file_name"`
	FileType   string    `json:"file_type"`
	FileSize   int64     `json:"file_size"`
	S3Key      string    `json:"s3_key"`
	UploadedBy uuid.UUID `json:"uploaded_by"`
	CreatedAt  time.Time `json:"created_at"`
}

func ToClaimDocumentResponse(d *claimEntity.ClaimDocument) ClaimDocumentResponse {
	return ClaimDocumentResponse{
		ID: d.ID, ClaimID: d.ClaimID, FileName: d.FileName,
		FileType: d.FileType, FileSize: d.FileSize, S3Key: d.S3Key,
		UploadedBy: d.UploadedBy, CreatedAt: d.CreatedAt,
	}
}
