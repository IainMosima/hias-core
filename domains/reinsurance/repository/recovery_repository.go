package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/reinsurance/entity"
	"github.com/google/uuid"
)

type AgedRecoveryBucket struct {
	Bucket           string `json:"bucket"`
	Count            int64  `json:"count"`
	TotalOutstanding int64  `json:"total_outstanding"`
}

type RecoveryRepository interface {
	Create(ctx context.Context, recovery *entity.ReinsuranceRecovery) (*entity.ReinsuranceRecovery, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ReinsuranceRecovery, error)
	GetByNumber(ctx context.Context, number string) (*entity.ReinsuranceRecovery, error)
	ListByClaim(ctx context.Context, claimID uuid.UUID, limit, offset int) ([]*entity.ReinsuranceRecovery, error)
	ListByTreaty(ctx context.Context, treatyID uuid.UUID, limit, offset int) ([]*entity.ReinsuranceRecovery, error)
	ListOutstanding(ctx context.Context, limit, offset int) ([]*entity.ReinsuranceRecovery, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status, workflowStatus string) (*entity.ReinsuranceRecovery, error)
	UpdateRecoveredAmount(ctx context.Context, id uuid.UUID, recoveredAmount, outstandingAmount int64) (*entity.ReinsuranceRecovery, error)
	GetTotalRecoverableByTreaty(ctx context.Context, treatyID uuid.UUID) (int64, error)
	GetTotalRecoveredByTreaty(ctx context.Context, treatyID uuid.UUID) (int64, error)
	Count(ctx context.Context) (int64, error)
	CountOutstanding(ctx context.Context) (int64, error)
	GetAgedAnalysis(ctx context.Context) ([]AgedRecoveryBucket, error)
	GetTotalRecoverableAmountAll(ctx context.Context) (int64, error)
	GetTotalRecoveredAmountAll(ctx context.Context) (int64, error)
}
