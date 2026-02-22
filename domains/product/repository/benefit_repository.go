package repository

import (
	"context"
	"github.com/bitbiz/hias-core/domains/product/entity"
	"github.com/google/uuid"
)

type BenefitRepository interface {
	Create(ctx context.Context, benefit *entity.Benefit) (*entity.Benefit, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Benefit, error)
	ListByPlan(ctx context.Context, planID uuid.UUID) ([]*entity.Benefit, error)
	ListByCategory(ctx context.Context, planID uuid.UUID, category string) ([]*entity.Benefit, error)
	Update(ctx context.Context, benefit *entity.Benefit) (*entity.Benefit, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
