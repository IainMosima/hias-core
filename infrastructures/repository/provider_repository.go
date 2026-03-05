package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/bitbiz/hias-core/domains/provider/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/provider/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type providerRepository struct {
	store db.Store
}

func NewProviderRepository(store db.Store) domainRepo.ProviderRepository {
	return &providerRepository{store: store}
}

func (r *providerRepository) Create(ctx context.Context, provider *entity.Provider) (*entity.Provider, error) {
	dbProvider, err := r.store.CreateProvider(ctx, db.CreateProviderParams{
		Name:          provider.Name,
		Type:          provider.Type,
		LicenseNumber: provider.LicenseNumber,
		Status:        provider.Status,
		County:        stringToPgtypeText(provider.County),
		Address:       stringToPgtypeText(provider.Address),
		Phone:         stringToPgtypeText(provider.Phone),
		Email:         stringToPgtypeText(provider.Email),
		ContactPerson: stringToPgtypeText(provider.ContactPerson),
		UserID:        uuidToPgtype(provider.UserID),
		CreatedBy:     uuidToPgtype(provider.CreatedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}
	return sqlcProviderToDomain(dbProvider), nil
}

func (r *providerRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Provider, error) {
	dbProvider, err := r.store.GetProviderByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider by ID: %w", err)
	}
	return sqlcProviderToDomain(dbProvider), nil
}

func (r *providerRepository) GetByLicense(ctx context.Context, license string) (*entity.Provider, error) {
	dbProvider, err := r.store.GetProviderByLicense(ctx, license)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider by license: %w", err)
	}
	return sqlcProviderToDomain(dbProvider), nil
}

func (r *providerRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Provider, error) {
	dbProvider, err := r.store.GetProviderByUserID(ctx, uuidToPgtype(userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get provider by user ID: %w", err)
	}
	return sqlcProviderToDomain(dbProvider), nil
}

func (r *providerRepository) List(ctx context.Context, limit, offset int) ([]*entity.Provider, error) {
	dbProviders, err := r.store.ListProviders(ctx, db.ListProvidersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list providers: %w", err)
	}
	providers := make([]*entity.Provider, len(dbProviders))
	for i, p := range dbProviders {
		providers[i] = sqlcProviderToDomain(p)
	}
	return providers, nil
}

func (r *providerRepository) ListByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.Provider, error) {
	dbProviders, err := r.store.ListProvidersByStatus(ctx, db.ListProvidersByStatusParams{
		Status: status,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list providers by status: %w", err)
	}
	providers := make([]*entity.Provider, len(dbProviders))
	for i, p := range dbProviders {
		providers[i] = sqlcProviderToDomain(p)
	}
	return providers, nil
}

func (r *providerRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.store.CountProviders(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count providers: %w", err)
	}
	return count, nil
}

func (r *providerRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Provider, error) {
	dbProvider, err := r.store.UpdateProviderStatus(ctx, db.UpdateProviderStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update provider status: %w", err)
	}
	return sqlcProviderToDomain(dbProvider), nil
}

func (r *providerRepository) Update(ctx context.Context, provider *entity.Provider) (*entity.Provider, error) {
	dbProvider, err := r.store.UpdateProvider(ctx, db.UpdateProviderParams{
		ID:            provider.ID,
		Name:          stringToPgtypeText(provider.Name),
		County:        stringToPgtypeText(provider.County),
		Address:       stringToPgtypeText(provider.Address),
		Phone:         stringToPgtypeText(provider.Phone),
		Email:         stringToPgtypeText(provider.Email),
		ContactPerson: stringToPgtypeText(provider.ContactPerson),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update provider: %w", err)
	}
	return sqlcProviderToDomain(dbProvider), nil
}

func (r *providerRepository) UpdateTier(ctx context.Context, id uuid.UUID, tier string) (*entity.Provider, error) {
	dbProvider, err := r.store.UpdateProviderTier(ctx, db.UpdateProviderTierParams{
		ID:   id,
		Tier: tier,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update provider tier: %w", err)
	}
	return sqlcProviderToDomain(dbProvider), nil
}

func (r *providerRepository) ListByTier(ctx context.Context, tier string, limit, offset int) ([]*entity.Provider, error) {
	dbProviders, err := r.store.ListProvidersByTier(ctx, db.ListProvidersByTierParams{
		Tier:   tier,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list providers by tier: %w", err)
	}
	providers := make([]*entity.Provider, len(dbProviders))
	for i, p := range dbProviders {
		providers[i] = sqlcProviderToDomain(p)
	}
	return providers, nil
}

func (r *providerRepository) UpdateAccreditation(ctx context.Context, id uuid.UUID, status string, expiry *time.Time, body string) (*entity.Provider, error) {
	dbProvider, err := r.store.UpdateAccreditation(ctx, db.UpdateAccreditationParams{
		ID:                  id,
		AccreditationStatus: stringToPgtypeText(status),
		AccreditationExpiry: timePtrToPgtypeTimestamptz(expiry),
		AccreditationBody:   stringToPgtypeText(body),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update accreditation: %w", err)
	}
	return sqlcProviderToDomain(dbProvider), nil
}

func (r *providerRepository) ListByAccreditationStatus(ctx context.Context, status string, limit, offset int) ([]*entity.Provider, error) {
	dbProviders, err := r.store.ListProvidersByAccreditationStatus(ctx, db.ListProvidersByAccreditationStatusParams{
		AccreditationStatus: stringToPgtypeText(status),
		Limit:               int32(limit),
		Offset:              int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list providers by accreditation status: %w", err)
	}
	providers := make([]*entity.Provider, len(dbProviders))
	for i, p := range dbProviders {
		providers[i] = sqlcProviderToDomain(p)
	}
	return providers, nil
}

func (r *providerRepository) ListExpiringAccreditations(ctx context.Context, days, limit, offset int) ([]*entity.Provider, error) {
	dbProviders, err := r.store.ListExpiringAccreditations(ctx, db.ListExpiringAccreditationsParams{
		Days:   int32(days),
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list expiring accreditations: %w", err)
	}
	providers := make([]*entity.Provider, len(dbProviders))
	for i, p := range dbProviders {
		providers[i] = sqlcProviderToDomain(p)
	}
	return providers, nil
}

func sqlcProviderToDomain(p db.Provider) *entity.Provider {
	return &entity.Provider{
		ID:                  p.ID,
		Name:                p.Name,
		Type:                p.Type,
		LicenseNumber:       p.LicenseNumber,
		Status:              p.Status,
		Tier:                p.Tier,
		County:              p.County.String,
		Address:             p.Address.String,
		Phone:               p.Phone.String,
		Email:               p.Email.String,
		ContactPerson:       p.ContactPerson.String,
		AccreditationStatus: p.AccreditationStatus.String,
		AccreditationExpiry: pgtypeTimestamptzToTimePtr(p.AccreditationExpiry),
		AccreditationBody:   p.AccreditationBody.String,
		UserID:              pgtypeToUUID(p.UserID),
		CreatedBy:           pgtypeToUUID(p.CreatedBy),
		CreatedAt:           p.CreatedAt,
		UpdatedAt:           p.UpdatedAt,
	}
}
