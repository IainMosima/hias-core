package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/sales/entity"
	"github.com/google/uuid"
)

type LeadRepository interface {
	Create(ctx context.Context, lead *entity.Lead) (*entity.Lead, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Lead, error)
	GetByNumber(ctx context.Context, number string) (*entity.Lead, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Lead, error)
	ListByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.Lead, error)
	ListByAssignedTo(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Lead, error)
	ListDueFollowUps(ctx context.Context, limit, offset int) ([]*entity.Lead, error)
	Update(ctx context.Context, lead *entity.Lead) (*entity.Lead, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Lead, error)
	Count(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, status string) (int64, error)
}
