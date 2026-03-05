package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/billing/entity"
	"github.com/google/uuid"
)

type CreditNoteRepository interface {
	Create(ctx context.Context, cn *entity.CreditNote) (*entity.CreditNote, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.CreditNote, error)
	GetByNumber(ctx context.Context, number string) (*entity.CreditNote, error)
	ListByPolicy(ctx context.Context, policyID uuid.UUID) ([]*entity.CreditNote, error)
	ListByStatus(ctx context.Context, status string, limit, offset int32) ([]*entity.CreditNote, error)
	Approve(ctx context.Context, id uuid.UUID, approvedBy uuid.UUID) (*entity.CreditNote, error)
	Apply(ctx context.Context, id uuid.UUID, invoiceID uuid.UUID) (*entity.CreditNote, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.CreditNote, error)
}
