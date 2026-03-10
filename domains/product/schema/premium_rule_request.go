package schema

type CreatePremiumRuleRequest struct {
	CalculationType string  `json:"calculation_type" binding:"required"`
	Relationship    string  `json:"relationship"`
	RateAmount      int64   `json:"rate_amount" binding:"required,min=1"`
	DiscountType    string  `json:"discount_type"`
	DiscountValue   int64   `json:"discount_value"`
	MinMembers      int     `json:"min_members"`
	MinAge          int     `json:"min_age"`
	MaxAge          int     `json:"max_age"`
	RuleType        string  `json:"rule_type"`
	EffectiveFrom   string  `json:"effective_from"`
	EffectiveTo     *string `json:"effective_to"`
	SortOrder       int     `json:"sort_order"`
}

type UpdatePremiumRuleRequest struct {
	CalculationType *string `json:"calculation_type"`
	Relationship    *string `json:"relationship"`
	RateAmount      *int64  `json:"rate_amount"`
	DiscountType    *string `json:"discount_type"`
	DiscountValue   *int64  `json:"discount_value"`
	MinMembers      *int    `json:"min_members"`
	MinAge          *int    `json:"min_age"`
	MaxAge          *int    `json:"max_age"`
	RuleType        *string `json:"rule_type"`
	EffectiveFrom   *string `json:"effective_from"`
	EffectiveTo     *string `json:"effective_to"`
	SortOrder       *int    `json:"sort_order"`
}
