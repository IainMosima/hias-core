package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/reinsurance/entity"
	"github.com/google/uuid"
)

type BordereauRepository interface {
	Create(ctx context.Context, bordereau *entity.Bordereau) (*entity.Bordereau, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Bordereau, error)
	GetByNumber(ctx context.Context, number string) (*entity.Bordereau, error)
	ListByTreaty(ctx context.Context, treatyID uuid.UUID, limit, offset int) ([]*entity.Bordereau, error)
	Update(ctx context.Context, bordereau *entity.Bordereau) (*entity.Bordereau, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Bordereau, error)
	Count(ctx context.Context) (int64, error)
}
