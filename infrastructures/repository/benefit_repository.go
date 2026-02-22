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

func sqlcBenefitToDomain(b db.Benefit) *entity.Benefit {
	return &entity.Benefit{
		ID:                b.ID,
		PlanID:            b.PlanID,
		Name:              b.Name,
		Category:          b.Category,
		AnnualLimit:       b.AnnualLimit,
		CoPayType:         b.CoPayType,
		CoPayValue:        b.CoPayValue,
		WaitingPeriodDays: int(b.WaitingPeriodDays),
		CreatedAt:         b.CreatedAt,
		UpdatedAt:         b.UpdatedAt,
	}
}
