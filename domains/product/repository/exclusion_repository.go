package repository

import (
	"context"
	"github.com/bitbiz/hias-core/domains/product/entity"
	"github.com/google/uuid"
)

type ExclusionRepository interface {
	Create(ctx context.Context, exclusion *entity.Exclusion) (*entity.Exclusion, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Exclusion, error)
	ListByPlan(ctx context.Context, planID uuid.UUID) ([]*entity.Exclusion, error)
	ListByType(ctx context.Context, planID uuid.UUID, exclusionType string) ([]*entity.Exclusion, error)
	Update(ctx context.Context, exclusion *entity.Exclusion) (*entity.Exclusion, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
