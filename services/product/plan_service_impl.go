package product

import (
	"context"
	"net/http"

	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/domains/product/entity"
	"github.com/bitbiz/hias-core/domains/product/repository"
	productSchema "github.com/bitbiz/hias-core/domains/product/schema"
	"github.com/bitbiz/hias-core/domains/product/service"
	"github.com/google/uuid"
)

type planServiceImpl struct {
	planRepo      repository.PlanRepository
	benefitRepo   repository.BenefitRepository
	exclusionRepo repository.ExclusionRepository
}

func NewPlanService(
	planRepo repository.PlanRepository,
	benefitRepo repository.BenefitRepository,
	exclusionRepo repository.ExclusionRepository,
) service.PlanService {
	return &planServiceImpl{
		planRepo:      planRepo,
		benefitRepo:   benefitRepo,
		exclusionRepo: exclusionRepo,
	}
}

func (s *planServiceImpl) CreatePlan(ctx context.Context, req productSchema.CreatePlanRequest, createdBy uuid.UUID) *schema.ServiceResponse[productSchema.PlanResponse] {
	currency := req.Currency
	if currency == "" {
		currency = "KES"
	}

	plan, err := s.planRepo.Create(ctx, &entity.Plan{
		Name:        req.Name,
		Type:        req.Type,
		BasePremium: req.BasePremium,
		Currency:    currency,
		Status:      "ACTIVE",
		Description: req.Description,
		CreatedBy:   createdBy,
	})
	if err != nil {
		return schema.NewServiceErrorResponse[productSchema.PlanResponse](http.StatusInternalServerError, "Failed to create plan", err)
	}

	return schema.NewServiceResponse(productSchema.ToPlanResponse(plan), http.StatusCreated, "Plan created successfully")
}

func (s *planServiceImpl) GetPlan(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[productSchema.PlanResponse] {
	plan, err := s.planRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[productSchema.PlanResponse](http.StatusNotFound, "Plan not found", err)
	}

	return schema.NewServiceResponse(productSchema.ToPlanResponse(plan), http.StatusOK, "Plan retrieved")
}

func (s *planServiceImpl) ListPlans(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]productSchema.PlanResponse] {
	offset := (page - 1) * pageSize
	plans, err := s.planRepo.List(ctx, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]productSchema.PlanResponse](http.StatusInternalServerError, "Failed to list plans", err)
	}

	responses := make([]productSchema.PlanResponse, len(plans))
	for i, p := range plans {
		responses[i] = productSchema.ToPlanResponse(p)
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Plans retrieved")
}

func (s *planServiceImpl) UpdatePlan(ctx context.Context, id uuid.UUID, req productSchema.UpdatePlanRequest) *schema.ServiceResponse[productSchema.PlanResponse] {
	existing, err := s.planRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[productSchema.PlanResponse](http.StatusNotFound, "Plan not found", err)
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Type != nil {
		existing.Type = *req.Type
	}
	if req.BasePremium != nil {
		existing.BasePremium = *req.BasePremium
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}

	updated, err := s.planRepo.Update(ctx, existing)
	if err != nil {
		return schema.NewServiceErrorResponse[productSchema.PlanResponse](http.StatusInternalServerError, "Failed to update plan", err)
	}

	return schema.NewServiceResponse(productSchema.ToPlanResponse(updated), http.StatusOK, "Plan updated")
}

func (s *planServiceImpl) GetTotalCount(ctx context.Context) *schema.ServiceResponse[int64] {
	count, err := s.planRepo.Count(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to count plans", err)
	}
	return schema.NewServiceResponse(count, http.StatusOK, "Count retrieved")
}
