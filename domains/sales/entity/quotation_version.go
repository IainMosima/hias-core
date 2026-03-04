package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type QuotationVersion struct {
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
