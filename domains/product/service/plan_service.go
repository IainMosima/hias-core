package service

import (
	"context"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	productSchema "github.com/bitbiz/hias-core/domains/product/schema"
	"github.com/google/uuid"
)

type PlanService interface {
	CreatePlan(ctx context.Context, req productSchema.CreatePlanRequest, createdBy uuid.UUID) *schema.ServiceResponse[productSchema.PlanResponse]
	GetPlan(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[productSchema.PlanResponse]
	ListPlans(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]productSchema.PlanResponse]
	ListPlansBySegment(ctx context.Context, segment string, page, pageSize int) *schema.ServiceResponse[[]productSchema.PlanResponse]
	UpdatePlan(ctx context.Context, id uuid.UUID, req productSchema.UpdatePlanRequest) *schema.ServiceResponse[productSchema.PlanResponse]
	GetTotalCount(ctx context.Context) *schema.ServiceResponse[int64]
}
