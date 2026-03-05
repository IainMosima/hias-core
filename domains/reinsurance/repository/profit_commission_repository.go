package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/reinsurance/entity"
	"github.com/google/uuid"
)

type ProfitCommissionRepository interface {
	Create(ctx context.Context, pc *entity.ProfitCommission) (*entity.ProfitCommission, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ProfitCommission, error)
	ListByTreaty(ctx context.Context, treatyID uuid.UUID) ([]*entity.ProfitCommission, error)
	Update(ctx context.Context, pc *entity.ProfitCommission) (*entity.ProfitCommission, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
