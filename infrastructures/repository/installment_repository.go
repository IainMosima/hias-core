package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/billing/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/billing/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type installmentRepository struct {
	store db.Store
}

func NewInstallmentRepository(store db.Store) domainRepo.InstallmentRepository {
	return &installmentRepository{store: store}
}

func (r *installmentRepository) Create(ctx context.Context, installment *entity.Installment) (*entity.Installment, error) {
	dbInstallment, err := r.store.CreateInstallment(ctx, db.CreateInstallmentParams{
		ScheduleID:        installment.ScheduleID,
		InstallmentNumber: int32(installment.InstallmentNumber),
		DueDate:           installment.DueDate,
		Amount:            installment.Amount,
		Status:            installment.Status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create installment: %w", err)
	}
	return sqlcInstallmentToDomain(dbInstallment), nil
}

func (r *installmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Installment, error) {
	dbInstallment, err := r.store.GetInstallmentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get installment by ID: %w", err)
	}
	return sqlcInstallmentToDomain(dbInstallment), nil
}

func (r *installmentRepository) ListBySchedule(ctx context.Context, scheduleID uuid.UUID) ([]*entity.Installment, error) {
	dbInstallments, err := r.store.ListInstallmentsBySchedule(ctx, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to list installments by schedule: %w", err)
	}
	installments := make([]*entity.Installment, len(dbInstallments))
	for i, inst := range dbInstallments {
		installments[i] = sqlcInstallmentToDomain(inst)
	}
	return installments, nil
}

func (r *installmentRepository) ListOverdue(ctx context.Context) ([]*entity.Installment, error) {
	dbInstallments, err := r.store.ListOverdueInstallments(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list overdue installments: %w", err)
	}
	installments := make([]*entity.Installment, len(dbInstallments))
	for i, inst := range dbInstallments {
		installments[i] = sqlcInstallmentToDomain(inst)
	}
	return installments, nil
}

func (r *installmentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Installment, error) {
	dbInstallment, err := r.store.UpdateInstallmentStatus(ctx, db.UpdateInstallmentStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update installment status: %w", err)
	}
	return sqlcInstallmentToDomain(dbInstallment), nil
}

func (r *installmentRepository) MarkPaid(ctx context.Context, id uuid.UUID, invoiceID uuid.UUID) (*entity.Installment, error) {
	dbInstallment, err := r.store.MarkInstallmentPaid(ctx, db.MarkInstallmentPaidParams{
		ID:        id,
		InvoiceID: uuidToPgtype(invoiceID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to mark installment paid: %w", err)
	}
	return sqlcInstallmentToDomain(dbInstallment), nil
}

func sqlcInstallmentToDomain(i db.Installment) *entity.Installment {
	return &entity.Installment{
		ID:                i.ID,
		ScheduleID:        i.ScheduleID,
		InstallmentNumber: int(i.InstallmentNumber),
		DueDate:           i.DueDate,
		Amount:            i.Amount,
		Status:            i.Status,
		PaidAt:            pgtypeTimestamptzToTimePtr(i.PaidAt),
		InvoiceID:         pgtypeToUUID(i.InvoiceID),
		CreatedAt:         i.CreatedAt,
		UpdatedAt:         i.UpdatedAt,
	}
}
