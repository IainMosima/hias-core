package entity

import (
	"time"

	"github.com/google/uuid"
)

type ProfitCommission struct {
	ID                  uuid.UUID  `json:"id"`
	TreatyID            uuid.UUID  `json:"treaty_id"`
	CommissionType      string     `json:"commission_type"`
	LossRatioFrom       float64    `json:"loss_ratio_from"`
	LossRatioTo         float64    `json:"loss_ratio_to"`
	CommissionRate      float64    `json:"commission_rate"`
	CarryForwardYears   int        `json:"carry_forward_years"`
	CarryForwardBalance int64      `json:"carry_forward_balance"`
	PeriodStart         *time.Time `json:"period_start,omitempty"`
	PeriodEnd           *time.Time `json:"period_end,omitempty"`
	CalculatedAmount    int64      `json:"calculated_amount"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}
