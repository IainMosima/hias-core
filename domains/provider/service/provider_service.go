package service

import (
	"context"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	providerSchema "github.com/bitbiz/hias-core/domains/provider/schema"
	"github.com/google/uuid"
)

type ProviderService interface {
	RegisterProvider(ctx context.Context, req providerSchema.RegisterProviderRequest, createdBy uuid.UUID) *schema.ServiceResponse[providerSchema.ProviderResponse]
	GetProvider(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[providerSchema.ProviderResponse]
	ListProviders(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]providerSchema.ProviderResponse]
	CredentialProvider(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[providerSchema.ProviderResponse]
	ActivateProvider(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[providerSchema.ProviderResponse]
	SuspendProvider(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[providerSchema.ProviderResponse]
	TerminateProvider(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[providerSchema.ProviderResponse]
	UpdateProvider(ctx context.Context, id uuid.UUID, req providerSchema.UpdateProviderRequest) *schema.ServiceResponse[providerSchema.ProviderResponse]
	GetTotalCount(ctx context.Context) *schema.ServiceResponse[int64]
}
