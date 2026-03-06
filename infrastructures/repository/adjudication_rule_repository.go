package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/claims/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type adjudicationRuleRepository struct {
	store db.Store
}

func NewAdjudicationRuleRepository(store db.Store) domainRepo.AdjudicationRuleRepository {
	return &adjudicationRuleRepository{store: store}
}

func (r *adjudicationRuleRepository) Create(ctx context.Context, rule *entity.AdjudicationRule) (*entity.AdjudicationRule, error) {
	dbRule, err := r.store.CreateAdjudicationRule(ctx, db.CreateAdjudicationRuleParams{
		Name:       rule.Name,
		RuleType:   rule.RuleType,
		Parameters: rule.Parameters,
		Priority:   int32(rule.Priority),
		IsActive:   rule.IsActive,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create adjudication rule: %w", err)
	}
	return sqlcAdjudicationRuleToDomain(dbRule), nil
}

func (r *adjudicationRuleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.AdjudicationRule, error) {
	dbRule, err := r.store.GetAdjudicationRuleByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get adjudication rule: %w", err)
	}
	return sqlcAdjudicationRuleToDomain(dbRule), nil
}

func (r *adjudicationRuleRepository) List(ctx context.Context, limit, offset int) ([]*entity.AdjudicationRule, error) {
	dbRules, err := r.store.ListAdjudicationRules(ctx, db.ListAdjudicationRulesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list adjudication rules: %w", err)
	}
	rules := make([]*entity.AdjudicationRule, len(dbRules))
	for i, r := range dbRules {
		rules[i] = sqlcAdjudicationRuleToDomain(r)
	}
	return rules, nil
}

func (r *adjudicationRuleRepository) ListActive(ctx context.Context) ([]*entity.AdjudicationRule, error) {
	dbRules, err := r.store.ListActiveAdjudicationRules(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list active adjudication rules: %w", err)
	}
	rules := make([]*entity.AdjudicationRule, len(dbRules))
	for i, r := range dbRules {
		rules[i] = sqlcAdjudicationRuleToDomain(r)
	}
	return rules, nil
}

func (r *adjudicationRuleRepository) Update(ctx context.Context, rule *entity.AdjudicationRule) (*entity.AdjudicationRule, error) {
	dbRule, err := r.store.UpdateAdjudicationRule(ctx, db.UpdateAdjudicationRuleParams{
		ID:         rule.ID,
		Name:       rule.Name,
		RuleType:   rule.RuleType,
		Parameters: rule.Parameters,
		Priority:   int32(rule.Priority),
		IsActive:   rule.IsActive,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update adjudication rule: %w", err)
	}
	return sqlcAdjudicationRuleToDomain(dbRule), nil
}

func (r *adjudicationRuleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.store.DeleteAdjudicationRule(ctx, id)
}

func sqlcAdjudicationRuleToDomain(r db.AdjudicationRule) *entity.AdjudicationRule {
	return &entity.AdjudicationRule{
		ID:         r.ID,
		Name:       r.Name,
		RuleType:   r.RuleType,
		Parameters: r.Parameters,
		Priority:   int(r.Priority),
		IsActive:   r.IsActive,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}
