package service

import (
	"context"

	"github.com/bitbiz/hias-core/domains/identity/schema"
	providerSchema "github.com/bitbiz/hias-core/domains/provider/schema"
	"github.com/google/uuid"
)

type RateCardService interface {
	CreateRateCard(ctx context.Context, providerID uuid.UUID, req providerSchema.CreateRateCardRequest) *schema.ServiceResponse[providerSchema.RateCardResponse]
	ListRateCards(ctx context.Context, providerID uuid.UUID) *schema.ServiceResponse[[]providerSchema.RateCardResponse]
	BulkCreateRateCards(ctx context.Context, providerID uuid.UUID, req providerSchema.BulkCreateRateCardRequest) *schema.ServiceResponse[[]providerSchema.RateCardResponse]
}
