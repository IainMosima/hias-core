package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/policy/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/policy/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type underwritingRuleRepository struct {
	store db.Store
}

func NewUnderwritingRuleRepository(store db.Store) domainRepo.UnderwritingRuleRepository {
	return &underwritingRuleRepository{store: store}
}

func (r *underwritingRuleRepository) Create(ctx context.Context, rule *entity.UnderwritingRule) (*entity.UnderwritingRule, error) {
	dbR, err := r.store.CreateUnderwritingRule(ctx, db.CreateUnderwritingRuleParams{
		PlanID:          rule.PlanID,
		RuleType:        rule.RuleType,
		Relationship:    stringToPgtypeText(rule.Relationship),
		ParameterKey:    rule.ParameterKey,
		ParameterValue:  rule.ParameterValue,
		Severity:        rule.Severity,
		RiskScoreWeight: int32(rule.RiskScoreWeight),
		IsBlocking:      rule.IsBlocking,
		IsActive:        rule.IsActive,
		Description:     stringToPgtypeText(rule.Description),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create underwriting rule: %w", err)
	}
	return sqlcUnderwritingRuleToDomain(dbR), nil
}

func (r *underwritingRuleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.UnderwritingRule, error) {
	dbR, err := r.store.GetUnderwritingRuleByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get underwriting rule: %w", err)
	}
	return sqlcUnderwritingRuleToDomain(dbR), nil
}

func (r *underwritingRuleRepository) ListByPlan(ctx context.Context, planID uuid.UUID) ([]*entity.UnderwritingRule, error) {
	dbRs, err := r.store.ListUnderwritingRulesByPlan(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("failed to list underwriting rules: %w", err)
	}
	return sqlcUnderwritingRulesToDomain(dbRs), nil
}

func (r *underwritingRuleRepository) ListActiveByPlan(ctx context.Context, planID uuid.UUID) ([]*entity.UnderwritingRule, error) {
	dbRs, err := r.store.ListActiveUnderwritingRulesByPlan(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("failed to list active underwriting rules: %w", err)
	}
	return sqlcUnderwritingRulesToDomain(dbRs), nil
}

func (r *underwritingRuleRepository) Update(ctx context.Context, rule *entity.UnderwritingRule) (*entity.UnderwritingRule, error) {
	dbR, err := r.store.UpdateUnderwritingRule(ctx, db.UpdateUnderwritingRuleParams{
		ID:              rule.ID,
		RuleType:        stringToPgtypeText(rule.RuleType),
		Relationship:    stringToPgtypeText(rule.Relationship),
		ParameterKey:    stringToPgtypeText(rule.ParameterKey),
		ParameterValue:  stringToPgtypeText(rule.ParameterValue),
		Severity:        stringToPgtypeText(rule.Severity),
		RiskScoreWeight: pgtype.Int4{Int32: int32(rule.RiskScoreWeight), Valid: true},
		IsBlocking:      pgtype.Bool{Bool: rule.IsBlocking, Valid: true},
		IsActive:        pgtype.Bool{Bool: rule.IsActive, Valid: true},
		Description:     stringToPgtypeText(rule.Description),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update underwriting rule: %w", err)
	}
	return sqlcUnderwritingRuleToDomain(dbR), nil
}

func (r *underwritingRuleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.store.DeleteUnderwritingRule(ctx, id); err != nil {
		return fmt.Errorf("failed to delete underwriting rule: %w", err)
	}
	return nil
}

func sqlcUnderwritingRuleToDomain(r db.UnderwritingRule) *entity.UnderwritingRule {
	return &entity.UnderwritingRule{
		ID:              r.ID,
		PlanID:          r.PlanID,
		RuleType:        r.RuleType,
		Relationship:    r.Relationship.String,
		ParameterKey:    r.ParameterKey,
		ParameterValue:  r.ParameterValue,
		Severity:        r.Severity,
		RiskScoreWeight: int(r.RiskScoreWeight),
		IsBlocking:      r.IsBlocking,
		IsActive:        r.IsActive,
		Description:     r.Description.String,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}

func sqlcUnderwritingRulesToDomain(rules []db.UnderwritingRule) []*entity.UnderwritingRule {
	result := make([]*entity.UnderwritingRule, len(rules))
	for i, r := range rules {
		result[i] = sqlcUnderwritingRuleToDomain(r)
	}
	return result
}
