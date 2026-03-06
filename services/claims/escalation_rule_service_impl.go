package claims

import (
	"context"
	"net/http"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	claimRepo "github.com/bitbiz/hias-core/domains/claims/repository"
	claimsSchema "github.com/bitbiz/hias-core/domains/claims/schema"
	"github.com/bitbiz/hias-core/domains/claims/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/google/uuid"
)

type escalationRuleServiceImpl struct {
	ruleRepo claimRepo.EscalationRuleRepository
}

func NewEscalationRuleService(ruleRepo claimRepo.EscalationRuleRepository) service.EscalationRuleService {
	return &escalationRuleServiceImpl{ruleRepo: ruleRepo}
}

func (s *escalationRuleServiceImpl) CreateRule(ctx context.Context, req claimsSchema.CreateEscalationRuleRequest) *schema.ServiceResponse[claimsSchema.EscalationRuleResponse] {
	rule := &entity.EscalationRule{
		Name:            req.Name,
		ConditionType:   req.ConditionType,
		ThresholdAmount: req.ThresholdAmount,
		EscalationRole:  req.EscalationRole,
		IsActive:        req.IsActive,
	}

	created, err := s.ruleRepo.Create(ctx, rule)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.EscalationRuleResponse](http.StatusInternalServerError, "Failed to create escalation rule", err)
	}
	return schema.NewServiceResponse(claimsSchema.ToEscalationRuleResponse(created), http.StatusCreated, "Escalation rule created")
}

func (s *escalationRuleServiceImpl) GetRule(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[claimsSchema.EscalationRuleResponse] {
	rule, err := s.ruleRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.EscalationRuleResponse](http.StatusNotFound, "Escalation rule not found", err)
	}
	return schema.NewServiceResponse(claimsSchema.ToEscalationRuleResponse(rule), http.StatusOK, "Escalation rule retrieved")
}

func (s *escalationRuleServiceImpl) ListRules(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]claimsSchema.EscalationRuleResponse] {
	offset := (page - 1) * pageSize
	rules, err := s.ruleRepo.List(ctx, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]claimsSchema.EscalationRuleResponse](http.StatusInternalServerError, "Failed to list escalation rules", err)
	}
	responses := make([]claimsSchema.EscalationRuleResponse, len(rules))
	for i, r := range rules {
		responses[i] = claimsSchema.ToEscalationRuleResponse(r)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Escalation rules retrieved")
}

func (s *escalationRuleServiceImpl) UpdateRule(ctx context.Context, id uuid.UUID, req claimsSchema.UpdateEscalationRuleRequest) *schema.ServiceResponse[claimsSchema.EscalationRuleResponse] {
	existing, err := s.ruleRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.EscalationRuleResponse](http.StatusNotFound, "Escalation rule not found", err)
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.ConditionType != "" {
		existing.ConditionType = req.ConditionType
	}
	if req.ThresholdAmount != 0 {
		existing.ThresholdAmount = req.ThresholdAmount
	}
	if req.EscalationRole != "" {
		existing.EscalationRole = req.EscalationRole
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	updated, err := s.ruleRepo.Update(ctx, existing)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.EscalationRuleResponse](http.StatusInternalServerError, "Failed to update escalation rule", err)
	}
	return schema.NewServiceResponse(claimsSchema.ToEscalationRuleResponse(updated), http.StatusOK, "Escalation rule updated")
}

func (s *escalationRuleServiceImpl) DeleteRule(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[string] {
	if err := s.ruleRepo.Delete(ctx, id); err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusInternalServerError, "Failed to delete escalation rule", err)
	}
	return schema.NewServiceResponse("deleted", http.StatusOK, "Escalation rule deleted")
}
