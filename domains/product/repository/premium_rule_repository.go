package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/product/entity"
	"github.com/google/uuid"
)

type PremiumRuleRepository interface {
	Create(ctx context.Context, rule *entity.PremiumRule) (*entity.PremiumRule, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.PremiumRule, error)
	ListByPlan(ctx context.Context, planID uuid.UUID) ([]*entity.PremiumRule, error)
	Update(ctx context.Context, rule *entity.PremiumRule) (*entity.PremiumRule, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
