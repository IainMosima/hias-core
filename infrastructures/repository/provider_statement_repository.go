package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/billing/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/billing/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type providerStatementRepository struct {
	store db.Store
}

func NewProviderStatementRepository(store db.Store) domainRepo.ProviderStatementRepository {
	return &providerStatementRepository{store: store}
}

func (r *providerStatementRepository) Create(ctx context.Context, stmt *entity.ProviderStatement) (*entity.ProviderStatement, error) {
	dbStmt, err := r.store.CreateProviderStatement(ctx, db.CreateProviderStatementParams{
		ProviderID:      stmt.ProviderID,
		StatementNumber: stmt.StatementNumber,
		PeriodStart:     timeToPgtypeDate(stmt.PeriodStart),
		PeriodEnd:       timeToPgtypeDate(stmt.PeriodEnd),
		TotalClaimed:    stmt.TotalClaimed,
		Status:          stmt.Status,
		FileName:        stringToPgtypeText(stmt.FileName),
		S3Key:           stringToPgtypeText(stmt.S3Key),
		CreatedBy:       stmt.CreatedBy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create provider statement: %w", err)
	}
	return sqlcProviderStatementToDomain(dbStmt), nil
}

func (r *providerStatementRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ProviderStatement, error) {
	dbStmt, err := r.store.GetProviderStatementByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider statement by ID: %w", err)
	}
	return sqlcProviderStatementToDomain(dbStmt), nil
}

func (r *providerStatementRepository) ListByProvider(ctx context.Context, providerID uuid.UUID, limit, offset int) ([]*entity.ProviderStatement, error) {
	dbStmts, err := r.store.ListProviderStatementsByProvider(ctx, db.ListProviderStatementsByProviderParams{
		ProviderID: providerID,
		Limit:      int32(limit),
		Offset:     int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list provider statements by provider: %w", err)
	}
	stmts := make([]*entity.ProviderStatement, len(dbStmts))
	for i, s := range dbStmts {
		stmts[i] = sqlcProviderStatementToDomain(s)
	}
	return stmts, nil
}

func (r *providerStatementRepository) Reconcile(ctx context.Context, id uuid.UUID, totalMatched, totalDiscrepancy int64, matchedCount, unmatchedCount int, reconciledBy uuid.UUID) (*entity.ProviderStatement, error) {
	dbStmt, err := r.store.ReconcileProviderStatement(ctx, db.ReconcileProviderStatementParams{
		ID:               id,
		TotalMatched:     totalMatched,
		TotalDiscrepancy: totalDiscrepancy,
		MatchedCount:     int32(matchedCount),
		UnmatchedCount:   int32(unmatchedCount),
		ReconciledBy:     uuidToPgtype(reconciledBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to reconcile provider statement: %w", err)
	}
	return sqlcProviderStatementToDomain(dbStmt), nil
}

func (r *providerStatementRepository) CreateLineItem(ctx context.Context, item *entity.StatementLineItem) (*entity.StatementLineItem, error) {
	dbItem, err := r.store.CreateStatementLineItem(ctx, db.CreateStatementLineItemParams{
		StatementID:   item.StatementID,
		ClaimNumber:   stringToPgtypeText(item.ClaimNumber),
		ServiceDate:   timePtrToPgtypeDate(item.ServiceDate),
		MemberName:    stringToPgtypeText(item.MemberName),
		ProcedureCode: stringToPgtypeText(item.ProcedureCode),
		ClaimedAmount: item.ClaimedAmount,
		MatchStatus:   item.MatchStatus,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create statement line item: %w", err)
	}
	return sqlcStatementLineItemToDomain(dbItem), nil
}

func (r *providerStatementRepository) ListLineItems(ctx context.Context, statementID uuid.UUID) ([]*entity.StatementLineItem, error) {
	dbItems, err := r.store.ListStatementLineItemsByStatement(ctx, statementID)
	if err != nil {
		return nil, fmt.Errorf("failed to list statement line items: %w", err)
	}
	items := make([]*entity.StatementLineItem, len(dbItems))
	for i, item := range dbItems {
		items[i] = sqlcStatementLineItemToDomain(item)
	}
	return items, nil
}

func (r *providerStatementRepository) MatchLineItem(ctx context.Context, id, matchedClaimID uuid.UUID, matchStatus string, discrepancy int64, notes string) (*entity.StatementLineItem, error) {
	dbItem, err := r.store.MatchStatementLineItem(ctx, db.MatchStatementLineItemParams{
		ID:                id,
		MatchedClaimID:    uuidToPgtype(matchedClaimID),
		MatchStatus:       matchStatus,
		DiscrepancyAmount: discrepancy,
		Notes:             stringToPgtypeText(notes),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to match statement line item: %w", err)
	}
	return sqlcStatementLineItemToDomain(dbItem), nil
}

func sqlcProviderStatementToDomain(s db.ProviderStatement) *entity.ProviderStatement {
	return &entity.ProviderStatement{
		ID:               s.ID,
		ProviderID:       s.ProviderID,
		StatementNumber:  s.StatementNumber,
		PeriodStart:      pgtypeDateToTime(s.PeriodStart),
		PeriodEnd:        pgtypeDateToTime(s.PeriodEnd),
		TotalClaimed:     s.TotalClaimed,
		TotalMatched:     s.TotalMatched,
		TotalDiscrepancy: s.TotalDiscrepancy,
		MatchedCount:     int(s.MatchedCount),
		UnmatchedCount:   int(s.UnmatchedCount),
		Status:           s.Status,
		FileName:         s.FileName.String,
		S3Key:            s.S3Key.String,
		ReconciledBy:     pgtypeToUUID(s.ReconciledBy),
		ReconciledAt:     pgtypeTimestamptzToTimePtr(s.ReconciledAt),
		CreatedBy:        s.CreatedBy,
		CreatedAt:        s.CreatedAt,
		UpdatedAt:        s.UpdatedAt,
	}
}

func sqlcStatementLineItemToDomain(i db.StatementLineItem) *entity.StatementLineItem {
	return &entity.StatementLineItem{
		ID:                i.ID,
		StatementID:       i.StatementID,
		ClaimNumber:       i.ClaimNumber.String,
		ServiceDate:       pgtypeDateToTimePtr(i.ServiceDate),
		MemberName:        i.MemberName.String,
		ProcedureCode:     i.ProcedureCode.String,
		ClaimedAmount:     i.ClaimedAmount,
		MatchedClaimID:    pgtypeToUUID(i.MatchedClaimID),
		MatchStatus:       i.MatchStatus,
		DiscrepancyAmount: i.DiscrepancyAmount,
		Notes:             i.Notes.String,
		CreatedAt:         i.CreatedAt,
	}
}
