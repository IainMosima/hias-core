package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/bitbiz/hias-core/domains/policy/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/policy/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type policyRepository struct {
	store db.Store
}

func NewPolicyRepository(store db.Store) domainRepo.PolicyRepository {
	return &policyRepository{store: store}
}

func (r *policyRepository) Create(ctx context.Context, policy *entity.Policy) (*entity.Policy, error) {
	dbPolicy, err := r.store.CreatePolicy(ctx, db.CreatePolicyParams{
		PlanID:            policy.PlanID,
		PolicyholderName:  policy.PolicyholderName,
		PolicyholderEmail: policy.PolicyholderEmail,
		PolicyholderPhone: policy.PolicyholderPhone,
		PolicyNumber:      policy.PolicyNumber,
		Status:            policy.Status,
		StartDate:         timeToPgtypeTimestamptz(policy.StartDate),
		EndDate:           timeToPgtypeTimestamptz(policy.EndDate),
		PremiumAmount:     policy.PremiumAmount,
		Currency:          policy.Currency,
		CreatedBy:         uuidToPgtype(policy.CreatedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create policy: %w", err)
	}
	return sqlcPolicyToDomain(dbPolicy), nil
}

func (r *policyRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Policy, error) {
	dbPolicy, err := r.store.GetPolicyByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get policy by ID: %w", err)
	}
	return sqlcPolicyToDomain(dbPolicy), nil
}

func (r *policyRepository) GetByNumber(ctx context.Context, number string) (*entity.Policy, error) {
	dbPolicy, err := r.store.GetPolicyByNumber(ctx, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get policy by number: %w", err)
	}
	return sqlcPolicyToDomain(dbPolicy), nil
}

func (r *policyRepository) List(ctx context.Context, limit, offset int) ([]*entity.Policy, error) {
	dbPolicies, err := r.store.ListPolicies(ctx, db.ListPoliciesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list policies: %w", err)
	}
	policies := make([]*entity.Policy, len(dbPolicies))
	for i, p := range dbPolicies {
		policies[i] = sqlcPolicyToDomain(p)
	}
	return policies, nil
}

func (r *policyRepository) ListByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.Policy, error) {
	dbPolicies, err := r.store.ListPoliciesByStatus(ctx, db.ListPoliciesByStatusParams{
		Status: status,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list policies by status: %w", err)
	}
	policies := make([]*entity.Policy, len(dbPolicies))
	for i, p := range dbPolicies {
		policies[i] = sqlcPolicyToDomain(p)
	}
	return policies, nil
}

func (r *policyRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.store.CountPolicies(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count policies: %w", err)
	}
	return count, nil
}

func (r *policyRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	count, err := r.store.CountPoliciesByStatus(ctx, status)
	if err != nil {
		return 0, fmt.Errorf("failed to count policies by status: %w", err)
	}
	return count, nil
}

func (r *policyRepository) GetActivePoliciesForBilling(ctx context.Context) ([]*entity.Policy, error) {
	dbPolicies, err := r.store.GetActivePoliciesForBilling(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active policies for billing: %w", err)
	}
	policies := make([]*entity.Policy, len(dbPolicies))
	for i, p := range dbPolicies {
		policies[i] = sqlcPolicyToDomain(p)
	}
	return policies, nil
}

func (r *policyRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Policy, error) {
	dbPolicy, err := r.store.UpdatePolicyStatus(ctx, db.UpdatePolicyStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update policy status: %w", err)
	}
	return sqlcPolicyToDomain(dbPolicy), nil
}

func (r *policyRepository) ActivateWithTimestamp(ctx context.Context, id uuid.UUID) (*entity.Policy, error) {
	dbPolicy, err := r.store.ActivatePolicyWithTimestamp(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to activate policy with timestamp: %w", err)
	}
	return sqlcPolicyToDomain(dbPolicy), nil
}

func (r *policyRepository) Update(ctx context.Context, policy *entity.Policy) (*entity.Policy, error) {
	dbPolicy, err := r.store.UpdatePolicy(ctx, db.UpdatePolicyParams{
		ID:                policy.ID,
		PolicyholderName:  stringToPgtypeText(policy.PolicyholderName),
		PolicyholderEmail: stringToPgtypeText(policy.PolicyholderEmail),
		PolicyholderPhone: stringToPgtypeText(policy.PolicyholderPhone),
		StartDate:         timeToPgtypeTimestamptz(policy.StartDate),
		EndDate:           timeToPgtypeTimestamptz(policy.EndDate),
		PremiumAmount:     int64ToPgtypeInt8(policy.PremiumAmount),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update policy: %w", err)
	}
	return sqlcPolicyToDomain(dbPolicy), nil
}

func (r *policyRepository) GetLapsedForTermination(ctx context.Context) ([]*entity.Policy, error) {
	dbPolicies, err := r.store.GetLapsedPoliciesForTermination(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get lapsed policies for termination: %w", err)
	}
	policies := make([]*entity.Policy, len(dbPolicies))
	for i, p := range dbPolicies {
		policies[i] = sqlcPolicyToDomain(p)
	}
	return policies, nil
}

func (r *policyRepository) GetOverdueForLapse(ctx context.Context) ([]*entity.Policy, error) {
	dbPolicies, err := r.store.GetOverduePoliciesForLapse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue policies for lapse: %w", err)
	}
	policies := make([]*entity.Policy, len(dbPolicies))
	for i, p := range dbPolicies {
		policies[i] = sqlcPolicyToDomain(p)
	}
	return policies, nil
}

func (r *policyRepository) ListExpiringSoon(ctx context.Context, days int) ([]*entity.Policy, error) {
	dbPolicies, err := r.store.ListPoliciesExpiringSoon(ctx, int32(days))
	if err != nil {
		return nil, fmt.Errorf("failed to list expiring policies: %w", err)
	}
	policies := make([]*entity.Policy, len(dbPolicies))
	for i, p := range dbPolicies {
		policies[i] = sqlcPolicyToDomain(p)
	}
	return policies, nil
}

func (r *policyRepository) UpdatePlanAndPremium(ctx context.Context, id uuid.UUID, planID uuid.UUID, premiumAmount int64) (*entity.Policy, error) {
	dbPolicy, err := r.store.UpdatePolicyPlanAndPremium(ctx, db.UpdatePolicyPlanAndPremiumParams{
		ID:            id,
		PlanID:        planID,
		PremiumAmount: premiumAmount,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update policy plan and premium: %w", err)
	}
	return sqlcPolicyToDomain(dbPolicy), nil
}

func sqlcPolicyToDomain(p db.Policy) *entity.Policy {
	var renewedFromID *uuid.UUID
	if p.RenewedFromID.Valid {
		id := uuid.UUID(p.RenewedFromID.Bytes)
		renewedFromID = &id
	}
	return &entity.Policy{
		ID:                p.ID,
		PlanID:            p.PlanID,
		PolicyholderName:  p.PolicyholderName,
		PolicyholderEmail: p.PolicyholderEmail,
		PolicyholderPhone: p.PolicyholderPhone,
		PolicyNumber:      p.PolicyNumber,
		Status:            p.Status,
		StartDate:         pgtypeTimestamptzToTime(p.StartDate),
		EndDate:           pgtypeTimestamptzToTime(p.EndDate),
		PremiumAmount:     p.PremiumAmount,
		Currency:          p.Currency,
		RenewedFromID:     renewedFromID,
		ActivatedAt:       pgtypeTimestamptzToTimePtr(p.ActivatedAt),
		CreatedBy:         pgtypeToUUID(p.CreatedBy),
		CreatedAt:         p.CreatedAt,
		UpdatedAt:         p.UpdatedAt,
	}
}

func (r *policyRepository) ListFiltered(ctx context.Context, dateFrom, dateTo *time.Time, search string, limit, offset int) ([]*entity.Policy, error) {
	dbPolicies, err := r.store.ListPoliciesFiltered(ctx, db.ListPoliciesFilteredParams{
		DateFrom:    timePtrToPgtypeTimestamptz(dateFrom),
		DateTo:      timePtrToPgtypeTimestamptz(dateTo),
		Search:      search,
		QueryLimit:  int32(limit),
		QueryOffset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list policies filtered: %w", err)
	}
	policies := make([]*entity.Policy, len(dbPolicies))
	for i, p := range dbPolicies {
		policies[i] = sqlcPolicyToDomain(p)
	}
	return policies, nil
}

func (r *policyRepository) CountFiltered(ctx context.Context, dateFrom, dateTo *time.Time, search string) (int64, error) {
	count, err := r.store.CountPoliciesFiltered(ctx, db.CountPoliciesFilteredParams{
		DateFrom: timePtrToPgtypeTimestamptz(dateFrom),
		DateTo:   timePtrToPgtypeTimestamptz(dateTo),
		Search:   search,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to count policies filtered: %w", err)
	}
	return count, nil
}
