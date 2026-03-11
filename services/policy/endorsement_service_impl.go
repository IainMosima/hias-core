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
	productRepo "github.com/bitbiz/hias-core/domains/product/repository"
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
	planRepo        productRepo.PlanRepository
	policyDocSvc    service.PolicyDocumentService
	auditSvc        auditService.AuditService
}

func NewEndorsementService(
	endorsementRepo repository.EndorsementRepository,
	policyRepo repository.PolicyRepository,
	memberRepo repository.MemberRepository,
	memberSvc service.MemberService,
	policySvc service.PolicyService,
	premiumRuleSvc productService.PremiumRuleService,
	planRepo productRepo.PlanRepository,
	policyDocSvc service.PolicyDocumentService,
	auditSvc auditService.AuditService,
) service.EndorsementService {
	return &endorsementServiceImpl{
		endorsementRepo: endorsementRepo,
		policyRepo:      policyRepo,
		memberRepo:      memberRepo,
		memberSvc:       memberSvc,
		policySvc:       policySvc,
		premiumRuleSvc:  premiumRuleSvc,
		planRepo:        planRepo,
		policyDocSvc:    policyDocSvc,
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

	// Effective date must be within policy term
	if effectiveDate.Before(pol.StartDate.Truncate(24*time.Hour)) || effectiveDate.After(pol.EndDate.Truncate(24*time.Hour)) {
		return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "Effective date must be within policy term (start_date to end_date)", nil)
	}

	// Effective date must not be more than 90 days in the past
	ninetyDaysAgo := time.Now().AddDate(0, 0, -90).Truncate(24 * time.Hour)
	if effectiveDate.Before(ninetyDaysAgo) {
		return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "Effective date cannot be more than 90 days in the past", nil)
	}

	// Validate changes payload based on endorsement type
	if errResp := s.validateChangesPayload(ctx, req.EndorsementType, req.Changes, policyID); errResp != nil {
		return errResp
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

// validateChangesPayload validates the changes JSON based on endorsement type
func (s *endorsementServiceImpl) validateChangesPayload(ctx context.Context, endorsementType string, changes json.RawMessage, policyID uuid.UUID) *schema.ServiceResponse[policySchema.EndorsementResponse] {
	switch endorsementType {
	case string(shared.EndorsementTypeAddMember):
		var req policySchema.EnrollMemberRequest
		if err := json.Unmarshal(changes, &req); err != nil {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "Invalid ADD_MEMBER payload: "+err.Error(), err)
		}
		if req.Name == "" {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "ADD_MEMBER: name is required", nil)
		}
		if req.DateOfBirth == "" {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "ADD_MEMBER: date_of_birth is required", nil)
		}
		if _, err := time.Parse("2006-01-02", req.DateOfBirth); err != nil {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "ADD_MEMBER: invalid date_of_birth format (YYYY-MM-DD)", err)
		}
		if req.Gender == "" {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "ADD_MEMBER: gender is required", nil)
		}
		validGenders := map[string]bool{"male": true, "female": true, "other": true}
		if !validGenders[req.Gender] {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "ADD_MEMBER: gender must be male, female, or other", nil)
		}
		if req.Relationship == "" {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "ADD_MEMBER: relationship is required", nil)
		}
		validRelationships := map[string]bool{"principal": true, "spouse": true, "child": true, "parent": true}
		if !validRelationships[req.Relationship] {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "ADD_MEMBER: relationship must be principal, spouse, child, or parent", nil)
		}

	case string(shared.EndorsementTypeRemoveMember):
		var req policySchema.RemoveMemberChanges
		if err := json.Unmarshal(changes, &req); err != nil {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "Invalid REMOVE_MEMBER payload: "+err.Error(), err)
		}
		if req.MemberID == "" {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "REMOVE_MEMBER: member_id is required", nil)
		}
		memberID, err := uuid.Parse(req.MemberID)
		if err != nil {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "REMOVE_MEMBER: invalid member_id (must be UUID)", err)
		}
		member, err := s.memberRepo.GetByID(ctx, memberID)
		if err != nil {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "REMOVE_MEMBER: member not found", err)
		}
		if member.PolicyID != policyID {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "REMOVE_MEMBER: member does not belong to this policy", nil)
		}
		if member.Status == string(shared.MemberStatusRemoved) {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "REMOVE_MEMBER: member is already REMOVED", nil)
		}

	case string(shared.EndorsementTypeUpdateMember):
		var req policySchema.UpdateMemberChanges
		if err := json.Unmarshal(changes, &req); err != nil {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "Invalid UPDATE_MEMBER payload: "+err.Error(), err)
		}
		if req.MemberID == "" {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "UPDATE_MEMBER: member_id is required", nil)
		}
		memberID, err := uuid.Parse(req.MemberID)
		if err != nil {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "UPDATE_MEMBER: invalid member_id (must be UUID)", err)
		}
		member, err := s.memberRepo.GetByID(ctx, memberID)
		if err != nil {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "UPDATE_MEMBER: member not found", err)
		}
		if member.PolicyID != policyID {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "UPDATE_MEMBER: member does not belong to this policy", nil)
		}
		if member.Status == string(shared.MemberStatusRemoved) {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "UPDATE_MEMBER: cannot update a REMOVED member", nil)
		}
		// At least one field in updates must be non-nil
		u := req.Updates
		if u.Name == nil && u.Phone == nil && u.Email == nil && u.KRAPin == nil && u.County == nil && u.Address == nil {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "UPDATE_MEMBER: at least one field in updates must be provided", nil)
		}

	case string(shared.EndorsementTypePlanChange):
		var req policySchema.ChangePlanRequest
		if err := json.Unmarshal(changes, &req); err != nil {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "Invalid PLAN_CHANGE payload: "+err.Error(), err)
		}
		if req.NewPlanID == "" {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "PLAN_CHANGE: new_plan_id is required", nil)
		}
		newPlanID, err := uuid.Parse(req.NewPlanID)
		if err != nil {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "PLAN_CHANGE: invalid new_plan_id (must be UUID)", err)
		}
		// Check plan exists and is active
		if s.planRepo != nil {
			plan, planErr := s.planRepo.GetByID(ctx, newPlanID)
			if planErr != nil {
				return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "PLAN_CHANGE: plan not found", planErr)
			}
			if plan.Status != string(shared.PlanStatusActive) {
				return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "PLAN_CHANGE: plan is not ACTIVE", nil)
			}
		}
		// Check new plan differs from current
		pol, polErr := s.policyRepo.GetByID(ctx, policyID)
		if polErr == nil && pol.PlanID == newPlanID {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest, "PLAN_CHANGE: new_plan_id must differ from current plan", nil)
		}
	}

	return nil
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
		resp := s.memberSvc.EnrollMember(ctx, endorsement.PolicyID, enrollReq)
		if resp.Error != nil {
			return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusInternalServerError, fmt.Sprintf("Failed to add member: %s", resp.Message), resp.Error)
		}
		// Set coverage_start_date on newly created member
		if s.memberRepo != nil {
			effectiveDate := endorsement.EffectiveDate
			if _, err := s.memberRepo.UpdateCoverageDates(ctx, resp.Data.ID, &effectiveDate, nil); err != nil {
				log.Printf("Warning: failed to set coverage_start_date for member %s: %v", resp.Data.ID, err)
			}
		}

	case string(shared.EndorsementTypeRemoveMember):
		var removePayload policySchema.RemoveMemberChanges
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
		// Set coverage_end_date on removed member
		if s.memberRepo != nil {
			effectiveDate := endorsement.EffectiveDate
			if _, err := s.memberRepo.UpdateCoverageDates(ctx, memberID, nil, &effectiveDate); err != nil {
				log.Printf("Warning: failed to set coverage_end_date for member %s: %v", memberID, err)
			}
		}

	case string(shared.EndorsementTypeUpdateMember):
		var updatePayload policySchema.UpdateMemberChanges
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

	// Apply premium adjustment to the policy — skip for REMOVE_MEMBER since RemoveMember already recalculates
	if endorsement.PremiumAdjustment != 0 && endorsement.EndorsementType != string(shared.EndorsementTypeRemoveMember) {
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

	// Generate endorsement letter
	if s.policyDocSvc != nil {
		go func() {
			bgCtx := context.Background()
			s.policyDocSvc.GenerateWelcomeLetter(bgCtx, endorsement.PolicyID, endorsement.ApprovedBy)
		}()
	}

	return schema.NewServiceResponse(policySchema.ToEndorsementResponse(updated), http.StatusOK, "Endorsement applied")
}

func (s *endorsementServiceImpl) CancelEndorsement(ctx context.Context, id uuid.UUID, reason string, cancelledBy uuid.UUID) *schema.ServiceResponse[policySchema.EndorsementResponse] {
	endorsement, err := s.endorsementRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusNotFound, "Endorsement not found", err)
	}

	// Only PENDING or APPROVED can be cancelled
	if endorsement.Status != string(shared.EndorsementStatusPending) && endorsement.Status != string(shared.EndorsementStatusApproved) {
		return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusBadRequest,
			fmt.Sprintf("Cannot cancel endorsement in %s status (must be PENDING or APPROVED)", endorsement.Status), nil)
	}

	endorsement.Status = string(shared.EndorsementStatusCancelled)
	endorsement.Reason = reason

	updated, err := s.endorsementRepo.Update(ctx, endorsement)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.EndorsementResponse](http.StatusInternalServerError, "Failed to cancel endorsement", err)
	}

	s.logAudit(ctx, cancelledBy, string(shared.AuditEntityTypeEndorsement), id, string(shared.AuditActionStateChange))

	return schema.NewServiceResponse(policySchema.ToEndorsementResponse(updated), http.StatusOK, "Endorsement cancelled")
}

func (s *endorsementServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}
