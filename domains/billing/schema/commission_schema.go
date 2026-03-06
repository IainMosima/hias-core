package schema

import (
	"github.com/bitbiz/hias-core/domains/billing/entity"
	"github.com/google/uuid"
	"time"
)

type CreateCommissionRuleRequest struct {
	PlanID         string     `json:"plan_id" binding:"required,uuid"`
	IntermediaryID string     `json:"intermediary_id" binding:"required,uuid"`
	RatePct        float64    `json:"rate_pct"`
	FlatAmount     int64      `json:"flat_amount"`
	EffectiveFrom  time.Time  `json:"effective_from" binding:"required"`
	EffectiveTo    *time.Time `json:"effective_to"`
}

type CalculateCommissionRequest struct {
	PlanID         string `json:"plan_id" binding:"required,uuid"`
	IntermediaryID string `json:"intermediary_id" binding:"required,uuid"`
	PremiumAmount  int64  `json:"premium_amount" binding:"required,min=1"`
}

type CommissionRuleResponse struct {
	ID             uuid.UUID  `json:"id"`
	PlanID         uuid.UUID  `json:"plan_id"`
	IntermediaryID uuid.UUID  `json:"intermediary_id"`
	RatePct        float64    `json:"rate_pct"`
	FlatAmount     int64      `json:"flat_amount"`
	EffectiveFrom  time.Time  `json:"effective_from"`
	EffectiveTo    *time.Time `json:"effective_to,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

type CommissionPaymentResponse struct {
	ID               uuid.UUID  `json:"id"`
	PolicyID         uuid.UUID  `json:"policy_id"`
	IntermediaryID   uuid.UUID  `json:"intermediary_id"`
	CommissionRuleID uuid.UUID  `json:"commission_rule_id"`
	Amount           int64      `json:"amount"`
	Currency         string     `json:"currency"`
	Status           string     `json:"status"`
	PeriodStart      time.Time  `json:"period_start"`
	PeriodEnd        time.Time  `json:"period_end"`
	PaidAt           *time.Time `json:"paid_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}

type CalculateCommissionResponse struct {
	CommissionAmount int64   `json:"commission_amount"`
	RatePct          float64 `json:"rate_pct"`
	FlatAmount       int64   `json:"flat_amount"`
	RuleID           string  `json:"rule_id"`
}

func ToCommissionRuleResponse(r *entity.CommissionRule) CommissionRuleResponse {
	return CommissionRuleResponse{
		ID: r.ID, PlanID: r.PlanID, IntermediaryID: r.IntermediaryID,
		RatePct: r.RatePct, FlatAmount: r.FlatAmount,
		EffectiveFrom: r.EffectiveFrom, EffectiveTo: r.EffectiveTo,
		CreatedAt: r.CreatedAt,
	}
}

func ToCommissionPaymentResponse(p *entity.CommissionPayment) CommissionPaymentResponse {
	return CommissionPaymentResponse{
		ID: p.ID, PolicyID: p.PolicyID, IntermediaryID: p.IntermediaryID,
		CommissionRuleID: p.CommissionRuleID, Amount: p.Amount,
		Currency: p.Currency, Status: p.Status,
		PeriodStart: p.PeriodStart, PeriodEnd: p.PeriodEnd,
		PaidAt: p.PaidAt, CreatedAt: p.CreatedAt,
	}
}
