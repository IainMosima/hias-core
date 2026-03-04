package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/billing/entity"
	"github.com/google/uuid"
)

type InstallmentRepository interface {
	Create(ctx context.Context, installment *entity.Installment) (*entity.Installment, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Installment, error)
	ListBySchedule(ctx context.Context, scheduleID uuid.UUID) ([]*entity.Installment, error)
	ListOverdue(ctx context.Context) ([]*entity.Installment, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Installment, error)
	MarkPaid(ctx context.Context, id uuid.UUID, invoiceID uuid.UUID) (*entity.Installment, error)
}
