package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	"github.com/google/uuid"
)

type APIPartnerRepository interface {
	Create(ctx context.Context, partner *entity.APIPartner) (*entity.APIPartner, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.APIPartner, error)
	GetByAPIKey(ctx context.Context, apiKey string) (*entity.APIPartner, error)
	List(ctx context.Context, limit, offset int) ([]*entity.APIPartner, error)
	Update(ctx context.Context, partner *entity.APIPartner) (*entity.APIPartner, error)
	Deactivate(ctx context.Context, id uuid.UUID) error
	UpdateAPIKey(ctx context.Context, id uuid.UUID, apiKey, apiSecretHash string) (*entity.APIPartner, error)
}
