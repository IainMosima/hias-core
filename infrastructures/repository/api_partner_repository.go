package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/claims/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type apiPartnerRepository struct {
	store db.Store
}

func NewAPIPartnerRepository(store db.Store) domainRepo.APIPartnerRepository {
	return &apiPartnerRepository{store: store}
}

func (r *apiPartnerRepository) Create(ctx context.Context, partner *entity.APIPartner) (*entity.APIPartner, error) {
	metadata, _ := json.Marshal(partner.Metadata)
	if partner.Metadata == nil {
		metadata = []byte("{}")
	}
	dbPartner, err := r.store.CreateAPIPartner(ctx, db.CreateAPIPartnerParams{
		Name:               partner.Name,
		PartnerType:        partner.PartnerType,
		ApiKey:             partner.APIKey,
		ApiSecretHash:      partner.APISecretHash,
		ProviderID:         uuidToPgtype(partner.ProviderID),
		IsActive:           partner.IsActive,
		RateLimitPerMinute: int32(partner.RateLimitPerMinute),
		AllowedClaimTypes:  partner.AllowedClaimTypes,
		WebhookUrl:         stringToPgtypeText(partner.WebhookURL),
		ContactEmail:       stringToPgtypeText(partner.ContactEmail),
		Metadata:           metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create API partner: %w", err)
	}
	return sqlcAPIPartnerToDomain(dbPartner), nil
}

func (r *apiPartnerRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.APIPartner, error) {
	dbPartner, err := r.store.GetAPIPartnerByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get API partner by ID: %w", err)
	}
	return sqlcAPIPartnerToDomain(dbPartner), nil
}

func (r *apiPartnerRepository) GetByAPIKey(ctx context.Context, apiKey string) (*entity.APIPartner, error) {
	dbPartner, err := r.store.GetAPIPartnerByAPIKey(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get API partner by API key: %w", err)
	}
	return sqlcAPIPartnerToDomain(dbPartner), nil
}

func (r *apiPartnerRepository) List(ctx context.Context, limit, offset int) ([]*entity.APIPartner, error) {
	dbPartners, err := r.store.ListAPIPartners(ctx, db.ListAPIPartnersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list API partners: %w", err)
	}
	partners := make([]*entity.APIPartner, len(dbPartners))
	for i, p := range dbPartners {
		partners[i] = sqlcAPIPartnerToDomain(p)
	}
	return partners, nil
}

func (r *apiPartnerRepository) Update(ctx context.Context, partner *entity.APIPartner) (*entity.APIPartner, error) {
	metadata, _ := json.Marshal(partner.Metadata)
	if partner.Metadata == nil {
		metadata = []byte("{}")
	}
	dbPartner, err := r.store.UpdateAPIPartner(ctx, db.UpdateAPIPartnerParams{
		ID:                 partner.ID,
		Name:               partner.Name,
		PartnerType:        partner.PartnerType,
		ProviderID:         uuidToPgtype(partner.ProviderID),
		RateLimitPerMinute: int32(partner.RateLimitPerMinute),
		AllowedClaimTypes:  partner.AllowedClaimTypes,
		WebhookUrl:         stringToPgtypeText(partner.WebhookURL),
		ContactEmail:       stringToPgtypeText(partner.ContactEmail),
		Metadata:           metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update API partner: %w", err)
	}
	return sqlcAPIPartnerToDomain(dbPartner), nil
}

func (r *apiPartnerRepository) Deactivate(ctx context.Context, id uuid.UUID) error {
	return r.store.DeactivateAPIPartner(ctx, id)
}

func (r *apiPartnerRepository) UpdateAPIKey(ctx context.Context, id uuid.UUID, apiKey, apiSecretHash string) (*entity.APIPartner, error) {
	dbPartner, err := r.store.UpdateAPIPartnerAPIKey(ctx, db.UpdateAPIPartnerAPIKeyParams{
		ID:            id,
		ApiKey:        apiKey,
		ApiSecretHash: apiSecretHash,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update API partner key: %w", err)
	}
	return sqlcAPIPartnerToDomain(dbPartner), nil
}

func sqlcAPIPartnerToDomain(p db.ApiPartner) *entity.APIPartner {
	return &entity.APIPartner{
		ID:                 p.ID,
		Name:               p.Name,
		PartnerType:        p.PartnerType,
		APIKey:             p.ApiKey,
		APISecretHash:      p.ApiSecretHash,
		ProviderID:         pgtypeToUUID(p.ProviderID),
		IsActive:           p.IsActive,
		RateLimitPerMinute: int(p.RateLimitPerMinute),
		AllowedClaimTypes:  p.AllowedClaimTypes,
		WebhookURL:         p.WebhookUrl.String,
		ContactEmail:       p.ContactEmail.String,
		Metadata:           p.Metadata,
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
	}
}
