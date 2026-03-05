package repository

import (
	"context"
	"time"

	"github.com/bitbiz/hias-core/domains/reinsurance/entity"
	"github.com/google/uuid"
)

type CessionRepository interface {
	Create(ctx context.Context, cession *entity.Cession) (*entity.Cession, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Cession, error)
	GetByNumber(ctx context.Context, number string) (*entity.Cession, error)
	ListByTreaty(ctx context.Context, treatyID uuid.UUID, limit, offset int) ([]*entity.Cession, error)
	ListByPolicy(ctx context.Context, policyID uuid.UUID, limit, offset int) ([]*entity.Cession, error)
	ListByTreatyAndPeriod(ctx context.Context, treatyID uuid.UUID, start, end time.Time, limit, offset int) ([]*entity.Cession, error)
	ListBookedByTreatyAndPeriod(ctx context.Context, treatyID uuid.UUID, start, end time.Time) ([]*entity.Cession, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Cession, error)
	GetTotalCededByTreaty(ctx context.Context, treatyID uuid.UUID) (int64, error)
	GetTotalCededByTreatyAndPeriod(ctx context.Context, treatyID uuid.UUID, start, end time.Time) (int64, error)
	Count(ctx context.Context) (int64, error)
	GetTotalCededAmountAll(ctx context.Context) (int64, error)
	GetTotalGrossAmountAll(ctx context.Context) (int64, error)
}
