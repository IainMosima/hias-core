package schema

import "time"

type CreateTreatyRequest struct {
	Name           string    `json:"name" binding:"required"`
	TreatyType     string    `json:"treaty_type" binding:"required,oneof=QUOTA_SHARE XOL"`
	EffectiveDate  time.Time `json:"effective_date" binding:"required"`
	ExpiryDate     time.Time `json:"expiry_date" binding:"required"`
	RetentionLimit int64     `json:"retention_limit"`
	Currency       string    `json:"currency"`
	Notes          string    `json:"notes"`
}

type UpdateTreatyRequest struct {
	Name           string    `json:"name"`
	EffectiveDate  time.Time `json:"effective_date"`
	ExpiryDate     time.Time `json:"expiry_date"`
	RetentionLimit int64     `json:"retention_limit"`
	Currency       string    `json:"currency"`
	Notes          string    `json:"notes"`
}

type AddParticipantRequest struct {
	ReinsurerName   string  `json:"reinsurer_name" binding:"required"`
	SharePercentage float64 `json:"share_percentage" binding:"required,gt=0,lte=100"`
	CommissionRate  float64 `json:"commission_rate" binding:"gte=0"`
	IsLead          bool    `json:"is_lead"`
}

type UpdateParticipantRequest struct {
	ReinsurerName   string  `json:"reinsurer_name"`
	SharePercentage float64 `json:"share_percentage" binding:"gte=0,lte=100"`
	CommissionRate  float64 `json:"commission_rate" binding:"gte=0"`
	IsLead          bool    `json:"is_lead"`
}

type AddLayerRequest struct {
	LayerNumber      int     `json:"layer_number" binding:"required,gte=1"`
	AttachmentPoint  int64   `json:"attachment_point" binding:"required,gte=0"`
	LayerLimit       int64   `json:"layer_limit" binding:"required,gt=0"`
	DeductibleAmount int64   `json:"deductible_amount" binding:"gte=0"`
	PremiumRate      float64 `json:"premium_rate" binding:"gte=0"`
	AggregateLimit   *int64  `json:"aggregate_limit"`
}

type UpdateLayerRequest struct {
	AttachmentPoint  int64   `json:"attachment_point"`
	LayerLimit       int64   `json:"layer_limit"`
	DeductibleAmount int64   `json:"deductible_amount"`
	PremiumRate      float64 `json:"premium_rate"`
	AggregateLimit   *int64  `json:"aggregate_limit"`
}

type CedePremiumRequest struct {
	TreatyID string `json:"treaty_id" binding:"required,uuid"`
	PolicyID string `json:"policy_id" binding:"required,uuid"`
	Amount   int64  `json:"amount" binding:"required,gt=0"`
}

type AutoCedePolicyPremiumRequest struct {
	PolicyID string `json:"policy_id" binding:"required,uuid"`
	Amount   int64  `json:"amount" binding:"required,gt=0"`
}

type CreateRecoveryRequest struct {
	ClaimID           string `json:"claim_id" binding:"required,uuid"`
	TreatyID          string `json:"treaty_id" binding:"required,uuid"`
	TreatyLayerID     string `json:"treaty_layer_id,omitempty"`
	CessionID         string `json:"cession_id,omitempty"`
	GrossAmount       int64  `json:"gross_amount" binding:"required,gt=0"`
	RecoverableAmount int64  `json:"recoverable_amount" binding:"required,gt=0"`
	Notes             string `json:"notes"`
}

type ApplyRecoveryForClaimRequest struct {
	ApprovedAmount int64 `json:"approved_amount" binding:"required,gt=0"`
}

type RecoveryWorkflowRequest struct {
	Notes string `json:"notes"`
}

type RecordPaymentRequest struct {
	Amount int64  `json:"amount" binding:"required,gt=0"`
	Notes  string `json:"notes"`
}

type GenerateBordereauRequest struct {
	TreatyID    string    `json:"treaty_id" binding:"required,uuid"`
	PeriodStart time.Time `json:"period_start" binding:"required"`
	PeriodEnd   time.Time `json:"period_end" binding:"required"`
}

type GenerateStatementRequest struct {
	TreatyID      string    `json:"treaty_id" binding:"required,uuid"`
	ParticipantID string    `json:"participant_id" binding:"required,uuid"`
	PeriodStart   time.Time `json:"period_start" binding:"required"`
	PeriodEnd     time.Time `json:"period_end" binding:"required"`
}

type CalculateProfitCommissionRequest struct {
	TreatyID    string    `json:"treaty_id" binding:"required,uuid"`
	PeriodStart time.Time `json:"period_start" binding:"required"`
	PeriodEnd   time.Time `json:"period_end" binding:"required"`
}

type AddProfitCommissionRuleRequest struct {
	CommissionType    string  `json:"commission_type" binding:"required,oneof=SLIDING_SCALE FLAT CARRY_FORWARD"`
	LossRatioFrom     float64 `json:"loss_ratio_from" binding:"gte=0"`
	LossRatioTo       float64 `json:"loss_ratio_to" binding:"gte=0"`
	CommissionRate    float64 `json:"commission_rate" binding:"required,gte=0"`
	CarryForwardYears int     `json:"carry_forward_years"`
}
