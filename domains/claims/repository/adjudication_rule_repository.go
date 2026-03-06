package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	"github.com/google/uuid"
)

type AdjudicationRuleRepository interface {
	Create(ctx context.Context, rule *entity.AdjudicationRule) (*entity.AdjudicationRule, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.AdjudicationRule, error)
	List(ctx context.Context, limit, offset int) ([]*entity.AdjudicationRule, error)
	ListActive(ctx context.Context) ([]*entity.AdjudicationRule, error)
	Update(ctx context.Context, rule *entity.AdjudicationRule) (*entity.AdjudicationRule, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
