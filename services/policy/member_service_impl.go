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
	"github.com/bitbiz/hias-core/shared"
	"github.com/google/uuid"
)

type memberServiceImpl struct {
	memberRepo repository.MemberRepository
	policyRepo repository.PolicyRepository
}

func NewMemberService(
	memberRepo repository.MemberRepository,
	policyRepo repository.PolicyRepository,
) service.MemberService {
	return &memberServiceImpl{
		memberRepo: memberRepo,
		policyRepo: policyRepo,
	}
}

func (s *memberServiceImpl) EnrollMember(ctx context.Context, policyID uuid.UUID, req policySchema.EnrollMemberRequest) *schema.ServiceResponse[policySchema.MemberResponse] {
	// Verify the policy exists
	policy, err := s.policyRepo.GetByID(ctx, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.MemberResponse](http.StatusNotFound, "Policy not found", err)
	}

	// Only allow enrollment on DRAFT or ACTIVE policies
	if policy.Status != string(shared.PolicyStatusDraft) && policy.Status != string(shared.PolicyStatusActive) {
		return schema.NewServiceErrorResponse[policySchema.MemberResponse](
			http.StatusConflict,
			fmt.Sprintf("Cannot enroll member: policy status is %s, expected DRAFT or ACTIVE", policy.Status),
			nil,
		)
	}

	// Parse date of birth
	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.MemberResponse](http.StatusBadRequest, "Invalid date_of_birth format, expected YYYY-MM-DD", err)
	}

	// Generate member number: MEM-{policyNumber suffix}-{sequential}
	count, err := s.memberRepo.CountByPolicy(ctx, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.MemberResponse](http.StatusInternalServerError, "Failed to count existing members", err)
	}
	memberNumber := fmt.Sprintf("MEM-%s-%03d", policy.PolicyNumber, count+1)

	member, err := s.memberRepo.Create(ctx, &entity.Member{
		PolicyID:     policyID,
		NationalID:   req.NationalID,
		Name:         req.Name,
		DateOfBirth:  dob,
		Gender:       req.Gender,
		Relationship: req.Relationship,
		MemberNumber: memberNumber,
		Phone:        req.Phone,
		Email:        req.Email,
		Verified:     false,
	})
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.MemberResponse](http.StatusInternalServerError, "Failed to enroll member", err)
	}

	return schema.NewServiceResponse(policySchema.ToMemberResponse(member), http.StatusCreated, "Member enrolled successfully")
}

func (s *memberServiceImpl) VerifyMember(ctx context.Context, memberID uuid.UUID) *schema.ServiceResponse[policySchema.MemberResponse] {
	member, err := s.memberRepo.GetByID(ctx, memberID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.MemberResponse](http.StatusNotFound, "Member not found", err)
	}

	if member.Verified {
		return schema.NewServiceResponse(policySchema.ToMemberResponse(member), http.StatusOK, "Member already verified")
	}

	// TODO: Integrate with IPRS (Integrated Population Registration System) for national ID verification.
	// For now, stub the integration and set verified=true directly.
	verified, err := s.memberRepo.Verify(ctx, memberID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.MemberResponse](http.StatusInternalServerError, "Failed to verify member", err)
	}

	return schema.NewServiceResponse(policySchema.ToMemberResponse(verified), http.StatusOK, "Member verified successfully")
}

func (s *memberServiceImpl) GetMemberEligibility(ctx context.Context, memberID uuid.UUID) *schema.ServiceResponse[bool] {
	member, err := s.memberRepo.GetByID(ctx, memberID)
	if err != nil {
		return schema.NewServiceErrorResponse[bool](http.StatusNotFound, "Member not found", err)
	}

	// Check if the member's policy is active
	policy, err := s.policyRepo.GetByID(ctx, member.PolicyID)
	if err != nil {
		return schema.NewServiceErrorResponse[bool](http.StatusInternalServerError, "Failed to retrieve policy", err)
	}

	if policy.Status != string(shared.PolicyStatusActive) {
		return schema.NewServiceResponse(false, http.StatusOK, fmt.Sprintf("Member not eligible: policy status is %s", policy.Status))
	}

	// Check if the policy is within the coverage period
	now := time.Now()
	if now.Before(policy.StartDate) || now.After(policy.EndDate) {
		return schema.NewServiceResponse(false, http.StatusOK, "Member not eligible: outside policy coverage period")
	}

	return schema.NewServiceResponse(true, http.StatusOK, "Member is eligible")
}

func (s *memberServiceImpl) ListMembers(ctx context.Context, policyID uuid.UUID) *schema.ServiceResponse[[]policySchema.MemberResponse] {
	// Verify the policy exists
	_, err := s.policyRepo.GetByID(ctx, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]policySchema.MemberResponse](http.StatusNotFound, "Policy not found", err)
	}

	members, err := s.memberRepo.ListByPolicy(ctx, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]policySchema.MemberResponse](http.StatusInternalServerError, "Failed to list members", err)
	}

	responses := make([]policySchema.MemberResponse, len(members))
	for i, m := range members {
		responses[i] = policySchema.ToMemberResponse(m)
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Members retrieved")
}
