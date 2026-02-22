package repository

import (
	"context"
	"github.com/bitbiz/hias-core/domains/claims/entity"
	"github.com/google/uuid"
)

type ClaimLineItemRepository interface {
	Create(ctx context.Context, item *entity.ClaimLineItem) (*entity.ClaimLineItem, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ClaimLineItem, error)
	ListByClaim(ctx context.Context, claimID uuid.UUID) ([]*entity.ClaimLineItem, error)
	UpdateApprovedAmount(ctx context.Context, id uuid.UUID, amount int64) (*entity.ClaimLineItem, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
