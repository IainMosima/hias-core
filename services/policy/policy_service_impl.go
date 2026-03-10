package policy

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	billingService "github.com/bitbiz/hias-core/domains/billing/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/domains/policy/entity"
	"github.com/bitbiz/hias-core/domains/policy/repository"
	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/bitbiz/hias-core/domains/policy/service"
	planRepo "github.com/bitbiz/hias-core/domains/product/repository"
	productService "github.com/bitbiz/hias-core/domains/product/service"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/google/uuid"
)

type policyServiceImpl struct {
	policyRepo     repository.PolicyRepository
	planRepo       planRepo.PlanRepository
	memberRepo     repository.MemberRepository
	premiumRuleSvc productService.PremiumRuleService
	policyDocSvc   service.PolicyDocumentService
	creditNoteSvc  billingService.CreditNoteService
	auditSvc       auditService.AuditService
}

func NewPolicyService(
	policyRepo repository.PolicyRepository,
	planRepo planRepo.PlanRepository,
	memberRepo repository.MemberRepository,
	premiumRuleSvc productService.PremiumRuleService,
	policyDocSvc service.PolicyDocumentService,
	creditNoteSvc billingService.CreditNoteService,
	auditSvc auditService.AuditService,
) service.PolicyService {
	return &policyServiceImpl{
		policyRepo:     policyRepo,
		planRepo:       planRepo,
		memberRepo:     memberRepo,
		premiumRuleSvc: premiumRuleSvc,
		policyDocSvc:   policyDocSvc,
		creditNoteSvc:  creditNoteSvc,
		auditSvc:       auditSvc,
	}
}

func (s *policyServiceImpl) CreatePolicy(ctx context.Context, req policySchema.CreatePolicyRequest, createdBy uuid.UUID) *schema.ServiceResponse[policySchema.PolicyResponse] {
	planID, err := uuid.Parse(req.PlanID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusBadRequest, "Invalid plan ID", err)
	}

	plan, err := s.planRepo.GetByID(ctx, planID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusNotFound, "Plan not found", err)
	}

	startDate := req.StartDate
	if startDate.IsZero() {
		startDate = time.Now()
	}
	endDate := req.EndDate
	if endDate.IsZero() {
		endDate = startDate.AddDate(1, 0, 0)
	}

	policy := &entity.Policy{
		PlanID:            planID,
		PolicyholderName:  req.PolicyholderName,
		PolicyholderEmail: req.PolicyholderEmail,
		PolicyholderPhone: req.PolicyholderPhone,
		PolicyNumber:      utils.GeneratePolicyNumber(),
		Status:            string(shared.PolicyStatusDraft),
		StartDate:         startDate,
		EndDate:           endDate,
		PremiumAmount:     plan.BasePremium,
		Currency:          plan.Currency,
		CreatedBy:         createdBy,
	}

	created, err := s.policyRepo.Create(ctx, policy)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusInternalServerError, "Failed to create policy", err)
	}

	s.logAudit(ctx, createdBy, string(shared.AuditEntityTypePolicy), created.ID, string(shared.AuditActionCreate))

	return schema.NewServiceResponse(policySchema.ToPolicyResponse(created), http.StatusCreated, "Policy created")
}

func (s *policyServiceImpl) GetPolicy(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[policySchema.PolicyResponse] {
	policy, err := s.policyRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusNotFound, "Policy not found", err)
	}
	return schema.NewServiceResponse(policySchema.ToPolicyResponse(policy), http.StatusOK, "Policy retrieved")
}

func (s *policyServiceImpl) ListPolicies(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]policySchema.PolicyResponse] {
	offset := (page - 1) * pageSize
	policies, err := s.policyRepo.List(ctx, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]policySchema.PolicyResponse](http.StatusInternalServerError, "Failed to list policies", err)
	}

	responses := make([]policySchema.PolicyResponse, len(policies))
	for i, p := range policies {
		responses[i] = policySchema.ToPolicyResponse(p)
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Policies retrieved")
}

func (s *policyServiceImpl) ListPoliciesByStatus(ctx context.Context, status string, page, pageSize int) *schema.ServiceResponse[[]policySchema.PolicyResponse] {
	offset := (page - 1) * pageSize
	policies, err := s.policyRepo.ListByStatus(ctx, status, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]policySchema.PolicyResponse](http.StatusInternalServerError, "Failed to list policies by status", err)
	}

	responses := make([]policySchema.PolicyResponse, len(policies))
	for i, p := range policies {
		responses[i] = policySchema.ToPolicyResponse(p)
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Policies retrieved")
}

func (s *policyServiceImpl) ActivatePolicy(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[policySchema.PolicyResponse] {
	policy, err := s.policyRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusNotFound, "Policy not found", err)
	}

	if policy.Status != string(shared.PolicyStatusDraft) {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusBadRequest, fmt.Sprintf("Cannot activate policy in %s status", policy.Status), nil)
	}

	updated, err := s.policyRepo.ActivateWithTimestamp(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusInternalServerError, "Failed to activate policy", err)
	}

	// Activate all pending members linked to this policy
	if err := s.memberRepo.ActivatePendingByPolicy(ctx, id); err != nil {
		log.Printf("Failed to activate pending members for policy %s: %v", id, err)
	}

	// Auto-generate welcome letter and member cards (non-blocking)
	if s.policyDocSvc != nil {
		go func() {
			bgCtx := context.Background()
			s.policyDocSvc.GenerateWelcomeLetter(bgCtx, id, uuid.Nil)
			s.policyDocSvc.BulkGenerateMemberCards(bgCtx, id, uuid.Nil)
		}()
	}

	return schema.NewServiceResponse(policySchema.ToPolicyResponse(updated), http.StatusOK, "Policy activated")
}

func (s *policyServiceImpl) LapsePolicy(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[policySchema.PolicyResponse] {
	policy, err := s.policyRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusNotFound, "Policy not found", err)
	}

	if policy.Status != string(shared.PolicyStatusActive) {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusBadRequest, fmt.Sprintf("Cannot lapse policy in %s status; must be ACTIVE", policy.Status), nil)
	}

	updated, err := s.policyRepo.UpdateStatus(ctx, id, string(shared.PolicyStatusLapsed))
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusInternalServerError, "Failed to lapse policy", err)
	}
	return schema.NewServiceResponse(policySchema.ToPolicyResponse(updated), http.StatusOK, "Policy lapsed")
}

func (s *policyServiceImpl) TerminatePolicy(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[policySchema.PolicyResponse] {
	policy, err := s.policyRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusNotFound, "Policy not found", err)
	}

	if policy.Status != string(shared.PolicyStatusActive) && policy.Status != string(shared.PolicyStatusLapsed) {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusBadRequest, fmt.Sprintf("Cannot terminate policy in %s status; must be ACTIVE or LAPSED", policy.Status), nil)
	}

	updated, err := s.policyRepo.UpdateStatus(ctx, id, string(shared.PolicyStatusTerminated))
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusInternalServerError, "Failed to terminate policy", err)
	}
	return schema.NewServiceResponse(policySchema.ToPolicyResponse(updated), http.StatusOK, "Policy terminated")
}

func (s *policyServiceImpl) ReinstatePolicy(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[policySchema.PolicyResponse] {
	policy, err := s.policyRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusNotFound, "Policy not found", err)
	}

	if policy.Status != string(shared.PolicyStatusLapsed) && policy.Status != string(shared.PolicyStatusSuspended) {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusBadRequest, "Only lapsed or suspended policies can be reinstated", nil)
	}

	updated, err := s.policyRepo.UpdateStatus(ctx, id, string(shared.PolicyStatusActive))
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusInternalServerError, "Failed to reinstate policy", err)
	}

	return schema.NewServiceResponse(policySchema.ToPolicyResponse(updated), http.StatusOK, "Policy reinstated")
}

func (s *policyServiceImpl) SuspendPolicy(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[policySchema.PolicyResponse] {
	policy, err := s.policyRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusNotFound, "Policy not found", err)
	}

	if policy.Status != string(shared.PolicyStatusActive) {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusBadRequest, fmt.Sprintf("Cannot suspend policy in %s status; must be ACTIVE", policy.Status), nil)
	}

	updated, err := s.policyRepo.UpdateStatus(ctx, id, string(shared.PolicyStatusSuspended))
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusInternalServerError, "Failed to suspend policy", err)
	}

	s.logAudit(ctx, uuid.Nil, string(shared.AuditEntityTypePolicy), id, string(shared.AuditActionStateChange))

	return schema.NewServiceResponse(policySchema.ToPolicyResponse(updated), http.StatusOK, "Policy suspended")
}

func (s *policyServiceImpl) UpdatePolicy(ctx context.Context, id uuid.UUID, req policySchema.UpdatePolicyRequest) *schema.ServiceResponse[policySchema.PolicyResponse] {
	policy, err := s.policyRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusNotFound, "Policy not found", err)
	}

	if req.PolicyholderName != nil {
		policy.PolicyholderName = *req.PolicyholderName
	}
	if req.PolicyholderEmail != nil {
		policy.PolicyholderEmail = *req.PolicyholderEmail
	}
	if req.PolicyholderPhone != nil {
		policy.PolicyholderPhone = *req.PolicyholderPhone
	}
	if req.StartDate != nil {
		policy.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		policy.EndDate = *req.EndDate
	}

	updated, err := s.policyRepo.Update(ctx, policy)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusInternalServerError, "Failed to update policy", err)
	}

	s.logAudit(ctx, uuid.Nil, string(shared.AuditEntityTypePolicy), id, string(shared.AuditActionUpdate))

	return schema.NewServiceResponse(policySchema.ToPolicyResponse(updated), http.StatusOK, "Policy updated")
}

func (s *policyServiceImpl) ChangePlan(ctx context.Context, policyID uuid.UUID, req policySchema.ChangePlanRequest, userID uuid.UUID) *schema.ServiceResponse[policySchema.PolicyResponse] {
	pol, err := s.policyRepo.GetByID(ctx, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusNotFound, "Policy not found", err)
	}

	if pol.Status != string(shared.PolicyStatusActive) {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusBadRequest, "Can only change plan for ACTIVE policies", nil)
	}

	newPlanID, err := uuid.Parse(req.NewPlanID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusBadRequest, "Invalid new plan ID", err)
	}

	newPlan, err := s.planRepo.GetByID(ctx, newPlanID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusNotFound, "New plan not found", err)
	}

	if newPlan.Status != string(shared.PlanStatusActive) {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusBadRequest, "New plan is not active", nil)
	}

	// Calculate new premium using premium rules if available
	newPremium := newPlan.BasePremium
	if s.premiumRuleSvc != nil && s.memberRepo != nil {
		memberCount, _ := s.memberRepo.CountActiveByPolicy(ctx, policyID)
		if memberCount > 0 {
			members, _ := s.memberRepo.ListActiveByPolicy(ctx, policyID)
			membersJSON, _ := json.Marshal(members)
			premResp := s.premiumRuleSvc.CalculatePremiumWithMembers(ctx, newPlanID, int(memberCount), membersJSON)
			if premResp.Error == nil && premResp.Data > 0 {
				newPremium = premResp.Data
			}
		}
	}

	// Prorate the premium difference for mid-term plan changes
	now := time.Now()
	totalDays := pol.EndDate.Sub(pol.StartDate).Hours() / 24
	remainingDays := pol.EndDate.Sub(now).Hours() / 24
	if totalDays > 0 && remainingDays > 0 {
		premiumDiff := newPremium - pol.PremiumAmount
		proratedAdjustment := int64(float64(premiumDiff) * remainingDays / totalDays)
		newPremium = pol.PremiumAmount + proratedAdjustment
	}

	oldPremium := pol.PremiumAmount

	updated, err := s.policyRepo.UpdatePlanAndPremium(ctx, policyID, newPlanID, newPremium)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusInternalServerError, "Failed to change plan", err)
	}

	// Create credit note for plan downgrade (new premium < old premium)
	if newPremium < oldPremium && s.creditNoteSvc != nil {
		refundAmount := oldPremium - newPremium
		reason := fmt.Sprintf("Plan downgrade from %s to %s — pro-rata refund", pol.PlanID, newPlanID)
		s.creditNoteSvc.CreateCreditNote(ctx, policyID, uuid.Nil, refundAmount, reason, userID)
	}

	s.logAudit(ctx, userID, string(shared.AuditEntityTypePolicy), policyID, string(shared.AuditActionUpdate))

	return schema.NewServiceResponse(policySchema.ToPolicyResponse(updated), http.StatusOK, "Plan changed successfully")
}

func (s *policyServiceImpl) BulkActivate(ctx context.Context, ids []uuid.UUID) *schema.ServiceResponse[policySchema.BulkResultResponse] {
	result := policySchema.BulkResultResponse{}
	for _, id := range ids {
		resp := s.ActivatePolicy(ctx, id)
		if resp.Error != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Policy %s: %s", id, resp.Message))
		} else {
			result.Succeeded++
		}
	}
	return schema.NewServiceResponse(result, http.StatusOK, "Bulk activation completed")
}

func (s *policyServiceImpl) BulkLapse(ctx context.Context, ids []uuid.UUID) *schema.ServiceResponse[policySchema.BulkResultResponse] {
	result := policySchema.BulkResultResponse{}
	for _, id := range ids {
		resp := s.LapsePolicy(ctx, id)
		if resp.Error != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Policy %s: %s", id, resp.Message))
		} else {
			result.Succeeded++
		}
	}
	return schema.NewServiceResponse(result, http.StatusOK, "Bulk lapse completed")
}

func (s *policyServiceImpl) CalculateProratedPremium(ctx context.Context, policyID uuid.UUID) *schema.ServiceResponse[int64] {
	pol, err := s.policyRepo.GetByID(ctx, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusNotFound, "Policy not found", err)
	}

	now := time.Now()
	totalDays := pol.EndDate.Sub(pol.StartDate).Hours() / 24
	if totalDays <= 0 {
		return schema.NewServiceResponse(pol.PremiumAmount, http.StatusOK, "Prorated premium calculated")
	}

	remainingDays := pol.EndDate.Sub(now).Hours() / 24
	if remainingDays < 0 {
		remainingDays = 0
	}

	prorated := int64(float64(pol.PremiumAmount) * remainingDays / totalDays)
	return schema.NewServiceResponse(prorated, http.StatusOK, "Prorated premium calculated")
}

func (s *policyServiceImpl) GetTotalCount(ctx context.Context) *schema.ServiceResponse[int64] {
	count, err := s.policyRepo.Count(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to get count", err)
	}
	return schema.NewServiceResponse(count, http.StatusOK, "Count retrieved")
}

func (s *policyServiceImpl) ListPoliciesFiltered(ctx context.Context, dateFrom, dateTo *time.Time, search string, page, pageSize int) *schema.ServiceResponse[[]policySchema.PolicyResponse] {
	offset := (page - 1) * pageSize
	policies, err := s.policyRepo.ListFiltered(ctx, dateFrom, dateTo, search, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]policySchema.PolicyResponse](http.StatusInternalServerError, "Failed to list policies", err)
	}
	responses := make([]policySchema.PolicyResponse, len(policies))
	for i, p := range policies {
		responses[i] = policySchema.ToPolicyResponse(p)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Policies retrieved")
}

func (s *policyServiceImpl) CountPoliciesFiltered(ctx context.Context, dateFrom, dateTo *time.Time, search string) *schema.ServiceResponse[int64] {
	count, err := s.policyRepo.CountFiltered(ctx, dateFrom, dateTo, search)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to count policies", err)
	}
	return schema.NewServiceResponse(count, http.StatusOK, "Policies counted")
}

func (s *policyServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}
