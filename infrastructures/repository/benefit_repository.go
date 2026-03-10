package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/product/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/product/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type benefitRepository struct {
	store db.Store
}

func NewBenefitRepository(store db.Store) domainRepo.BenefitRepository {
	return &benefitRepository{store: store}
}

func (r *benefitRepository) Create(ctx context.Context, benefit *entity.Benefit) (*entity.Benefit, error) {
	dbBenefit, err := r.store.CreateBenefit(ctx, db.CreateBenefitParams{
		PlanID:            benefit.PlanID,
		Name:              benefit.Name,
		Category:          benefit.Category,
		AnnualLimit:       benefit.AnnualLimit,
		CoPayType:         benefit.CoPayType,
		CoPayValue:        benefit.CoPayValue,
		WaitingPeriodDays: int32(benefit.WaitingPeriodDays),
		SubLimitType:      benefit.SubLimitType,
		SubLimitValue:     benefit.SubLimitValue,
		MinAge:            int32(benefit.MinAge),
		MaxAge:            int32(benefit.MaxAge),
		WaitingPeriodType: benefit.WaitingPeriodType,
		DeductibleAmount:  benefit.DeductibleAmount,
		IsOptional:        benefit.IsOptional,
		AddonPremium:      benefit.AddonPremium,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create benefit: %w", err)
	}
	return sqlcBenefitToDomain(dbBenefit), nil
}

func (r *benefitRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Benefit, error) {
	dbBenefit, err := r.store.GetBenefitByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get benefit by ID: %w", err)
	}
	return sqlcBenefitToDomain(dbBenefit), nil
}

func (r *benefitRepository) ListByPlan(ctx context.Context, planID uuid.UUID) ([]*entity.Benefit, error) {
	dbBenefits, err := r.store.ListBenefitsByPlan(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("failed to list benefits by plan: %w", err)
	}
	benefits := make([]*entity.Benefit, len(dbBenefits))
	for i, b := range dbBenefits {
		benefits[i] = sqlcBenefitToDomain(b)
	}
	return benefits, nil
}

func (r *benefitRepository) ListByCategory(ctx context.Context, planID uuid.UUID, category string) ([]*entity.Benefit, error) {
	dbBenefits, err := r.store.ListBenefitsByCategory(ctx, db.ListBenefitsByCategoryParams{
		PlanID:   planID,
		Category: category,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list benefits by category: %w", err)
	}
	benefits := make([]*entity.Benefit, len(dbBenefits))
	for i, b := range dbBenefits {
		benefits[i] = sqlcBenefitToDomain(b)
	}
	return benefits, nil
}

func (r *benefitRepository) Update(ctx context.Context, benefit *entity.Benefit) (*entity.Benefit, error) {
	dbBenefit, err := r.store.UpdateBenefit(ctx, db.UpdateBenefitParams{
		ID:                benefit.ID,
		Name:              stringToPgtypeText(benefit.Name),
		Category:          stringToPgtypeText(benefit.Category),
		AnnualLimit:       int64ToPgtypeInt8(benefit.AnnualLimit),
		CoPayType:         stringToPgtypeText(benefit.CoPayType),
		CoPayValue:        int64ToPgtypeInt8(benefit.CoPayValue),
		WaitingPeriodDays: intToPgtypeInt4(benefit.WaitingPeriodDays),
		SubLimitType:      stringToPgtypeText(benefit.SubLimitType),
		SubLimitValue:     int64ToPgtypeInt8(benefit.SubLimitValue),
		MinAge:            intToPgtypeInt4(benefit.MinAge),
		MaxAge:            intToPgtypeInt4(benefit.MaxAge),
		WaitingPeriodType: stringToPgtypeText(benefit.WaitingPeriodType),
		IsOptional:        boolToPgtypeBool(benefit.IsOptional),
		AddonPremium:      int64ToPgtypeInt8(benefit.AddonPremium),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update benefit: %w", err)
	}
	return sqlcBenefitToDomain(dbBenefit), nil
}

func (r *benefitRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.store.DeleteBenefit(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete benefit: %w", err)
	}
	return nil
}

func (r *benefitRepository) CreateWithParent(ctx context.Context, benefit *entity.Benefit) (*entity.Benefit, error) {
	dbBenefit, err := r.store.CreateBenefitWithParent(ctx, db.CreateBenefitWithParentParams{
		PlanID:            benefit.PlanID,
		ParentBenefitID:   uuidPtrToPgtype(benefit.ParentBenefitID),
		Name:              benefit.Name,
		Category:          benefit.Category,
		AnnualLimit:       benefit.AnnualLimit,
		CoPayType:         benefit.CoPayType,
		CoPayValue:        benefit.CoPayValue,
		WaitingPeriodDays: int32(benefit.WaitingPeriodDays),
		SubLimitType:      benefit.SubLimitType,
		SubLimitValue:     benefit.SubLimitValue,
		MinAge:            int32(benefit.MinAge),
		MaxAge:            int32(benefit.MaxAge),
		WaitingPeriodType: benefit.WaitingPeriodType,
		DeductibleAmount:  benefit.DeductibleAmount,
		IsOptional:        benefit.IsOptional,
		AddonPremium:      benefit.AddonPremium,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create sub-benefit: %w", err)
	}
	return sqlcBenefitToDomain(dbBenefit), nil
}

func (r *benefitRepository) ListSubBenefits(ctx context.Context, parentID uuid.UUID) ([]*entity.Benefit, error) {
	dbBenefits, err := r.store.ListSubBenefits(ctx, uuidToPgtype(parentID))
	if err != nil {
		return nil, fmt.Errorf("failed to list sub-benefits: %w", err)
	}
	benefits := make([]*entity.Benefit, len(dbBenefits))
	for i, b := range dbBenefits {
		benefits[i] = sqlcBenefitToDomain(b)
	}
	return benefits, nil
}

func sqlcBenefitToDomain(b db.Benefit) *entity.Benefit {
	var parentBenefitID *uuid.UUID
	if b.ParentBenefitID.Valid {
		id := uuid.UUID(b.ParentBenefitID.Bytes)
		parentBenefitID = &id
	}
	return &entity.Benefit{
		ID:                b.ID,
		PlanID:            b.PlanID,
		ParentBenefitID:   parentBenefitID,
		Name:              b.Name,
		Category:          b.Category,
		AnnualLimit:       b.AnnualLimit,
		CoPayType:         b.CoPayType,
		CoPayValue:        b.CoPayValue,
		WaitingPeriodDays: int(b.WaitingPeriodDays),
		SubLimitType:      b.SubLimitType,
		SubLimitValue:     b.SubLimitValue,
		MinAge:            int(b.MinAge),
		MaxAge:            int(b.MaxAge),
		WaitingPeriodType: b.WaitingPeriodType,
		DeductibleAmount:  b.DeductibleAmount,
		IsOptional:        b.IsOptional,
		AddonPremium:      b.AddonPremium,
		CreatedAt:         b.CreatedAt,
		UpdatedAt:         b.UpdatedAt,
	}
}
