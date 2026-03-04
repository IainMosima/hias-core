package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/sales/entity"
	"github.com/google/uuid"
)

type QuotationRepository interface {
	Create(ctx context.Context, quotation *entity.Quotation) (*entity.Quotation, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Quotation, error)
	GetByNumber(ctx context.Context, number string) (*entity.Quotation, error)
	ListByLead(ctx context.Context, leadID uuid.UUID, limit, offset int) ([]*entity.Quotation, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Quotation, error)
	ListByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.Quotation, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Quotation, error)
	UpdateCurrentVersion(ctx context.Context, id uuid.UUID, version int) (*entity.Quotation, error)
	SetPolicyID(ctx context.Context, id uuid.UUID, policyID uuid.UUID) (*entity.Quotation, error)
	Count(ctx context.Context) (int64, error)
	ListExpired(ctx context.Context) ([]*entity.Quotation, error)
}
