package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/reinsurance/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/reinsurance/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type treatyLayerRepository struct {
	store db.Store
}

func NewTreatyLayerRepository(store db.Store) domainRepo.TreatyLayerRepository {
	return &treatyLayerRepository{store: store}
}

func (r *treatyLayerRepository) Create(ctx context.Context, layer *entity.TreatyLayer) (*entity.TreatyLayer, error) {
	dbLayer, err := r.store.CreateTreatyLayer(ctx, db.CreateTreatyLayerParams{
		TreatyID:         layer.TreatyID,
		LayerNumber:      int32(layer.LayerNumber),
		AttachmentPoint:  layer.AttachmentPoint,
		LayerLimit:       layer.LayerLimit,
		DeductibleAmount: layer.DeductibleAmount,
		PremiumRate:      float64ToPgNumeric(layer.PremiumRate),
		AggregateLimit:   int64PtrToPgtypeInt8(layer.AggregateLimit),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create treaty layer: %w", err)
	}
	return sqlcTreatyLayerToDomain(dbLayer), nil
}

func (r *treatyLayerRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.TreatyLayer, error) {
	dbLayer, err := r.store.GetTreatyLayerByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get treaty layer by ID: %w", err)
	}
	return sqlcTreatyLayerToDomain(dbLayer), nil
}

func (r *treatyLayerRepository) ListByTreaty(ctx context.Context, treatyID uuid.UUID) ([]*entity.TreatyLayer, error) {
	dbLayers, err := r.store.ListTreatyLayersByTreaty(ctx, treatyID)
	if err != nil {
		return nil, fmt.Errorf("failed to list treaty layers by treaty: %w", err)
	}
	layers := make([]*entity.TreatyLayer, len(dbLayers))
	for i, l := range dbLayers {
		layers[i] = sqlcTreatyLayerToDomain(l)
	}
	return layers, nil
}

func (r *treatyLayerRepository) Update(ctx context.Context, layer *entity.TreatyLayer) (*entity.TreatyLayer, error) {
	dbLayer, err := r.store.UpdateTreatyLayer(ctx, db.UpdateTreatyLayerParams{
		ID:               layer.ID,
		AttachmentPoint:  layer.AttachmentPoint,
		LayerLimit:       layer.LayerLimit,
		DeductibleAmount: layer.DeductibleAmount,
		PremiumRate:      float64ToPgNumeric(layer.PremiumRate),
		AggregateLimit:   int64PtrToPgtypeInt8(layer.AggregateLimit),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update treaty layer: %w", err)
	}
	return sqlcTreatyLayerToDomain(dbLayer), nil
}

func (r *treatyLayerRepository) UpdateAggregateUsed(ctx context.Context, id uuid.UUID, aggregateUsed int64) (*entity.TreatyLayer, error) {
	dbLayer, err := r.store.UpdateTreatyLayerAggregateUsed(ctx, db.UpdateTreatyLayerAggregateUsedParams{
		ID:            id,
		AggregateUsed: aggregateUsed,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update treaty layer aggregate used: %w", err)
	}
	return sqlcTreatyLayerToDomain(dbLayer), nil
}

func (r *treatyLayerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.store.DeleteTreatyLayer(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete treaty layer: %w", err)
	}
	return nil
}

func sqlcTreatyLayerToDomain(l db.TreatyLayer) *entity.TreatyLayer {
	return &entity.TreatyLayer{
		ID:               l.ID,
		TreatyID:         l.TreatyID,
		LayerNumber:      int(l.LayerNumber),
		AttachmentPoint:  l.AttachmentPoint,
		LayerLimit:       l.LayerLimit,
		DeductibleAmount: l.DeductibleAmount,
		PremiumRate:      pgNumericToFloat64(l.PremiumRate),
		AggregateLimit:   pgtypeInt8ToInt64Ptr(l.AggregateLimit),
		AggregateUsed:    l.AggregateUsed,
		CreatedAt:        l.CreatedAt,
		UpdatedAt:        l.UpdatedAt,
	}
}
