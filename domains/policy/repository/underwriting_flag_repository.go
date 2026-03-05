package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/policy/entity"
	"github.com/google/uuid"
)

type UnderwritingFlagRepository interface {
	Create(ctx context.Context, flag *entity.UnderwritingFlag) (*entity.UnderwritingFlag, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.UnderwritingFlag, error)
	ListByPolicy(ctx context.Context, policyID uuid.UUID) ([]*entity.UnderwritingFlag, error)
	ListByMember(ctx context.Context, memberID uuid.UUID) ([]*entity.UnderwritingFlag, error)
	ListByAssessment(ctx context.Context, assessmentID uuid.UUID) ([]*entity.UnderwritingFlag, error)
	ListOpen(ctx context.Context, limit, offset int32) ([]*entity.UnderwritingFlag, error)
	CountOpen(ctx context.Context) (int64, error)
	Resolve(ctx context.Context, id uuid.UUID, resolvedBy uuid.UUID, resolution string) (*entity.UnderwritingFlag, error)
	Override(ctx context.Context, id uuid.UUID, resolvedBy uuid.UUID, resolution string) (*entity.UnderwritingFlag, error)
}
