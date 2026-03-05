package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/bitbiz/hias-core/domains/reinsurance/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/reinsurance/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type cessionRepository struct {
	store db.Store
}

func NewCessionRepository(store db.Store) domainRepo.CessionRepository {
	return &cessionRepository{store: store}
}

func (r *cessionRepository) Create(ctx context.Context, cession *entity.Cession) (*entity.Cession, error) {
	dbCession, err := r.store.CreateCession(ctx, db.CreateCessionParams{
		CessionNumber:    cession.CessionNumber,
		TreatyID:         cession.TreatyID,
		PolicyID:         cession.PolicyID,
		TreatyLayerID:    uuidToPgtype(cession.TreatyLayerID),
		CessionType:      cession.CessionType,
		GrossAmount:      cession.GrossAmount,
		CededAmount:      cession.CededAmount,
		RetainedAmount:   cession.RetainedAmount,
		CommissionAmount: cession.CommissionAmount,
		SharePercentage:  float64ToPgNumeric(cession.SharePercentage),
		Status:           cession.Status,
		CreatedBy:        cession.CreatedBy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create cession: %w", err)
	}
	return sqlcCessionToDomain(dbCession), nil
}

func (r *cessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Cession, error) {
	dbCession, err := r.store.GetCessionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get cession by ID: %w", err)
	}
	return sqlcCessionToDomain(dbCession), nil
}

func (r *cessionRepository) GetByNumber(ctx context.Context, number string) (*entity.Cession, error) {
	dbCession, err := r.store.GetCessionByNumber(ctx, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get cession by number: %w", err)
	}
	return sqlcCessionToDomain(dbCession), nil
}

func (r *cessionRepository) ListByTreaty(ctx context.Context, treatyID uuid.UUID, limit, offset int) ([]*entity.Cession, error) {
	dbCessions, err := r.store.ListCessionsByTreaty(ctx, db.ListCessionsByTreatyParams{
		TreatyID: treatyID,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list cessions by treaty: %w", err)
	}
	cessions := make([]*entity.Cession, len(dbCessions))
	for i, c := range dbCessions {
		cessions[i] = sqlcCessionToDomain(c)
	}
	return cessions, nil
}

func (r *cessionRepository) ListByPolicy(ctx context.Context, policyID uuid.UUID, limit, offset int) ([]*entity.Cession, error) {
	dbCessions, err := r.store.ListCessionsByPolicy(ctx, db.ListCessionsByPolicyParams{
		PolicyID: policyID,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list cessions by policy: %w", err)
	}
	cessions := make([]*entity.Cession, len(dbCessions))
	for i, c := range dbCessions {
		cessions[i] = sqlcCessionToDomain(c)
	}
	return cessions, nil
}

func (r *cessionRepository) ListByTreatyAndPeriod(ctx context.Context, treatyID uuid.UUID, start, end time.Time, limit, offset int) ([]*entity.Cession, error) {
	dbCessions, err := r.store.ListCessionsByTreatyAndPeriod(ctx, db.ListCessionsByTreatyAndPeriodParams{
		TreatyID:    treatyID,
		CreatedAt:   start,
		CreatedAt_2: end,
		Limit:       int32(limit),
		Offset:      int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list cessions by treaty and period: %w", err)
	}
	cessions := make([]*entity.Cession, len(dbCessions))
	for i, c := range dbCessions {
		cessions[i] = sqlcCessionToDomain(c)
	}
	return cessions, nil
}

func (r *cessionRepository) ListBookedByTreatyAndPeriod(ctx context.Context, treatyID uuid.UUID, start, end time.Time) ([]*entity.Cession, error) {
	dbCessions, err := r.store.ListBookedCessionsByTreatyAndPeriod(ctx, db.ListBookedCessionsByTreatyAndPeriodParams{
		TreatyID:    treatyID,
		CreatedAt:   start,
		CreatedAt_2: end,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list booked cessions by treaty and period: %w", err)
	}
	cessions := make([]*entity.Cession, len(dbCessions))
	for i, c := range dbCessions {
		cessions[i] = sqlcCessionToDomain(c)
	}
	return cessions, nil
}

func (r *cessionRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Cession, error) {
	dbCession, err := r.store.UpdateCessionStatus(ctx, db.UpdateCessionStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update cession status: %w", err)
	}
	return sqlcCessionToDomain(dbCession), nil
}

func (r *cessionRepository) GetTotalCededByTreaty(ctx context.Context, treatyID uuid.UUID) (int64, error) {
	total, err := r.store.GetTotalCededByTreaty(ctx, treatyID)
	if err != nil {
		return 0, fmt.Errorf("failed to get total ceded by treaty: %w", err)
	}
	return total, nil
}

func (r *cessionRepository) GetTotalCededByTreatyAndPeriod(ctx context.Context, treatyID uuid.UUID, start, end time.Time) (int64, error) {
	total, err := r.store.GetTotalCededByTreatyAndPeriod(ctx, db.GetTotalCededByTreatyAndPeriodParams{
		TreatyID:    treatyID,
		CreatedAt:   start,
		CreatedAt_2: end,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get total ceded by treaty and period: %w", err)
	}
	return total, nil
}

func (r *cessionRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.store.CountCessions(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count cessions: %w", err)
	}
	return count, nil
}

func (r *cessionRepository) GetTotalCededAmountAll(ctx context.Context) (int64, error) {
	total, err := r.store.GetTotalCededAmountAll(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get total ceded amount: %w", err)
	}
	return total, nil
}

func (r *cessionRepository) GetTotalGrossAmountAll(ctx context.Context) (int64, error) {
	total, err := r.store.GetTotalGrossAmountAll(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get total gross amount: %w", err)
	}
	return total, nil
}

func sqlcCessionToDomain(c db.Cession) *entity.Cession {
	return &entity.Cession{
		ID:               c.ID,
		CessionNumber:    c.CessionNumber,
		TreatyID:         c.TreatyID,
		PolicyID:         c.PolicyID,
		TreatyLayerID:    pgtypeToUUID(c.TreatyLayerID),
		CessionType:      c.CessionType,
		GrossAmount:      c.GrossAmount,
		CededAmount:      c.CededAmount,
		RetainedAmount:   c.RetainedAmount,
		CommissionAmount: c.CommissionAmount,
		SharePercentage:  pgNumericToFloat64(c.SharePercentage),
		Status:           c.Status,
		CreatedBy:        c.CreatedBy,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
	}
}
