package repository

import (
	"context"
	"github.com/bitbiz/hias-core/domains/provider/entity"
	"github.com/google/uuid"
)

type ContractRepository interface {
	Create(ctx context.Context, contract *entity.Contract) (*entity.Contract, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Contract, error)
	ListByProvider(ctx context.Context, providerID uuid.UUID) ([]*entity.Contract, error)
	GetActiveByProvider(ctx context.Context, providerID uuid.UUID) (*entity.Contract, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Contract, error)
}
