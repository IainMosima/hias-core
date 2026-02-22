package repository

import (
	"context"
	"github.com/bitbiz/hias-core/domains/billing/entity"
	"github.com/google/uuid"
)

type RemittanceRepository interface {
	Create(ctx context.Context, remittance *entity.Remittance) (*entity.Remittance, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Remittance, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Remittance, error)
	ListByProvider(ctx context.Context, providerID uuid.UUID, limit, offset int) ([]*entity.Remittance, error)
	ListByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.Remittance, error)
	Count(ctx context.Context) (int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Remittance, error)
	MarkAdviceSent(ctx context.Context, id uuid.UUID) (*entity.Remittance, error)
	SetPayment(ctx context.Context, id uuid.UUID, paymentID uuid.UUID) (*entity.Remittance, error)
}
