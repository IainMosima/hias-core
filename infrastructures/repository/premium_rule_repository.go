package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/product/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/product/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type premiumRuleRepository struct {
	store db.Store
}

func NewPremiumRuleRepository(store db.Store) domainRepo.PremiumRuleRepository {
	return &premiumRuleRepository{store: store}
}

func (r *premiumRuleRepository) Create(ctx context.Context, rule *entity.PremiumRule) (*entity.PremiumRule, error) {
	dbRule, err := r.store.CreatePremiumRule(ctx, db.CreatePremiumRuleParams{
		PlanID:          rule.PlanID,
		CalculationType: rule.CalculationType,
		Relationship:    stringToPgtypeText(rule.Relationship),
		RateAmount:      rule.RateAmount,
		DiscountType:    stringToPgtypeText(rule.DiscountType),
		DiscountValue:   rule.DiscountValue,
		MinMembers:      int32(rule.MinMembers),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create premium rule: %w", err)
	}
	return sqlcPremiumRuleToDomain(dbRule), nil
}

func (r *premiumRuleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.PremiumRule, error) {
	dbRule, err := r.store.GetPremiumRuleByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get premium rule by ID: %w", err)
	}
	return sqlcPremiumRuleToDomain(dbRule), nil
}

func (r *premiumRuleRepository) ListByPlan(ctx context.Context, planID uuid.UUID) ([]*entity.PremiumRule, error) {
	dbRules, err := r.store.ListPremiumRulesByPlan(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("failed to list premium rules by plan: %w", err)
	}
	rules := make([]*entity.PremiumRule, len(dbRules))
	for i, r := range dbRules {
		rules[i] = sqlcPremiumRuleToDomain(r)
	}
	return rules, nil
}

func (r *premiumRuleRepository) Update(ctx context.Context, rule *entity.PremiumRule) (*entity.PremiumRule, error) {
	dbRule, err := r.store.UpdatePremiumRule(ctx, db.UpdatePremiumRuleParams{
		ID:              rule.ID,
		CalculationType: stringToPgtypeText(rule.CalculationType),
		Relationship:    stringToPgtypeText(rule.Relationship),
		RateAmount:      int64ToPgtypeInt8(rule.RateAmount),
		DiscountType:    stringToPgtypeText(rule.DiscountType),
		DiscountValue:   int64ToPgtypeInt8(rule.DiscountValue),
		MinMembers:      intToPgtypeInt4(rule.MinMembers),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update premium rule: %w", err)
	}
	return sqlcPremiumRuleToDomain(dbRule), nil
}

func (r *premiumRuleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.store.DeletePremiumRule(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete premium rule: %w", err)
	}
	return nil
}

func sqlcPremiumRuleToDomain(r db.PremiumRule) *entity.PremiumRule {
	return &entity.PremiumRule{
		ID:              r.ID,
		PlanID:          r.PlanID,
		CalculationType: r.CalculationType,
		Relationship:    r.Relationship.String,
		RateAmount:      r.RateAmount,
		DiscountType:    r.DiscountType.String,
		DiscountValue:   r.DiscountValue,
		MinMembers:      int(r.MinMembers),
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}
