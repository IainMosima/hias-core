package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/reinsurance/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/reinsurance/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type reinsurerStatementRepository struct {
	store db.Store
}

func NewReinsurerStatementRepository(store db.Store) domainRepo.ReinsurerStatementRepository {
	return &reinsurerStatementRepository{store: store}
}

func (r *reinsurerStatementRepository) Create(ctx context.Context, statement *entity.ReinsurerStatement) (*entity.ReinsurerStatement, error) {
	dbStatement, err := r.store.CreateReinsurerStatement(ctx, db.CreateReinsurerStatementParams{
		StatementNumber:  statement.StatementNumber,
		TreatyID:         statement.TreatyID,
		ParticipantID:    statement.ParticipantID,
		PeriodStart:      timeToPgtypeDate(statement.PeriodStart),
		PeriodEnd:        timeToPgtypeDate(statement.PeriodEnd),
		PremiumCeded:     statement.PremiumCeded,
		ClaimsRecovered:  statement.ClaimsRecovered,
		CommissionDue:    statement.CommissionDue,
		ProfitCommission: statement.ProfitCommission,
		NetBalance:       statement.NetBalance,
		Status:           statement.Status,
		CreatedBy:        statement.CreatedBy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create reinsurer statement: %w", err)
	}
	return sqlcReinsurerStatementToDomain(dbStatement), nil
}

func (r *reinsurerStatementRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ReinsurerStatement, error) {
	dbStatement, err := r.store.GetReinsurerStatementByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get reinsurer statement by ID: %w", err)
	}
	return sqlcReinsurerStatementToDomain(dbStatement), nil
}

func (r *reinsurerStatementRepository) GetByNumber(ctx context.Context, number string) (*entity.ReinsurerStatement, error) {
	dbStatement, err := r.store.GetReinsurerStatementByNumber(ctx, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get reinsurer statement by number: %w", err)
	}
	return sqlcReinsurerStatementToDomain(dbStatement), nil
}

func (r *reinsurerStatementRepository) ListByTreaty(ctx context.Context, treatyID uuid.UUID, limit, offset int) ([]*entity.ReinsurerStatement, error) {
	dbStatements, err := r.store.ListReinsurerStatementsByTreaty(ctx, db.ListReinsurerStatementsByTreatyParams{
		TreatyID: treatyID,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list reinsurer statements by treaty: %w", err)
	}
	statements := make([]*entity.ReinsurerStatement, len(dbStatements))
	for i, s := range dbStatements {
		statements[i] = sqlcReinsurerStatementToDomain(s)
	}
	return statements, nil
}

func (r *reinsurerStatementRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.ReinsurerStatement, error) {
	dbStatement, err := r.store.UpdateReinsurerStatementStatus(ctx, db.UpdateReinsurerStatementStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update reinsurer statement status: %w", err)
	}
	return sqlcReinsurerStatementToDomain(dbStatement), nil
}

func (r *reinsurerStatementRepository) Update(ctx context.Context, statement *entity.ReinsurerStatement) (*entity.ReinsurerStatement, error) {
	dbStatement, err := r.store.UpdateReinsurerStatement(ctx, db.UpdateReinsurerStatementParams{
		ID:               statement.ID,
		PremiumCeded:     statement.PremiumCeded,
		ClaimsRecovered:  statement.ClaimsRecovered,
		CommissionDue:    statement.CommissionDue,
		ProfitCommission: statement.ProfitCommission,
		NetBalance:       statement.NetBalance,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update reinsurer statement: %w", err)
	}
	return sqlcReinsurerStatementToDomain(dbStatement), nil
}

func (r *reinsurerStatementRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.store.CountReinsurerStatements(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count reinsurer statements: %w", err)
	}
	return count, nil
}

func sqlcReinsurerStatementToDomain(s db.ReinsurerStatement) *entity.ReinsurerStatement {
	return &entity.ReinsurerStatement{
		ID:               s.ID,
		StatementNumber:  s.StatementNumber,
		TreatyID:         s.TreatyID,
		ParticipantID:    s.ParticipantID,
		PeriodStart:      pgtypeDateToTime(s.PeriodStart),
		PeriodEnd:        pgtypeDateToTime(s.PeriodEnd),
		PremiumCeded:     s.PremiumCeded,
		ClaimsRecovered:  s.ClaimsRecovered,
		CommissionDue:    s.CommissionDue,
		ProfitCommission: s.ProfitCommission,
		NetBalance:       s.NetBalance,
		Status:           s.Status,
		CreatedBy:        s.CreatedBy,
		CreatedAt:        s.CreatedAt,
		UpdatedAt:        s.UpdatedAt,
	}
}
