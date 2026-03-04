package schema

import (
	"encoding/json"
	"time"

	"github.com/bitbiz/hias-core/domains/sales/entity"
	"github.com/google/uuid"
)

type QuotationResponse struct {
	ID              uuid.UUID  `json:"id"`
	QuotationNumber string     `json:"quotation_number"`
	LeadID          uuid.UUID  `json:"lead_id"`
	PlanID          uuid.UUID  `json:"plan_id"`
	QuotationType   string     `json:"quotation_type"`
	Status          string     `json:"status"`
	CurrentVersion  int        `json:"current_version"`
	PolicyID        uuid.UUID  `json:"policy_id,omitempty"`
	ValidFrom       *time.Time `json:"valid_from,omitempty"`
	ValidUntil      *time.Time `json:"valid_until,omitempty"`
	ClientName      string     `json:"client_name"`
	ClientEmail     string     `json:"client_email"`
	ClientPhone     string     `json:"client_phone"`
	Currency        string     `json:"currency"`
	CreatedBy       uuid.UUID  `json:"created_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type QuotationVersionResponse struct {
	ID               uuid.UUID       `json:"id"`
	QuotationID      uuid.UUID       `json:"quotation_id"`
	VersionNumber    int             `json:"version_number"`
	BasePremium      int64           `json:"base_premium"`
	DiscountType     string          `json:"discount_type"`
	DiscountValue    int64           `json:"discount_value"`
	DiscountReason   string          `json:"discount_reason"`
	LoadingType      string          `json:"loading_type"`
	LoadingValue     int64           `json:"loading_value"`
	LoadingReason    string          `json:"loading_reason"`
	FinalPremium     int64           `json:"final_premium"`
	MemberCount      int             `json:"member_count"`
	ProposedMembers  json.RawMessage `json:"proposed_members"`
	BillingFrequency string          `json:"billing_frequency"`
	RequiresApproval bool            `json:"requires_approval"`
	ApprovalStatus   string          `json:"approval_status"`
	ApprovedBy       uuid.UUID       `json:"approved_by,omitempty"`
	ApprovedAt       *time.Time      `json:"approved_at,omitempty"`
	RejectionReason  string          `json:"rejection_reason,omitempty"`
	PricingBreakdown json.RawMessage `json:"pricing_breakdown"`
	CreatedBy        uuid.UUID       `json:"created_by"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

type QuotationDetailResponse struct {
	QuotationResponse
	Versions  []QuotationVersionResponse  `json:"versions"`
	Documents []QuotationDocumentResponse `json:"documents"`
}

type QuotationDocumentResponse struct {
	ID             uuid.UUID       `json:"id"`
	QuotationID    uuid.UUID       `json:"quotation_id"`
	VersionNumber  int             `json:"version_number"`
	FileName       string          `json:"file_name"`
	FileType       string          `json:"file_type"`
	FileSize       int64           `json:"file_size"`
	UploadedBy     uuid.UUID       `json:"uploaded_by"`
	CanEditRoles   json.RawMessage `json:"can_edit_roles"`
	CanDeleteRoles json.RawMessage `json:"can_delete_roles"`
	CreatedAt      time.Time       `json:"created_at"`
}

type VersionComparisonResponse struct {
	VersionA    QuotationVersionResponse `json:"version_a"`
	VersionB    QuotationVersionResponse `json:"version_b"`
	PricingDiff PricingDiff              `json:"pricing_diff"`
}

type PricingDiff struct {
	BasePremiumDiff  int64 `json:"base_premium_diff"`
	DiscountDiff     int64 `json:"discount_diff"`
	LoadingDiff      int64 `json:"loading_diff"`
	FinalPremiumDiff int64 `json:"final_premium_diff"`
	MemberCountDiff  int   `json:"member_count_diff"`
}

type ApprovalLimitResponse struct {
	ID                    uuid.UUID `json:"id"`
	RoleName              string    `json:"role_name"`
	MaxDiscountPercentage int64     `json:"max_discount_percentage"`
	MaxDiscountAmount     int64     `json:"max_discount_amount"`
	MaxLoadingPercentage  int64     `json:"max_loading_percentage"`
	MaxLoadingAmount      int64     `json:"max_loading_amount"`
	EscalationRole        string    `json:"escalation_role"`
	IsActive              bool      `json:"is_active"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type ConversionResultResponse struct {
	QuotationID     uuid.UUID `json:"quotation_id"`
	PolicyID        uuid.UUID `json:"policy_id"`
	QuotationNumber string    `json:"quotation_number"`
	PolicyNumber    string    `json:"policy_number"`
	Message         string    `json:"message"`
}

func ToQuotationResponse(q *entity.Quotation) QuotationResponse {
	return QuotationResponse{
		ID: q.ID, QuotationNumber: q.QuotationNumber, LeadID: q.LeadID,
		PlanID: q.PlanID, QuotationType: q.QuotationType, Status: q.Status,
		CurrentVersion: q.CurrentVersion, PolicyID: q.PolicyID,
		ValidFrom: q.ValidFrom, ValidUntil: q.ValidUntil,
		ClientName: q.ClientName, ClientEmail: q.ClientEmail,
		ClientPhone: q.ClientPhone, Currency: q.Currency,
		CreatedBy: q.CreatedBy, CreatedAt: q.CreatedAt, UpdatedAt: q.UpdatedAt,
	}
}

func ToQuotationVersionResponse(v *entity.QuotationVersion) QuotationVersionResponse {
	return QuotationVersionResponse{
		ID: v.ID, QuotationID: v.QuotationID, VersionNumber: v.VersionNumber,
		BasePremium: v.BasePremium, DiscountType: v.DiscountType,
		DiscountValue: v.DiscountValue, DiscountReason: v.DiscountReason,
		LoadingType: v.LoadingType, LoadingValue: v.LoadingValue,
		LoadingReason: v.LoadingReason, FinalPremium: v.FinalPremium,
		MemberCount: v.MemberCount, ProposedMembers: v.ProposedMembers,
		BillingFrequency: v.BillingFrequency, RequiresApproval: v.RequiresApproval,
		ApprovalStatus: v.ApprovalStatus, ApprovedBy: v.ApprovedBy,
		ApprovedAt: v.ApprovedAt, RejectionReason: v.RejectionReason,
		PricingBreakdown: v.PricingBreakdown, CreatedBy: v.CreatedBy,
		CreatedAt: v.CreatedAt, UpdatedAt: v.UpdatedAt,
	}
}

func ToQuotationDocumentResponse(d *entity.QuotationDocument) QuotationDocumentResponse {
	return QuotationDocumentResponse{
		ID: d.ID, QuotationID: d.QuotationID, VersionNumber: d.VersionNumber,
		FileName: d.FileName, FileType: d.FileType, FileSize: d.FileSize,
		UploadedBy: d.UploadedBy, CanEditRoles: d.CanEditRoles,
		CanDeleteRoles: d.CanDeleteRoles, CreatedAt: d.CreatedAt,
	}
}

func ToApprovalLimitResponse(a *entity.ApprovalLimit) ApprovalLimitResponse {
	return ApprovalLimitResponse{
		ID: a.ID, RoleName: a.RoleName,
		MaxDiscountPercentage: a.MaxDiscountPercentage, MaxDiscountAmount: a.MaxDiscountAmount,
		MaxLoadingPercentage: a.MaxLoadingPercentage, MaxLoadingAmount: a.MaxLoadingAmount,
		EscalationRole: a.EscalationRole, IsActive: a.IsActive,
		CreatedAt: a.CreatedAt, UpdatedAt: a.UpdatedAt,
	}
}
