package policy

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/domains/policy/entity"
	"github.com/bitbiz/hias-core/domains/policy/repository"
	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/bitbiz/hias-core/domains/policy/service"
	planRepo "github.com/bitbiz/hias-core/domains/product/repository"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/google/uuid"
)

type policyServiceImpl struct {
	policyRepo repository.PolicyRepository
	planRepo   planRepo.PlanRepository
	auditSvc   auditService.AuditService
}

func NewPolicyService(
	policyRepo repository.PolicyRepository,
	planRepo planRepo.PlanRepository,
	auditSvc auditService.AuditService,
) service.PolicyService {
	return &policyServiceImpl{
		policyRepo: policyRepo,
		planRepo:   planRepo,
		auditSvc:   auditSvc,
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

func (s *policyServiceImpl) ActivatePolicy(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[policySchema.PolicyResponse] {
	policy, err := s.policyRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusNotFound, "Policy not found", err)
	}

	if policy.Status != string(shared.PolicyStatusDraft) {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusBadRequest, fmt.Sprintf("Cannot activate policy in %s status", policy.Status), nil)
	}

	updated, err := s.policyRepo.UpdateStatus(ctx, id, string(shared.PolicyStatusActive))
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusInternalServerError, "Failed to activate policy", err)
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

	if policy.Status != string(shared.PolicyStatusLapsed) {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusBadRequest, "Only lapsed policies can be reinstated", nil)
	}

	updated, err := s.policyRepo.UpdateStatus(ctx, id, string(shared.PolicyStatusActive))
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusInternalServerError, "Failed to reinstate policy", err)
	}

	return schema.NewServiceResponse(policySchema.ToPolicyResponse(updated), http.StatusOK, "Policy reinstated")
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

func (s *policyServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}
