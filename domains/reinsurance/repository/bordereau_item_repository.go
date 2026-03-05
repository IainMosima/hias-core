package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/reinsurance/entity"
	"github.com/google/uuid"
)

type BordereauItemRepository interface {
	Create(ctx context.Context, item *entity.BordereauItem) (*entity.BordereauItem, error)
	ListByBordereau(ctx context.Context, bordereauID uuid.UUID) ([]*entity.BordereauItem, error)
	DeleteByBordereau(ctx context.Context, bordereauID uuid.UUID) error
}
