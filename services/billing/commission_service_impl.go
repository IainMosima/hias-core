package billing

import (
	"context"
	"net/http"
	"time"

	"github.com/bitbiz/hias-core/domains/billing/entity"
	billingRepo "github.com/bitbiz/hias-core/domains/billing/repository"
	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	"github.com/bitbiz/hias-core/domains/billing/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/google/uuid"
)

type commissionServiceImpl struct {
	ruleRepo    billingRepo.CommissionRuleRepository
	paymentRepo billingRepo.CommissionPaymentRepository
}

func NewCommissionService(ruleRepo billingRepo.CommissionRuleRepository, paymentRepo billingRepo.CommissionPaymentRepository) service.CommissionService {
	return &commissionServiceImpl{ruleRepo: ruleRepo, paymentRepo: paymentRepo}
}

func (s *commissionServiceImpl) CreateRule(ctx context.Context, req billingSchema.CreateCommissionRuleRequest, createdBy uuid.UUID) *schema.ServiceResponse[billingSchema.CommissionRuleResponse] {
	planID, _ := uuid.Parse(req.PlanID)
	intermediaryID, _ := uuid.Parse(req.IntermediaryID)

	rule := &entity.CommissionRule{
		PlanID: planID, IntermediaryID: intermediaryID,
		RatePct: req.RatePct, FlatAmount: req.FlatAmount,
		EffectiveFrom: req.EffectiveFrom, EffectiveTo: req.EffectiveTo,
		CreatedBy: createdBy,
	}

	created, err := s.ruleRepo.Create(ctx, rule)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.CommissionRuleResponse](http.StatusInternalServerError, "Failed to create commission rule", err)
	}
	return schema.NewServiceResponse(billingSchema.ToCommissionRuleResponse(created), http.StatusCreated, "Commission rule created")
}

func (s *commissionServiceImpl) ListRulesByPlan(ctx context.Context, planID uuid.UUID) *schema.ServiceResponse[[]billingSchema.CommissionRuleResponse] {
	rules, err := s.ruleRepo.ListByPlan(ctx, planID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]billingSchema.CommissionRuleResponse](http.StatusInternalServerError, "Failed to list commission rules", err)
	}
	responses := make([]billingSchema.CommissionRuleResponse, len(rules))
	for i, r := range rules {
		responses[i] = billingSchema.ToCommissionRuleResponse(r)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Commission rules retrieved")
}

func (s *commissionServiceImpl) CalculateCommission(ctx context.Context, req billingSchema.CalculateCommissionRequest) *schema.ServiceResponse[billingSchema.CalculateCommissionResponse] {
	planID, _ := uuid.Parse(req.PlanID)
	rules, err := s.ruleRepo.ListByPlan(ctx, planID)
	if err != nil || len(rules) == 0 {
		return schema.NewServiceErrorResponse[billingSchema.CalculateCommissionResponse](http.StatusNotFound, "No commission rules found", err)
	}

	// Filter rules by effective date range
	now := time.Now()
	var activeRules []*entity.CommissionRule
	for _, r := range rules {
		if now.Before(r.EffectiveFrom) {
			continue
		}
		if r.EffectiveTo != nil && now.After(*r.EffectiveTo) {
			continue
		}
		activeRules = append(activeRules, r)
	}

	if len(activeRules) == 0 {
		return schema.NewServiceErrorResponse[billingSchema.CalculateCommissionResponse](http.StatusNotFound, "No active commission rule for current date", nil)
	}

	// Use first active rule
	rule := activeRules[0]
	commission := int64(float64(req.PremiumAmount)*rule.RatePct/100) + rule.FlatAmount

	return schema.NewServiceResponse(billingSchema.CalculateCommissionResponse{
		CommissionAmount: commission, RatePct: rule.RatePct, FlatAmount: rule.FlatAmount,
		RuleID: rule.ID.String(),
	}, http.StatusOK, "Commission calculated")
}

func (s *commissionServiceImpl) ListPayments(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]billingSchema.CommissionPaymentResponse] {
	offset := (page - 1) * pageSize
	payments, err := s.paymentRepo.List(ctx, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]billingSchema.CommissionPaymentResponse](http.StatusInternalServerError, "Failed to list commission payments", err)
	}
	responses := make([]billingSchema.CommissionPaymentResponse, len(payments))
	for i, p := range payments {
		responses[i] = billingSchema.ToCommissionPaymentResponse(p)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Commission payments retrieved")
}

func (s *commissionServiceImpl) ProcessPayments(ctx context.Context) *schema.ServiceResponse[int] {
	// Process pending commission payments — mark as PROCESSED
	payments, err := s.paymentRepo.List(ctx, 100, 0)
	if err != nil {
		return schema.NewServiceErrorResponse[int](http.StatusInternalServerError, "Failed to list payments", err)
	}

	processed := 0
	for _, p := range payments {
		if p.Status == "PENDING" {
			_, err := s.paymentRepo.UpdateStatus(ctx, p.ID, "PROCESSED")
			if err == nil {
				processed++
			}
		}
	}
	return schema.NewServiceResponse(processed, http.StatusOK, "Commission payments processed")
}
