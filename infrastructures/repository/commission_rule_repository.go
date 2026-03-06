package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/bitbiz/hias-core/domains/billing/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/billing/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type commissionRuleRepository struct {
	store db.Store
}

func NewCommissionRuleRepository(store db.Store) domainRepo.CommissionRuleRepository {
	return &commissionRuleRepository{store: store}
}

func (r *commissionRuleRepository) Create(ctx context.Context, rule *entity.CommissionRule) (*entity.CommissionRule, error) {
	var createdBy pgtype.UUID
	if rule.CreatedBy != uuid.Nil {
		createdBy = pgtype.UUID{Bytes: rule.CreatedBy, Valid: true}
	}
	var effectiveTo pgtype.Timestamptz
	if rule.EffectiveTo != nil {
		effectiveTo = pgtype.Timestamptz{Time: *rule.EffectiveTo, Valid: true}
	}
	dbRule, err := r.store.CreateCommissionRule(ctx, db.CreateCommissionRuleParams{
		PlanID:         rule.PlanID,
		IntermediaryID: rule.IntermediaryID,
		RatePct:        pgtype.Numeric{Valid: true},
		FlatAmount:     rule.FlatAmount,
		EffectiveFrom:  rule.EffectiveFrom,
		EffectiveTo:    effectiveTo,
		CreatedBy:      createdBy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create commission rule: %w", err)
	}
	return sqlcCommissionRuleToDomain(dbRule), nil
}

func (r *commissionRuleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.CommissionRule, error) {
	dbRule, err := r.store.GetCommissionRuleByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get commission rule: %w", err)
	}
	return sqlcCommissionRuleToDomain(dbRule), nil
}

func (r *commissionRuleRepository) ListByPlan(ctx context.Context, planID uuid.UUID) ([]*entity.CommissionRule, error) {
	dbRules, err := r.store.ListCommissionRulesByPlan(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("failed to list commission rules: %w", err)
	}
	rules := make([]*entity.CommissionRule, len(dbRules))
	for i, r := range dbRules {
		rules[i] = sqlcCommissionRuleToDomain(r)
	}
	return rules, nil
}

func sqlcCommissionRuleToDomain(r db.CommissionRule) *entity.CommissionRule {
	var createdBy uuid.UUID
	if r.CreatedBy.Valid {
		createdBy = r.CreatedBy.Bytes
	}
	var effectiveTo *time.Time
	if r.EffectiveTo.Valid {
		effectiveTo = &r.EffectiveTo.Time
	}
	return &entity.CommissionRule{
		ID: r.ID, PlanID: r.PlanID, IntermediaryID: r.IntermediaryID,
		FlatAmount: r.FlatAmount, EffectiveFrom: r.EffectiveFrom,
		EffectiveTo: effectiveTo, CreatedBy: createdBy,
		CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
}
