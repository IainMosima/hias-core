package product

import (
	"context"
	"encoding/json"
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

type exclusionServiceImpl struct {
	exclusionRepo repository.ExclusionRepository
	planRepo      repository.PlanRepository
	auditSvc      auditService.AuditService
}

func NewExclusionService(exclusionRepo repository.ExclusionRepository, planRepo repository.PlanRepository, auditSvc auditService.AuditService) service.ExclusionService {
	return &exclusionServiceImpl{
		exclusionRepo: exclusionRepo,
		planRepo:      planRepo,
		auditSvc:      auditSvc,
	}
}

func (s *exclusionServiceImpl) CreateExclusion(ctx context.Context, planID uuid.UUID, req productSchema.CreateExclusionRequest) *schema.ServiceResponse[productSchema.ExclusionResponse] {
	// Validate plan exists
	_, err := s.planRepo.GetByID(ctx, planID)
	if err != nil {
		return schema.NewServiceErrorResponse[productSchema.ExclusionResponse](http.StatusNotFound, "Plan not found", err)
	}

	icdJSON, _ := json.Marshal(req.ICDCodes)

	exclusion := &entity.Exclusion{
		PlanID:      planID,
		Description: req.Description,
		Type:        req.Type,
		ICDCodes:    icdJSON,
	}

	created, err := s.exclusionRepo.Create(ctx, exclusion)
	if err != nil {
		return schema.NewServiceErrorResponse[productSchema.ExclusionResponse](http.StatusInternalServerError, "Failed to create exclusion", err)
	}

	s.logAudit(ctx, uuid.Nil, string(shared.AuditEntityTypeExclusion), created.ID, string(shared.AuditActionCreate))

	return schema.NewServiceResponse(productSchema.ToExclusionResponse(created), http.StatusCreated, "Exclusion created")
}

func (s *exclusionServiceImpl) GetExclusion(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[productSchema.ExclusionResponse] {
	exclusion, err := s.exclusionRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[productSchema.ExclusionResponse](http.StatusNotFound, "Exclusion not found", err)
	}

	return schema.NewServiceResponse(productSchema.ToExclusionResponse(exclusion), http.StatusOK, "Exclusion retrieved")
}

func (s *exclusionServiceImpl) ListExclusionsByPlan(ctx context.Context, planID uuid.UUID) *schema.ServiceResponse[[]productSchema.ExclusionResponse] {
	exclusions, err := s.exclusionRepo.ListByPlan(ctx, planID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]productSchema.ExclusionResponse](http.StatusInternalServerError, "Failed to list exclusions", err)
	}

	responses := make([]productSchema.ExclusionResponse, len(exclusions))
	for i, e := range exclusions {
		responses[i] = productSchema.ToExclusionResponse(e)
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Exclusions retrieved")
}

func (s *exclusionServiceImpl) UpdateExclusion(ctx context.Context, id uuid.UUID, req productSchema.UpdateExclusionRequest) *schema.ServiceResponse[productSchema.ExclusionResponse] {
	existing, err := s.exclusionRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[productSchema.ExclusionResponse](http.StatusNotFound, "Exclusion not found", err)
	}

	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Type != nil {
		existing.Type = *req.Type
	}
	if req.ICDCodes != nil {
		icdJSON, _ := json.Marshal(req.ICDCodes)
		existing.ICDCodes = icdJSON
	}

	updated, err := s.exclusionRepo.Update(ctx, existing)
	if err != nil {
		return schema.NewServiceErrorResponse[productSchema.ExclusionResponse](http.StatusInternalServerError, "Failed to update exclusion", err)
	}

	return schema.NewServiceResponse(productSchema.ToExclusionResponse(updated), http.StatusOK, "Exclusion updated")
}

func (s *exclusionServiceImpl) DeleteExclusion(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[string] {
	err := s.exclusionRepo.Delete(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusInternalServerError, "Failed to delete exclusion", err)
	}

	return schema.NewServiceResponse("Exclusion deleted", http.StatusOK, "Exclusion deleted")
}

func (s *exclusionServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}
