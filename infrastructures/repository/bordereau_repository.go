package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/reinsurance/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/reinsurance/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type bordereauRepository struct {
	store db.Store
}

func NewBordereauRepository(store db.Store) domainRepo.BordereauRepository {
	return &bordereauRepository{store: store}
}

func (r *bordereauRepository) Create(ctx context.Context, bordereau *entity.Bordereau) (*entity.Bordereau, error) {
	dbBordereau, err := r.store.CreateBordereau(ctx, db.CreateBordereauParams{
		BordereauNumber: bordereau.BordereauNumber,
		TreatyID:        bordereau.TreatyID,
		BordereauType:   bordereau.BordereauType,
		PeriodStart:     timeToPgtypeDate(bordereau.PeriodStart),
		PeriodEnd:       timeToPgtypeDate(bordereau.PeriodEnd),
		TotalGross:      bordereau.TotalGross,
		TotalCeded:      bordereau.TotalCeded,
		TotalCommission: bordereau.TotalCommission,
		ItemCount:       int32(bordereau.ItemCount),
		Status:          bordereau.Status,
		CreatedBy:       bordereau.CreatedBy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create bordereau: %w", err)
	}
	return sqlcBordereauToDomain(dbBordereau), nil
}

func (r *bordereauRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Bordereau, error) {
	dbBordereau, err := r.store.GetBordereauByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get bordereau by ID: %w", err)
	}
	return sqlcBordereauToDomain(dbBordereau), nil
}

func (r *bordereauRepository) GetByNumber(ctx context.Context, number string) (*entity.Bordereau, error) {
	dbBordereau, err := r.store.GetBordereauByNumber(ctx, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get bordereau by number: %w", err)
	}
	return sqlcBordereauToDomain(dbBordereau), nil
}

func (r *bordereauRepository) ListByTreaty(ctx context.Context, treatyID uuid.UUID, limit, offset int) ([]*entity.Bordereau, error) {
	dbBordereauxList, err := r.store.ListBordereauxByTreaty(ctx, db.ListBordereauxByTreatyParams{
		TreatyID: treatyID,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list bordereauxby treaty: %w", err)
	}
	bordereauxList := make([]*entity.Bordereau, len(dbBordereauxList))
	for i, b := range dbBordereauxList {
		bordereauxList[i] = sqlcBordereauToDomain(b)
	}
	return bordereauxList, nil
}

func (r *bordereauRepository) Update(ctx context.Context, bordereau *entity.Bordereau) (*entity.Bordereau, error) {
	dbBordereau, err := r.store.UpdateBordereau(ctx, db.UpdateBordereauParams{
		ID:              bordereau.ID,
		TotalGross:      bordereau.TotalGross,
		TotalCeded:      bordereau.TotalCeded,
		TotalCommission: bordereau.TotalCommission,
		ItemCount:       int32(bordereau.ItemCount),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update bordereau: %w", err)
	}
	return sqlcBordereauToDomain(dbBordereau), nil
}

func (r *bordereauRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Bordereau, error) {
	dbBordereau, err := r.store.UpdateBordereauStatus(ctx, db.UpdateBordereauStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update bordereau status: %w", err)
	}
	return sqlcBordereauToDomain(dbBordereau), nil
}

func (r *bordereauRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.store.CountBordereaux(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count bordereaux: %w", err)
	}
	return count, nil
}

func sqlcBordereauToDomain(b db.Bordereaux) *entity.Bordereau {
	return &entity.Bordereau{
		ID:              b.ID,
		BordereauNumber: b.BordereauNumber,
		TreatyID:        b.TreatyID,
		BordereauType:   b.BordereauType,
		PeriodStart:     pgtypeDateToTime(b.PeriodStart),
		PeriodEnd:       pgtypeDateToTime(b.PeriodEnd),
		TotalGross:      b.TotalGross,
		TotalCeded:      b.TotalCeded,
		TotalCommission: b.TotalCommission,
		ItemCount:       int(b.ItemCount),
		Status:          b.Status,
		CreatedBy:       b.CreatedBy,
		CreatedAt:       b.CreatedAt,
		UpdatedAt:       b.UpdatedAt,
	}
}
