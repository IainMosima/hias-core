package schema

import (
	"encoding/json"
	"github.com/bitbiz/hias-core/domains/product/entity"
	"github.com/google/uuid"
	"time"
)

type PlanResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Segment     string    `json:"segment"`
	BasePremium int64     `json:"base_premium"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BenefitResponse struct {
	ID                uuid.UUID `json:"id"`
	PlanID            uuid.UUID `json:"plan_id"`
	Name              string    `json:"name"`
	Category          string    `json:"category"`
	AnnualLimit       int64     `json:"annual_limit"`
	CoPayType         string    `json:"co_pay_type"`
	CoPayValue        int64     `json:"co_pay_value"`
	WaitingPeriodDays int       `json:"waiting_period_days"`
	SubLimitType      string    `json:"sub_limit_type"`
	SubLimitValue     int64     `json:"sub_limit_value"`
	MinAge            int       `json:"min_age"`
	MaxAge            int       `json:"max_age"`
	WaitingPeriodType string    `json:"waiting_period_type"`
	CreatedAt         time.Time `json:"created_at"`
}

type ExclusionResponse struct {
	ID          uuid.UUID       `json:"id"`
	PlanID      uuid.UUID       `json:"plan_id"`
	Description string          `json:"description"`
	Type        string          `json:"type"`
	ICDCodes    json.RawMessage `json:"icd_codes"`
	CreatedAt   time.Time       `json:"created_at"`
}

func ToPlanResponse(p *entity.Plan) PlanResponse {
	return PlanResponse{
		ID: p.ID, Name: p.Name, Type: p.Type, Segment: p.Segment,
		BasePremium: p.BasePremium, Currency: p.Currency, Status: p.Status,
		Description: p.Description, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}

func ToBenefitResponse(b *entity.Benefit) BenefitResponse {
	return BenefitResponse{
		ID: b.ID, PlanID: b.PlanID, Name: b.Name, Category: b.Category,
		AnnualLimit: b.AnnualLimit, CoPayType: b.CoPayType, CoPayValue: b.CoPayValue,
		WaitingPeriodDays: b.WaitingPeriodDays, SubLimitType: b.SubLimitType,
		SubLimitValue: b.SubLimitValue, MinAge: b.MinAge, MaxAge: b.MaxAge,
		WaitingPeriodType: b.WaitingPeriodType, CreatedAt: b.CreatedAt,
	}
}

func ToExclusionResponse(e *entity.Exclusion) ExclusionResponse {
	return ExclusionResponse{
		ID: e.ID, PlanID: e.PlanID, Description: e.Description,
		Type: e.Type, ICDCodes: e.ICDCodes, CreatedAt: e.CreatedAt,
	}
}
