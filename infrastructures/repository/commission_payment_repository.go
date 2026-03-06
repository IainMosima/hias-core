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

type commissionPaymentRepository struct {
	store db.Store
}

func NewCommissionPaymentRepository(store db.Store) domainRepo.CommissionPaymentRepository {
	return &commissionPaymentRepository{store: store}
}

func (r *commissionPaymentRepository) Create(ctx context.Context, payment *entity.CommissionPayment) (*entity.CommissionPayment, error) {
	var createdBy pgtype.UUID
	if payment.CreatedBy != uuid.Nil {
		createdBy = pgtype.UUID{Bytes: payment.CreatedBy, Valid: true}
	}
	dbPayment, err := r.store.CreateCommissionPayment(ctx, db.CreateCommissionPaymentParams{
		PolicyID:         payment.PolicyID,
		IntermediaryID:   payment.IntermediaryID,
		CommissionRuleID: payment.CommissionRuleID,
		Amount:           payment.Amount,
		Currency:         payment.Currency,
		Status:           payment.Status,
		PeriodStart:      payment.PeriodStart,
		PeriodEnd:        payment.PeriodEnd,
		CreatedBy:        createdBy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create commission payment: %w", err)
	}
	return sqlcCommissionPaymentToDomain(dbPayment), nil
}

func (r *commissionPaymentRepository) List(ctx context.Context, limit, offset int) ([]*entity.CommissionPayment, error) {
	dbPayments, err := r.store.ListCommissionPayments(ctx, db.ListCommissionPaymentsParams{
		Limit: int32(limit), Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list commission payments: %w", err)
	}
	payments := make([]*entity.CommissionPayment, len(dbPayments))
	for i, p := range dbPayments {
		payments[i] = sqlcCommissionPaymentToDomain(p)
	}
	return payments, nil
}

func (r *commissionPaymentRepository) ListByIntermediary(ctx context.Context, intermediaryID uuid.UUID, limit, offset int) ([]*entity.CommissionPayment, error) {
	dbPayments, err := r.store.ListCommissionPaymentsByIntermediary(ctx, db.ListCommissionPaymentsByIntermediaryParams{
		IntermediaryID: intermediaryID, Limit: int32(limit), Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list commission payments: %w", err)
	}
	payments := make([]*entity.CommissionPayment, len(dbPayments))
	for i, p := range dbPayments {
		payments[i] = sqlcCommissionPaymentToDomain(p)
	}
	return payments, nil
}

func (r *commissionPaymentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.CommissionPayment, error) {
	dbPayment, err := r.store.UpdateCommissionPaymentStatus(ctx, db.UpdateCommissionPaymentStatusParams{
		ID: id, Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update commission payment status: %w", err)
	}
	return sqlcCommissionPaymentToDomain(dbPayment), nil
}

func sqlcCommissionPaymentToDomain(p db.CommissionPayment) *entity.CommissionPayment {
	var createdBy uuid.UUID
	if p.CreatedBy.Valid {
		createdBy = p.CreatedBy.Bytes
	}
	var paidAt *time.Time
	if p.PaidAt.Valid {
		paidAt = &p.PaidAt.Time
	}
	return &entity.CommissionPayment{
		ID: p.ID, PolicyID: p.PolicyID, IntermediaryID: p.IntermediaryID,
		CommissionRuleID: p.CommissionRuleID, Amount: p.Amount,
		Currency: p.Currency, Status: p.Status,
		PeriodStart: p.PeriodStart, PeriodEnd: p.PeriodEnd,
		PaidAt: paidAt, CreatedBy: createdBy,
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}
