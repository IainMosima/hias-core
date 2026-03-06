package repository

import (
	"context"
	"github.com/bitbiz/hias-core/domains/billing/entity"
	"github.com/google/uuid"
)

type PremiumLedgerRepository interface {
	Create(ctx context.Context, entry *entity.PremiumLedgerEntry) (*entity.PremiumLedgerEntry, error)
	ListByPolicy(ctx context.Context, policyID uuid.UUID, limit, offset int) ([]*entity.PremiumLedgerEntry, error)
	GetBalanceByPolicy(ctx context.Context, policyID uuid.UUID) (int64, error)
}
