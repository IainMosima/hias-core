package policy

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	claimRepo "github.com/bitbiz/hias-core/domains/claims/repository"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/domains/policy/entity"
	"github.com/bitbiz/hias-core/domains/policy/repository"
	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/bitbiz/hias-core/domains/policy/service"
	productEntity "github.com/bitbiz/hias-core/domains/product/entity"
	productRepo "github.com/bitbiz/hias-core/domains/product/repository"
	productService "github.com/bitbiz/hias-core/domains/product/service"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/google/uuid"
)

type renewalServiceImpl struct {
	renewalRepo          repository.PolicyRenewalRepository
	policyRepo           repository.PolicyRepository
	memberRepo           repository.MemberRepository
	claimRepo            claimRepo.ClaimRepository
	premiumRuleSvc       productService.PremiumRuleService
	premiumRuleRepo      productRepo.PremiumRuleRepository
	planRepo             productRepo.PlanRepository
	underwritingFlagRepo repository.UnderwritingFlagRepository
	auditSvc             auditService.AuditService
	policyDocSvc         service.PolicyDocumentService
}

func NewRenewalService(
	renewalRepo repository.PolicyRenewalRepository,
	policyRepo repository.PolicyRepository,
	memberRepo repository.MemberRepository,
	claimRepo claimRepo.ClaimRepository,
	premiumRuleSvc productService.PremiumRuleService,
	premiumRuleRepo productRepo.PremiumRuleRepository,
	planRepo productRepo.PlanRepository,
	underwritingFlagRepo repository.UnderwritingFlagRepository,
	auditSvc auditService.AuditService,
	policyDocSvc service.PolicyDocumentService,
) service.RenewalService {
	return &renewalServiceImpl{
		renewalRepo:          renewalRepo,
		policyRepo:           policyRepo,
		memberRepo:           memberRepo,
		claimRepo:            claimRepo,
		premiumRuleSvc:       premiumRuleSvc,
		premiumRuleRepo:      premiumRuleRepo,
		planRepo:             planRepo,
		underwritingFlagRepo: underwritingFlagRepo,
		auditSvc:             auditSvc,
		policyDocSvc:         policyDocSvc,
	}
}

func (s *renewalServiceImpl) InitiateRenewal(ctx context.Context, req policySchema.InitiateRenewalRequest, createdBy uuid.UUID) *schema.ServiceResponse[policySchema.RenewalResponse] {
	policyID, err := uuid.Parse(req.PolicyID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.RenewalResponse](http.StatusBadRequest, "Invalid policy ID", err)
	}

	pol, err := s.policyRepo.GetByID(ctx, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.RenewalResponse](http.StatusNotFound, "Policy not found", err)
	}

	if pol.Status != string(shared.PolicyStatusActive) {
		return schema.NewServiceErrorResponse[policySchema.RenewalResponse](http.StatusBadRequest, "Can only renew ACTIVE policies", nil)
	}

	renewalDate, err := utils.ParseFlexibleDate(req.RenewalDate)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.RenewalResponse](http.StatusBadRequest, "Invalid renewal date format (YYYY-MM-DD or ISO 8601)", err)
	}

	var expiresAt *time.Time
	if req.ExpiresAt != "" {
		t, err := utils.ParseFlexibleDate(req.ExpiresAt)
		if err != nil {
			return schema.NewServiceErrorResponse[policySchema.RenewalResponse](http.StatusBadRequest, "Invalid expires_at format (YYYY-MM-DD or ISO 8601)", err)
		}
		expiresAt = &t
	}

	var newPlanID uuid.UUID
	if req.NewPlanID != "" {
		newPlanID, err = uuid.Parse(req.NewPlanID)
		if err != nil {
			return schema.NewServiceErrorResponse[policySchema.RenewalResponse](http.StatusBadRequest, "Invalid new plan ID", err)
		}
	}

	renewal := &entity.PolicyRenewal{
		PolicyID:    policyID,
		Status:      string(shared.RenewalStatusPending),
		RenewalDate: renewalDate,
		NewPremium:  pol.PremiumAmount,
		NewPlanID:   newPlanID,
		ExpiresAt:   expiresAt,
		CreatedBy:   createdBy,
	}

	created, err := s.renewalRepo.Create(ctx, renewal)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.RenewalResponse](http.StatusInternalServerError, "Failed to initiate renewal", err)
	}

	s.logAudit(ctx, createdBy, string(shared.AuditEntityTypeRenewal), created.ID, string(shared.AuditActionCreate))

	// Auto-generate renewal notice (non-blocking)
	if s.policyDocSvc != nil {
		go func() {
			bgCtx := context.Background()
			s.policyDocSvc.GenerateDocument(bgCtx, policySchema.GenerateDocumentRequest{
				EntityType:     string(shared.DocumentEntityTypeRenewal),
				EntityID:       created.ID.String(),
				DocumentType:   string(shared.PolicyDocumentTypeRenewalNotice),
				GenerationMode: string(shared.GenerationModeAuto),
				GeneratedBy:    createdBy,
			})
		}()
	}

	return schema.NewServiceResponse(policySchema.ToRenewalResponse(created), http.StatusCreated, "Renewal initiated")
}

func (s *renewalServiceImpl) GetRenewal(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[policySchema.RenewalResponse] {
	renewal, err := s.renewalRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.RenewalResponse](http.StatusNotFound, "Renewal not found", err)
	}
	return schema.NewServiceResponse(policySchema.ToRenewalResponse(renewal), http.StatusOK, "Renewal retrieved")
}

func (s *renewalServiceImpl) ListByPolicy(ctx context.Context, policyID uuid.UUID) *schema.ServiceResponse[[]policySchema.RenewalResponse] {
	// Get the latest renewal for this policy
	renewal, err := s.renewalRepo.GetByPolicyID(ctx, policyID)
	if err != nil {
		return schema.NewServiceResponse([]policySchema.RenewalResponse{}, http.StatusOK, "No renewals found")
	}
	responses := []policySchema.RenewalResponse{policySchema.ToRenewalResponse(renewal)}
	return schema.NewServiceResponse(responses, http.StatusOK, "Renewals retrieved")
}

func (s *renewalServiceImpl) ApproveRenewal(ctx context.Context, id uuid.UUID, approvedBy uuid.UUID) *schema.ServiceResponse[policySchema.RenewalResponse] {
	renewal, err := s.renewalRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.RenewalResponse](http.StatusNotFound, "Renewal not found", err)
	}

	if renewal.Status != string(shared.RenewalStatusPending) {
		return schema.NewServiceErrorResponse[policySchema.RenewalResponse](http.StatusBadRequest, fmt.Sprintf("Cannot approve renewal in %s status", renewal.Status), nil)
	}

	now := time.Now()
	renewal.Status = string(shared.RenewalStatusApproved)
	renewal.ApprovedBy = approvedBy
	renewal.ApprovedAt = &now

	updated, err := s.renewalRepo.Update(ctx, renewal)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.RenewalResponse](http.StatusInternalServerError, "Failed to approve renewal", err)
	}

	s.logAudit(ctx, approvedBy, string(shared.AuditEntityTypeRenewal), id, string(shared.AuditActionStateChange))

	return schema.NewServiceResponse(policySchema.ToRenewalResponse(updated), http.StatusOK, "Renewal approved")
}

func (s *renewalServiceImpl) RejectRenewal(ctx context.Context, id uuid.UUID, reason string) *schema.ServiceResponse[policySchema.RenewalResponse] {
	renewal, err := s.renewalRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.RenewalResponse](http.StatusNotFound, "Renewal not found", err)
	}

	if renewal.Status != string(shared.RenewalStatusPending) {
		return schema.NewServiceErrorResponse[policySchema.RenewalResponse](http.StatusBadRequest, fmt.Sprintf("Cannot reject renewal in %s status", renewal.Status), nil)
	}

	renewal.Status = string(shared.RenewalStatusRejected)
	renewal.PremiumChangeReason = reason

	updated, err := s.renewalRepo.Update(ctx, renewal)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.RenewalResponse](http.StatusInternalServerError, "Failed to reject renewal", err)
	}

	return schema.NewServiceResponse(policySchema.ToRenewalResponse(updated), http.StatusOK, "Renewal rejected")
}

func (s *renewalServiceImpl) CompleteRenewal(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[policySchema.RenewalResponse] {
	renewal, err := s.renewalRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.RenewalResponse](http.StatusNotFound, "Renewal not found", err)
	}

	if renewal.Status != string(shared.RenewalStatusApproved) {
		return schema.NewServiceErrorResponse[policySchema.RenewalResponse](http.StatusBadRequest, "Only approved renewals can be completed", nil)
	}

	originalPolicy, err := s.policyRepo.GetByID(ctx, renewal.PolicyID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.RenewalResponse](http.StatusInternalServerError, "Failed to get original policy", err)
	}

	// Create new policy for next year
	planID := originalPolicy.PlanID
	if renewal.NewPlanID != uuid.Nil {
		planID = renewal.NewPlanID
	}
	premium := originalPolicy.PremiumAmount
	if renewal.NewPremium > 0 {
		premium = renewal.NewPremium
	}

	// Claims experience loading — adjust premium based on loss ratio
	if s.claimRepo != nil {
		claims, claimErr := s.claimRepo.ListByPolicy(ctx, renewal.PolicyID, 10000, 0)
		if claimErr == nil && len(claims) > 0 {
			var totalApproved int64
			for _, c := range claims {
				totalApproved += c.ApprovedAmount
			}
			if originalPolicy.PremiumAmount > 0 {
				lossRatio := float64(totalApproved) / float64(originalPolicy.PremiumAmount) * 100
				var loading float64
				var reason string
				switch {
				case lossRatio > 100:
					loading = 0.25
					reason = fmt.Sprintf("Claims loading +25%% (loss ratio %.1f%%)", lossRatio)
				case lossRatio > 75:
					loading = 0.15
					reason = fmt.Sprintf("Claims loading +15%% (loss ratio %.1f%%)", lossRatio)
				case lossRatio > 50:
					loading = 0.10
					reason = fmt.Sprintf("Claims loading +10%% (loss ratio %.1f%%)", lossRatio)
				case lossRatio < 30:
					loading = -0.05
					reason = fmt.Sprintf("Good claims discount -5%% (loss ratio %.1f%%)", lossRatio)
				}
				if loading != 0 {
					premium += int64(float64(premium) * loading)
					renewal.PremiumChangeReason = reason
				}
			}
		}
	}

	// Recalculate premium with premium rules if available
	if s.premiumRuleSvc != nil {
		activeMembers, memErr := s.memberRepo.ListActiveByPolicy(ctx, originalPolicy.ID)
		if memErr == nil && len(activeMembers) > 0 {
			membersJSON, _ := json.Marshal(activeMembers)
			premResp := s.premiumRuleSvc.CalculatePremiumWithMembers(ctx, planID, len(activeMembers), membersJSON)
			if premResp.Error == nil && premResp.Data > 0 {
				premium = premResp.Data
			}
		}
	}

	newPolicy := &entity.Policy{
		PlanID:            planID,
		PolicyholderName:  originalPolicy.PolicyholderName,
		PolicyholderEmail: originalPolicy.PolicyholderEmail,
		PolicyholderPhone: originalPolicy.PolicyholderPhone,
		PolicyNumber:      utils.GeneratePolicyNumber(),
		Status:            string(shared.PolicyStatusDraft),
		StartDate:         originalPolicy.EndDate,
		EndDate:           originalPolicy.EndDate.AddDate(1, 0, 0),
		PremiumAmount:     premium,
		Currency:          originalPolicy.Currency,
		RenewedFromID:     &originalPolicy.ID,
		CreatedBy:         renewal.CreatedBy,
	}

	createdPolicy, err := s.policyRepo.Create(ctx, newPolicy)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.RenewalResponse](http.StatusInternalServerError, "Failed to create renewed policy", err)
	}

	// Get premium rules for age validation
	var premiumRules []*productEntity.PremiumRule
	if s.premiumRuleRepo != nil {
		if rules, rErr := s.premiumRuleRepo.ListEffectiveByPlan(ctx, planID, time.Now()); rErr == nil {
			premiumRules = rules
		}
	}

	// Copy active members to new policy (with re-validation)
	members, err := s.memberRepo.ListActiveByPolicy(ctx, originalPolicy.ID)
	if err != nil {
		log.Printf("Warning: Failed to copy members during renewal: %v", err)
	} else {
		for _, m := range members {
			// Overage re-validation
			memberAge := renewalCalculateMemberAge(m.DateOfBirth)
			if len(premiumRules) > 0 && !isAgeValidForPlan(premiumRules, m.Relationship, memberAge) {
				log.Printf("Warning: Skipping member %s during renewal — age %d outside allowed range for %s", m.Name, memberAge, m.Relationship)
				s.createRenewalFlag(ctx, createdPolicy.ID, m.ID,
					string(shared.UnderwritingFlagRenewalSkip), "MEDIUM",
					fmt.Sprintf("Member %s skipped during renewal: age %d outside allowed range for %s", m.Name, memberAge, m.Relationship))
				continue
			}

			// Double insurance re-validation
			if m.NationalID != "" {
				existing, natErr := s.memberRepo.GetByNationalID(ctx, m.NationalID)
				if natErr == nil && existing != nil && existing.PolicyID != originalPolicy.ID {
					existingPol, polErr := s.policyRepo.GetByID(ctx, existing.PolicyID)
					if polErr == nil && existingPol.Status == string(shared.PolicyStatusActive) && existingPol.ID != createdPolicy.ID {
						log.Printf("Warning: Skipping member %s during renewal — double insurance detected in policy %s", m.Name, existingPol.PolicyNumber)
						s.createRenewalFlag(ctx, createdPolicy.ID, m.ID,
							string(shared.UnderwritingFlagRenewalSkip), "HIGH",
							fmt.Sprintf("Member %s skipped during renewal: double insurance detected in policy %s", m.Name, existingPol.PolicyNumber))
						continue
					}
				}
			}

			newMember := &entity.Member{
				PolicyID:     createdPolicy.ID,
				NationalID:   m.NationalID,
				Name:         m.Name,
				DateOfBirth:  m.DateOfBirth,
				Gender:       m.Gender,
				Relationship: m.Relationship,
				MemberNumber: utils.GenerateMemberNumber(),
				Phone:        m.Phone,
				Email:        m.Email,
				KRAPin:       m.KRAPin,
				County:       m.County,
				Address:      m.Address,
				Status:       string(shared.MemberStatusActive),
				Verified:     m.Verified,
			}
			if _, err := s.memberRepo.Create(ctx, newMember); err != nil {
				log.Printf("Warning: Failed to copy member %s during renewal: %v", m.Name, err)
			}
		}
	}

	// Update renewal record
	now := time.Now()
	renewal.Status = string(shared.RenewalStatusCompleted)
	renewal.RenewedPolicyID = createdPolicy.ID
	renewal.CompletedAt = &now

	updated, err := s.renewalRepo.Update(ctx, renewal)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.RenewalResponse](http.StatusInternalServerError, "Failed to complete renewal", err)
	}

	s.logAudit(ctx, renewal.CreatedBy, string(shared.AuditEntityTypeRenewal), id, string(shared.AuditActionStateChange))

	return schema.NewServiceResponse(policySchema.ToRenewalResponse(updated), http.StatusOK, "Renewal completed — new policy created")
}

func (s *renewalServiceImpl) ExpirePendingRenewals(ctx context.Context) *schema.ServiceResponse[int] {
	expired, err := s.renewalRepo.ListExpired(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int](http.StatusInternalServerError, "Failed to find expired renewals", err)
	}

	count := 0
	for _, r := range expired {
		if _, err := s.renewalRepo.UpdateStatus(ctx, r.ID, string(shared.RenewalStatusExpired)); err != nil {
			log.Printf("Failed to expire renewal %s: %v", r.ID, err)
			continue
		}
		count++
	}

	return schema.NewServiceResponse(count, http.StatusOK, fmt.Sprintf("Expired %d renewals", count))
}

func (s *renewalServiceImpl) BulkInitiateRenewals(ctx context.Context, policyIDs []uuid.UUID, createdBy uuid.UUID) *schema.ServiceResponse[policySchema.BulkResultResponse] {
	result := policySchema.BulkResultResponse{}
	for _, policyID := range policyIDs {
		req := policySchema.InitiateRenewalRequest{
			PolicyID:    policyID.String(),
			RenewalDate: time.Now().AddDate(0, 0, 30).Format("2006-01-02"),
		}
		resp := s.InitiateRenewal(ctx, req, createdBy)
		if resp.Error != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Policy %s: %s", policyID, resp.Message))
		} else {
			result.Succeeded++
		}
	}

	return schema.NewServiceResponse(result, http.StatusOK, "Bulk renewals processed")
}

func renewalCalculateMemberAge(dob time.Time) int {
	now := time.Now()
	age := now.Year() - dob.Year()
	if now.YearDay() < dob.YearDay() {
		age--
	}
	return age
}

func isAgeValidForPlan(rules []*productEntity.PremiumRule, relationship string, age int) bool {
	for _, rule := range rules {
		if rule.Relationship == "" || strings.EqualFold(rule.Relationship, relationship) {
			if (rule.MinAge == 0 || age >= rule.MinAge) && (rule.MaxAge == 0 || age <= rule.MaxAge) {
				return true
			}
		}
	}
	return false
}

func (s *renewalServiceImpl) createRenewalFlag(ctx context.Context, policyID, memberID uuid.UUID, flagType, severity, details string) {
	if s.underwritingFlagRepo == nil {
		return
	}
	flag := &entity.UnderwritingFlag{
		PolicyID: policyID,
		MemberID: memberID,
		FlagType: flagType,
		Severity: severity,
		Details:  details,
		Status:   string(shared.UnderwritingFlagStatusOpen),
	}
	if _, err := s.underwritingFlagRepo.Create(ctx, flag); err != nil {
		log.Printf("Warning: failed to create renewal flag: %v", err)
	}
}

func (s *renewalServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}
