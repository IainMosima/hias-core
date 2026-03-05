package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/reinsurance/entity"
	"github.com/google/uuid"
)

type ReinsurerStatementRepository interface {
	Create(ctx context.Context, statement *entity.ReinsurerStatement) (*entity.ReinsurerStatement, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ReinsurerStatement, error)
	GetByNumber(ctx context.Context, number string) (*entity.ReinsurerStatement, error)
	ListByTreaty(ctx context.Context, treatyID uuid.UUID, limit, offset int) ([]*entity.ReinsurerStatement, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.ReinsurerStatement, error)
	Update(ctx context.Context, statement *entity.ReinsurerStatement) (*entity.ReinsurerStatement, error)
	Count(ctx context.Context) (int64, error)
}
