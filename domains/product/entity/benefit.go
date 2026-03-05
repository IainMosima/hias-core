package entity

import (
	"github.com/google/uuid"
	"time"
)

type Benefit struct {
	ID                uuid.UUID `json:"id"`
	PlanID            uuid.UUID `json:"plan_id"`
	Name              string    `json:"name"`
	Category          string    `json:"category"` // outpatient, inpatient, dental, optical, maternity
	AnnualLimit       int64     `json:"annual_limit"`
	CoPayType         string    `json:"co_pay_type"` // percentage, fixed
	CoPayValue        int64     `json:"co_pay_value"`
	WaitingPeriodDays int       `json:"waiting_period_days"`
	SubLimitType      string    `json:"sub_limit_type"` // none, per_visit, per_item
	SubLimitValue     int64     `json:"sub_limit_value"`
	MinAge            int       `json:"min_age"`
	MaxAge            int       `json:"max_age"`
	WaitingPeriodType string    `json:"waiting_period_type"` // general, maternity, pre_existing
	DeductibleAmount  int64     `json:"deductible_amount"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
