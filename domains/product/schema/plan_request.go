package schema

type CreatePlanRequest struct {
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required,oneof=individual group"`
	BasePremium int64  `json:"base_premium" binding:"required,min=1"`
	Currency    string `json:"currency"`
	Description string `json:"description"`
}

type UpdatePlanRequest struct {
	Name        *string `json:"name"`
	Type        *string `json:"type"`
	BasePremium *int64  `json:"base_premium"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
}

type CreateBenefitRequest struct {
	Name              string `json:"name" binding:"required"`
	Category          string `json:"category" binding:"required,oneof=outpatient inpatient dental optical maternity"`
	AnnualLimit       int64  `json:"annual_limit" binding:"required,min=1"`
	CoPayType         string `json:"co_pay_type" binding:"required,oneof=percentage fixed"`
	CoPayValue        int64  `json:"co_pay_value" binding:"min=0"`
	WaitingPeriodDays int    `json:"waiting_period_days" binding:"min=0"`
}

type CreateExclusionRequest struct {
	Description string   `json:"description" binding:"required"`
	Type        string   `json:"type" binding:"required,oneof=pre_existing cosmetic experimental"`
	ICDCodes    []string `json:"icd_codes"`
}
