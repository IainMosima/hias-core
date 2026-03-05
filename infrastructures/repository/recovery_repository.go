package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/reinsurance/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/reinsurance/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type recoveryRepository struct {
	store db.Store
}

func NewRecoveryRepository(store db.Store) domainRepo.RecoveryRepository {
	return &recoveryRepository{store: store}
}

func (r *recoveryRepository) Create(ctx context.Context, recovery *entity.ReinsuranceRecovery) (*entity.ReinsuranceRecovery, error) {
	dbRecovery, err := r.store.CreateReinsuranceRecovery(ctx, db.CreateReinsuranceRecoveryParams{
		RecoveryNumber:    recovery.RecoveryNumber,
		ClaimID:           recovery.ClaimID,
		TreatyID:          recovery.TreatyID,
		TreatyLayerID:     uuidToPgtype(recovery.TreatyLayerID),
		CessionID:         uuidToPgtype(recovery.CessionID),
		GrossClaimAmount:  recovery.GrossClaimAmount,
		RecoverableAmount: recovery.RecoverableAmount,
		RecoveredAmount:   recovery.RecoveredAmount,
		OutstandingAmount: recovery.OutstandingAmount,
		Status:            recovery.Status,
		WorkflowStatus:    recovery.WorkflowStatus,
		Notes:             stringToPgtypeText(recovery.Notes),
		CreatedBy:         recovery.CreatedBy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create reinsurance recovery: %w", err)
	}
	return sqlcReinsuranceRecoveryToDomain(dbRecovery), nil
}

func (r *recoveryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ReinsuranceRecovery, error) {
	dbRecovery, err := r.store.GetReinsuranceRecoveryByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get reinsurance recovery by ID: %w", err)
	}
	return sqlcReinsuranceRecoveryToDomain(dbRecovery), nil
}

func (r *recoveryRepository) GetByNumber(ctx context.Context, number string) (*entity.ReinsuranceRecovery, error) {
	dbRecovery, err := r.store.GetReinsuranceRecoveryByNumber(ctx, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get reinsurance recovery by number: %w", err)
	}
	return sqlcReinsuranceRecoveryToDomain(dbRecovery), nil
}

func (r *recoveryRepository) ListByClaim(ctx context.Context, claimID uuid.UUID, limit, offset int) ([]*entity.ReinsuranceRecovery, error) {
	dbRecoveries, err := r.store.ListReinsuranceRecoveriesByClaim(ctx, db.ListReinsuranceRecoveriesByClaimParams{
		ClaimID: claimID,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list reinsurance recoveries by claim: %w", err)
	}
	recoveries := make([]*entity.ReinsuranceRecovery, len(dbRecoveries))
	for i, rec := range dbRecoveries {
		recoveries[i] = sqlcReinsuranceRecoveryToDomain(rec)
	}
	return recoveries, nil
}

func (r *recoveryRepository) ListByTreaty(ctx context.Context, treatyID uuid.UUID, limit, offset int) ([]*entity.ReinsuranceRecovery, error) {
	dbRecoveries, err := r.store.ListReinsuranceRecoveriesByTreaty(ctx, db.ListReinsuranceRecoveriesByTreatyParams{
		TreatyID: treatyID,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list reinsurance recoveries by treaty: %w", err)
	}
	recoveries := make([]*entity.ReinsuranceRecovery, len(dbRecoveries))
	for i, rec := range dbRecoveries {
		recoveries[i] = sqlcReinsuranceRecoveryToDomain(rec)
	}
	return recoveries, nil
}

func (r *recoveryRepository) ListOutstanding(ctx context.Context, limit, offset int) ([]*entity.ReinsuranceRecovery, error) {
	dbRecoveries, err := r.store.ListOutstandingRecoveries(ctx, db.ListOutstandingRecoveriesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list outstanding recoveries: %w", err)
	}
	recoveries := make([]*entity.ReinsuranceRecovery, len(dbRecoveries))
	for i, rec := range dbRecoveries {
		recoveries[i] = sqlcReinsuranceRecoveryToDomain(rec)
	}
	return recoveries, nil
}

func (r *recoveryRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status, workflowStatus string) (*entity.ReinsuranceRecovery, error) {
	dbRecovery, err := r.store.UpdateReinsuranceRecoveryStatus(ctx, db.UpdateReinsuranceRecoveryStatusParams{
		ID:             id,
		Status:         status,
		WorkflowStatus: workflowStatus,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update reinsurance recovery status: %w", err)
	}
	return sqlcReinsuranceRecoveryToDomain(dbRecovery), nil
}

func (r *recoveryRepository) UpdateRecoveredAmount(ctx context.Context, id uuid.UUID, recoveredAmount, outstandingAmount int64) (*entity.ReinsuranceRecovery, error) {
	dbRecovery, err := r.store.UpdateReinsuranceRecoveryAmounts(ctx, db.UpdateReinsuranceRecoveryAmountsParams{
		ID:                id,
		RecoveredAmount:   recoveredAmount,
		OutstandingAmount: outstandingAmount,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update reinsurance recovery amounts: %w", err)
	}
	return sqlcReinsuranceRecoveryToDomain(dbRecovery), nil
}

func (r *recoveryRepository) GetTotalRecoverableByTreaty(ctx context.Context, treatyID uuid.UUID) (int64, error) {
	total, err := r.store.GetTotalRecoverableByTreaty(ctx, treatyID)
	if err != nil {
		return 0, fmt.Errorf("failed to get total recoverable by treaty: %w", err)
	}
	return total, nil
}

func (r *recoveryRepository) GetTotalRecoveredByTreaty(ctx context.Context, treatyID uuid.UUID) (int64, error) {
	total, err := r.store.GetTotalRecoveredByTreaty(ctx, treatyID)
	if err != nil {
		return 0, fmt.Errorf("failed to get total recovered by treaty: %w", err)
	}
	return total, nil
}

func (r *recoveryRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.store.CountReinsuranceRecoveries(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count reinsurance recoveries: %w", err)
	}
	return count, nil
}

func (r *recoveryRepository) CountOutstanding(ctx context.Context) (int64, error) {
	count, err := r.store.CountOutstandingRecoveries(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count outstanding recoveries: %w", err)
	}
	return count, nil
}

func (r *recoveryRepository) GetAgedAnalysis(ctx context.Context) ([]domainRepo.AgedRecoveryBucket, error) {
	rows, err := r.store.GetAgedRecoveryAnalysis(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get aged recovery analysis: %w", err)
	}
	buckets := make([]domainRepo.AgedRecoveryBucket, len(rows))
	for i, row := range rows {
		buckets[i] = domainRepo.AgedRecoveryBucket{
			Bucket:           row.Bucket,
			Count:            row.Count,
			TotalOutstanding: row.TotalOutstanding,
		}
	}
	return buckets, nil
}

func (r *recoveryRepository) GetTotalRecoverableAmountAll(ctx context.Context) (int64, error) {
	total, err := r.store.GetTotalRecoverableAmountAll(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get total recoverable amount: %w", err)
	}
	return total, nil
}

func (r *recoveryRepository) GetTotalRecoveredAmountAll(ctx context.Context) (int64, error) {
	total, err := r.store.GetTotalRecoveredAmountAll(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get total recovered amount: %w", err)
	}
	return total, nil
}

func sqlcReinsuranceRecoveryToDomain(r db.ReinsuranceRecovery) *entity.ReinsuranceRecovery {
	return &entity.ReinsuranceRecovery{
		ID:                r.ID,
		RecoveryNumber:    r.RecoveryNumber,
		ClaimID:           r.ClaimID,
		TreatyID:          r.TreatyID,
		TreatyLayerID:     pgtypeToUUID(r.TreatyLayerID),
		CessionID:         pgtypeToUUID(r.CessionID),
		GrossClaimAmount:  r.GrossClaimAmount,
		RecoverableAmount: r.RecoverableAmount,
		RecoveredAmount:   r.RecoveredAmount,
		OutstandingAmount: r.OutstandingAmount,
		Status:            r.Status,
		WorkflowStatus:    r.WorkflowStatus,
		Notes:             r.Notes.String,
		CreatedBy:         r.CreatedBy,
		CreatedAt:         r.CreatedAt,
		UpdatedAt:         r.UpdatedAt,
	}
}
