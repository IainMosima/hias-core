package repository

import (
	"context"
	"github.com/bitbiz/hias-core/domains/claims/entity"
	"github.com/google/uuid"
)

type ClaimRepository interface {
	Create(ctx context.Context, claim *entity.Claim) (*entity.Claim, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Claim, error)
	GetByNumber(ctx context.Context, number string) (*entity.Claim, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Claim, error)
	ListByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.Claim, error)
	ListByProvider(ctx context.Context, providerID uuid.UUID, limit, offset int) ([]*entity.Claim, error)
	ListByMember(ctx context.Context, memberID uuid.UUID, limit, offset int) ([]*entity.Claim, error)
	ListByPolicy(ctx context.Context, policyID uuid.UUID, limit, offset int) ([]*entity.Claim, error)
	Count(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, status string) (int64, error)
	GetForAdjudication(ctx context.Context, limit int) ([]*entity.Claim, error)
	GetApprovedForRemittance(ctx context.Context, providerID uuid.UUID) ([]*entity.Claim, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Claim, error)
	UpdateAmounts(ctx context.Context, id uuid.UUID, approved, copay, memberResp int64) (*entity.Claim, error)
	Reject(ctx context.Context, id uuid.UUID, reason string) (*entity.Claim, error)
	GetApprovedAmountForBenefitThisYear(ctx context.Context, memberID uuid.UUID, category string) (int64, error)
}
