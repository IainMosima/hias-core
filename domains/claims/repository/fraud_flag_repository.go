package repository

import (
	"context"
	"github.com/bitbiz/hias-core/domains/claims/entity"
	"github.com/google/uuid"
)

type FraudFlagRepository interface {
	Create(ctx context.Context, flag *entity.FraudFlag) (*entity.FraudFlag, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.FraudFlag, error)
	ListByClaim(ctx context.Context, claimID uuid.UUID) ([]*entity.FraudFlag, error)
	ListUnresolved(ctx context.Context, limit, offset int) ([]*entity.FraudFlag, error)
	Resolve(ctx context.Context, id uuid.UUID, resolvedBy uuid.UUID) (*entity.FraudFlag, error)
	CheckDuplicate(ctx context.Context, claimNumber string, excludeID uuid.UUID) (int64, error)
	CheckFrequency(ctx context.Context, memberID, providerID uuid.UUID, procedureCode string, excludeID uuid.UUID) (int64, error)
}
