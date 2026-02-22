package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/product/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/product/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type exclusionRepository struct {
	store db.Store
}

func NewExclusionRepository(store db.Store) domainRepo.ExclusionRepository {
	return &exclusionRepository{store: store}
}

func (r *exclusionRepository) Create(ctx context.Context, exclusion *entity.Exclusion) (*entity.Exclusion, error) {
	dbExclusion, err := r.store.CreateExclusion(ctx, db.CreateExclusionParams{
		PlanID:      exclusion.PlanID,
		Description: exclusion.Description,
		Type:        exclusion.Type,
		IcdCodes:    exclusion.ICDCodes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create exclusion: %w", err)
	}
	return sqlcExclusionToDomain(dbExclusion), nil
}

func (r *exclusionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Exclusion, error) {
	dbExclusion, err := r.store.GetExclusionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get exclusion by ID: %w", err)
	}
	return sqlcExclusionToDomain(dbExclusion), nil
}

func (r *exclusionRepository) ListByPlan(ctx context.Context, planID uuid.UUID) ([]*entity.Exclusion, error) {
	dbExclusions, err := r.store.ListExclusionsByPlan(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("failed to list exclusions by plan: %w", err)
	}
	exclusions := make([]*entity.Exclusion, len(dbExclusions))
	for i, e := range dbExclusions {
		exclusions[i] = sqlcExclusionToDomain(e)
	}
	return exclusions, nil
}

func (r *exclusionRepository) ListByType(ctx context.Context, planID uuid.UUID, exclusionType string) ([]*entity.Exclusion, error) {
	dbExclusions, err := r.store.ListExclusionsByType(ctx, db.ListExclusionsByTypeParams{
		PlanID: planID,
		Type:   exclusionType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list exclusions by type: %w", err)
	}
	exclusions := make([]*entity.Exclusion, len(dbExclusions))
	for i, e := range dbExclusions {
		exclusions[i] = sqlcExclusionToDomain(e)
	}
	return exclusions, nil
}

func (r *exclusionRepository) Update(ctx context.Context, exclusion *entity.Exclusion) (*entity.Exclusion, error) {
	dbExclusion, err := r.store.UpdateExclusion(ctx, db.UpdateExclusionParams{
		ID:          exclusion.ID,
		Description: stringToPgtypeText(exclusion.Description),
		Type:        stringToPgtypeText(exclusion.Type),
		IcdCodes:    exclusion.ICDCodes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update exclusion: %w", err)
	}
	return sqlcExclusionToDomain(dbExclusion), nil
}

func (r *exclusionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.store.DeleteExclusion(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete exclusion: %w", err)
	}
	return nil
}

func sqlcExclusionToDomain(e db.Exclusion) *entity.Exclusion {
	return &entity.Exclusion{
		ID:          e.ID,
		PlanID:      e.PlanID,
		Description: e.Description,
		Type:        e.Type,
		ICDCodes:    e.IcdCodes,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
