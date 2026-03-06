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

type adjudicationRuleServiceImpl struct {
	ruleRepo claimRepo.AdjudicationRuleRepository
}

func NewAdjudicationRuleService(ruleRepo claimRepo.AdjudicationRuleRepository) service.AdjudicationRuleService {
	return &adjudicationRuleServiceImpl{ruleRepo: ruleRepo}
}

func (s *adjudicationRuleServiceImpl) CreateRule(ctx context.Context, req claimsSchema.CreateAdjudicationRuleRequest) *schema.ServiceResponse[claimsSchema.AdjudicationRuleResponse] {
	rule := &entity.AdjudicationRule{
		Name:       req.Name,
		RuleType:   req.RuleType,
		Parameters: req.Parameters,
		Priority:   req.Priority,
		IsActive:   req.IsActive,
	}

	created, err := s.ruleRepo.Create(ctx, rule)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.AdjudicationRuleResponse](http.StatusInternalServerError, "Failed to create adjudication rule", err)
	}
	return schema.NewServiceResponse(claimsSchema.ToAdjudicationRuleResponse(created), http.StatusCreated, "Adjudication rule created")
}

func (s *adjudicationRuleServiceImpl) GetRule(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[claimsSchema.AdjudicationRuleResponse] {
	rule, err := s.ruleRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.AdjudicationRuleResponse](http.StatusNotFound, "Adjudication rule not found", err)
	}
	return schema.NewServiceResponse(claimsSchema.ToAdjudicationRuleResponse(rule), http.StatusOK, "Adjudication rule retrieved")
}

func (s *adjudicationRuleServiceImpl) ListRules(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]claimsSchema.AdjudicationRuleResponse] {
	offset := (page - 1) * pageSize
	rules, err := s.ruleRepo.List(ctx, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]claimsSchema.AdjudicationRuleResponse](http.StatusInternalServerError, "Failed to list adjudication rules", err)
	}
	responses := make([]claimsSchema.AdjudicationRuleResponse, len(rules))
	for i, r := range rules {
		responses[i] = claimsSchema.ToAdjudicationRuleResponse(r)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Adjudication rules retrieved")
}

func (s *adjudicationRuleServiceImpl) UpdateRule(ctx context.Context, id uuid.UUID, req claimsSchema.UpdateAdjudicationRuleRequest) *schema.ServiceResponse[claimsSchema.AdjudicationRuleResponse] {
	existing, err := s.ruleRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.AdjudicationRuleResponse](http.StatusNotFound, "Adjudication rule not found", err)
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.RuleType != "" {
		existing.RuleType = req.RuleType
	}
	if req.Parameters != nil {
		existing.Parameters = req.Parameters
	}
	if req.Priority != 0 {
		existing.Priority = req.Priority
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	updated, err := s.ruleRepo.Update(ctx, existing)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.AdjudicationRuleResponse](http.StatusInternalServerError, "Failed to update adjudication rule", err)
	}
	return schema.NewServiceResponse(claimsSchema.ToAdjudicationRuleResponse(updated), http.StatusOK, "Adjudication rule updated")
}

func (s *adjudicationRuleServiceImpl) DeleteRule(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[string] {
	if err := s.ruleRepo.Delete(ctx, id); err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusInternalServerError, "Failed to delete adjudication rule", err)
	}
	return schema.NewServiceResponse("deleted", http.StatusOK, "Adjudication rule deleted")
}
