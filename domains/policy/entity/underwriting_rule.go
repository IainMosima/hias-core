package entity

import (
	"time"

	"github.com/google/uuid"
)

type UnderwritingRule struct {
	ID              uuid.UUID `json:"id"`
	PlanID          uuid.UUID `json:"plan_id"`
	RuleType        string    `json:"rule_type"`
	Relationship    string    `json:"relationship"`
	ParameterKey    string    `json:"parameter_key"`
	ParameterValue  string    `json:"parameter_value"`
	Severity        string    `json:"severity"`
	RiskScoreWeight int       `json:"risk_score_weight"`
	IsBlocking      bool      `json:"is_blocking"`
	IsActive        bool      `json:"is_active"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
