package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/sales/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/sales/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type approvalLimitRepository struct {
	store db.Store
}

func NewApprovalLimitRepository(store db.Store) domainRepo.ApprovalLimitRepository {
	return &approvalLimitRepository{store: store}
}

func (r *approvalLimitRepository) GetByRole(ctx context.Context, roleName string) (*entity.ApprovalLimit, error) {
	dbLimit, err := r.store.GetApprovalLimitByRole(ctx, roleName)
	if err != nil {
		return nil, fmt.Errorf("failed to get approval limit by role: %w", err)
	}
	return sqlcApprovalLimitToDomain(dbLimit), nil
}

func (r *approvalLimitRepository) List(ctx context.Context) ([]*entity.ApprovalLimit, error) {
	dbLimits, err := r.store.ListApprovalLimits(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list approval limits: %w", err)
	}
	limits := make([]*entity.ApprovalLimit, len(dbLimits))
	for i, l := range dbLimits {
		limits[i] = sqlcApprovalLimitToDomain(l)
	}
	return limits, nil
}

func (r *approvalLimitRepository) Create(ctx context.Context, limit *entity.ApprovalLimit) (*entity.ApprovalLimit, error) {
	dbLimit, err := r.store.CreateApprovalLimit(ctx, db.CreateApprovalLimitParams{
		RoleName:              limit.RoleName,
		MaxDiscountPercentage: limit.MaxDiscountPercentage,
		MaxDiscountAmount:     limit.MaxDiscountAmount,
		MaxLoadingPercentage:  limit.MaxLoadingPercentage,
		MaxLoadingAmount:      limit.MaxLoadingAmount,
		EscalationRole:        stringToPgtypeText(limit.EscalationRole),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create approval limit: %w", err)
	}
	return sqlcApprovalLimitToDomain(dbLimit), nil
}

func (r *approvalLimitRepository) Update(ctx context.Context, limit *entity.ApprovalLimit) (*entity.ApprovalLimit, error) {
	dbLimit, err := r.store.UpdateApprovalLimit(ctx, db.UpdateApprovalLimitParams{
		ID:                    limit.ID,
		MaxDiscountPercentage: limit.MaxDiscountPercentage,
		MaxDiscountAmount:     limit.MaxDiscountAmount,
		MaxLoadingPercentage:  limit.MaxLoadingPercentage,
		MaxLoadingAmount:      limit.MaxLoadingAmount,
		Column6:               limit.EscalationRole,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update approval limit: %w", err)
	}
	return sqlcApprovalLimitToDomain(dbLimit), nil
}

func (r *approvalLimitRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.store.DeleteApprovalLimit(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete approval limit: %w", err)
	}
	return nil
}

func sqlcApprovalLimitToDomain(l db.ApprovalLimit) *entity.ApprovalLimit {
	return &entity.ApprovalLimit{
		ID:                    l.ID,
		RoleName:              l.RoleName,
		MaxDiscountPercentage: l.MaxDiscountPercentage,
		MaxDiscountAmount:     l.MaxDiscountAmount,
		MaxLoadingPercentage:  l.MaxLoadingPercentage,
		MaxLoadingAmount:      l.MaxLoadingAmount,
		EscalationRole:        l.EscalationRole.String,
		IsActive:              l.IsActive,
		CreatedAt:             l.CreatedAt,
		UpdatedAt:             l.UpdatedAt,
	}
}
