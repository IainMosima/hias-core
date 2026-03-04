package schema

import (
	"time"

	"github.com/bitbiz/hias-core/domains/product/entity"
	"github.com/google/uuid"
)

type PremiumRuleResponse struct {
	ID              uuid.UUID `json:"id"`
	PlanID          uuid.UUID `json:"plan_id"`
	CalculationType string    `json:"calculation_type"`
	Relationship    string    `json:"relationship,omitempty"`
	RateAmount      int64     `json:"rate_amount"`
	DiscountType    string    `json:"discount_type,omitempty"`
	DiscountValue   int64     `json:"discount_value"`
	MinMembers      int       `json:"min_members"`
	MinAge          int       `json:"min_age"`
	MaxAge          int       `json:"max_age"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func ToPremiumRuleResponse(r *entity.PremiumRule) PremiumRuleResponse {
	return PremiumRuleResponse{
		ID: r.ID, PlanID: r.PlanID, CalculationType: r.CalculationType,
		Relationship: r.Relationship, RateAmount: r.RateAmount,
		DiscountType: r.DiscountType, DiscountValue: r.DiscountValue,
		MinMembers: r.MinMembers, MinAge: r.MinAge, MaxAge: r.MaxAge,
		CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
}
