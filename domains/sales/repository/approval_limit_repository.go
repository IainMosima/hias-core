package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/sales/entity"
	"github.com/google/uuid"
)

type ApprovalLimitRepository interface {
	GetByRole(ctx context.Context, roleName string) (*entity.ApprovalLimit, error)
	List(ctx context.Context) ([]*entity.ApprovalLimit, error)
	Create(ctx context.Context, limit *entity.ApprovalLimit) (*entity.ApprovalLimit, error)
	Update(ctx context.Context, limit *entity.ApprovalLimit) (*entity.ApprovalLimit, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
