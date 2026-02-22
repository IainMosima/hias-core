package repository

import (
	"context"
	"encoding/json"
	"github.com/bitbiz/hias-core/domains/billing/entity"
	"github.com/google/uuid"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *entity.Payment) (*entity.Payment, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error)
	GetByReference(ctx context.Context, reference string) (*entity.Payment, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Payment, error)
	ListByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.Payment, error)
	ListByInvoice(ctx context.Context, invoiceID uuid.UUID) ([]*entity.Payment, error)
	Count(ctx context.Context) (int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Payment, error)
	Confirm(ctx context.Context, id uuid.UUID, gatewayResponse json.RawMessage) (*entity.Payment, error)
	Reconcile(ctx context.Context, id uuid.UUID) (*entity.Payment, error)
	IncrementRetry(ctx context.Context, id uuid.UUID) (*entity.Payment, error)
	GetFailedForRetry(ctx context.Context, limit int) ([]*entity.Payment, error)
}
