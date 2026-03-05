package policy

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"encoding/json"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	billingService "github.com/bitbiz/hias-core/domains/billing/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/domains/policy/entity"
	"github.com/bitbiz/hias-core/domains/policy/repository"
	policyRepo "github.com/bitbiz/hias-core/domains/policy/repository"
	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/bitbiz/hias-core/domains/policy/service"
	productRepo "github.com/bitbiz/hias-core/domains/product/repository"
	productService "github.com/bitbiz/hias-core/domains/product/service"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/google/uuid"
)

type memberServiceImpl struct {
	memberRepo           repository.MemberRepository
	policyRepo           policyRepo.PolicyRepository
	planRepo             productRepo.PlanRepository
	premiumRuleRepo      productRepo.PremiumRuleRepository
	premiumRuleSvc       productService.PremiumRuleService
	underwritingFlagRepo repository.UnderwritingFlagRepository
	creditNoteSvc        billingService.CreditNoteService
	auditSvc             auditService.AuditService
}

func NewMemberService(
	memberRepo repository.MemberRepository,
	policyRepo policyRepo.PolicyRepository,
	planRepo productRepo.PlanRepository,
	premiumRuleRepo productRepo.PremiumRuleRepository,
	premiumRuleSvc productService.PremiumRuleService,
	underwritingFlagRepo repository.UnderwritingFlagRepository,
	creditNoteSvc billingService.CreditNoteService,
	auditSvc auditService.AuditService,
) service.MemberService {
	return &memberServiceImpl{
		memberRepo:           memberRepo,
		policyRepo:           policyRepo,
		planRepo:             planRepo,
		premiumRuleRepo:      premiumRuleRepo,
		premiumRuleSvc:       premiumRuleSvc,
		underwritingFlagRepo: underwritingFlagRepo,
		creditNoteSvc:        creditNoteSvc,
		auditSvc:             auditSvc,
	}
}

func (s *memberServiceImpl) EnrollMember(ctx context.Context, policyID uuid.UUID, req policySchema.EnrollMemberRequest) *schema.ServiceResponse[policySchema.MemberResponse] {
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

	// Underwriting: double insurance check
	if req.NationalID != "" {
		existing, err := s.memberRepo.GetByNationalID(ctx, req.NationalID)
		if err == nil && existing != nil {
			existingPol, polErr := s.policyRepo.GetByID(ctx, existing.PolicyID)
			if polErr == nil && existingPol.ID != policyID && existingPol.Status == string(shared.PolicyStatusActive) {
				s.createFlag(ctx, policyID, uuid.Nil, uuid.Nil,
					string(shared.UnderwritingFlagDoubleInsurance), "HIGH",
					fmt.Sprintf("Double insurance: NationalID %s already covered under policy %s", req.NationalID, existingPol.PolicyNumber))
				return schema.NewServiceErrorResponse[policySchema.MemberResponse](http.StatusBadRequest, "Double insurance detected: member already covered under another active policy", nil)
			}
		}
	}

	// Underwriting: overage check
	if err := s.validateMemberAge(ctx, pol.PlanID, req.Relationship, dob); err != nil {
		s.createFlag(ctx, policyID, uuid.Nil, uuid.Nil,
			string(shared.UnderwritingFlagMaxAge), "HIGH",
			fmt.Sprintf("Age violation for %s: %s", req.Name, err.Error()))
		return schema.NewServiceErrorResponse[policySchema.MemberResponse](http.StatusBadRequest, err.Error(), nil)
	}

	member := &entity.Member{
		PolicyID:     policyID,
		NationalID:   req.NationalID,
		Name:         req.Name,
		DateOfBirth:  dob,
		Gender:       req.Gender,
		Relationship: req.Relationship,
		MemberNumber: utils.GenerateMemberNumber(),
		Phone:        req.Phone,
		Email:        req.Email,
		KRAPin:       req.KRAPin,
		County:       req.County,
		Address:      req.Address,
		Status:       string(shared.MemberStatusActive),
		Verified:     false,
	}

	created, err := s.memberRepo.Create(ctx, member)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.MemberResponse](http.StatusInternalServerError, "Failed to enroll member", err)
	}

	s.logAudit(ctx, uuid.Nil, string(shared.AuditEntityTypeMember), created.ID, string(shared.AuditActionCreate))

	return schema.NewServiceResponse(policySchema.ToMemberResponse(created), http.StatusCreated, "Member enrolled")
}

func (s *memberServiceImpl) GetMember(ctx context.Context, memberID uuid.UUID) *schema.ServiceResponse[policySchema.MemberResponse] {
	member, err := s.memberRepo.GetByID(ctx, memberID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.MemberResponse](http.StatusNotFound, "Member not found", err)
	}
	return schema.NewServiceResponse(policySchema.ToMemberResponse(member), http.StatusOK, "Member retrieved")
}

func (s *memberServiceImpl) UpdateMember(ctx context.Context, memberID uuid.UUID, req policySchema.UpdateMemberRequest) *schema.ServiceResponse[policySchema.MemberResponse] {
	member, err := s.memberRepo.GetByID(ctx, memberID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.MemberResponse](http.StatusNotFound, "Member not found", err)
	}

	if member.Status == string(shared.MemberStatusRemoved) {
		return schema.NewServiceErrorResponse[policySchema.MemberResponse](http.StatusBadRequest, "Cannot update a removed member", nil)
	}

	if req.Name != nil {
		member.Name = *req.Name
	}
	if req.Phone != nil {
		member.Phone = *req.Phone
	}
	if req.Email != nil {
		member.Email = *req.Email
	}
	if req.KRAPin != nil {
		member.KRAPin = *req.KRAPin
	}
	if req.County != nil {
		member.County = *req.County
	}
	if req.Address != nil {
		member.Address = *req.Address
	}

	updated, err := s.memberRepo.Update(ctx, member)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.MemberResponse](http.StatusInternalServerError, "Failed to update member", err)
	}

	s.logAudit(ctx, uuid.Nil, string(shared.AuditEntityTypeMember), memberID, string(shared.AuditActionUpdate))

	return schema.NewServiceResponse(policySchema.ToMemberResponse(updated), http.StatusOK, "Member updated")
}

func (s *memberServiceImpl) RemoveMember(ctx context.Context, memberID uuid.UUID, reason string) *schema.ServiceResponse[string] {
	member, err := s.memberRepo.GetByID(ctx, memberID)
	if err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusNotFound, "Member not found", err)
	}

	if member.Status == string(shared.MemberStatusRemoved) {
		return schema.NewServiceErrorResponse[string](http.StatusBadRequest, "Member is already removed", nil)
	}

	pol, err := s.policyRepo.GetByID(ctx, member.PolicyID)
	if err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusInternalServerError, "Failed to get policy", err)
	}

	if pol.Status != string(shared.PolicyStatusActive) && pol.Status != string(shared.PolicyStatusDraft) {
		return schema.NewServiceErrorResponse[string](http.StatusBadRequest, fmt.Sprintf("Cannot remove members from %s policy", pol.Status), nil)
	}

	// Capture old premium before recalc
	oldPremium := pol.PremiumAmount

	if _, err := s.memberRepo.UpdateStatus(ctx, memberID, string(shared.MemberStatusRemoved)); err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusInternalServerError, "Failed to remove member", err)
	}

	// Recalculate premium after member removal
	s.recalculatePolicyPremium(ctx, member.PolicyID, pol.PlanID)

	// Create credit note if premium decreased
	if s.creditNoteSvc != nil {
		updatedPol, polErr := s.policyRepo.GetByID(ctx, member.PolicyID)
		if polErr == nil && updatedPol.PremiumAmount < oldPremium {
			totalDays := pol.EndDate.Sub(pol.StartDate).Hours() / 24
			remainingDays := pol.EndDate.Sub(time.Now()).Hours() / 24
			if totalDays > 0 && remainingDays > 0 {
				premiumDiff := oldPremium - updatedPol.PremiumAmount
				refundAmount := int64(float64(premiumDiff) * remainingDays / totalDays)
				if refundAmount > 0 {
					cnReason := fmt.Sprintf("Pro-rata refund for member removal: %s", member.Name)
					if reason != "" {
						cnReason = fmt.Sprintf("%s (reason: %s)", cnReason, reason)
					}
					s.creditNoteSvc.CreateCreditNote(ctx, member.PolicyID, memberID, refundAmount, cnReason, uuid.Nil)
				}
			}
		}
	}

	s.logAudit(ctx, uuid.Nil, string(shared.AuditEntityTypeMember), memberID, string(shared.AuditActionDelete))

	msg := "Member removed"
	if reason != "" {
		msg = fmt.Sprintf("Member removed: %s", reason)
	}
	return schema.NewServiceResponse(msg, http.StatusOK, msg)
}

func (s *memberServiceImpl) SuspendMember(ctx context.Context, memberID uuid.UUID) *schema.ServiceResponse[policySchema.MemberResponse] {
	member, err := s.memberRepo.GetByID(ctx, memberID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.MemberResponse](http.StatusNotFound, "Member not found", err)
	}

	if member.Status != string(shared.MemberStatusActive) {
		return schema.NewServiceErrorResponse[policySchema.MemberResponse](http.StatusBadRequest, fmt.Sprintf("Cannot suspend member in %s status", member.Status), nil)
	}

	updated, err := s.memberRepo.UpdateStatus(ctx, memberID, string(shared.MemberStatusSuspended))
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.MemberResponse](http.StatusInternalServerError, "Failed to suspend member", err)
	}

	s.logAudit(ctx, uuid.Nil, string(shared.AuditEntityTypeMember), memberID, string(shared.AuditActionStateChange))

	return schema.NewServiceResponse(policySchema.ToMemberResponse(updated), http.StatusOK, "Member suspended")
}

func (s *memberServiceImpl) ReactivateMember(ctx context.Context, memberID uuid.UUID) *schema.ServiceResponse[policySchema.MemberResponse] {
	member, err := s.memberRepo.GetByID(ctx, memberID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.MemberResponse](http.StatusNotFound, "Member not found", err)
	}

	if member.Status != string(shared.MemberStatusSuspended) {
		return schema.NewServiceErrorResponse[policySchema.MemberResponse](http.StatusBadRequest, "Only suspended members can be reactivated", nil)
	}

	updated, err := s.memberRepo.UpdateStatus(ctx, memberID, string(shared.MemberStatusActive))
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.MemberResponse](http.StatusInternalServerError, "Failed to reactivate member", err)
	}

	s.logAudit(ctx, uuid.Nil, string(shared.AuditEntityTypeMember), memberID, string(shared.AuditActionStateChange))

	return schema.NewServiceResponse(policySchema.ToMemberResponse(updated), http.StatusOK, "Member reactivated")
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

	eligible := pol.Status == string(shared.PolicyStatusActive) &&
		member.Status == string(shared.MemberStatusActive) &&
		time.Now().Before(pol.EndDate)
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

func (s *memberServiceImpl) BulkEnrollMembers(ctx context.Context, policyID uuid.UUID, reqs []policySchema.EnrollMemberRequest) *schema.ServiceResponse[policySchema.BulkMemberResultResponse] {
	result := policySchema.BulkMemberResultResponse{}

	for _, req := range reqs {
		resp := s.EnrollMember(ctx, policyID, req)
		if resp.Error != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Member %s: %s", req.Name, resp.Message))
		} else {
			result.Succeeded++
			result.Members = append(result.Members, resp.Data)
		}
	}

	return schema.NewServiceResponse(result, http.StatusOK, fmt.Sprintf("Bulk enrollment: %d succeeded, %d failed", result.Succeeded, result.Failed))
}

func (s *memberServiceImpl) ImportMembersCSV(ctx context.Context, policyID uuid.UUID, csvData []byte) *schema.ServiceResponse[policySchema.BulkMemberResultResponse] {
	reader := csv.NewReader(bytes.NewReader(csvData))

	// Read header row
	header, err := reader.Read()
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.BulkMemberResultResponse](http.StatusBadRequest, "Failed to read CSV header", err)
	}

	// Build column index map
	colIndex := make(map[string]int)
	for i, col := range header {
		colIndex[strings.TrimSpace(strings.ToLower(col))] = i
	}

	// Required columns
	requiredCols := []string{"name", "date_of_birth", "gender", "relationship"}
	for _, col := range requiredCols {
		if _, ok := colIndex[col]; !ok {
			return schema.NewServiceErrorResponse[policySchema.BulkMemberResultResponse](
				http.StatusBadRequest,
				fmt.Sprintf("Missing required CSV column: %s", col),
				nil,
			)
		}
	}

	var reqs []policySchema.EnrollMemberRequest
	lineNum := 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		lineNum++
		if err != nil {
			return schema.NewServiceErrorResponse[policySchema.BulkMemberResultResponse](
				http.StatusBadRequest,
				fmt.Sprintf("CSV parse error at line %d", lineNum),
				err,
			)
		}

		req := policySchema.EnrollMemberRequest{
			Name:         getCSVField(record, colIndex, "name"),
			DateOfBirth:  getCSVField(record, colIndex, "date_of_birth"),
			Gender:       getCSVField(record, colIndex, "gender"),
			Relationship: getCSVField(record, colIndex, "relationship"),
			NationalID:   getCSVField(record, colIndex, "national_id"),
			Phone:        getCSVField(record, colIndex, "phone"),
			Email:        getCSVField(record, colIndex, "email"),
			KRAPin:       getCSVField(record, colIndex, "kra_pin"),
			County:       getCSVField(record, colIndex, "county"),
			Address:      getCSVField(record, colIndex, "address"),
		}
		reqs = append(reqs, req)
	}

	if len(reqs) == 0 {
		return schema.NewServiceErrorResponse[policySchema.BulkMemberResultResponse](http.StatusBadRequest, "CSV contains no data rows", nil)
	}

	return s.BulkEnrollMembers(ctx, policyID, reqs)
}

func (s *memberServiceImpl) BulkRemoveMembers(ctx context.Context, policyID uuid.UUID, memberIDs []uuid.UUID, reason string) *schema.ServiceResponse[policySchema.BulkMemberResultResponse] {
	result := policySchema.BulkMemberResultResponse{}

	for _, memberID := range memberIDs {
		resp := s.RemoveMember(ctx, memberID, reason)
		if resp.Error != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Member %s: %s", memberID, resp.Message))
		} else {
			result.Succeeded++
		}
	}

	return schema.NewServiceResponse(result, http.StatusOK, fmt.Sprintf("Bulk removal: %d succeeded, %d failed", result.Succeeded, result.Failed))
}

func (s *memberServiceImpl) validateMemberAge(ctx context.Context, planID uuid.UUID, relationship string, dob time.Time) error {
	if s.premiumRuleRepo == nil {
		return nil
	}

	rules, err := s.premiumRuleRepo.ListByPlan(ctx, planID)
	if err != nil || len(rules) == 0 {
		return nil
	}

	age := calculateMemberAge(dob)

	for _, rule := range rules {
		if rule.Relationship == "" || strings.EqualFold(rule.Relationship, relationship) {
			if (rule.MinAge == 0 || age >= rule.MinAge) && (rule.MaxAge == 0 || age <= rule.MaxAge) {
				return nil
			}
		}
	}

	return fmt.Errorf("member age %d outside allowed range for %s relationship", age, relationship)
}

func (s *memberServiceImpl) recalculatePolicyPremium(ctx context.Context, policyID, planID uuid.UUID) {
	if s.premiumRuleSvc == nil {
		return
	}

	members, err := s.memberRepo.ListActiveByPolicy(ctx, policyID)
	if err != nil {
		log.Printf("Warning: failed to list members for premium recalculation: %v", err)
		return
	}

	membersJSON, err := json.Marshal(members)
	if err != nil {
		log.Printf("Warning: failed to marshal members for premium recalculation: %v", err)
		return
	}

	premResp := s.premiumRuleSvc.CalculatePremiumWithMembers(ctx, planID, len(members), membersJSON)
	if premResp.Error != nil || premResp.Data <= 0 {
		log.Printf("Warning: premium recalculation returned no result")
		return
	}

	if _, err := s.policyRepo.UpdatePlanAndPremium(ctx, policyID, planID, premResp.Data); err != nil {
		log.Printf("Warning: failed to update policy premium after member removal: %v", err)
	}
}

func calculateMemberAge(dob time.Time) int {
	now := time.Now()
	age := now.Year() - dob.Year()
	if now.YearDay() < dob.YearDay() {
		age--
	}
	return age
}

func getCSVField(record []string, colIndex map[string]int, field string) string {
	if idx, ok := colIndex[field]; ok && idx < len(record) {
		return strings.TrimSpace(record[idx])
	}
	return ""
}

func (s *memberServiceImpl) createFlag(ctx context.Context, policyID, memberID, assessmentID uuid.UUID, flagType, severity, details string) {
	if s.underwritingFlagRepo == nil {
		return
	}
	flag := &entity.UnderwritingFlag{
		AssessmentID: assessmentID,
		PolicyID:     policyID,
		MemberID:     memberID,
		FlagType:     flagType,
		Severity:     severity,
		Details:      details,
		Status:       string(shared.UnderwritingFlagStatusOpen),
	}
	if _, err := s.underwritingFlagRepo.Create(ctx, flag); err != nil {
		log.Printf("Warning: failed to create underwriting flag: %v", err)
	}
}

func (s *memberServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}
