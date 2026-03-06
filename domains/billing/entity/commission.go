package entity

import (
	"github.com/google/uuid"
	"time"
)

type CommissionRule struct {
	ID             uuid.UUID  `json:"id"`
	PlanID         uuid.UUID  `json:"plan_id"`
	IntermediaryID uuid.UUID  `json:"intermediary_id"`
	RatePct        float64    `json:"rate_pct"`
	FlatAmount     int64      `json:"flat_amount"`
	EffectiveFrom  time.Time  `json:"effective_from"`
	EffectiveTo    *time.Time `json:"effective_to,omitempty"`
	CreatedBy      uuid.UUID  `json:"created_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type CommissionPayment struct {
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
	CreatedBy        uuid.UUID  `json:"created_by,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}
