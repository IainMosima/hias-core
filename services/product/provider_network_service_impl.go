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

type providerNetworkServiceImpl struct {
	networkRepo repository.ProviderNetworkRepository
	planRepo    repository.PlanRepository
	auditSvc    auditService.AuditService
}

func NewProviderNetworkService(
	networkRepo repository.ProviderNetworkRepository,
	planRepo repository.PlanRepository,
	auditSvc auditService.AuditService,
) service.ProviderNetworkService {
	return &providerNetworkServiceImpl{
		networkRepo: networkRepo,
		planRepo:    planRepo,
		auditSvc:    auditSvc,
	}
}

func (s *providerNetworkServiceImpl) CreateProviderNetwork(ctx context.Context, planID uuid.UUID, req productSchema.CreateProviderNetworkRequest) *schema.ServiceResponse[productSchema.ProviderNetworkResponse] {
	_, err := s.planRepo.GetByID(ctx, planID)
	if err != nil {
		return schema.NewServiceErrorResponse[productSchema.ProviderNetworkResponse](http.StatusNotFound, "Plan not found", err)
	}

	providerID, err := uuid.Parse(req.ProviderID)
	if err != nil {
		return schema.NewServiceErrorResponse[productSchema.ProviderNetworkResponse](http.StatusBadRequest, "Invalid provider ID", err)
	}

	network := &entity.ProviderNetwork{
		PlanID:          planID,
		ProviderID:      providerID,
		BenefitCategory: req.BenefitCategory,
		Status:          string(shared.ProviderNetworkStatusActive),
	}

	created, err := s.networkRepo.Create(ctx, network)
	if err != nil {
		return schema.NewServiceErrorResponse[productSchema.ProviderNetworkResponse](http.StatusInternalServerError, "Failed to create provider network", err)
	}

	s.logAudit(ctx, uuid.Nil, string(shared.AuditEntityTypeProviderNetwork), created.ID, string(shared.AuditActionCreate))

	return schema.NewServiceResponse(productSchema.ToProviderNetworkResponse(created), http.StatusCreated, "Provider network created")
}

func (s *providerNetworkServiceImpl) ListProviderNetworksByPlan(ctx context.Context, planID uuid.UUID) *schema.ServiceResponse[[]productSchema.ProviderNetworkResponse] {
	networks, err := s.networkRepo.ListByPlan(ctx, planID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]productSchema.ProviderNetworkResponse](http.StatusInternalServerError, "Failed to list provider networks", err)
	}

	responses := make([]productSchema.ProviderNetworkResponse, len(networks))
	for i, n := range networks {
		responses[i] = productSchema.ToProviderNetworkResponse(n)
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Provider networks retrieved")
}

func (s *providerNetworkServiceImpl) UpdateProviderNetworkStatus(ctx context.Context, id uuid.UUID, status string) *schema.ServiceResponse[productSchema.ProviderNetworkResponse] {
	updated, err := s.networkRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		return schema.NewServiceErrorResponse[productSchema.ProviderNetworkResponse](http.StatusInternalServerError, "Failed to update provider network status", err)
	}

	return schema.NewServiceResponse(productSchema.ToProviderNetworkResponse(updated), http.StatusOK, "Provider network status updated")
}

func (s *providerNetworkServiceImpl) CheckEligibility(ctx context.Context, planID, providerID uuid.UUID, category string) *schema.ServiceResponse[bool] {
	eligible, err := s.networkRepo.CheckEligibility(ctx, planID, providerID, category)
	if err != nil {
		return schema.NewServiceErrorResponse[bool](http.StatusInternalServerError, "Failed to check eligibility", err)
	}
	return schema.NewServiceResponse(eligible, http.StatusOK, "Eligibility checked")
}

func (s *providerNetworkServiceImpl) DeleteProviderNetwork(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[string] {
	err := s.networkRepo.Delete(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusInternalServerError, "Failed to delete provider network", err)
	}
	return schema.NewServiceResponse("Provider network deleted", http.StatusOK, "Provider network deleted")
}

func (s *providerNetworkServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}
