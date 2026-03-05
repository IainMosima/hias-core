package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/reinsurance/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/reinsurance/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type bordereauItemRepository struct {
	store db.Store
}

func NewBordereauItemRepository(store db.Store) domainRepo.BordereauItemRepository {
	return &bordereauItemRepository{store: store}
}

func (r *bordereauItemRepository) Create(ctx context.Context, item *entity.BordereauItem) (*entity.BordereauItem, error) {
	dbItem, err := r.store.CreateBordereauItem(ctx, db.CreateBordereauItemParams{
		BordereauID:      item.BordereauID,
		CessionID:        uuidToPgtype(item.CessionID),
		RecoveryID:       uuidToPgtype(item.RecoveryID),
		PolicyNumber:     stringToPgtypeText(item.PolicyNumber),
		ClaimNumber:      stringToPgtypeText(item.ClaimNumber),
		GrossAmount:      item.GrossAmount,
		CededAmount:      item.CededAmount,
		CommissionAmount: item.CommissionAmount,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create bordereau item: %w", err)
	}
	return sqlcBordereauItemToDomain(dbItem), nil
}

func (r *bordereauItemRepository) ListByBordereau(ctx context.Context, bordereauID uuid.UUID) ([]*entity.BordereauItem, error) {
	dbItems, err := r.store.ListBordereauItemsByBordereau(ctx, bordereauID)
	if err != nil {
		return nil, fmt.Errorf("failed to list bordereau items by bordereau: %w", err)
	}
	items := make([]*entity.BordereauItem, len(dbItems))
	for i, item := range dbItems {
		items[i] = sqlcBordereauItemToDomain(item)
	}
	return items, nil
}

func (r *bordereauItemRepository) DeleteByBordereau(ctx context.Context, bordereauID uuid.UUID) error {
	err := r.store.DeleteBordereauItemsByBordereau(ctx, bordereauID)
	if err != nil {
		return fmt.Errorf("failed to delete bordereau items by bordereau: %w", err)
	}
	return nil
}

func sqlcBordereauItemToDomain(i db.BordereauItem) *entity.BordereauItem {
	return &entity.BordereauItem{
		ID:               i.ID,
		BordereauID:      i.BordereauID,
		CessionID:        pgtypeToUUID(i.CessionID),
		RecoveryID:       pgtypeToUUID(i.RecoveryID),
		PolicyNumber:     i.PolicyNumber.String,
		ClaimNumber:      i.ClaimNumber.String,
		GrossAmount:      i.GrossAmount,
		CededAmount:      i.CededAmount,
		CommissionAmount: i.CommissionAmount,
		CreatedAt:        i.CreatedAt,
	}
}
