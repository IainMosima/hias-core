package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/policy/entity"
	"github.com/google/uuid"
)

type UnderwritingRuleRepository interface {
	Create(ctx context.Context, rule *entity.UnderwritingRule) (*entity.UnderwritingRule, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.UnderwritingRule, error)
	ListByPlan(ctx context.Context, planID uuid.UUID) ([]*entity.UnderwritingRule, error)
	ListActiveByPlan(ctx context.Context, planID uuid.UUID) ([]*entity.UnderwritingRule, error)
	Update(ctx context.Context, rule *entity.UnderwritingRule) (*entity.UnderwritingRule, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
