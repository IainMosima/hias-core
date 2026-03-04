package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/product/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/product/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type providerNetworkRepository struct {
	store db.Store
}

func NewProviderNetworkRepository(store db.Store) domainRepo.ProviderNetworkRepository {
	return &providerNetworkRepository{store: store}
}

func (r *providerNetworkRepository) Create(ctx context.Context, network *entity.ProviderNetwork) (*entity.ProviderNetwork, error) {
	dbNetwork, err := r.store.CreateProviderNetwork(ctx, db.CreateProviderNetworkParams{
		PlanID:          network.PlanID,
		ProviderID:      network.ProviderID,
		BenefitCategory: stringToPgtypeText(network.BenefitCategory),
		Status:          network.Status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create provider network: %w", err)
	}
	return sqlcProviderNetworkToDomain(dbNetwork), nil
}

func (r *providerNetworkRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ProviderNetwork, error) {
	dbNetwork, err := r.store.GetProviderNetworkByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider network by ID: %w", err)
	}
	return sqlcProviderNetworkToDomain(dbNetwork), nil
}

func (r *providerNetworkRepository) ListByPlan(ctx context.Context, planID uuid.UUID) ([]*entity.ProviderNetwork, error) {
	dbNetworks, err := r.store.ListProviderNetworksByPlan(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("failed to list provider networks by plan: %w", err)
	}
	networks := make([]*entity.ProviderNetwork, len(dbNetworks))
	for i, n := range dbNetworks {
		networks[i] = sqlcProviderNetworkToDomain(n)
	}
	return networks, nil
}

func (r *providerNetworkRepository) ListByProvider(ctx context.Context, providerID uuid.UUID) ([]*entity.ProviderNetwork, error) {
	dbNetworks, err := r.store.ListProviderNetworksByProvider(ctx, providerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list provider networks by provider: %w", err)
	}
	networks := make([]*entity.ProviderNetwork, len(dbNetworks))
	for i, n := range dbNetworks {
		networks[i] = sqlcProviderNetworkToDomain(n)
	}
	return networks, nil
}

func (r *providerNetworkRepository) CheckEligibility(ctx context.Context, planID, providerID uuid.UUID, category string) (bool, error) {
	count, err := r.store.CheckProviderNetworkEligibility(ctx, db.CheckProviderNetworkEligibilityParams{
		PlanID:          planID,
		ProviderID:      providerID,
		BenefitCategory: stringToPgtypeText(category),
	})
	if err != nil {
		return false, fmt.Errorf("failed to check provider network eligibility: %w", err)
	}
	return count > 0, nil
}

func (r *providerNetworkRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.ProviderNetwork, error) {
	dbNetwork, err := r.store.UpdateProviderNetworkStatus(ctx, db.UpdateProviderNetworkStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update provider network status: %w", err)
	}
	return sqlcProviderNetworkToDomain(dbNetwork), nil
}

func (r *providerNetworkRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.store.DeleteProviderNetwork(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete provider network: %w", err)
	}
	return nil
}

func sqlcProviderNetworkToDomain(n db.ProviderNetwork) *entity.ProviderNetwork {
	return &entity.ProviderNetwork{
		ID:              n.ID,
		PlanID:          n.PlanID,
		ProviderID:      n.ProviderID,
		BenefitCategory: n.BenefitCategory.String,
		Status:          n.Status,
		CreatedAt:       n.CreatedAt,
		UpdatedAt:       n.UpdatedAt,
	}
}
