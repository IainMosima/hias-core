package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bitbiz/hias-core/domains/billing/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/billing/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type paymentRepository struct {
	store db.Store
}

func NewPaymentRepository(store db.Store) domainRepo.PaymentRepository {
	return &paymentRepository{store: store}
}

func (r *paymentRepository) Create(ctx context.Context, payment *entity.Payment) (*entity.Payment, error) {
	dbPayment, err := r.store.CreatePayment(ctx, db.CreatePaymentParams{
		InvoiceID:       uuidToPgtype(payment.InvoiceID),
		ClaimID:         uuidToPgtype(payment.ClaimID),
		Type:            payment.Type,
		Amount:          payment.Amount,
		Currency:        payment.Currency,
		Method:          payment.Method,
		ReferenceNumber: stringToPgtypeText(payment.ReferenceNumber),
		Status:          payment.Status,
		CreatedBy:       uuidToPgtype(payment.CreatedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}
	return sqlcPaymentToDomain(dbPayment), nil
}

func (r *paymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
	dbPayment, err := r.store.GetPaymentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment by ID: %w", err)
	}
	return sqlcPaymentToDomain(dbPayment), nil
}

func (r *paymentRepository) GetByReference(ctx context.Context, reference string) (*entity.Payment, error) {
	dbPayment, err := r.store.GetPaymentByReference(ctx, stringToPgtypeText(reference))
	if err != nil {
		return nil, fmt.Errorf("failed to get payment by reference: %w", err)
	}
	return sqlcPaymentToDomain(dbPayment), nil
}

func (r *paymentRepository) List(ctx context.Context, limit, offset int) ([]*entity.Payment, error) {
	dbPayments, err := r.store.ListPayments(ctx, db.ListPaymentsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}
	payments := make([]*entity.Payment, len(dbPayments))
	for i, p := range dbPayments {
		payments[i] = sqlcPaymentToDomain(p)
	}
	return payments, nil
}

func (r *paymentRepository) ListByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.Payment, error) {
	dbPayments, err := r.store.ListPaymentsByStatus(ctx, db.ListPaymentsByStatusParams{
		Status: status,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list payments by status: %w", err)
	}
	payments := make([]*entity.Payment, len(dbPayments))
	for i, p := range dbPayments {
		payments[i] = sqlcPaymentToDomain(p)
	}
	return payments, nil
}

func (r *paymentRepository) ListByInvoice(ctx context.Context, invoiceID uuid.UUID) ([]*entity.Payment, error) {
	dbPayments, err := r.store.ListPaymentsByInvoice(ctx, uuidToPgtype(invoiceID))
	if err != nil {
		return nil, fmt.Errorf("failed to list payments by invoice: %w", err)
	}
	payments := make([]*entity.Payment, len(dbPayments))
	for i, p := range dbPayments {
		payments[i] = sqlcPaymentToDomain(p)
	}
	return payments, nil
}

func (r *paymentRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.store.CountPayments(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count payments: %w", err)
	}
	return count, nil
}

func (r *paymentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Payment, error) {
	dbPayment, err := r.store.UpdatePaymentStatus(ctx, db.UpdatePaymentStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update payment status: %w", err)
	}
	return sqlcPaymentToDomain(dbPayment), nil
}

func (r *paymentRepository) Confirm(ctx context.Context, id uuid.UUID, gatewayResponse json.RawMessage) (*entity.Payment, error) {
	dbPayment, err := r.store.ConfirmPayment(ctx, db.ConfirmPaymentParams{
		ID:              id,
		GatewayResponse: gatewayResponse,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to confirm payment: %w", err)
	}
	return sqlcPaymentToDomain(dbPayment), nil
}

func (r *paymentRepository) Reconcile(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
	dbPayment, err := r.store.ReconcilePayment(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to reconcile payment: %w", err)
	}
	return sqlcPaymentToDomain(dbPayment), nil
}

func (r *paymentRepository) IncrementRetry(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
	dbPayment, err := r.store.IncrementPaymentRetry(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to increment payment retry: %w", err)
	}
	return sqlcPaymentToDomain(dbPayment), nil
}

func (r *paymentRepository) GetFailedForRetry(ctx context.Context, limit int) ([]*entity.Payment, error) {
	dbPayments, err := r.store.GetFailedPaymentsForRetry(ctx, int32(limit))
	if err != nil {
		return nil, fmt.Errorf("failed to get failed payments for retry: %w", err)
	}
	payments := make([]*entity.Payment, len(dbPayments))
	for i, p := range dbPayments {
		payments[i] = sqlcPaymentToDomain(p)
	}
	return payments, nil
}

func sqlcPaymentToDomain(p db.Payment) *entity.Payment {
	return &entity.Payment{
		ID:              p.ID,
		InvoiceID:       pgtypeToUUID(p.InvoiceID),
		ClaimID:         pgtypeToUUID(p.ClaimID),
		Type:            p.Type,
		Amount:          p.Amount,
		Currency:        p.Currency,
		Method:          p.Method,
		ReferenceNumber: p.ReferenceNumber.String,
		Status:          p.Status,
		RetryCount:      int(p.RetryCount),
		MaxRetries:      int(p.MaxRetries),
		GatewayResponse: p.GatewayResponse,
		PaidAt:          pgtypeTimestamptzToTimePtr(p.PaidAt),
		ReconciledAt:    pgtypeTimestamptzToTimePtr(p.ReconciledAt),
		CreatedBy:       pgtypeToUUID(p.CreatedBy),
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}
