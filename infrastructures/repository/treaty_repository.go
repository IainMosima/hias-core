package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/reinsurance/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/reinsurance/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type treatyRepository struct {
	store db.Store
}

func NewTreatyRepository(store db.Store) domainRepo.TreatyRepository {
	return &treatyRepository{store: store}
}

func (r *treatyRepository) Create(ctx context.Context, treaty *entity.Treaty) (*entity.Treaty, error) {
	dbTreaty, err := r.store.CreateTreaty(ctx, db.CreateTreatyParams{
		TreatyNumber:   treaty.TreatyNumber,
		Name:           treaty.Name,
		TreatyType:     treaty.TreatyType,
		Status:         treaty.Status,
		EffectiveDate:  timeToPgtypeDate(treaty.EffectiveDate),
		ExpiryDate:     timeToPgtypeDate(treaty.ExpiryDate),
		RetentionLimit: treaty.RetentionLimit,
		Currency:       treaty.Currency,
		Notes:          stringToPgtypeText(treaty.Notes),
		CreatedBy:      treaty.CreatedBy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create treaty: %w", err)
	}
	return sqlcTreatyToDomain(dbTreaty), nil
}

func (r *treatyRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Treaty, error) {
	dbTreaty, err := r.store.GetTreatyByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get treaty by ID: %w", err)
	}
	return sqlcTreatyToDomain(dbTreaty), nil
}

func (r *treatyRepository) GetByNumber(ctx context.Context, number string) (*entity.Treaty, error) {
	dbTreaty, err := r.store.GetTreatyByNumber(ctx, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get treaty by number: %w", err)
	}
	return sqlcTreatyToDomain(dbTreaty), nil
}

func (r *treatyRepository) List(ctx context.Context, limit, offset int) ([]*entity.Treaty, error) {
	dbTreaties, err := r.store.ListTreaties(ctx, db.ListTreatiesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list treaties: %w", err)
	}
	treaties := make([]*entity.Treaty, len(dbTreaties))
	for i, t := range dbTreaties {
		treaties[i] = sqlcTreatyToDomain(t)
	}
	return treaties, nil
}

func (r *treatyRepository) ListActive(ctx context.Context, limit, offset int) ([]*entity.Treaty, error) {
	dbTreaties, err := r.store.ListActiveTreaties(ctx, db.ListActiveTreatiesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list active treaties: %w", err)
	}
	treaties := make([]*entity.Treaty, len(dbTreaties))
	for i, t := range dbTreaties {
		treaties[i] = sqlcTreatyToDomain(t)
	}
	return treaties, nil
}

func (r *treatyRepository) ListByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.Treaty, error) {
	dbTreaties, err := r.store.ListTreatiesByStatus(ctx, db.ListTreatiesByStatusParams{
		Status: status,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list treaties by status: %w", err)
	}
	treaties := make([]*entity.Treaty, len(dbTreaties))
	for i, t := range dbTreaties {
		treaties[i] = sqlcTreatyToDomain(t)
	}
	return treaties, nil
}

func (r *treatyRepository) ListByType(ctx context.Context, treatyType string, limit, offset int) ([]*entity.Treaty, error) {
	dbTreaties, err := r.store.ListTreatiesByType(ctx, db.ListTreatiesByTypeParams{
		TreatyType: treatyType,
		Limit:      int32(limit),
		Offset:     int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list treaties by type: %w", err)
	}
	treaties := make([]*entity.Treaty, len(dbTreaties))
	for i, t := range dbTreaties {
		treaties[i] = sqlcTreatyToDomain(t)
	}
	return treaties, nil
}

func (r *treatyRepository) Update(ctx context.Context, treaty *entity.Treaty) (*entity.Treaty, error) {
	dbTreaty, err := r.store.UpdateTreaty(ctx, db.UpdateTreatyParams{
		ID:             treaty.ID,
		Name:           treaty.Name,
		EffectiveDate:  timeToPgtypeDate(treaty.EffectiveDate),
		ExpiryDate:     timeToPgtypeDate(treaty.ExpiryDate),
		RetentionLimit: treaty.RetentionLimit,
		Currency:       treaty.Currency,
		Notes:          stringToPgtypeText(treaty.Notes),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update treaty: %w", err)
	}
	return sqlcTreatyToDomain(dbTreaty), nil
}

func (r *treatyRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Treaty, error) {
	dbTreaty, err := r.store.UpdateTreatyStatus(ctx, db.UpdateTreatyStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update treaty status: %w", err)
	}
	return sqlcTreatyToDomain(dbTreaty), nil
}

func (r *treatyRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.store.CountTreaties(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count treaties: %w", err)
	}
	return count, nil
}

func (r *treatyRepository) ListExpiring(ctx context.Context, withinDays int, limit, offset int) ([]*entity.Treaty, error) {
	dbTreaties, err := r.store.ListExpiringTreaties(ctx, db.ListExpiringTreatiesParams{
		Column1: int32(withinDays),
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list expiring treaties: %w", err)
	}
	treaties := make([]*entity.Treaty, len(dbTreaties))
	for i, t := range dbTreaties {
		treaties[i] = sqlcTreatyToDomain(t)
	}
	return treaties, nil
}

func sqlcTreatyToDomain(t db.Treaty) *entity.Treaty {
	return &entity.Treaty{
		ID:             t.ID,
		TreatyNumber:   t.TreatyNumber,
		Name:           t.Name,
		TreatyType:     t.TreatyType,
		Status:         t.Status,
		EffectiveDate:  pgtypeDateToTime(t.EffectiveDate),
		ExpiryDate:     pgtypeDateToTime(t.ExpiryDate),
		RetentionLimit: t.RetentionLimit,
		Currency:       t.Currency,
		Notes:          t.Notes.String,
		CreatedBy:      t.CreatedBy,
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
	}
}
