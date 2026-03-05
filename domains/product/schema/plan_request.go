package schema

type CreatePlanRequest struct {
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required,oneof=individual group"`
	Segment     string `json:"segment"`
	BasePremium int64  `json:"base_premium" binding:"required,min=1"`
	Currency    string `json:"currency"`
	Description string `json:"description"`
}

type UpdatePlanRequest struct {
	Name        *string `json:"name"`
	Type        *string `json:"type"`
	Segment     *string `json:"segment"`
	BasePremium *int64  `json:"base_premium"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
}

type CreateBenefitRequest struct {
	ParentBenefitID   string `json:"parent_benefit_id,omitempty"`
	Name              string `json:"name" binding:"required"`
	Category          string `json:"category" binding:"required"`
	AnnualLimit       int64  `json:"annual_limit" binding:"required,min=1"`
	CoPayType         string `json:"co_pay_type" binding:"required"`
	CoPayValue        int64  `json:"co_pay_value" binding:"min=0"`
	WaitingPeriodDays int    `json:"waiting_period_days" binding:"min=0"`
	SubLimitType      string `json:"sub_limit_type"`
	SubLimitValue     int64  `json:"sub_limit_value"`
	MinAge            int    `json:"min_age"`
	MaxAge            int    `json:"max_age"`
	WaitingPeriodType string `json:"waiting_period_type"`
	DeductibleAmount  int64  `json:"deductible_amount" binding:"min=0"`
}

type CreateExclusionRequest struct {
	Description string   `json:"description" binding:"required"`
	Type        string   `json:"type" binding:"required,oneof=pre_existing cosmetic experimental"`
	ICDCodes    []string `json:"icd_codes"`
}

type UpdateExclusionRequest struct {
	Description *string  `json:"description"`
	Type        *string  `json:"type"`
	ICDCodes    []string `json:"icd_codes"`
}
