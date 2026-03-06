package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/claims/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type escalationRuleRepository struct {
	store db.Store
}

func NewEscalationRuleRepository(store db.Store) domainRepo.EscalationRuleRepository {
	return &escalationRuleRepository{store: store}
}

func (r *escalationRuleRepository) Create(ctx context.Context, rule *entity.EscalationRule) (*entity.EscalationRule, error) {
	dbRule, err := r.store.CreateEscalationRule(ctx, db.CreateEscalationRuleParams{
		Name:            rule.Name,
		ConditionType:   rule.ConditionType,
		ThresholdAmount: rule.ThresholdAmount,
		EscalationRole:  rule.EscalationRole,
		IsActive:        rule.IsActive,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create escalation rule: %w", err)
	}
	return sqlcEscalationRuleToDomain(dbRule), nil
}

func (r *escalationRuleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.EscalationRule, error) {
	dbRule, err := r.store.GetEscalationRuleByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get escalation rule: %w", err)
	}
	return sqlcEscalationRuleToDomain(dbRule), nil
}

func (r *escalationRuleRepository) List(ctx context.Context, limit, offset int) ([]*entity.EscalationRule, error) {
	dbRules, err := r.store.ListEscalationRules(ctx, db.ListEscalationRulesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list escalation rules: %w", err)
	}
	rules := make([]*entity.EscalationRule, len(dbRules))
	for i, r := range dbRules {
		rules[i] = sqlcEscalationRuleToDomain(r)
	}
	return rules, nil
}

func (r *escalationRuleRepository) ListActive(ctx context.Context) ([]*entity.EscalationRule, error) {
	dbRules, err := r.store.ListActiveEscalationRules(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list active escalation rules: %w", err)
	}
	rules := make([]*entity.EscalationRule, len(dbRules))
	for i, r := range dbRules {
		rules[i] = sqlcEscalationRuleToDomain(r)
	}
	return rules, nil
}

func (r *escalationRuleRepository) Update(ctx context.Context, rule *entity.EscalationRule) (*entity.EscalationRule, error) {
	dbRule, err := r.store.UpdateEscalationRule(ctx, db.UpdateEscalationRuleParams{
		ID:              rule.ID,
		Name:            rule.Name,
		ConditionType:   rule.ConditionType,
		ThresholdAmount: rule.ThresholdAmount,
		EscalationRole:  rule.EscalationRole,
		IsActive:        rule.IsActive,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update escalation rule: %w", err)
	}
	return sqlcEscalationRuleToDomain(dbRule), nil
}

func (r *escalationRuleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.store.DeleteEscalationRule(ctx, id)
}

func sqlcEscalationRuleToDomain(r db.EscalationRule) *entity.EscalationRule {
	return &entity.EscalationRule{
		ID:              r.ID,
		Name:            r.Name,
		ConditionType:   r.ConditionType,
		ThresholdAmount: r.ThresholdAmount,
		EscalationRole:  r.EscalationRole,
		IsActive:        r.IsActive,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}
