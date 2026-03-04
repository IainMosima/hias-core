package product

import (
	"context"
	"log"
	"net/http"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/domains/product/entity"
	"github.com/bitbiz/hias-core/domains/product/repository"
	productSchema "github.com/bitbiz/hias-core/domains/product/schema"
	"github.com/bitbiz/hias-core/domains/product/service"
	"github.com/bitbiz/hias-core/shared"
	"github.com/google/uuid"
)

type planServiceImpl struct {
	planRepo repository.PlanRepository
	auditSvc auditService.AuditService
}

func NewPlanService(
	planRepo repository.PlanRepository,
	auditSvc auditService.AuditService,
) service.PlanService {
	return &planServiceImpl{
		planRepo: planRepo,
		auditSvc: auditSvc,
	}
}

func (s *planServiceImpl) CreatePlan(ctx context.Context, req productSchema.CreatePlanRequest, createdBy uuid.UUID) *schema.ServiceResponse[productSchema.PlanResponse] {
	currency := req.Currency
	if currency == "" {
		currency = string(shared.CurrencyKES)
	}
	segment := req.Segment
	if segment == "" {
		segment = string(shared.PlanSegmentRetail)
	}

	plan := &entity.Plan{
		Name:        req.Name,
		Type:        req.Type,
		Segment:     segment,
		BasePremium: req.BasePremium,
		Currency:    currency,
		Status:      string(shared.PlanStatusActive),
		Description: req.Description,
		CreatedBy:   createdBy,
	}

	created, err := s.planRepo.Create(ctx, plan)
	if err != nil {
		return schema.NewServiceErrorResponse[productSchema.PlanResponse](http.StatusInternalServerError, "Failed to create plan", err)
	}

	s.logAudit(ctx, createdBy, string(shared.AuditEntityTypePlan), created.ID, string(shared.AuditActionCreate))

	return schema.NewServiceResponse(productSchema.ToPlanResponse(created), http.StatusCreated, "Plan created")
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

func (s *planServiceImpl) ListPlansBySegment(ctx context.Context, segment string, page, pageSize int) *schema.ServiceResponse[[]productSchema.PlanResponse] {
	offset := (page - 1) * pageSize
	plans, err := s.planRepo.ListBySegment(ctx, segment, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]productSchema.PlanResponse](http.StatusInternalServerError, "Failed to list plans by segment", err)
	}

	responses := make([]productSchema.PlanResponse, len(plans))
	for i, p := range plans {
		responses[i] = productSchema.ToPlanResponse(p)
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Plans retrieved")
}

func (s *planServiceImpl) UpdatePlan(ctx context.Context, id uuid.UUID, req productSchema.UpdatePlanRequest) *schema.ServiceResponse[productSchema.PlanResponse] {
	plan, err := s.planRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[productSchema.PlanResponse](http.StatusNotFound, "Plan not found", err)
	}

	if req.Name != nil {
		plan.Name = *req.Name
	}
	if req.Type != nil {
		plan.Type = *req.Type
	}
	if req.Segment != nil {
		plan.Segment = *req.Segment
	}
	if req.BasePremium != nil {
		plan.BasePremium = *req.BasePremium
	}
	if req.Description != nil {
		plan.Description = *req.Description
	}
	if req.Status != nil {
		plan.Status = *req.Status
	}

	updated, err := s.planRepo.Update(ctx, plan)
	if err != nil {
		return schema.NewServiceErrorResponse[productSchema.PlanResponse](http.StatusInternalServerError, "Failed to update plan", err)
	}

	return schema.NewServiceResponse(productSchema.ToPlanResponse(updated), http.StatusOK, "Plan updated")
}

func (s *planServiceImpl) GetTotalCount(ctx context.Context) *schema.ServiceResponse[int64] {
	count, err := s.planRepo.Count(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to get count", err)
	}
	return schema.NewServiceResponse(count, http.StatusOK, "Count retrieved")
}

func (s *planServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}
