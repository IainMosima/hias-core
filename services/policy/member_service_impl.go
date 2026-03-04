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
	policyRepo "github.com/bitbiz/hias-core/domains/policy/repository"
	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/bitbiz/hias-core/domains/policy/service"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/google/uuid"
)

type memberServiceImpl struct {
	memberRepo repository.MemberRepository
	policyRepo policyRepo.PolicyRepository
	auditSvc   auditService.AuditService
}

func NewMemberService(
	memberRepo repository.MemberRepository,
	policyRepo policyRepo.PolicyRepository,
	auditSvc auditService.AuditService,
) service.MemberService {
	return &memberServiceImpl{
		memberRepo: memberRepo,
		policyRepo: policyRepo,
		auditSvc:   auditSvc,
	}
}

func (s *memberServiceImpl) EnrollMember(ctx context.Context, policyID uuid.UUID, req policySchema.EnrollMemberRequest) *schema.ServiceResponse[policySchema.MemberResponse] {
	// Verify policy exists
	pol, err := s.policyRepo.GetByID(ctx, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.MemberResponse](http.StatusNotFound, "Policy not found", err)
	}

	if pol.Status != string(shared.PolicyStatusActive) && pol.Status != string(shared.PolicyStatusDraft) {
		return schema.NewServiceErrorResponse[policySchema.MemberResponse](http.StatusBadRequest, fmt.Sprintf("Cannot enroll members in %s policy", pol.Status), nil)
	}

	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.MemberResponse](http.StatusBadRequest, "Invalid date of birth format (YYYY-MM-DD)", err)
	}

	memberNumber := utils.GenerateMemberNumber()

	member := &entity.Member{
		PolicyID:     policyID,
		NationalID:   req.NationalID,
		Name:         req.Name,
		DateOfBirth:  dob,
		Gender:       req.Gender,
		Relationship: req.Relationship,
		MemberNumber: memberNumber,
		Phone:        req.Phone,
		Email:        req.Email,
		KRAPin:       req.KRAPin,
		County:       req.County,
		Address:      req.Address,
		Verified:     false,
	}

	created, err := s.memberRepo.Create(ctx, member)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.MemberResponse](http.StatusInternalServerError, "Failed to enroll member", err)
	}

	s.logAudit(ctx, uuid.Nil, string(shared.AuditEntityTypeMember), created.ID, string(shared.AuditActionCreate))

	return schema.NewServiceResponse(policySchema.ToMemberResponse(created), http.StatusCreated, "Member enrolled")
}

func (s *memberServiceImpl) VerifyMember(ctx context.Context, memberID uuid.UUID) *schema.ServiceResponse[policySchema.MemberResponse] {
	verified, err := s.memberRepo.Verify(ctx, memberID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.MemberResponse](http.StatusInternalServerError, "Failed to verify member", err)
	}
	return schema.NewServiceResponse(policySchema.ToMemberResponse(verified), http.StatusOK, "Member verified")
}

func (s *memberServiceImpl) GetMemberEligibility(ctx context.Context, memberID uuid.UUID) *schema.ServiceResponse[bool] {
	member, err := s.memberRepo.GetByID(ctx, memberID)
	if err != nil {
		return schema.NewServiceErrorResponse[bool](http.StatusNotFound, "Member not found", err)
	}

	pol, err := s.policyRepo.GetByID(ctx, member.PolicyID)
	if err != nil {
		return schema.NewServiceErrorResponse[bool](http.StatusNotFound, "Policy not found", err)
	}

	eligible := pol.Status == string(shared.PolicyStatusActive) && time.Now().Before(pol.EndDate)
	return schema.NewServiceResponse(eligible, http.StatusOK, "Eligibility checked")
}

func (s *memberServiceImpl) ListMembers(ctx context.Context, policyID uuid.UUID) *schema.ServiceResponse[[]policySchema.MemberResponse] {
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

func (s *memberServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}
