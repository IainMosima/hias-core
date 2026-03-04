package service

import (
	"context"

	"github.com/bitbiz/hias-core/domains/identity/schema"
	productSchema "github.com/bitbiz/hias-core/domains/product/schema"
	"github.com/google/uuid"
)

type ExclusionService interface {
	CreateExclusion(ctx context.Context, planID uuid.UUID, req productSchema.CreateExclusionRequest) *schema.ServiceResponse[productSchema.ExclusionResponse]
	GetExclusion(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[productSchema.ExclusionResponse]
	ListExclusionsByPlan(ctx context.Context, planID uuid.UUID) *schema.ServiceResponse[[]productSchema.ExclusionResponse]
	UpdateExclusion(ctx context.Context, id uuid.UUID, req productSchema.UpdateExclusionRequest) *schema.ServiceResponse[productSchema.ExclusionResponse]
	DeleteExclusion(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[string]
}
