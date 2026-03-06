package repository

import (
	"context"
	"time"

	"github.com/bitbiz/hias-core/domains/provider/entity"
	"github.com/google/uuid"
)

type ProviderRepository interface {
	Create(ctx context.Context, provider *entity.Provider) (*entity.Provider, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Provider, error)
	GetByLicense(ctx context.Context, license string) (*entity.Provider, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Provider, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Provider, error)
	ListByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.Provider, error)
	Count(ctx context.Context) (int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Provider, error)
	Update(ctx context.Context, provider *entity.Provider) (*entity.Provider, error)
	UpdateTier(ctx context.Context, id uuid.UUID, tier string) (*entity.Provider, error)
	ListByTier(ctx context.Context, tier string, limit, offset int) ([]*entity.Provider, error)
	UpdateAccreditation(ctx context.Context, id uuid.UUID, status string, expiry *time.Time, body string) (*entity.Provider, error)
	ListByAccreditationStatus(ctx context.Context, status string, limit, offset int) ([]*entity.Provider, error)
	ListExpiringAccreditations(ctx context.Context, days, limit, offset int) ([]*entity.Provider, error)
	ListFiltered(ctx context.Context, search string, limit, offset int) ([]*entity.Provider, error)
	CountFiltered(ctx context.Context, search string) (int64, error)
}
