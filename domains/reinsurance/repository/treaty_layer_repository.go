package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/reinsurance/entity"
	"github.com/google/uuid"
)

type TreatyLayerRepository interface {
	Create(ctx context.Context, layer *entity.TreatyLayer) (*entity.TreatyLayer, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.TreatyLayer, error)
	ListByTreaty(ctx context.Context, treatyID uuid.UUID) ([]*entity.TreatyLayer, error)
	Update(ctx context.Context, layer *entity.TreatyLayer) (*entity.TreatyLayer, error)
	UpdateAggregateUsed(ctx context.Context, id uuid.UUID, aggregateUsed int64) (*entity.TreatyLayer, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
