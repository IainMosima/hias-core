package entity

import (
	"time"
	"github.com/google/uuid"
)

type Benefit struct {
	ID               uuid.UUID `json:"id"`
	PlanID           uuid.UUID `json:"plan_id"`
	Name             string    `json:"name"`
	Category         string    `json:"category"` // outpatient, inpatient, dental, optical, maternity
	AnnualLimit      int64     `json:"annual_limit"`
	CoPayType        string    `json:"co_pay_type"` // percentage, fixed
	CoPayValue       int64     `json:"co_pay_value"`
	WaitingPeriodDays int      `json:"waiting_period_days"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
