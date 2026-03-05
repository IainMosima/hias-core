package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/policy/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/policy/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type policyRenewalRepository struct {
	store db.Store
}

func NewPolicyRenewalRepository(store db.Store) domainRepo.PolicyRenewalRepository {
	return &policyRenewalRepository{store: store}
}

func (r *policyRenewalRepository) Create(ctx context.Context, renewal *entity.PolicyRenewal) (*entity.PolicyRenewal, error) {
	dbR, err := r.store.CreatePolicyRenewal(ctx, db.CreatePolicyRenewalParams{
		PolicyID:            renewal.PolicyID,
		Status:              renewal.Status,
		RenewalDate:         renewal.RenewalDate,
		NewPremium:          int64ToPgtypeInt8(renewal.NewPremium),
		PremiumChangeReason: stringToPgtypeText(renewal.PremiumChangeReason),
		NewPlanID:           uuidToPgtype(renewal.NewPlanID),
		ExpiresAt:           timePtrToPgtypeTimestamptz(renewal.ExpiresAt),
		CreatedBy:           renewal.CreatedBy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create policy renewal: %w", err)
	}
	return sqlcPolicyRenewalToDomain(dbR), nil
}

func (r *policyRenewalRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.PolicyRenewal, error) {
	dbR, err := r.store.GetPolicyRenewalByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get policy renewal: %w", err)
	}
	return sqlcPolicyRenewalToDomain(dbR), nil
}

func (r *policyRenewalRepository) GetByPolicyID(ctx context.Context, policyID uuid.UUID) (*entity.PolicyRenewal, error) {
	dbR, err := r.store.GetPolicyRenewalByPolicyID(ctx, policyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get policy renewal by policy: %w", err)
	}
	return sqlcPolicyRenewalToDomain(dbR), nil
}

func (r *policyRenewalRepository) ListPending(ctx context.Context) ([]*entity.PolicyRenewal, error) {
	dbRs, err := r.store.ListPendingRenewals(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list pending renewals: %w", err)
	}
	result := make([]*entity.PolicyRenewal, len(dbRs))
	for i, rn := range dbRs {
		result[i] = sqlcPolicyRenewalToDomain(rn)
	}
	return result, nil
}

func (r *policyRenewalRepository) ListExpired(ctx context.Context) ([]*entity.PolicyRenewal, error) {
	dbRs, err := r.store.ListExpiredRenewals(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list expired renewals: %w", err)
	}
	result := make([]*entity.PolicyRenewal, len(dbRs))
	for i, rn := range dbRs {
		result[i] = sqlcPolicyRenewalToDomain(rn)
	}
	return result, nil
}

func (r *policyRenewalRepository) Update(ctx context.Context, renewal *entity.PolicyRenewal) (*entity.PolicyRenewal, error) {
	dbR, err := r.store.UpdatePolicyRenewal(ctx, db.UpdatePolicyRenewalParams{
		ID:                  renewal.ID,
		Status:              stringToPgtypeText(renewal.Status),
		RenewedPolicyID:     uuidToPgtype(renewal.RenewedPolicyID),
		NewPremium:          int64ToPgtypeInt8(renewal.NewPremium),
		PremiumChangeReason: stringToPgtypeText(renewal.PremiumChangeReason),
		ApprovedBy:          uuidToPgtype(renewal.ApprovedBy),
		ApprovedAt:          timePtrToPgtypeTimestamptz(renewal.ApprovedAt),
		CompletedAt:         timePtrToPgtypeTimestamptz(renewal.CompletedAt),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update policy renewal: %w", err)
	}
	return sqlcPolicyRenewalToDomain(dbR), nil
}

func (r *policyRenewalRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.PolicyRenewal, error) {
	dbR, err := r.store.UpdatePolicyRenewalStatus(ctx, db.UpdatePolicyRenewalStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update renewal status: %w", err)
	}
	return sqlcPolicyRenewalToDomain(dbR), nil
}

func sqlcPolicyRenewalToDomain(r db.PolicyRenewal) *entity.PolicyRenewal {
	var newPremium int64
	if r.NewPremium.Valid {
		newPremium = r.NewPremium.Int64
	}
	return &entity.PolicyRenewal{
		ID:                  r.ID,
		PolicyID:            r.PolicyID,
		RenewedPolicyID:     pgtypeToUUID(r.RenewedPolicyID),
		Status:              r.Status,
		RenewalDate:         r.RenewalDate,
		NewPremium:          newPremium,
		PremiumChangeReason: r.PremiumChangeReason.String,
		NewPlanID:           pgtypeToUUID(r.NewPlanID),
		ApprovedBy:          pgtypeToUUID(r.ApprovedBy),
		ApprovedAt:          pgtypeTimestamptzToTimePtr(r.ApprovedAt),
		CompletedAt:         pgtypeTimestamptzToTimePtr(r.CompletedAt),
		ExpiresAt:           pgtypeTimestamptzToTimePtr(r.ExpiresAt),
		CreatedBy:           r.CreatedBy,
		CreatedAt:           r.CreatedAt,
		UpdatedAt:           r.UpdatedAt,
	}
}
