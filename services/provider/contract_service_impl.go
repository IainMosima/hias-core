package provider

import (
	"context"
	"fmt"
	"log"
	"net/http"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/domains/provider/entity"
	"github.com/bitbiz/hias-core/domains/provider/repository"
	providerSchema "github.com/bitbiz/hias-core/domains/provider/schema"
	"github.com/bitbiz/hias-core/domains/provider/service"
	"github.com/bitbiz/hias-core/shared"
	"github.com/google/uuid"
)

type contractServiceImpl struct {
	contractRepo repository.ContractRepository
	auditSvc     auditService.AuditService
}

func NewContractService(contractRepo repository.ContractRepository, auditSvc auditService.AuditService) service.ContractService {
	return &contractServiceImpl{
		contractRepo: contractRepo,
		auditSvc:     auditSvc,
	}
}

func (s *contractServiceImpl) CreateContract(ctx context.Context, providerID uuid.UUID, req providerSchema.CreateContractRequest) *schema.ServiceResponse[providerSchema.ContractResponse] {
	contract := &entity.Contract{
		ProviderID: providerID,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		Terms:      req.Terms,
		Status:     string(shared.ContractStatusActive),
	}

	created, err := s.contractRepo.Create(ctx, contract)
	if err != nil {
		return schema.NewServiceErrorResponse[providerSchema.ContractResponse](
			http.StatusInternalServerError, fmt.Sprintf("failed to create contract: %v", err), err,
		)
	}

	s.logAudit(ctx, uuid.Nil, string(shared.AuditEntityTypeContract), created.ID, "CREATE")

	resp := providerSchema.ToContractResponse(created)
	return schema.NewServiceResponse(resp, http.StatusCreated, "Contract created")
}

func (s *contractServiceImpl) ListContracts(ctx context.Context, providerID uuid.UUID) *schema.ServiceResponse[[]providerSchema.ContractResponse] {
	contracts, err := s.contractRepo.ListByProvider(ctx, providerID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]providerSchema.ContractResponse](
			http.StatusInternalServerError, fmt.Sprintf("failed to list contracts: %v", err), err,
		)
	}

	responses := make([]providerSchema.ContractResponse, len(contracts))
	for i, c := range contracts {
		responses[i] = providerSchema.ToContractResponse(c)
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Contracts retrieved")
}

func (s *contractServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}
