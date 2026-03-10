package repository

import (
	"context"
	"github.com/bitbiz/hias-core/domains/product/entity"
	"github.com/google/uuid"
)

type PlanRepository interface {
	Create(ctx context.Context, plan *entity.Plan) (*entity.Plan, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Plan, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Plan, error)
	ListByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.Plan, error)
	ListBySegment(ctx context.Context, segment string, limit, offset int) ([]*entity.Plan, error)
	Count(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, status string) (int64, error)
	Update(ctx context.Context, plan *entity.Plan) (*entity.Plan, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
