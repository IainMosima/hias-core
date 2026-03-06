package policy

import (
	"context"
	"encoding/json"
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
	productService "github.com/bitbiz/hias-core/domains/product/service"
	"github.com/bitbiz/hias-core/shared"
	"github.com/google/uuid"
)

type endorsementServiceImpl struct {
	endorsementRepo repository.EndorsementRepository
	policyRepo      repository.PolicyRepository
	memberRepo      repository.MemberRepository
	memberSvc       service.MemberService
	policySvc       service.PolicyService
	premiumRuleSvc  productService.PremiumRuleService
	auditSvc        auditService.AuditService
}

func NewEndorsementService(
	endorsementRepo repository.EndorsementRepository,
	policyRepo repository.PolicyRepository,
	memberRepo repository.MemberRepository,
	memberSvc service.MemberService,
	policySvc service.PolicyService,
	premiumRuleSvc productService.PremiumRuleService,
	auditSvc auditService.AuditService,
) service.EndorsementService {
	return &endorsementServiceImpl{
		endorsementRepo: endorsementRepo,
		policyRepo:      policyRepo,
		memberRepo:      memberRepo,
		memberSvc:       memberSvc,
		policySvc:       policySvc,
		premiumRuleSvc:  premiumRuleSvc,
		auditSvc:        auditSvc,
	}
}

func (s *endorsementServiceImpl) CreateEndorsement(ctx context.Context, req policySchema.CreateEndorsementRequest, requestedBy uuid.UUID) *schema.ServiceResponse[policySchema.EndorsementResponse] {
	policyID, err := uuid.Parse(req.PolicyID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "Invalid policy ID", err)
	}

	pol, err := s.policyRepo.GetByID(ctx, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusNotFound, "Policy not found", err)
	}

	if pol.Status != string(shared.PolicyStatusActive) {
		return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "Endorsements can only be created for ACTIVE policies", nil)
	}

	validTypes := map[string]bool{
		string(shared.EndorsementTypeAddMember):    true,
		string(shared.EndorsementTypeRemoveMember): true,
		string(shared.EndorsementTypeUpdateMember): true,
		string(shared.EndorsementTypePlanChange):   true,
	}
	if !validTypes[req.EndorsementType] {
		return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "Invalid endorsement type", nil)
	}

	effectiveDate, err := time.Parse("2006-01-02", req.EffectiveDate)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "Invalid effective date format (YYYY-MM-DD)", err)
	}

	endorsement := &entity.Endorsement{
		PolicyID:          policyID,
		EndorsementType:   req.EndorsementType,
		Status:            string(shared.EndorsementStatusPending),
		EffectiveDate:     effectiveDate,
		Changes:           req.Changes,
		Reason:            req.Reason,
		PremiumAdjustment: req.PremiumAdjustment,
		RequestedBy:       requestedBy,
	}

	created, err := s.endorsementRepo.Create(ctx, endorsement)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusInternalServerError, "Failed to create endorsement", err)
	}

	s.logAudit(ctx, requestedBy, string(shared.AuditEntityTypeEndorsement), created.ID, string(shared.AuditActionCreate))

	return schema.NewServiceResponse(policySchema.ToEndorsementResponse(created), http.StatusCreated, "Endorsement created")
}

func (s *endorsementServiceImpl) GetEndorsement(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[policySchema.EndorsementResponse] {
	endorsement, err := s.endorsementRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusNotFound, "Endorsement not found", err)
	}
	return schema.NewServiceResponse(policySchema.ToEndorsementResponse(endorsement), http.StatusOK, "Endorsement retrieved")
}

func (s *endorsementServiceImpl) ListByPolicy(ctx context.Context, policyID uuid.UUID) *schema.ServiceResponse[[]policySchema.EndorsementResponse] {
	endorsements, err := s.endorsementRepo.ListByPolicy(ctx, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]policySchema.EndorsementResponse](http.StatusInternalServerError, "Failed to list endorsements", err)
	}

	responses := make([]policySchema.EndorsementResponse, len(endorsements))
	for i, e := range endorsements {
		responses[i] = policySchema.ToEndorsementResponse(e)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Endorsements retrieved")
}

func (s *endorsementServiceImpl) ApproveEndorsement(ctx context.Context, id uuid.UUID, approvedBy uuid.UUID) *schema.ServiceResponse[policySchema.EndorsementResponse] {
	endorsement, err := s.endorsementRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusNotFound, "Endorsement not found", err)
	}

	if endorsement.Status != string(shared.EndorsementStatusPending) {
		return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, fmt.Sprintf("Cannot approve endorsement in %s status", endorsement.Status), nil)
	}

	now := time.Now()
	endorsement.Status = string(shared.EndorsementStatusApproved)
	endorsement.ApprovedBy = approvedBy
	endorsement.ApprovedAt = &now

	updated, err := s.endorsementRepo.Update(ctx, endorsement)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusInternalServerError, "Failed to approve endorsement", err)
	}

	s.logAudit(ctx, approvedBy, string(shared.AuditEntityTypeEndorsement), id, string(shared.AuditActionStateChange))

	return schema.NewServiceResponse(policySchema.ToEndorsementResponse(updated), http.StatusOK, "Endorsement approved")
}

func (s *endorsementServiceImpl) RejectEndorsement(ctx context.Context, id uuid.UUID, reason string, rejectedBy uuid.UUID) *schema.ServiceResponse[policySchema.EndorsementResponse] {
	endorsement, err := s.endorsementRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusNotFound, "Endorsement not found", err)
	}

	if endorsement.Status != string(shared.EndorsementStatusPending) {
		return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, fmt.Sprintf("Cannot reject endorsement in %s status", endorsement.Status), nil)
	}

	endorsement.Status = string(shared.EndorsementStatusRejected)
	endorsement.Reason = reason

	updated, err := s.endorsementRepo.Update(ctx, endorsement)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusInternalServerError, "Failed to reject endorsement", err)
	}

	s.logAudit(ctx, rejectedBy, string(shared.AuditEntityTypeEndorsement), id, string(shared.AuditActionStateChange))

	return schema.NewServiceResponse(policySchema.ToEndorsementResponse(updated), http.StatusOK, "Endorsement rejected")
}

func (s *endorsementServiceImpl) ApplyEndorsement(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[policySchema.EndorsementResponse] {
	endorsement, err := s.endorsementRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusNotFound, "Endorsement not found", err)
	}

	if endorsement.Status != string(shared.EndorsementStatusApproved) {
		return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "Only approved endorsements can be applied", nil)
	}

	// Calculate premium adjustment if not manually set
	if endorsement.PremiumAdjustment == 0 && s.premiumRuleSvc != nil && s.memberRepo != nil {
		pol, polErr := s.policyRepo.GetByID(ctx, endorsement.PolicyID)
		if polErr == nil {
			currentCount, countErr := s.memberRepo.CountActiveByPolicy(ctx, endorsement.PolicyID)
			if countErr == nil {
				var newCount int
				switch endorsement.EndorsementType {
				case string(shared.EndorsementTypeAddMember):
					newCount = int(currentCount) + 1
				case string(shared.EndorsementTypeRemoveMember):
					newCount = int(currentCount) - 1
					if newCount < 0 {
						newCount = 0
					}
				}
				if newCount > 0 && newCount != int(currentCount) {
					premResp := s.premiumRuleSvc.CalculatePremiumWithMembers(ctx, pol.PlanID, newCount, nil)
					if premResp.Error == nil && premResp.Data > 0 {
						fullDelta := premResp.Data - pol.PremiumAmount
						// Pro-rate for remaining days in policy period
						if pol.EndDate.After(time.Now()) && pol.StartDate.Before(pol.EndDate) {
							totalDays := pol.EndDate.Sub(pol.StartDate).Hours() / 24
							remainingDays := pol.EndDate.Sub(time.Now()).Hours() / 24
							if totalDays > 0 {
								endorsement.PremiumAdjustment = int64(float64(fullDelta) * remainingDays / totalDays)
							}
						} else {
							endorsement.PremiumAdjustment = fullDelta
						}
						log.Printf("Endorsement premium adjustment calculated: %d (full delta: %d)", endorsement.PremiumAdjustment, fullDelta)
					}
				}
			}
		}
	}

	// Dispatch action based on endorsement type
	switch endorsement.EndorsementType {
	case string(shared.EndorsementTypeAddMember):
		var enrollReq policySchema.EnrollMemberRequest
		if err := json.Unmarshal(endorsement.Changes, &enrollReq); err != nil {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "Invalid changes payload for ADD_MEMBER", err)
		}
		if resp := s.memberSvc.EnrollMember(ctx, endorsement.PolicyID, enrollReq); resp.Error != nil {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusInternalServerError, fmt.Sprintf("Failed to add member: %s", resp.Message), resp.Error)
		}

	case string(shared.EndorsementTypeRemoveMember):
		var removePayload struct {
			MemberID string `json:"member_id"`
			Reason   string `json:"reason"`
		}
		if err := json.Unmarshal(endorsement.Changes, &removePayload); err != nil {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "Invalid changes payload for REMOVE_MEMBER", err)
		}
		memberID, err := uuid.Parse(removePayload.MemberID)
		if err != nil {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "Invalid member ID in changes", err)
		}
		if resp := s.memberSvc.RemoveMember(ctx, memberID, removePayload.Reason); resp.Error != nil {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusInternalServerError, fmt.Sprintf("Failed to remove member: %s", resp.Message), resp.Error)
		}

	case string(shared.EndorsementTypeUpdateMember):
		var updatePayload struct {
			MemberID string                           `json:"member_id"`
			Updates  policySchema.UpdateMemberRequest `json:"updates"`
		}
		if err := json.Unmarshal(endorsement.Changes, &updatePayload); err != nil {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "Invalid changes payload for UPDATE_MEMBER", err)
		}
		memberID, err := uuid.Parse(updatePayload.MemberID)
		if err != nil {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "Invalid member ID in changes", err)
		}
		if resp := s.memberSvc.UpdateMember(ctx, memberID, updatePayload.Updates); resp.Error != nil {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusInternalServerError, fmt.Sprintf("Failed to update member: %s", resp.Message), resp.Error)
		}

	case string(shared.EndorsementTypePlanChange):
		var changePlanReq policySchema.ChangePlanRequest
		if err := json.Unmarshal(endorsement.Changes, &changePlanReq); err != nil {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "Invalid changes payload for PLAN_CHANGE", err)
		}
		if resp := s.policySvc.ChangePlan(ctx, endorsement.PolicyID, changePlanReq, endorsement.ApprovedBy); resp.Error != nil {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusInternalServerError, fmt.Sprintf("Failed to change plan: %s", resp.Message), resp.Error)
		}
	}

	// Apply premium adjustment to the policy if specified
	if endorsement.PremiumAdjustment != 0 {
		pol, err := s.policyRepo.GetByID(ctx, endorsement.PolicyID)
		if err != nil {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusInternalServerError, "Failed to get policy for adjustment", err)
		}
		pol.PremiumAmount += endorsement.PremiumAdjustment
		if _, err := s.policyRepo.Update(ctx, pol); err != nil {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusInternalServerError, "Failed to apply premium adjustment", err)
		}
	}

	now := time.Now()
	endorsement.Status = string(shared.EndorsementStatusApplied)
	endorsement.AppliedAt = &now

	updated, err := s.endorsementRepo.Update(ctx, endorsement)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusInternalServerError, "Failed to apply endorsement", err)
	}

	s.logAudit(ctx, endorsement.ApprovedBy, string(shared.AuditEntityTypeEndorsement), id, string(shared.AuditActionStateChange))

	return schema.NewServiceResponse(policySchema.ToEndorsementResponse(updated), http.StatusOK, "Endorsement applied")
}

func (s *endorsementServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}
