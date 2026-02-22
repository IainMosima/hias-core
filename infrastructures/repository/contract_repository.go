package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/provider/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/provider/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type contractRepository struct {
	store db.Store
}

func NewContractRepository(store db.Store) domainRepo.ContractRepository {
	return &contractRepository{store: store}
}

func (r *contractRepository) Create(ctx context.Context, contract *entity.Contract) (*entity.Contract, error) {
	dbContract, err := r.store.CreateContract(ctx, db.CreateContractParams{
		ProviderID: contract.ProviderID,
		StartDate:  contract.StartDate,
		EndDate:    contract.EndDate,
		Terms:      contract.Terms,
		Status:     contract.Status,
		CreatedBy:  uuidToPgtype(contract.CreatedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create contract: %w", err)
	}
	return sqlcContractToDomain(dbContract), nil
}

func (r *contractRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Contract, error) {
	dbContract, err := r.store.GetContractByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract by ID: %w", err)
	}
	return sqlcContractToDomain(dbContract), nil
}

func (r *contractRepository) ListByProvider(ctx context.Context, providerID uuid.UUID) ([]*entity.Contract, error) {
	dbContracts, err := r.store.ListContractsByProvider(ctx, providerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list contracts by provider: %w", err)
	}
	contracts := make([]*entity.Contract, len(dbContracts))
	for i, c := range dbContracts {
		contracts[i] = sqlcContractToDomain(c)
	}
	return contracts, nil
}

func (r *contractRepository) GetActiveByProvider(ctx context.Context, providerID uuid.UUID) (*entity.Contract, error) {
	dbContract, err := r.store.GetActiveContractByProvider(ctx, providerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active contract by provider: %w", err)
	}
	return sqlcContractToDomain(dbContract), nil
}

func (r *contractRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Contract, error) {
	dbContract, err := r.store.UpdateContractStatus(ctx, db.UpdateContractStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update contract status: %w", err)
	}
	return sqlcContractToDomain(dbContract), nil
}

func sqlcContractToDomain(c db.Contract) *entity.Contract {
	return &entity.Contract{
		ID:         c.ID,
		ProviderID: c.ProviderID,
		StartDate:  c.StartDate,
		EndDate:    c.EndDate,
		Terms:      c.Terms,
		Status:     c.Status,
		CreatedBy:  pgtypeToUUID(c.CreatedBy),
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}
}
