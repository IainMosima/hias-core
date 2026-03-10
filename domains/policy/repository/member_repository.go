package repository

import (
	"context"
	"github.com/bitbiz/hias-core/domains/policy/entity"
	"github.com/google/uuid"
)

type MemberRepository interface {
	Create(ctx context.Context, member *entity.Member) (*entity.Member, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Member, error)
	GetByNumber(ctx context.Context, number string) (*entity.Member, error)
	GetByNationalID(ctx context.Context, nationalID string) (*entity.Member, error)
	ListByPolicy(ctx context.Context, policyID uuid.UUID) ([]*entity.Member, error)
	ListActiveByPolicy(ctx context.Context, policyID uuid.UUID) ([]*entity.Member, error)
	CountByPolicy(ctx context.Context, policyID uuid.UUID) (int64, error)
	CountActiveByPolicy(ctx context.Context, policyID uuid.UUID) (int64, error)
	Verify(ctx context.Context, id uuid.UUID) (*entity.Member, error)
	Update(ctx context.Context, member *entity.Member) (*entity.Member, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Member, error)
	ActivatePendingByPolicy(ctx context.Context, policyID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListFiltered(ctx context.Context, search string, limit, offset int) ([]*entity.Member, error)
	CountFiltered(ctx context.Context, search string) (int64, error)
}
