package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/product/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/product/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type planRepository struct {
	store db.Store
}

func NewPlanRepository(store db.Store) domainRepo.PlanRepository {
	return &planRepository{store: store}
}

func (r *planRepository) Create(ctx context.Context, plan *entity.Plan) (*entity.Plan, error) {
	dbPlan, err := r.store.CreatePlan(ctx, db.CreatePlanParams{
		Name:        plan.Name,
		Type:        plan.Type,
		BasePremium: plan.BasePremium,
		Currency:    plan.Currency,
		Status:      plan.Status,
		Description: plan.Description,
		CreatedBy:   uuidToPgtype(plan.CreatedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create plan: %w", err)
	}
	return sqlcPlanToDomain(dbPlan), nil
}

func (r *planRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Plan, error) {
	dbPlan, err := r.store.GetPlanByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get plan by ID: %w", err)
	}
	return sqlcPlanToDomain(dbPlan), nil
}

func (r *planRepository) List(ctx context.Context, limit, offset int) ([]*entity.Plan, error) {
	dbPlans, err := r.store.ListPlans(ctx, db.ListPlansParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list plans: %w", err)
	}
	plans := make([]*entity.Plan, len(dbPlans))
	for i, p := range dbPlans {
		plans[i] = sqlcPlanToDomain(p)
	}
	return plans, nil
}

func (r *planRepository) ListByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.Plan, error) {
	dbPlans, err := r.store.ListPlansByStatus(ctx, db.ListPlansByStatusParams{
		Status: status,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list plans by status: %w", err)
	}
	plans := make([]*entity.Plan, len(dbPlans))
	for i, p := range dbPlans {
		plans[i] = sqlcPlanToDomain(p)
	}
	return plans, nil
}

func (r *planRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.store.CountPlans(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count plans: %w", err)
	}
	return count, nil
}

func (r *planRepository) Update(ctx context.Context, plan *entity.Plan) (*entity.Plan, error) {
	dbPlan, err := r.store.UpdatePlan(ctx, db.UpdatePlanParams{
		ID:          plan.ID,
		Name:        stringToPgtypeText(plan.Name),
		Type:        stringToPgtypeText(plan.Type),
		BasePremium: int64ToPgtypeInt8(plan.BasePremium),
		Description: stringToPgtypeText(plan.Description),
		Status:      stringToPgtypeText(plan.Status),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update plan: %w", err)
	}
	return sqlcPlanToDomain(dbPlan), nil
}

func sqlcPlanToDomain(p db.Plan) *entity.Plan {
	return &entity.Plan{
		ID:          p.ID,
		Name:        p.Name,
		Type:        p.Type,
		BasePremium: p.BasePremium,
		Currency:    p.Currency,
		Status:      p.Status,
		Description: p.Description,
		CreatedBy:   pgtypeToUUID(p.CreatedBy),
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
