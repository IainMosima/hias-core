package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/reinsurance/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/reinsurance/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type profitCommissionRepository struct {
	store db.Store
}

func NewProfitCommissionRepository(store db.Store) domainRepo.ProfitCommissionRepository {
	return &profitCommissionRepository{store: store}
}

func (r *profitCommissionRepository) Create(ctx context.Context, pc *entity.ProfitCommission) (*entity.ProfitCommission, error) {
	dbPC, err := r.store.CreateProfitCommission(ctx, db.CreateProfitCommissionParams{
		TreatyID:            pc.TreatyID,
		CommissionType:      pc.CommissionType,
		LossRatioFrom:       float64ToPgNumeric(pc.LossRatioFrom),
		LossRatioTo:         float64ToPgNumeric(pc.LossRatioTo),
		CommissionRate:      float64ToPgNumeric(pc.CommissionRate),
		CarryForwardYears:   int32(pc.CarryForwardYears),
		CarryForwardBalance: pc.CarryForwardBalance,
		PeriodStart:         timePtrToPgtypeDate(pc.PeriodStart),
		PeriodEnd:           timePtrToPgtypeDate(pc.PeriodEnd),
		CalculatedAmount:    pc.CalculatedAmount,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create profit commission: %w", err)
	}
	return sqlcProfitCommissionToDomain(dbPC), nil
}

func (r *profitCommissionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ProfitCommission, error) {
	dbPC, err := r.store.GetProfitCommissionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get profit commission by ID: %w", err)
	}
	return sqlcProfitCommissionToDomain(dbPC), nil
}

func (r *profitCommissionRepository) ListByTreaty(ctx context.Context, treatyID uuid.UUID) ([]*entity.ProfitCommission, error) {
	dbPCs, err := r.store.ListProfitCommissionsByTreaty(ctx, treatyID)
	if err != nil {
		return nil, fmt.Errorf("failed to list profit commissions by treaty: %w", err)
	}
	pcs := make([]*entity.ProfitCommission, len(dbPCs))
	for i, pc := range dbPCs {
		pcs[i] = sqlcProfitCommissionToDomain(pc)
	}
	return pcs, nil
}

func (r *profitCommissionRepository) Update(ctx context.Context, pc *entity.ProfitCommission) (*entity.ProfitCommission, error) {
	dbPC, err := r.store.UpdateProfitCommission(ctx, db.UpdateProfitCommissionParams{
		ID:                  pc.ID,
		CommissionType:      pc.CommissionType,
		LossRatioFrom:       float64ToPgNumeric(pc.LossRatioFrom),
		LossRatioTo:         float64ToPgNumeric(pc.LossRatioTo),
		CommissionRate:      float64ToPgNumeric(pc.CommissionRate),
		CarryForwardYears:   int32(pc.CarryForwardYears),
		CarryForwardBalance: pc.CarryForwardBalance,
		PeriodStart:         timePtrToPgtypeDate(pc.PeriodStart),
		PeriodEnd:           timePtrToPgtypeDate(pc.PeriodEnd),
		CalculatedAmount:    pc.CalculatedAmount,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update profit commission: %w", err)
	}
	return sqlcProfitCommissionToDomain(dbPC), nil
}

func (r *profitCommissionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.store.DeleteProfitCommission(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete profit commission: %w", err)
	}
	return nil
}

func sqlcProfitCommissionToDomain(pc db.ProfitCommission) *entity.ProfitCommission {
	return &entity.ProfitCommission{
		ID:                  pc.ID,
		TreatyID:            pc.TreatyID,
		CommissionType:      pc.CommissionType,
		LossRatioFrom:       pgNumericToFloat64(pc.LossRatioFrom),
		LossRatioTo:         pgNumericToFloat64(pc.LossRatioTo),
		CommissionRate:      pgNumericToFloat64(pc.CommissionRate),
		CarryForwardYears:   int(pc.CarryForwardYears),
		CarryForwardBalance: pc.CarryForwardBalance,
		PeriodStart:         pgtypeDateToTimePtr(pc.PeriodStart),
		PeriodEnd:           pgtypeDateToTimePtr(pc.PeriodEnd),
		CalculatedAmount:    pc.CalculatedAmount,
		CreatedAt:           pc.CreatedAt,
		UpdatedAt:           pc.UpdatedAt,
	}
}
