package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/billing/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/billing/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type premiumLedgerRepository struct {
	store db.Store
}

func NewPremiumLedgerRepository(store db.Store) domainRepo.PremiumLedgerRepository {
	return &premiumLedgerRepository{store: store}
}

func (r *premiumLedgerRepository) Create(ctx context.Context, entry *entity.PremiumLedgerEntry) (*entity.PremiumLedgerEntry, error) {
	var createdBy pgtype.UUID
	if entry.CreatedBy != uuid.Nil {
		createdBy = pgtype.UUID{Bytes: entry.CreatedBy, Valid: true}
	}
	dbEntry, err := r.store.CreatePremiumLedgerEntry(ctx, db.CreatePremiumLedgerEntryParams{
		PolicyID:        entry.PolicyID,
		EntryType:       entry.EntryType,
		Amount:          entry.Amount,
		Description:     entry.Description,
		ReferenceNumber: entry.ReferenceNumber,
		EffectiveDate:   entry.EffectiveDate,
		BalanceAfter:    entry.BalanceAfter,
		CreatedBy:       createdBy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create premium ledger entry: %w", err)
	}
	return sqlcPremiumLedgerToDomain(dbEntry), nil
}

func (r *premiumLedgerRepository) ListByPolicy(ctx context.Context, policyID uuid.UUID, limit, offset int) ([]*entity.PremiumLedgerEntry, error) {
	dbEntries, err := r.store.ListPremiumLedgerByPolicy(ctx, db.ListPremiumLedgerByPolicyParams{
		PolicyID: policyID,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list premium ledger: %w", err)
	}
	entries := make([]*entity.PremiumLedgerEntry, len(dbEntries))
	for i, e := range dbEntries {
		entries[i] = sqlcPremiumLedgerToDomain(e)
	}
	return entries, nil
}

func (r *premiumLedgerRepository) GetBalanceByPolicy(ctx context.Context, policyID uuid.UUID) (int64, error) {
	balance, err := r.store.GetPremiumBalanceByPolicy(ctx, policyID)
	if err != nil {
		return 0, fmt.Errorf("failed to get premium balance: %w", err)
	}
	return balance, nil
}

func sqlcPremiumLedgerToDomain(e db.PremiumLedgerEntry) *entity.PremiumLedgerEntry {
	var createdBy uuid.UUID
	if e.CreatedBy.Valid {
		createdBy = e.CreatedBy.Bytes
	}
	return &entity.PremiumLedgerEntry{
		ID: e.ID, PolicyID: e.PolicyID, EntryType: e.EntryType,
		Amount: e.Amount, Description: e.Description,
		ReferenceNumber: e.ReferenceNumber, EffectiveDate: e.EffectiveDate,
		BalanceAfter: e.BalanceAfter, CreatedBy: createdBy,
		CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt,
	}
}
