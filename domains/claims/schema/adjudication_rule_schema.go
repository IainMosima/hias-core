package schema

import (
	"encoding/json"
	"time"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	"github.com/google/uuid"
)

type CreateAdjudicationRuleRequest struct {
	Name       string          `json:"name" binding:"required"`
	RuleType   string          `json:"rule_type" binding:"required"`
	Parameters json.RawMessage `json:"parameters"`
	Priority   int             `json:"priority"`
	IsActive   bool            `json:"is_active"`
}

type UpdateAdjudicationRuleRequest struct {
	Name       string          `json:"name"`
	RuleType   string          `json:"rule_type"`
	Parameters json.RawMessage `json:"parameters"`
	Priority   int             `json:"priority"`
	IsActive   *bool           `json:"is_active"`
}

type AdjudicationRuleResponse struct {
	ID         uuid.UUID       `json:"id"`
	Name       string          `json:"name"`
	RuleType   string          `json:"rule_type"`
	Parameters json.RawMessage `json:"parameters"`
	Priority   int             `json:"priority"`
	IsActive   bool            `json:"is_active"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

func ToAdjudicationRuleResponse(r *entity.AdjudicationRule) AdjudicationRuleResponse {
	return AdjudicationRuleResponse{
		ID: r.ID, Name: r.Name, RuleType: r.RuleType,
		Parameters: r.Parameters, Priority: r.Priority,
		IsActive: r.IsActive, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
}
