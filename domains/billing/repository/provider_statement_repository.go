package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/billing/entity"
	"github.com/google/uuid"
)

type ProviderStatementRepository interface {
	Create(ctx context.Context, stmt *entity.ProviderStatement) (*entity.ProviderStatement, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ProviderStatement, error)
	ListByProvider(ctx context.Context, providerID uuid.UUID, limit, offset int) ([]*entity.ProviderStatement, error)
	Reconcile(ctx context.Context, id uuid.UUID, totalMatched, totalDiscrepancy int64, matchedCount, unmatchedCount int, reconciledBy uuid.UUID) (*entity.ProviderStatement, error)
	CreateLineItem(ctx context.Context, item *entity.StatementLineItem) (*entity.StatementLineItem, error)
	ListLineItems(ctx context.Context, statementID uuid.UUID) ([]*entity.StatementLineItem, error)
	MatchLineItem(ctx context.Context, id, matchedClaimID uuid.UUID, matchStatus string, discrepancy int64, notes string) (*entity.StatementLineItem, error)
}
