package entity

import (
	"time"

	"github.com/google/uuid"
)

type PremiumRule struct {
	ID              uuid.UUID `json:"id"`
	PlanID          uuid.UUID `json:"plan_id"`
	CalculationType string    `json:"calculation_type"` // flat, per_member, tiered
	Relationship    string    `json:"relationship,omitempty"`
	RateAmount      int64     `json:"rate_amount"`
	DiscountType    string    `json:"discount_type,omitempty"` // percentage, fixed
	DiscountValue   int64     `json:"discount_value"`
	MinMembers      int       `json:"min_members"`
	MinAge          int       `json:"min_age"`
	MaxAge          int       `json:"max_age"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
