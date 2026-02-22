package policy

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/domains/policy/entity"
	"github.com/bitbiz/hias-core/domains/policy/repository"
	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/bitbiz/hias-core/domains/policy/service"
	productRepo "github.com/bitbiz/hias-core/domains/product/repository"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/google/uuid"
)

type policyServiceImpl struct {
	policyRepo repository.PolicyRepository
	planRepo   productRepo.PlanRepository
}

func NewPolicyService(
	policyRepo repository.PolicyRepository,
	planRepo productRepo.PlanRepository,
) service.PolicyService {
	return &policyServiceImpl{
		policyRepo: policyRepo,
		planRepo:   planRepo,
	}
}

func (s *policyServiceImpl) CreatePolicy(ctx context.Context, req policySchema.CreatePolicyRequest, createdBy uuid.UUID) *schema.ServiceResponse[policySchema.PolicyResponse] {
	planID, err := uuid.Parse(req.PlanID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusBadRequest, "Invalid plan ID", err)
	}

	// Verify the plan exists
	plan, err := s.planRepo.GetByID(ctx, planID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusNotFound, "Plan not found", err)
	}

	// Set default dates if not provided
	startDate := req.StartDate
	if startDate.IsZero() {
		startDate = time.Now()
	}
	endDate := req.EndDate
	if endDate.IsZero() {
		endDate = startDate.AddDate(1, 0, 0) // Default: 1 year from start
	}

	policyNumber := utils.GeneratePolicyNumber()

	policy, err := s.policyRepo.Create(ctx, &entity.Policy{
		PlanID:            planID,
		PolicyholderName:  req.PolicyholderName,
		PolicyholderEmail: req.PolicyholderEmail,
		PolicyholderPhone: req.PolicyholderPhone,
		PolicyNumber:      policyNumber,
		Status:            string(shared.PolicyStatusDraft),
		StartDate:         startDate,
		EndDate:           endDate,
		PremiumAmount:     plan.BasePremium,
		Currency:          plan.Currency,
		CreatedBy:         createdBy,
	})
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusInternalServerError, "Failed to create policy", err)
	}

	return schema.NewServiceResponse(policySchema.ToPolicyResponse(policy), http.StatusCreated, "Policy created successfully")
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

	// State machine: only DRAFT → ACTIVE
	if policy.Status != string(shared.PolicyStatusDraft) {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](
			http.StatusConflict,
			fmt.Sprintf("Cannot activate policy: current status is %s, expected DRAFT", policy.Status),
			nil,
		)
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

	// State machine: only ACTIVE → LAPSED
	if policy.Status != string(shared.PolicyStatusActive) {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](
			http.StatusConflict,
			fmt.Sprintf("Cannot lapse policy: current status is %s, expected ACTIVE", policy.Status),
			nil,
		)
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

	// State machine: ACTIVE or LAPSED → TERMINATED
	if policy.Status != string(shared.PolicyStatusActive) && policy.Status != string(shared.PolicyStatusLapsed) {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](
			http.StatusConflict,
			fmt.Sprintf("Cannot terminate policy: current status is %s, expected ACTIVE or LAPSED", policy.Status),
			nil,
		)
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

	// State machine: only LAPSED → ACTIVE
	if policy.Status != string(shared.PolicyStatusLapsed) {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](
			http.StatusConflict,
			fmt.Sprintf("Cannot reinstate policy: current status is %s, expected LAPSED", policy.Status),
			nil,
		)
	}

	updated, err := s.policyRepo.UpdateStatus(ctx, id, string(shared.PolicyStatusActive))
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyResponse](http.StatusInternalServerError, "Failed to reinstate policy", err)
	}

	return schema.NewServiceResponse(policySchema.ToPolicyResponse(updated), http.StatusOK, "Policy reinstated")
}

func (s *policyServiceImpl) GetTotalCount(ctx context.Context) *schema.ServiceResponse[int64] {
	count, err := s.policyRepo.Count(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to count policies", err)
	}
	return schema.NewServiceResponse(count, http.StatusOK, "Count retrieved")
}
