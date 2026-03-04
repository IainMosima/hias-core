package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/product/entity"
	"github.com/google/uuid"
)

type ProviderNetworkRepository interface {
	Create(ctx context.Context, network *entity.ProviderNetwork) (*entity.ProviderNetwork, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ProviderNetwork, error)
	ListByPlan(ctx context.Context, planID uuid.UUID) ([]*entity.ProviderNetwork, error)
	ListByProvider(ctx context.Context, providerID uuid.UUID) ([]*entity.ProviderNetwork, error)
	CheckEligibility(ctx context.Context, planID, providerID uuid.UUID, category string) (bool, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.ProviderNetwork, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
