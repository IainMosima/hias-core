package policy

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/domains/policy/entity"
	"github.com/bitbiz/hias-core/domains/policy/repository"
	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/bitbiz/hias-core/domains/policy/service"
	"github.com/bitbiz/hias-core/shared"
	"github.com/google/uuid"
)

type underwritingServiceImpl struct {
	underwritingRepo     repository.UnderwritingRepository
	policyRepo           repository.PolicyRepository
	memberRepo           repository.MemberRepository
	underwritingRuleRepo repository.UnderwritingRuleRepository
	underwritingFlagRepo repository.UnderwritingFlagRepository
	auditSvc             auditService.AuditService
}

func NewUnderwritingService(
	underwritingRepo repository.UnderwritingRepository,
	policyRepo repository.PolicyRepository,
	memberRepo repository.MemberRepository,
	underwritingRuleRepo repository.UnderwritingRuleRepository,
	underwritingFlagRepo repository.UnderwritingFlagRepository,
	auditSvc auditService.AuditService,
) service.UnderwritingService {
	return &underwritingServiceImpl{
		underwritingRepo:     underwritingRepo,
		policyRepo:           policyRepo,
		memberRepo:           memberRepo,
		underwritingRuleRepo: underwritingRuleRepo,
		underwritingFlagRepo: underwritingFlagRepo,
		auditSvc:             auditSvc,
	}
}

func (s *underwritingServiceImpl) SubmitAssessment(ctx context.Context, req policySchema.SubmitAssessmentRequest, createdBy uuid.UUID) *schema.ServiceResponse[policySchema.UnderwritingResponse] {
	policyID, err := uuid.Parse(req.PolicyID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.UnderwritingResponse](http.StatusBadRequest, "Invalid policy ID", err)
	}

	pol, err := s.policyRepo.GetByID(ctx, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.UnderwritingResponse](http.StatusNotFound, "Policy not found", err)
	}

	var memberID uuid.UUID
	var member *entity.Member
	if req.MemberID != "" {
		memberID, err = uuid.Parse(req.MemberID)
		if err != nil {
			return schema.NewServiceErrorResponse[policySchema.UnderwritingResponse](http.StatusBadRequest, "Invalid member ID", err)
		}
		member, err = s.memberRepo.GetByID(ctx, memberID)
		if err != nil {
			return schema.NewServiceErrorResponse[policySchema.UnderwritingResponse](http.StatusNotFound, "Member not found", err)
		}
	}

	medDecl := json.RawMessage("{}")
	if req.MedicalDeclarations != nil {
		medDecl = req.MedicalDeclarations
	}

	assessment := &entity.UnderwritingAssessment{
		PolicyID:            policyID,
		MemberID:            memberID,
		Status:              string(shared.UnderwritingStatusPending),
		Questionnaire:       req.Questionnaire,
		MedicalDeclarations: medDecl,
		CreatedBy:           createdBy,
	}

	// Set policy underwriting status to PENDING
	if _, uwErr := s.policyRepo.UpdateUnderwritingStatus(ctx, policyID, string(shared.PolicyUWStatusPending)); uwErr != nil {
		log.Printf("Warning: failed to set policy underwriting status to PENDING: %v", uwErr)
	}

	created, err := s.underwritingRepo.Create(ctx, assessment)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.UnderwritingResponse](http.StatusInternalServerError, "Failed to submit assessment", err)
	}

	// Evaluate rules engine
	s.evaluateRules(ctx, created, pol, member, req.Questionnaire)

	// Re-fetch to get updated assessment
	updated, err := s.underwritingRepo.GetByID(ctx, created.ID)
	if err != nil {
		updated = created
	}

	s.logAudit(ctx, createdBy, string(shared.AuditEntityTypeUnderwriting), created.ID, string(shared.AuditActionCreate))

	return schema.NewServiceResponse(policySchema.ToUnderwritingResponse(updated), http.StatusCreated, "Assessment submitted")
}

func (s *underwritingServiceImpl) evaluateRules(ctx context.Context, assessment *entity.UnderwritingAssessment, pol *entity.Policy, member *entity.Member, questionnaire json.RawMessage) {
	if s.underwritingRuleRepo == nil || s.underwritingFlagRepo == nil {
		return
	}

	rules, err := s.underwritingRuleRepo.ListActiveByPlan(ctx, pol.PlanID)
	if err != nil || len(rules) == 0 {
		return
	}

	// Parse questionnaire
	var qData map[string]interface{}
	if questionnaire != nil {
		_ = json.Unmarshal(questionnaire, &qData)
	}
	if qData == nil {
		qData = make(map[string]interface{})
	}

	var memberAge int
	var memberRelationship string
	var memberNationalID string
	if member != nil {
		memberAge = underwritingCalculateMemberAge(member.DateOfBirth)
		memberRelationship = member.Relationship
		memberNationalID = member.NationalID
	}

	totalRiskScore := 0
	hasBlocker := false
	var flagDetails []map[string]interface{}

	for _, rule := range rules {
		// Skip if rule is for a specific relationship and member doesn't match
		if rule.Relationship != "" && !strings.EqualFold(rule.Relationship, memberRelationship) {
			continue
		}

		triggered := false
		details := ""

		switch rule.RuleType {
		case string(shared.UnderwritingRuleMaxAge):
			if member != nil {
				maxAge, _ := strconv.Atoi(rule.ParameterValue)
				if maxAge > 0 && memberAge > maxAge {
					triggered = true
					details = fmt.Sprintf("Member age %d exceeds max age %d for %s", memberAge, maxAge, rule.Relationship)
				}
			}
		case string(shared.UnderwritingRuleMinAge):
			if member != nil {
				minAge, _ := strconv.Atoi(rule.ParameterValue)
				if minAge > 0 && memberAge < minAge {
					triggered = true
					details = fmt.Sprintf("Member age %d below min age %d for %s", memberAge, minAge, rule.Relationship)
				}
			}
		case string(shared.UnderwritingRuleDoubleInsurance):
			if member != nil && memberNationalID != "" {
				existing, natErr := s.memberRepo.GetByNationalID(ctx, memberNationalID)
				if natErr == nil && existing != nil && existing.PolicyID != pol.ID {
					existingPol, polErr := s.policyRepo.GetByID(ctx, existing.PolicyID)
					if polErr == nil && existingPol.Status == string(shared.PolicyStatusActive) {
						triggered = true
						details = fmt.Sprintf("Double insurance: NationalID %s already covered under policy %s", memberNationalID, existingPol.PolicyNumber)
					}
				}
			}
		case string(shared.UnderwritingRulePreExisting):
			val, ok := qData[rule.ParameterKey]
			if ok {
				valStr := fmt.Sprintf("%v", val)
				if strings.EqualFold(valStr, rule.ParameterValue) || strings.EqualFold(valStr, "yes") || strings.EqualFold(valStr, "true") {
					triggered = true
					details = fmt.Sprintf("Pre-existing condition flagged: %s = %s", rule.ParameterKey, valStr)
				}
			}
		case string(shared.UnderwritingRuleBMIThreshold):
			bmiVal, ok := qData["bmi"]
			if ok {
				bmiStr := fmt.Sprintf("%v", bmiVal)
				bmi, parseErr := strconv.ParseFloat(bmiStr, 64)
				threshold, threshErr := strconv.ParseFloat(rule.ParameterValue, 64)
				if parseErr == nil && threshErr == nil && bmi > threshold {
					triggered = true
					details = fmt.Sprintf("BMI %.1f exceeds threshold %.1f", bmi, threshold)
				}
			}
		case string(shared.UnderwritingRuleWaitingPeriod):
			// Informational — always create a flag if the condition exists in questionnaire
			val, ok := qData[rule.ParameterKey]
			if ok {
				valStr := fmt.Sprintf("%v", val)
				if strings.EqualFold(valStr, "yes") || strings.EqualFold(valStr, "true") {
					triggered = true
					details = fmt.Sprintf("Waiting period applies: %s days for %s", rule.ParameterValue, rule.ParameterKey)
				}
			}
		}

		if triggered {
			// Create flag
			flag := &entity.UnderwritingFlag{
				AssessmentID: assessment.ID,
				PolicyID:     assessment.PolicyID,
				MemberID:     assessment.MemberID,
				FlagType:     rule.RuleType,
				Severity:     rule.Severity,
				Details:      details,
				Status:       string(shared.UnderwritingFlagStatusOpen),
			}
			if _, flagErr := s.underwritingFlagRepo.Create(ctx, flag); flagErr != nil {
				log.Printf("Warning: failed to create underwriting flag: %v", flagErr)
			}

			totalRiskScore += rule.RiskScoreWeight
			if rule.IsBlocking {
				hasBlocker = true
			}

			flagDetails = append(flagDetails, map[string]interface{}{
				"rule_type": rule.RuleType,
				"severity":  rule.Severity,
				"details":   details,
				"blocking":  rule.IsBlocking,
				"weight":    rule.RiskScoreWeight,
			})
		}
	}

	// Determine auto-decision
	var status string
	var reason string
	if hasBlocker {
		status = string(shared.UnderwritingStatusDeclined)
		reason = "Declined: blocking rule triggered"
	} else if totalRiskScore > shared.UnderwritingReferThreshold {
		status = string(shared.UnderwritingStatusDeclined)
		reason = fmt.Sprintf("Declined: risk score %d exceeds threshold %d", totalRiskScore, shared.UnderwritingReferThreshold)
	} else if totalRiskScore > shared.UnderwritingAutoApproveThreshold {
		status = string(shared.UnderwritingStatusRefer)
		reason = fmt.Sprintf("Referred: risk score %d exceeds auto-approve threshold %d", totalRiskScore, shared.UnderwritingAutoApproveThreshold)
	} else {
		status = string(shared.UnderwritingStatusApproved)
		reason = "Auto-approved: risk score within acceptable range"
	}

	// Update assessment with results
	now := time.Now()
	riskFlagsJSON, _ := json.Marshal(flagDetails)
	assessment.Status = status
	assessment.RiskScore = totalRiskScore
	assessment.RiskFlags = riskFlagsJSON
	assessment.DecisionReason = reason
	assessment.AssessedAt = &now

	if _, err := s.underwritingRepo.Update(ctx, assessment); err != nil {
		log.Printf("Warning: failed to update assessment with rules result: %v", err)
	}

	// Sync underwriting status to policy
	policyUWStatus := assessmentStatusToPolicyUWStatus(status)
	if _, err := s.policyRepo.UpdateUnderwritingStatus(ctx, assessment.PolicyID, policyUWStatus); err != nil {
		log.Printf("Warning: failed to sync underwriting status to policy: %v", err)
	}
}

func (s *underwritingServiceImpl) GetAssessment(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[policySchema.UnderwritingResponse] {
	assessment, err := s.underwritingRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.UnderwritingResponse](http.StatusNotFound, "Assessment not found", err)
	}
	return schema.NewServiceResponse(policySchema.ToUnderwritingResponse(assessment), http.StatusOK, "Assessment retrieved")
}

func (s *underwritingServiceImpl) ListByPolicy(ctx context.Context, policyID uuid.UUID) *schema.ServiceResponse[[]policySchema.UnderwritingResponse] {
	assessments, err := s.underwritingRepo.GetByPolicyID(ctx, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]policySchema.UnderwritingResponse](http.StatusInternalServerError, "Failed to list assessments", err)
	}

	responses := make([]policySchema.UnderwritingResponse, len(assessments))
	for i, a := range assessments {
		responses[i] = policySchema.ToUnderwritingResponse(a)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Assessments retrieved")
}

func (s *underwritingServiceImpl) ReviewAssessment(ctx context.Context, id uuid.UUID, req policySchema.ReviewAssessmentRequest, assessedBy uuid.UUID) *schema.ServiceResponse[policySchema.UnderwritingResponse] {
	assessment, err := s.underwritingRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.UnderwritingResponse](http.StatusNotFound, "Assessment not found", err)
	}

	if assessment.Status != string(shared.UnderwritingStatusPending) && assessment.Status != string(shared.UnderwritingStatusRefer) {
		return schema.NewServiceErrorResponse[policySchema.UnderwritingResponse](http.StatusBadRequest, fmt.Sprintf("Cannot review assessment in %s status", assessment.Status), nil)
	}

	validStatuses := map[string]bool{
		string(shared.UnderwritingStatusApproved):              true,
		string(shared.UnderwritingStatusDeclined):              true,
		string(shared.UnderwritingStatusRefer):                 true,
		string(shared.UnderwritingStatusApprovedWithLoading):   true,
		string(shared.UnderwritingStatusApprovedWithExclusion): true,
	}
	if !validStatuses[req.Status] {
		return schema.NewServiceErrorResponse[policySchema.UnderwritingResponse](http.StatusBadRequest, "Invalid status; must be APPROVED, DECLINED, REFER, APPROVED_WITH_LOADING, or APPROVED_WITH_EXCLUSION", nil)
	}

	now := time.Now()
	assessment.Status = req.Status
	assessment.RiskScore = req.RiskScore
	assessment.DecisionReason = req.DecisionReason
	assessment.AssessedBy = assessedBy
	assessment.AssessedAt = &now

	updated, err := s.underwritingRepo.Update(ctx, assessment)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.UnderwritingResponse](http.StatusInternalServerError, "Failed to review assessment", err)
	}

	// Sync underwriting status to policy
	policyUWStatus := assessmentStatusToPolicyUWStatus(req.Status)
	if _, uwErr := s.policyRepo.UpdateUnderwritingStatus(ctx, assessment.PolicyID, policyUWStatus); uwErr != nil {
		log.Printf("Warning: failed to sync underwriting status to policy on review: %v", uwErr)
	}

	s.logAudit(ctx, assessedBy, string(shared.AuditEntityTypeUnderwriting), id, string(shared.AuditActionUpdate))

	return schema.NewServiceResponse(policySchema.ToUnderwritingResponse(updated), http.StatusOK, "Assessment reviewed")
}

// assessmentStatusToPolicyUWStatus maps an assessment status to the corresponding policy underwriting status.
func assessmentStatusToPolicyUWStatus(assessmentStatus string) string {
	switch assessmentStatus {
	case string(shared.UnderwritingStatusApproved):
		return string(shared.PolicyUWStatusApproved)
	case string(shared.UnderwritingStatusDeclined):
		return string(shared.PolicyUWStatusDeclined)
	case string(shared.UnderwritingStatusRefer):
		return string(shared.PolicyUWStatusInReview)
	case string(shared.UnderwritingStatusApprovedWithLoading):
		return string(shared.PolicyUWStatusApprovedWithLoading)
	case string(shared.UnderwritingStatusApprovedWithExclusion):
		return string(shared.PolicyUWStatusApprovedWithExclusion)
	default:
		return string(shared.PolicyUWStatusPending)
	}
}

func underwritingCalculateMemberAge(dob time.Time) int {
	now := time.Now()
	age := now.Year() - dob.Year()
	if now.YearDay() < dob.YearDay() {
		age--
	}
	return age
}

func (s *underwritingServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}
