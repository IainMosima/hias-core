package repository

import (
	"context"
	"github.com/bitbiz/hias-core/domains/billing/entity"
	"github.com/google/uuid"
)

type InvoiceRepository interface {
	Create(ctx context.Context, invoice *entity.Invoice) (*entity.Invoice, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Invoice, error)
	GetByNumber(ctx context.Context, number string) (*entity.Invoice, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Invoice, error)
	ListByPolicy(ctx context.Context, policyID uuid.UUID, limit, offset int) ([]*entity.Invoice, error)
	ListByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.Invoice, error)
	ListOverdue(ctx context.Context) ([]*entity.Invoice, error)
	Count(ctx context.Context) (int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Invoice, error)
}
