package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/billing/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/billing/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type remittanceRepository struct {
	store db.Store
}

func NewRemittanceRepository(store db.Store) domainRepo.RemittanceRepository {
	return &remittanceRepository{store: store}
}

func (r *remittanceRepository) Create(ctx context.Context, remittance *entity.Remittance) (*entity.Remittance, error) {
	dbRemittance, err := r.store.CreateRemittance(ctx, db.CreateRemittanceParams{
		ProviderID:  remittance.ProviderID,
		ClaimIds:    remittance.ClaimIDs,
		TotalAmount: remittance.TotalAmount,
		Currency:    remittance.Currency,
		Status:      remittance.Status,
		PeriodStart: remittance.PeriodStart,
		PeriodEnd:   remittance.PeriodEnd,
		CreatedBy:   uuidToPgtype(remittance.CreatedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create remittance: %w", err)
	}
	return sqlcRemittanceToDomain(dbRemittance), nil
}

func (r *remittanceRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Remittance, error) {
	dbRemittance, err := r.store.GetRemittanceByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get remittance by ID: %w", err)
	}
	return sqlcRemittanceToDomain(dbRemittance), nil
}

func (r *remittanceRepository) List(ctx context.Context, limit, offset int) ([]*entity.Remittance, error) {
	dbRemittances, err := r.store.ListRemittances(ctx, db.ListRemittancesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list remittances: %w", err)
	}
	remittances := make([]*entity.Remittance, len(dbRemittances))
	for i, rem := range dbRemittances {
		remittances[i] = sqlcRemittanceToDomain(rem)
	}
	return remittances, nil
}

func (r *remittanceRepository) ListByProvider(ctx context.Context, providerID uuid.UUID, limit, offset int) ([]*entity.Remittance, error) {
	dbRemittances, err := r.store.ListRemittancesByProvider(ctx, db.ListRemittancesByProviderParams{
		ProviderID: providerID,
		Limit:      int32(limit),
		Offset:     int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list remittances by provider: %w", err)
	}
	remittances := make([]*entity.Remittance, len(dbRemittances))
	for i, rem := range dbRemittances {
		remittances[i] = sqlcRemittanceToDomain(rem)
	}
	return remittances, nil
}

func (r *remittanceRepository) ListByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.Remittance, error) {
	dbRemittances, err := r.store.ListRemittancesByStatus(ctx, db.ListRemittancesByStatusParams{
		Status: status,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list remittances by status: %w", err)
	}
	remittances := make([]*entity.Remittance, len(dbRemittances))
	for i, rem := range dbRemittances {
		remittances[i] = sqlcRemittanceToDomain(rem)
	}
	return remittances, nil
}

func (r *remittanceRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.store.CountRemittances(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count remittances: %w", err)
	}
	return count, nil
}

func (r *remittanceRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Remittance, error) {
	dbRemittance, err := r.store.UpdateRemittanceStatus(ctx, db.UpdateRemittanceStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update remittance status: %w", err)
	}
	return sqlcRemittanceToDomain(dbRemittance), nil
}

func (r *remittanceRepository) MarkAdviceSent(ctx context.Context, id uuid.UUID) (*entity.Remittance, error) {
	dbRemittance, err := r.store.MarkRemittanceAdviceSent(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to mark remittance advice sent: %w", err)
	}
	return sqlcRemittanceToDomain(dbRemittance), nil
}

func (r *remittanceRepository) SetPayment(ctx context.Context, id uuid.UUID, paymentID uuid.UUID) (*entity.Remittance, error) {
	dbRemittance, err := r.store.SetRemittancePayment(ctx, db.SetRemittancePaymentParams{
		ID:        id,
		PaymentID: uuidToPgtype(paymentID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to set remittance payment: %w", err)
	}
	return sqlcRemittanceToDomain(dbRemittance), nil
}

func sqlcRemittanceToDomain(rem db.Remittance) *entity.Remittance {
	return &entity.Remittance{
		ID:                   rem.ID,
		ProviderID:           rem.ProviderID,
		ClaimIDs:             rem.ClaimIds,
		TotalAmount:          rem.TotalAmount,
		Currency:             rem.Currency,
		Status:               rem.Status,
		RemittanceAdviceSent: rem.RemittanceAdviceSent,
		PeriodStart:          rem.PeriodStart,
		PeriodEnd:            rem.PeriodEnd,
		PaymentID:            pgtypeToUUID(rem.PaymentID),
		CreatedBy:            pgtypeToUUID(rem.CreatedBy),
		CreatedAt:            rem.CreatedAt,
		UpdatedAt:            rem.UpdatedAt,
	}
}
