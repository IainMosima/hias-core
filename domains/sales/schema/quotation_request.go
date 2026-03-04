package schema

import "encoding/json"

type CreateQuotationRequest struct {
	LeadID           string          `json:"lead_id" binding:"required,uuid"`
	PlanID           string          `json:"plan_id" binding:"required,uuid"`
	QuotationType    string          `json:"quotation_type" binding:"required,oneof=standard tailor_made"`
	ClientName       string          `json:"client_name" binding:"required"`
	ClientEmail      string          `json:"client_email"`
	ClientPhone      string          `json:"client_phone"`
	MemberCount      int             `json:"member_count" binding:"required,min=1"`
	ProposedMembers  json.RawMessage `json:"proposed_members"`
	BillingFrequency string          `json:"billing_frequency" binding:"required,oneof=monthly quarterly semi_annual annual"`
	DiscountType     string          `json:"discount_type"`
	DiscountValue    int64           `json:"discount_value"`
	DiscountReason   string          `json:"discount_reason"`
	LoadingType      string          `json:"loading_type"`
	LoadingValue     int64           `json:"loading_value"`
	LoadingReason    string          `json:"loading_reason"`
}

type CreateQuotationVersionRequest struct {
	MemberCount      int             `json:"member_count" binding:"required,min=1"`
	ProposedMembers  json.RawMessage `json:"proposed_members"`
	BillingFrequency string          `json:"billing_frequency" binding:"required,oneof=monthly quarterly semi_annual annual"`
	DiscountType     string          `json:"discount_type"`
	DiscountValue    int64           `json:"discount_value"`
	DiscountReason   string          `json:"discount_reason"`
	LoadingType      string          `json:"loading_type"`
	LoadingValue     int64           `json:"loading_value"`
	LoadingReason    string          `json:"loading_reason"`
}

type ApproveVersionRequest struct {
	Notes string `json:"notes"`
}

type RejectVersionRequest struct {
	Reason string `json:"reason" binding:"required"`
}

type SendQuotationRequest struct {
	Channel string `json:"channel" binding:"required,oneof=SMS EMAIL"`
	Message string `json:"message"`
}

type ConvertToPolicyRequest struct {
	StartDate string `json:"start_date" binding:"required"`
	Notes     string `json:"notes"`
}

type UploadDocumentMeta struct {
	FileName       string          `json:"file_name" binding:"required"`
	FileType       string          `json:"file_type" binding:"required"`
	FileSize       int64           `json:"file_size" binding:"required"`
	VersionNumber  int             `json:"version_number"`
	CanEditRoles   json.RawMessage `json:"can_edit_roles"`
	CanDeleteRoles json.RawMessage `json:"can_delete_roles"`
}

type CreateApprovalLimitRequest struct {
	RoleName              string `json:"role_name" binding:"required"`
	MaxDiscountPercentage int64  `json:"max_discount_percentage"`
	MaxDiscountAmount     int64  `json:"max_discount_amount"`
	MaxLoadingPercentage  int64  `json:"max_loading_percentage"`
	MaxLoadingAmount      int64  `json:"max_loading_amount"`
	EscalationRole        string `json:"escalation_role"`
}

type UpdateDocumentMeta struct {
	FileName       string          `json:"file_name"`
	CanEditRoles   json.RawMessage `json:"can_edit_roles"`
	CanDeleteRoles json.RawMessage `json:"can_delete_roles"`
}

type UpdateApprovalLimitRequest struct {
	MaxDiscountPercentage *int64 `json:"max_discount_percentage"`
	MaxDiscountAmount     *int64 `json:"max_discount_amount"`
	MaxLoadingPercentage  *int64 `json:"max_loading_percentage"`
	MaxLoadingAmount      *int64 `json:"max_loading_amount"`
	EscalationRole        string `json:"escalation_role"`
}
