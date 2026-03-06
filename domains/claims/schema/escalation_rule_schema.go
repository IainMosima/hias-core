package schema

import (
	"time"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	"github.com/google/uuid"
)

type CreateEscalationRuleRequest struct {
	Name            string `json:"name" binding:"required"`
	ConditionType   string `json:"condition_type" binding:"required"`
	ThresholdAmount int64  `json:"threshold_amount"`
	EscalationRole  string `json:"escalation_role" binding:"required"`
	IsActive        bool   `json:"is_active"`
}

type UpdateEscalationRuleRequest struct {
	Name            string `json:"name"`
	ConditionType   string `json:"condition_type"`
	ThresholdAmount int64  `json:"threshold_amount"`
	EscalationRole  string `json:"escalation_role"`
	IsActive        *bool  `json:"is_active"`
}

type EscalationRuleResponse struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	ConditionType   string    `json:"condition_type"`
	ThresholdAmount int64     `json:"threshold_amount"`
	EscalationRole  string    `json:"escalation_role"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func ToEscalationRuleResponse(r *entity.EscalationRule) EscalationRuleResponse {
	return EscalationRuleResponse{
		ID: r.ID, Name: r.Name, ConditionType: r.ConditionType,
		ThresholdAmount: r.ThresholdAmount, EscalationRole: r.EscalationRole,
		IsActive: r.IsActive, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
}
