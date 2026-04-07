package service

import (
	"context"

	"github.com/bitbiz/hias-core/domains/identity/schema"
	productSchema "github.com/bitbiz/hias-core/domains/product/schema"
	"github.com/google/uuid"
)

type ProviderNetworkService interface {
	CreateProviderNetwork(ctx context.Context, planID uuid.UUID, req productSchema.CreateProviderNetworkRequest) *schema.ServiceResponse[productSchema.ProviderNetworkResponse]
	GetProviderNetwork(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[productSchema.ProviderNetworkResponse]
	ListProviderNetworksByPlan(ctx context.Context, planID uuid.UUID) *schema.ServiceResponse[[]productSchema.ProviderNetworkResponse]
	UpdateProviderNetworkStatus(ctx context.Context, id uuid.UUID, status string) *schema.ServiceResponse[productSchema.ProviderNetworkResponse]
	CheckEligibility(ctx context.Context, planID, providerID uuid.UUID, category string) *schema.ServiceResponse[bool]
	DeleteProviderNetwork(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[string]
}
