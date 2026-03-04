package sales

import (
	"context"
	"log"
	"net/http"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/domains/sales/entity"
	salesRepo "github.com/bitbiz/hias-core/domains/sales/repository"
	salesSchema "github.com/bitbiz/hias-core/domains/sales/schema"
	"github.com/bitbiz/hias-core/domains/sales/service"
	"github.com/bitbiz/hias-core/shared"
	"github.com/google/uuid"
)

type approvalLimitServiceImpl struct {
	approvalRepo salesRepo.ApprovalLimitRepository
	auditSvc     auditService.AuditService
}

func NewApprovalLimitService(
	approvalRepo salesRepo.ApprovalLimitRepository,
	auditSvc auditService.AuditService,
) service.ApprovalLimitService {
	return &approvalLimitServiceImpl{
		approvalRepo: approvalRepo,
		auditSvc:     auditSvc,
	}
}

func (s *approvalLimitServiceImpl) GetLimits(ctx context.Context) *schema.ServiceResponse[[]salesSchema.ApprovalLimitResponse] {
	limits, err := s.approvalRepo.List(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[[]salesSchema.ApprovalLimitResponse](http.StatusInternalServerError, "Failed to list approval limits", err)
	}
	responses := make([]salesSchema.ApprovalLimitResponse, len(limits))
	for i, l := range limits {
		responses[i] = salesSchema.ToApprovalLimitResponse(l)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Approval limits retrieved")
}

func (s *approvalLimitServiceImpl) GetLimitByRole(ctx context.Context, roleName string) *schema.ServiceResponse[salesSchema.ApprovalLimitResponse] {
	limit, err := s.approvalRepo.GetByRole(ctx, roleName)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.ApprovalLimitResponse](http.StatusNotFound, "Approval limit not found", err)
	}
	return schema.NewServiceResponse(salesSchema.ToApprovalLimitResponse(limit), http.StatusOK, "Approval limit retrieved")
}

func (s *approvalLimitServiceImpl) CreateLimit(ctx context.Context, req salesSchema.CreateApprovalLimitRequest, createdBy uuid.UUID) *schema.ServiceResponse[salesSchema.ApprovalLimitResponse] {
	limit := &entity.ApprovalLimit{
		RoleName:              req.RoleName,
		MaxDiscountPercentage: req.MaxDiscountPercentage,
		MaxDiscountAmount:     req.MaxDiscountAmount,
		MaxLoadingPercentage:  req.MaxLoadingPercentage,
		MaxLoadingAmount:      req.MaxLoadingAmount,
		EscalationRole:        req.EscalationRole,
		IsActive:              true,
	}

	created, err := s.approvalRepo.Create(ctx, limit)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.ApprovalLimitResponse](http.StatusInternalServerError, "Failed to create approval limit", err)
	}

	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, createdBy, string(shared.AuditEntityTypeApprovalLimit), created.ID, string(shared.AuditActionCreate), nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}

	return schema.NewServiceResponse(salesSchema.ToApprovalLimitResponse(created), http.StatusCreated, "Approval limit created")
}

func (s *approvalLimitServiceImpl) UpdateLimit(ctx context.Context, id uuid.UUID, req salesSchema.UpdateApprovalLimitRequest, updatedBy uuid.UUID) *schema.ServiceResponse[salesSchema.ApprovalLimitResponse] {
	// Get existing limits to find the one we need
	limits, err := s.approvalRepo.List(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.ApprovalLimitResponse](http.StatusInternalServerError, "Failed to get limits", err)
	}

	var existing *entity.ApprovalLimit
	for _, l := range limits {
		if l.ID == id {
			existing = l
			break
		}
	}
	if existing == nil {
		return schema.NewServiceErrorResponse[salesSchema.ApprovalLimitResponse](http.StatusNotFound, "Approval limit not found", nil)
	}

	// Apply updates
	if req.MaxDiscountPercentage != nil {
		existing.MaxDiscountPercentage = *req.MaxDiscountPercentage
	}
	if req.MaxDiscountAmount != nil {
		existing.MaxDiscountAmount = *req.MaxDiscountAmount
	}
	if req.MaxLoadingPercentage != nil {
		existing.MaxLoadingPercentage = *req.MaxLoadingPercentage
	}
	if req.MaxLoadingAmount != nil {
		existing.MaxLoadingAmount = *req.MaxLoadingAmount
	}
	if req.EscalationRole != "" {
		existing.EscalationRole = req.EscalationRole
	}

	result, err := s.approvalRepo.Update(ctx, existing)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.ApprovalLimitResponse](http.StatusInternalServerError, "Failed to update approval limit", err)
	}

	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, updatedBy, string(shared.AuditEntityTypeApprovalLimit), id, string(shared.AuditActionUpdate), nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}

	return schema.NewServiceResponse(salesSchema.ToApprovalLimitResponse(result), http.StatusOK, "Approval limit updated")
}
