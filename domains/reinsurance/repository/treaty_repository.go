package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/reinsurance/entity"
	"github.com/google/uuid"
)

type TreatyRepository interface {
	Create(ctx context.Context, treaty *entity.Treaty) (*entity.Treaty, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Treaty, error)
	GetByNumber(ctx context.Context, number string) (*entity.Treaty, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Treaty, error)
	ListActive(ctx context.Context, limit, offset int) ([]*entity.Treaty, error)
	ListByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.Treaty, error)
	ListByType(ctx context.Context, treatyType string, limit, offset int) ([]*entity.Treaty, error)
	Update(ctx context.Context, treaty *entity.Treaty) (*entity.Treaty, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Treaty, error)
	Count(ctx context.Context) (int64, error)
	ListExpiring(ctx context.Context, withinDays int, limit, offset int) ([]*entity.Treaty, error)
}
