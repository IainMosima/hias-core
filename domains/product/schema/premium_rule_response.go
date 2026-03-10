package schema

import (
	"time"

	"github.com/bitbiz/hias-core/domains/product/entity"
	"github.com/google/uuid"
)

type PremiumRuleResponse struct {
	ID              uuid.UUID  `json:"id"`
	PlanID          uuid.UUID  `json:"plan_id"`
	CalculationType string     `json:"calculation_type"`
	Relationship    string     `json:"relationship,omitempty"`
	RateAmount      int64      `json:"rate_amount"`
	DiscountType    string     `json:"discount_type,omitempty"`
	DiscountValue   int64      `json:"discount_value"`
	MinMembers      int        `json:"min_members"`
	MinAge          int        `json:"min_age"`
	MaxAge          int        `json:"max_age"`
	RuleType        string     `json:"rule_type"`
	EffectiveFrom   time.Time  `json:"effective_from"`
	EffectiveTo     *time.Time `json:"effective_to,omitempty"`
	SortOrder       int        `json:"sort_order"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type MemberPremiumBreakdown struct {
	Relationship string `json:"relationship"`
	Age          int    `json:"age"`
	AgeBand      string `json:"age_band,omitempty"`
	RateAmount   int64  `json:"rate_amount"`
	RuleName     string `json:"rule_name,omitempty"`
}

type AddonPremiumBreakdown struct {
	BenefitID   uuid.UUID `json:"benefit_id"`
	BenefitName string    `json:"benefit_name"`
	Category    string    `json:"category"`
	Premium     int64     `json:"premium"`
}

type PremiumBreakdownResponse struct {
	TotalPremium     int64                    `json:"total_premium"`
	MemberBreakdown  []MemberPremiumBreakdown `json:"member_breakdown"`
	DiscountApplied  int64                    `json:"discount_applied"`
	LoadingApplied   int64                    `json:"loading_applied"`
	CalculationType  string                   `json:"calculation_type"`
	AddonBreakdown   []AddonPremiumBreakdown  `json:"addon_breakdown,omitempty"`
	PremiumFrequency string                   `json:"premium_frequency,omitempty"`
}

func ToPremiumRuleResponse(r *entity.PremiumRule) PremiumRuleResponse {
	return PremiumRuleResponse{
		ID: r.ID, PlanID: r.PlanID, CalculationType: r.CalculationType,
		Relationship: r.Relationship, RateAmount: r.RateAmount,
		DiscountType: r.DiscountType, DiscountValue: r.DiscountValue,
		MinMembers: r.MinMembers, MinAge: r.MinAge, MaxAge: r.MaxAge,
		RuleType: r.RuleType, EffectiveFrom: r.EffectiveFrom,
		EffectiveTo: r.EffectiveTo, SortOrder: r.SortOrder,
		CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
}
