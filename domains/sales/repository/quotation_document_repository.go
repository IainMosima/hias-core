package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/sales/entity"
	"github.com/google/uuid"
)

type QuotationDocumentRepository interface {
	Create(ctx context.Context, doc *entity.QuotationDocument) (*entity.QuotationDocument, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.QuotationDocument, error)
	ListByQuotation(ctx context.Context, quotationID uuid.UUID) ([]*entity.QuotationDocument, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, doc *entity.QuotationDocument) (*entity.QuotationDocument, error)
}
