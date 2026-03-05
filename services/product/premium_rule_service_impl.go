package product

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"time"

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

	maxAge := req.MaxAge
	if maxAge == 0 {
		maxAge = 150
	}

	rule := &entity.PremiumRule{
		PlanID:          planID,
		CalculationType: req.CalculationType,
		Relationship:    req.Relationship,
		RateAmount:      req.RateAmount,
		DiscountType:    req.DiscountType,
		DiscountValue:   req.DiscountValue,
		MinMembers:      req.MinMembers,
		MinAge:          req.MinAge,
		MaxAge:          maxAge,
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
		return schema.NewServiceResponse(plan.BasePremium, http.StatusOK, "Premium calculated (base)")
	}

	// Check if any rule uses per_family calculation
	for _, r := range rules {
		if r.CalculationType == string(shared.PremiumCalculationTypePerFamily) {
			// per_family: find the best matching rule for the family size
			matchedRule := findFamilyRule(rules, memberCount)
			if matchedRule != nil {
				return schema.NewServiceResponse(matchedRule.RateAmount, http.StatusOK, "Premium calculated (per_family)")
			}
			return schema.NewServiceResponse(plan.BasePremium, http.StatusOK, "Premium calculated (base)")
		}
	}

	var totalPremium int64

	for _, rel := range relationships {
		matchedRule := findMatchingRule(rules, rel, 0)
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
				totalPremium -= totalPremium * r.DiscountValue / 10000
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

// CalculatePremiumWithMembers calculates premium using proposed_members JSON for age-band matching
func (s *premiumRuleServiceImpl) CalculatePremiumWithMembers(ctx context.Context, planID uuid.UUID, memberCount int, proposedMembers json.RawMessage) *schema.ServiceResponse[int64] {
	plan, err := s.planRepo.GetByID(ctx, planID)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusNotFound, "Plan not found", err)
	}

	rules, err := s.ruleRepo.ListByPlan(ctx, planID)
	if err != nil || len(rules) == 0 {
		return schema.NewServiceResponse(plan.BasePremium, http.StatusOK, "Premium calculated (base)")
	}

	// Check if any rule uses per_family calculation
	for _, r := range rules {
		if r.CalculationType == string(shared.PremiumCalculationTypePerFamily) {
			matchedRule := findFamilyRule(rules, memberCount)
			if matchedRule != nil {
				return schema.NewServiceResponse(matchedRule.RateAmount, http.StatusOK, "Premium calculated (per_family)")
			}
			return schema.NewServiceResponse(plan.BasePremium, http.StatusOK, "Premium calculated (base)")
		}
	}

	// Parse members for age-band matching
	var members []struct {
		Relationship string `json:"relationship"`
		DateOfBirth  string `json:"date_of_birth"`
	}
	hasAgeData := false
	if err := json.Unmarshal(proposedMembers, &members); err == nil && len(members) > 0 {
		for _, m := range members {
			if m.DateOfBirth != "" {
				hasAgeData = true
				break
			}
		}
	}

	var totalPremium int64

	if hasAgeData {
		for _, m := range members {
			age := calculateAge(m.DateOfBirth)
			matchedRule := findMatchingRule(rules, m.Relationship, age)
			if matchedRule != nil {
				totalPremium += matchedRule.RateAmount
			} else {
				totalPremium += plan.BasePremium
			}
		}
	} else {
		// Fallback: extract relationships only
		for _, m := range members {
			matchedRule := findMatchingRule(rules, m.Relationship, 0)
			if matchedRule != nil {
				totalPremium += matchedRule.RateAmount
			} else {
				totalPremium += plan.BasePremium
			}
		}
		// If no members parsed, use memberCount with base premium
		if len(members) == 0 {
			totalPremium = plan.BasePremium * int64(memberCount)
		}
	}

	// Apply group discount if applicable
	for _, r := range rules {
		if r.DiscountType != "" && r.MinMembers > 0 && memberCount >= r.MinMembers {
			if r.DiscountType == string(shared.DiscountTypePercentage) {
				totalPremium -= totalPremium * r.DiscountValue / 10000
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

// findMatchingRule finds the best matching rule for a relationship and age.
// Priority: relationship+age match > relationship-only match > generic (no relationship) match.
func findMatchingRule(rules []*entity.PremiumRule, relationship string, age int) *entity.PremiumRule {
	var relAgeMatch, relOnlyMatch, genericMatch *entity.PremiumRule

	for _, r := range rules {
		matchesRel := r.Relationship == relationship
		matchesAge := age >= r.MinAge && age <= r.MaxAge
		isGeneric := r.Relationship == ""

		if matchesRel && age > 0 && matchesAge {
			relAgeMatch = r
			break
		}
		if matchesRel && relOnlyMatch == nil {
			relOnlyMatch = r
		}
		if isGeneric && age > 0 && matchesAge && genericMatch == nil {
			genericMatch = r
		}
		if isGeneric && genericMatch == nil {
			genericMatch = r
		}
	}

	if relAgeMatch != nil {
		return relAgeMatch
	}
	if relOnlyMatch != nil {
		return relOnlyMatch
	}
	return genericMatch
}

// findFamilyRule finds the per_family rule best matching the given family size
func findFamilyRule(rules []*entity.PremiumRule, memberCount int) *entity.PremiumRule {
	var best *entity.PremiumRule
	for _, r := range rules {
		if r.CalculationType != string(shared.PremiumCalculationTypePerFamily) {
			continue
		}
		if memberCount >= r.MinMembers {
			if best == nil || r.MinMembers > best.MinMembers {
				best = r
			}
		}
	}
	if best == nil {
		// Fallback to any per_family rule
		for _, r := range rules {
			if r.CalculationType == string(shared.PremiumCalculationTypePerFamily) {
				return r
			}
		}
	}
	return best
}

// calculateAge computes age from a date_of_birth string (YYYY-MM-DD format)
func calculateAge(dob string) int {
	t, err := time.Parse("2006-01-02", dob)
	if err != nil {
		return 0
	}
	now := time.Now()
	age := now.Year() - t.Year()
	if now.YearDay() < t.YearDay() {
		age--
	}
	return int(math.Max(0, float64(age)))
}

func (s *premiumRuleServiceImpl) GetRateSheet(ctx context.Context, planID uuid.UUID) *schema.ServiceResponse[[]productSchema.PremiumRuleResponse] {
	_, err := s.planRepo.GetByID(ctx, planID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]productSchema.PremiumRuleResponse](http.StatusNotFound, "Plan not found", err)
	}

	rules, err := s.ruleRepo.ListByPlan(ctx, planID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]productSchema.PremiumRuleResponse](http.StatusInternalServerError, "Failed to get rate sheet", err)
	}

	responses := make([]productSchema.PremiumRuleResponse, len(rules))
	for i, r := range rules {
		responses[i] = productSchema.ToPremiumRuleResponse(r)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Rate sheet retrieved")
}

func (s *premiumRuleServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}
