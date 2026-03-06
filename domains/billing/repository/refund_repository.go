package repository

import (
	"context"
	"github.com/bitbiz/hias-core/domains/billing/entity"
	"github.com/google/uuid"
)

type RefundRepository interface {
	Create(ctx context.Context, refund *entity.Refund) (*entity.Refund, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Refund, error)
	ListByPolicy(ctx context.Context, policyID uuid.UUID, limit, offset int) ([]*entity.Refund, error)
	Approve(ctx context.Context, id uuid.UUID, approvedBy uuid.UUID) (*entity.Refund, error)
	Process(ctx context.Context, id uuid.UUID) (*entity.Refund, error)
}
