package policy

import (
	"context"
	"log"
	"net/http"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/domains/policy/entity"
	"github.com/bitbiz/hias-core/domains/policy/repository"
	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/bitbiz/hias-core/domains/policy/service"
	productRepo "github.com/bitbiz/hias-core/domains/product/repository"
	"github.com/bitbiz/hias-core/shared"
	"github.com/google/uuid"
)

type underwritingRuleServiceImpl struct {
	ruleRepo repository.UnderwritingRuleRepository
	planRepo productRepo.PlanRepository
	auditSvc auditService.AuditService
}

func NewUnderwritingRuleService(
	ruleRepo repository.UnderwritingRuleRepository,
	planRepo productRepo.PlanRepository,
	auditSvc auditService.AuditService,
) service.UnderwritingRuleService {
	return &underwritingRuleServiceImpl{
		ruleRepo: ruleRepo,
		planRepo: planRepo,
		auditSvc: auditSvc,
	}
}

func (s *underwritingRuleServiceImpl) CreateRule(ctx context.Context, req policySchema.CreateUnderwritingRuleRequest, createdBy uuid.UUID) *schema.ServiceResponse[policySchema.UnderwritingRuleResponse] {
	planID, err := uuid.Parse(req.PlanID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.UnderwritingRuleResponse](http.StatusBadRequest, "Invalid plan ID", err)
	}

	if _, err := s.planRepo.GetByID(ctx, planID); err != nil {
		return schema.NewServiceErrorResponse[policySchema.UnderwritingRuleResponse](http.StatusNotFound, "Plan not found", err)
	}

	severity := req.Severity
	if severity == "" {
		severity = "MEDIUM"
	}
	riskWeight := req.RiskScoreWeight
	if riskWeight == 0 {
		riskWeight = 20
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	rule := &entity.UnderwritingRule{
		PlanID:          planID,
		RuleType:        req.RuleType,
		Relationship:    req.Relationship,
		ParameterKey:    req.ParameterKey,
		ParameterValue:  req.ParameterValue,
		Severity:        severity,
		RiskScoreWeight: riskWeight,
		IsBlocking:      req.IsBlocking,
		IsActive:        isActive,
		Description:     req.Description,
	}

	created, err := s.ruleRepo.Create(ctx, rule)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.UnderwritingRuleResponse](http.StatusInternalServerError, "Failed to create rule", err)
	}

	s.logAudit(ctx, createdBy, string(shared.AuditEntityTypeUnderwritingRule), created.ID, string(shared.AuditActionCreate))

	return schema.NewServiceResponse(policySchema.ToUnderwritingRuleResponse(created), http.StatusCreated, "Underwriting rule created")
}

func (s *underwritingRuleServiceImpl) ListByPlan(ctx context.Context, planID uuid.UUID) *schema.ServiceResponse[[]policySchema.UnderwritingRuleResponse] {
	rules, err := s.ruleRepo.ListByPlan(ctx, planID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]policySchema.UnderwritingRuleResponse](http.StatusInternalServerError, "Failed to list rules", err)
	}
	responses := make([]policySchema.UnderwritingRuleResponse, len(rules))
	for i, r := range rules {
		responses[i] = policySchema.ToUnderwritingRuleResponse(r)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Rules retrieved")
}

func (s *underwritingRuleServiceImpl) UpdateRule(ctx context.Context, id uuid.UUID, req policySchema.UpdateUnderwritingRuleRequest) *schema.ServiceResponse[policySchema.UnderwritingRuleResponse] {
	existing, err := s.ruleRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.UnderwritingRuleResponse](http.StatusNotFound, "Rule not found", err)
	}

	if req.RuleType != nil {
		existing.RuleType = *req.RuleType
	}
	if req.Relationship != nil {
		existing.Relationship = *req.Relationship
	}
	if req.ParameterKey != nil {
		existing.ParameterKey = *req.ParameterKey
	}
	if req.ParameterValue != nil {
		existing.ParameterValue = *req.ParameterValue
	}
	if req.Severity != nil {
		existing.Severity = *req.Severity
	}
	if req.RiskScoreWeight != nil {
		existing.RiskScoreWeight = *req.RiskScoreWeight
	}
	if req.IsBlocking != nil {
		existing.IsBlocking = *req.IsBlocking
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}

	updated, err := s.ruleRepo.Update(ctx, existing)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.UnderwritingRuleResponse](http.StatusInternalServerError, "Failed to update rule", err)
	}

	s.logAudit(ctx, uuid.Nil, string(shared.AuditEntityTypeUnderwritingRule), id, string(shared.AuditActionUpdate))

	return schema.NewServiceResponse(policySchema.ToUnderwritingRuleResponse(updated), http.StatusOK, "Rule updated")
}

func (s *underwritingRuleServiceImpl) DeleteRule(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[string] {
	if _, err := s.ruleRepo.GetByID(ctx, id); err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusNotFound, "Rule not found", err)
	}
	if err := s.ruleRepo.Delete(ctx, id); err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusInternalServerError, "Failed to delete rule", err)
	}
	return schema.NewServiceResponse("Rule deleted", http.StatusOK, "Rule deleted")
}

func (s *underwritingRuleServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}
