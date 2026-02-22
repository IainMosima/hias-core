package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/claims/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type claimLineItemRepository struct {
	store db.Store
}

func NewClaimLineItemRepository(store db.Store) domainRepo.ClaimLineItemRepository {
	return &claimLineItemRepository{store: store}
}

func (r *claimLineItemRepository) Create(ctx context.Context, item *entity.ClaimLineItem) (*entity.ClaimLineItem, error) {
	dbItem, err := r.store.CreateClaimLineItem(ctx, db.CreateClaimLineItemParams{
		ClaimID:       item.ClaimID,
		ProcedureCode: item.ProcedureCode,
		ProcedureName: item.ProcedureName,
		DiagnosisCode: stringToPgtypeText(item.DiagnosisCode),
		Quantity:      int32(item.Quantity),
		UnitPrice:     item.UnitPrice,
		TotalPrice:    item.TotalPrice,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create claim line item: %w", err)
	}
	return sqlcClaimLineItemToDomain(dbItem), nil
}

func (r *claimLineItemRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ClaimLineItem, error) {
	dbItem, err := r.store.GetClaimLineItemByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get claim line item by ID: %w", err)
	}
	return sqlcClaimLineItemToDomain(dbItem), nil
}

func (r *claimLineItemRepository) ListByClaim(ctx context.Context, claimID uuid.UUID) ([]*entity.ClaimLineItem, error) {
	dbItems, err := r.store.ListClaimLineItems(ctx, claimID)
	if err != nil {
		return nil, fmt.Errorf("failed to list claim line items by claim: %w", err)
	}
	items := make([]*entity.ClaimLineItem, len(dbItems))
	for i, item := range dbItems {
		items[i] = sqlcClaimLineItemToDomain(item)
	}
	return items, nil
}

func (r *claimLineItemRepository) UpdateApprovedAmount(ctx context.Context, id uuid.UUID, amount int64) (*entity.ClaimLineItem, error) {
	dbItem, err := r.store.UpdateClaimLineItemApprovedAmount(ctx, db.UpdateClaimLineItemApprovedAmountParams{
		ID:             id,
		ApprovedAmount: amount,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update claim line item approved amount: %w", err)
	}
	return sqlcClaimLineItemToDomain(dbItem), nil
}

func (r *claimLineItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.store.DeleteClaimLineItem(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete claim line item: %w", err)
	}
	return nil
}

func sqlcClaimLineItemToDomain(i db.ClaimLineItem) *entity.ClaimLineItem {
	return &entity.ClaimLineItem{
		ID:             i.ID,
		ClaimID:        i.ClaimID,
		ProcedureCode:  i.ProcedureCode,
		ProcedureName:  i.ProcedureName,
		DiagnosisCode:  i.DiagnosisCode.String,
		Quantity:        int(i.Quantity),
		UnitPrice:      i.UnitPrice,
		TotalPrice:     i.TotalPrice,
		ApprovedAmount: i.ApprovedAmount,
		CreatedAt:      i.CreatedAt,
		UpdatedAt:      i.UpdatedAt,
	}
}
