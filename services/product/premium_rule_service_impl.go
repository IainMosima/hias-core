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

type premiumRuleServiceImpl struct {
	ruleRepo repository.PremiumRuleRepository
	planRepo repository.PlanRepository
	auditSvc auditService.AuditService
}

func NewPremiumRuleService(
	ruleRepo repository.PremiumRuleRepository,
	planRepo repository.PlanRepository,
	auditSvc auditService.AuditService,
) service.PremiumRuleService {
	return &premiumRuleServiceImpl{
		ruleRepo: ruleRepo,
		planRepo: planRepo,
		auditSvc: auditSvc,
	}
}

func (s *premiumRuleServiceImpl) CreatePremiumRule(ctx context.Context, planID uuid.UUID, req productSchema.CreatePremiumRuleRequest) *schema.ServiceResponse[productSchema.PremiumRuleResponse] {
	_, err := s.planRepo.GetByID(ctx, planID)
	if err != nil {
		return schema.NewServiceErrorResponse[productSchema.PremiumRuleResponse](http.StatusNotFound, "Plan not found", err)
	}

	rule := &entity.PremiumRule{
		PlanID:          planID,
		CalculationType: req.CalculationType,
		Relationship:    req.Relationship,
		RateAmount:      req.RateAmount,
		DiscountType:    req.DiscountType,
		DiscountValue:   req.DiscountValue,
		MinMembers:      req.MinMembers,
	}

	created, err := s.ruleRepo.Create(ctx, rule)
	if err != nil {
		return schema.NewServiceErrorResponse[productSchema.PremiumRuleResponse](http.StatusInternalServerError, "Failed to create premium rule", err)
	}

	s.logAudit(ctx, uuid.Nil, string(shared.AuditEntityTypePremiumRule), created.ID, string(shared.AuditActionCreate))

	return schema.NewServiceResponse(productSchema.ToPremiumRuleResponse(created), http.StatusCreated, "Premium rule created")
}

func (s *premiumRuleServiceImpl) ListPremiumRulesByPlan(ctx context.Context, planID uuid.UUID) *schema.ServiceResponse[[]productSchema.PremiumRuleResponse] {
	rules, err := s.ruleRepo.ListByPlan(ctx, planID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]productSchema.PremiumRuleResponse](http.StatusInternalServerError, "Failed to list premium rules", err)
	}

	responses := make([]productSchema.PremiumRuleResponse, len(rules))
	for i, r := range rules {
		responses[i] = productSchema.ToPremiumRuleResponse(r)
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Premium rules retrieved")
}

func (s *premiumRuleServiceImpl) DeletePremiumRule(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[string] {
	err := s.ruleRepo.Delete(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusInternalServerError, "Failed to delete premium rule", err)
	}
	return schema.NewServiceResponse("Premium rule deleted", http.StatusOK, "Premium rule deleted")
}

func (s *premiumRuleServiceImpl) CalculatePremium(ctx context.Context, planID uuid.UUID, memberCount int, relationships []string) *schema.ServiceResponse[int64] {
	plan, err := s.planRepo.GetByID(ctx, planID)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusNotFound, "Plan not found", err)
	}

	rules, err := s.ruleRepo.ListByPlan(ctx, planID)
	if err != nil || len(rules) == 0 {
		// Fallback to base premium
		return schema.NewServiceResponse(plan.BasePremium, http.StatusOK, "Premium calculated (base)")
	}

	var totalPremium int64

	for _, rel := range relationships {
		// Find matching rule for relationship
		var matchedRule *entity.PremiumRule
		for _, r := range rules {
			if r.Relationship == rel {
				matchedRule = r
				break
			}
		}
		if matchedRule == nil {
			// Use first rule without relationship filter
			for _, r := range rules {
				if r.Relationship == "" {
					matchedRule = r
					break
				}
			}
		}

		if matchedRule != nil {
			totalPremium += matchedRule.RateAmount
		} else {
			totalPremium += plan.BasePremium
		}
	}

	// Apply group discount if applicable
	for _, r := range rules {
		if r.DiscountType != "" && r.MinMembers > 0 && memberCount >= r.MinMembers {
			if r.DiscountType == string(shared.DiscountTypePercentage) {
				totalPremium -= totalPremium * r.DiscountValue / 100
			} else if r.DiscountType == string(shared.DiscountTypeFixed) {
				totalPremium -= r.DiscountValue
			}
			break
		}
	}

	if totalPremium < 0 {
		totalPremium = 0
	}

	return schema.NewServiceResponse(totalPremium, http.StatusOK, "Premium calculated")
}

func (s *premiumRuleServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}
