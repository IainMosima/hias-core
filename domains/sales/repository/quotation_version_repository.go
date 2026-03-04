package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/sales/entity"
	"github.com/google/uuid"
)

type QuotationVersionRepository interface {
	Create(ctx context.Context, version *entity.QuotationVersion) (*entity.QuotationVersion, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.QuotationVersion, error)
	GetByQuotationAndVersion(ctx context.Context, quotationID uuid.UUID, versionNumber int) (*entity.QuotationVersion, error)
	ListByQuotation(ctx context.Context, quotationID uuid.UUID) ([]*entity.QuotationVersion, error)
	GetLatestByQuotation(ctx context.Context, quotationID uuid.UUID) (*entity.QuotationVersion, error)
	UpdateApprovalStatus(ctx context.Context, id uuid.UUID, status string, approvedBy uuid.UUID) (*entity.QuotationVersion, error)
	RejectVersion(ctx context.Context, id uuid.UUID, reason string, rejectedBy uuid.UUID) (*entity.QuotationVersion, error)
}
