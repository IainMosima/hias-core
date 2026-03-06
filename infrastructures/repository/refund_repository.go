package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/bitbiz/hias-core/domains/billing/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/billing/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type refundRepository struct {
	store db.Store
}

func NewRefundRepository(store db.Store) domainRepo.RefundRepository {
	return &refundRepository{store: store}
}

func (r *refundRepository) Create(ctx context.Context, refund *entity.Refund) (*entity.Refund, error) {
	var createdBy pgtype.UUID
	if refund.CreatedBy != uuid.Nil {
		createdBy = pgtype.UUID{Bytes: refund.CreatedBy, Valid: true}
	}
	var creditNoteID pgtype.UUID
	if refund.CreditNoteID != uuid.Nil {
		creditNoteID = pgtype.UUID{Bytes: refund.CreditNoteID, Valid: true}
	}
	dbRefund, err := r.store.CreateRefund(ctx, db.CreateRefundParams{
		PolicyID:     refund.PolicyID,
		CreditNoteID: creditNoteID,
		Amount:       refund.Amount,
		Currency:     refund.Currency,
		Status:       refund.Status,
		Reason:       refund.Reason,
		CreatedBy:    createdBy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create refund: %w", err)
	}
	return sqlcRefundToDomain(dbRefund), nil
}

func (r *refundRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Refund, error) {
	dbRefund, err := r.store.GetRefundByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get refund: %w", err)
	}
	return sqlcRefundToDomain(dbRefund), nil
}

func (r *refundRepository) ListByPolicy(ctx context.Context, policyID uuid.UUID, limit, offset int) ([]*entity.Refund, error) {
	dbRefunds, err := r.store.ListRefundsByPolicy(ctx, db.ListRefundsByPolicyParams{
		PolicyID: policyID, Limit: int32(limit), Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list refunds: %w", err)
	}
	refunds := make([]*entity.Refund, len(dbRefunds))
	for i, r := range dbRefunds {
		refunds[i] = sqlcRefundToDomain(r)
	}
	return refunds, nil
}

func (r *refundRepository) Approve(ctx context.Context, id uuid.UUID, approvedBy uuid.UUID) (*entity.Refund, error) {
	dbRefund, err := r.store.ApproveRefund(ctx, db.ApproveRefundParams{
		ID: id, ApprovedBy: pgtype.UUID{Bytes: approvedBy, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to approve refund: %w", err)
	}
	return sqlcRefundToDomain(dbRefund), nil
}

func (r *refundRepository) Process(ctx context.Context, id uuid.UUID) (*entity.Refund, error) {
	dbRefund, err := r.store.ProcessRefund(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to process refund: %w", err)
	}
	return sqlcRefundToDomain(dbRefund), nil
}

func sqlcRefundToDomain(r db.Refund) *entity.Refund {
	var createdBy, approvedBy, creditNoteID uuid.UUID
	if r.CreatedBy.Valid {
		createdBy = r.CreatedBy.Bytes
	}
	if r.ApprovedBy.Valid {
		approvedBy = r.ApprovedBy.Bytes
	}
	if r.CreditNoteID.Valid {
		creditNoteID = r.CreditNoteID.Bytes
	}
	var approvedAt, processedAt *time.Time
	if r.ApprovedAt.Valid {
		approvedAt = &r.ApprovedAt.Time
	}
	if r.ProcessedAt.Valid {
		processedAt = &r.ProcessedAt.Time
	}
	return &entity.Refund{
		ID: r.ID, PolicyID: r.PolicyID, CreditNoteID: creditNoteID,
		Amount: r.Amount, Currency: r.Currency, Status: r.Status,
		Reason: r.Reason, ApprovedBy: approvedBy,
		ApprovedAt: approvedAt, ProcessedAt: processedAt,
		CreatedBy: createdBy, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
}
