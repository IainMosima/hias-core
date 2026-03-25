package service

import (
	"context"

	claimsSchema "github.com/bitbiz/hias-core/domains/claims/schema"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/google/uuid"
)

type APIPartnerService interface {
	CreatePartner(ctx context.Context, req claimsSchema.CreateAPIPartnerRequest) *schema.ServiceResponse[claimsSchema.CreateAPIPartnerResponse]
	ListPartners(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]claimsSchema.APIPartnerResponse]
	GetPartner(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[claimsSchema.APIPartnerResponse]
	UpdatePartner(ctx context.Context, id uuid.UUID, req claimsSchema.UpdateAPIPartnerRequest) *schema.ServiceResponse[claimsSchema.APIPartnerResponse]
	DeactivatePartner(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[string]
	RegenerateAPIKey(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[claimsSchema.CreateAPIPartnerResponse]
}
