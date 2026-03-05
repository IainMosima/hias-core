package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/reinsurance/entity"
	"github.com/google/uuid"
)

type TreatyAlertRepository interface {
	Create(ctx context.Context, alert *entity.TreatyAlert) (*entity.TreatyAlert, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.TreatyAlert, error)
	List(ctx context.Context, limit, offset int) ([]*entity.TreatyAlert, error)
	ListByTreaty(ctx context.Context, treatyID uuid.UUID, limit, offset int) ([]*entity.TreatyAlert, error)
	ListUnacknowledged(ctx context.Context, limit, offset int) ([]*entity.TreatyAlert, error)
	Acknowledge(ctx context.Context, id uuid.UUID, acknowledgedBy uuid.UUID) (*entity.TreatyAlert, error)
	CountUnacknowledged(ctx context.Context) (int64, error)
}
