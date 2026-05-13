package claims

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	"github.com/bitbiz/hias-core/domains/claims/entity"
	claimRepo "github.com/bitbiz/hias-core/domains/claims/repository"
	claimsSchema "github.com/bitbiz/hias-core/domains/claims/schema"
	"github.com/bitbiz/hias-core/domains/claims/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	notifService "github.com/bitbiz/hias-core/domains/notification/service"
	preauthRepo "github.com/bitbiz/hias-core/domains/preauth/repository"
	salesRepo "github.com/bitbiz/hias-core/domains/sales/repository"
	"github.com/bitbiz/hias-core/infrastructures/queue"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/google/uuid"
)

var claimCounterInitOnce sync.Once

type claimServiceImpl struct {
	claimRepo          claimRepo.ClaimRepository
	lineItemRepo       claimRepo.ClaimLineItemRepository
	adjudicatorSvc     service.AdjudicatorService
	validatorSvc       service.ValidatorService
	fraudSvc           service.FraudService
	adjRepo            claimRepo.AdjudicationRepository
	fraudFlagRepo      claimRepo.FraudFlagRepository
	claimDocRepo       claimRepo.ClaimDocumentRepository
	preauthRepo        preauthRepo.PreAuthRepository
	auditSvc           auditService.AuditService
	notifSvc           notifService.NotificationService
	approvalLimitRepo  salesRepo.ApprovalLimitRepository
	escalationRuleRepo claimRepo.EscalationRuleRepository
	queueManager       queue.QueueManager
}

func NewClaimService(
	claimRepo claimRepo.ClaimRepository,
	lineItemRepo claimRepo.ClaimLineItemRepository,
	adjudicatorSvc service.AdjudicatorService,
	validatorSvc service.ValidatorService,
	fraudSvc service.FraudService,
	adjRepo claimRepo.AdjudicationRepository,
	fraudFlagRepo claimRepo.FraudFlagRepository,
	claimDocRepo claimRepo.ClaimDocumentRepository,
	preauthRepo preauthRepo.PreAuthRepository,
	auditSvc auditService.AuditService,
	notifSvc notifService.NotificationService,
	approvalLimitRepo salesRepo.ApprovalLimitRepository,
	escalationRuleRepo claimRepo.EscalationRuleRepository,
	queueManager queue.QueueManager,
) service.ClaimService {
	return &claimServiceImpl{
		claimRepo:          claimRepo,
		lineItemRepo:       lineItemRepo,
		adjudicatorSvc:     adjudicatorSvc,
		validatorSvc:       validatorSvc,
		fraudSvc:           fraudSvc,
		adjRepo:            adjRepo,
		fraudFlagRepo:      fraudFlagRepo,
		claimDocRepo:       claimDocRepo,
		preauthRepo:        preauthRepo,
		auditSvc:           auditSvc,
		notifSvc:           notifSvc,
		approvalLimitRepo:  approvalLimitRepo,
		escalationRuleRepo: escalationRuleRepo,
		queueManager:       queueManager,
	}
}

func (s *claimServiceImpl) SubmitClaim(ctx context.Context, req claimsSchema.SubmitClaimRequest, createdBy uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse] {
	log.Printf("[CLAIMS-SVC] SubmitClaim ENTRY: policy=%s member=%s provider=%s claim_type=%s line_items=%d diagnosis=%v created_by=%s",
		req.PolicyID, req.MemberID, req.ProviderID, req.ClaimType, len(req.LineItems), req.DiagnosisCodes, createdBy)

	// Prime the in-memory counter once per process from the DB max so we don't collide with seeded data.
	claimCounterInitOnce.Do(func() {
		maxCounter, err := s.claimRepo.GetMaxCounterForYear(ctx, time.Now().Year())
		if err != nil {
			log.Printf("[CLAIMS-SVC] Failed to prime claim counter (will start from 0): %v", err)
		} else {
			utils.InitClaimCounter(maxCounter)
			log.Printf("[CLAIMS-SVC] Claim counter primed to %d", maxCounter)
		}
	})

	policyID, _ := uuid.Parse(req.PolicyID)
	memberID, _ := uuid.Parse(req.MemberID)
	providerID, _ := uuid.Parse(req.ProviderID)

	diagJSON, _ := json.Marshal(req.DiagnosisCodes)

	// Calculate total amount from line items
	var totalAmount int64
	for _, li := range req.LineItems {
		totalAmount += li.UnitPrice * int64(li.Quantity)
	}
	log.Printf("[CLAIMS-SVC] total_amount=%d (from %d line items)", totalAmount, len(req.LineItems))

	// Set claim type (default DIRECT)
	claimType := string(shared.ClaimTypeDirect)
	if req.ClaimType != "" {
		claimType = req.ClaimType
	}

	// Calculate SLA breach time
	slaBreachAt := time.Now().Add(time.Duration(shared.ClaimSLAHours) * time.Hour)

	claim := &entity.Claim{
		PolicyID:       policyID,
		MemberID:       memberID,
		ProviderID:     providerID,
		Status:         string(shared.ClaimStatusReceived),
		TotalAmount:    totalAmount,
		DiagnosisCodes: diagJSON,
		ServiceDate:    req.ServiceDate,
		AdmissionDate:  req.AdmissionDate,
		DischargeDate:  req.DischargeDate,
		Notes:          req.Notes,
		ClaimType:      claimType,
		SLABreachAt:    &slaBreachAt,
		CreatedBy:      createdBy,
	}

	if req.PreAuthID != "" {
		preAuthID, _ := uuid.Parse(req.PreAuthID)
		claim.PreAuthID = preAuthID
	}

	// Retry on claim number collision (up to 3 attempts) — happens when in-memory counter
	// falls behind the DB (e.g. after a process restart with pre-seeded data).
	var created *entity.Claim
	var err error
	const maxRetries = 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		claim.ClaimNumber = utils.GenerateClaimNumber()
		log.Printf("[CLAIMS-SVC] Attempting claim create: attempt=%d claim_number=%s", attempt, claim.ClaimNumber)
		created, err = s.claimRepo.Create(ctx, claim)
		if err == nil {
			break
		}
		if errors.Is(err, claimRepo.ErrClaimNumberCollision) {
			log.Printf("[CLAIMS-SVC] Claim number collision on %s (attempt %d/%d), re-syncing counter",
				claim.ClaimNumber, attempt, maxRetries)
			if maxCounter, qErr := s.claimRepo.GetMaxCounterForYear(ctx, time.Now().Year()); qErr == nil {
				utils.ResetClaimCounterForCollision(maxCounter)
			}
			continue
		}
		// FK violation — bad policy/member/provider UUID — return 400
		var fkErr *claimRepo.ClaimFKViolationError
		if errors.As(err, &fkErr) {
			log.Printf("[CLAIMS-SVC] FK violation: constraint=%s detail=%s", fkErr.Constraint, fkErr.Detail)
			return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](
				http.StatusBadRequest,
				fmt.Sprintf("Invalid reference: check policy, member, or provider ID (%s)", fkErr.Constraint),
				err,
			)
		}
		// All other DB errors
		log.Printf("[CLAIMS-SVC] Create FAILED (attempt %d): %v", attempt, err)
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusInternalServerError, "Failed to create claim", err)
	}
	if err != nil {
		log.Printf("[CLAIMS-SVC] All %d create attempts exhausted, last error: %v", maxRetries, err)
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusInternalServerError, "Failed to create claim after retries", err)
	}
	log.Printf("[CLAIMS-SVC] Claim created: id=%s claim_number=%s status=%s", created.ID, created.ClaimNumber, created.Status)

	s.recordTimeline(ctx, created.ID, "", string(shared.ClaimStatusReceived), "Claim Submitted", "", createdBy)

	// Create line items
	for _, li := range req.LineItems {
		lineItem := &entity.ClaimLineItem{
			ClaimID:       created.ID,
			ProcedureCode: li.ProcedureCode,
			ProcedureName: li.ProcedureName,
			DiagnosisCode: li.DiagnosisCode,
			Quantity:      li.Quantity,
			UnitPrice:     li.UnitPrice,
			TotalPrice:    li.UnitPrice * int64(li.Quantity),
		}
		_, lineErr := s.lineItemRepo.Create(ctx, lineItem)
		if lineErr != nil {
			log.Printf("Failed to create line item: %v", lineErr)
		}
	}

	// --- Claims Pipeline ---

	// Step 1: Fetch line items for validation and adjudication
	lineItems, _ := s.lineItemRepo.ListByClaim(ctx, created.ID)

	// Step 2: Validate claim
	valid, validationErrors, _ := s.validatorSvc.ValidateClaim(ctx, created, lineItems)
	if !valid {
		reason := strings.Join(validationErrors, "; ")
		rejected, rejectErr := s.claimRepo.Reject(ctx, created.ID, reason)
		if rejectErr != nil {
			log.Printf("Failed to reject invalid claim: %v", rejectErr)
		} else {
			created = rejected
		}

		s.logAudit(ctx, createdBy, string(shared.AuditEntityTypeClaim), created.ID, string(shared.AuditActionCreate))
		return schema.NewServiceResponse(claimsSchema.ToClaimResponse(created), http.StatusCreated, "Claim submitted but rejected: "+reason)
	}

	// Update status to VALIDATED
	validated, err := s.claimRepo.UpdateStatus(ctx, created.ID, string(shared.ClaimStatusValidated))
	if err != nil {
		log.Printf("Failed to update claim status to VALIDATED: %v", err)
	} else {
		created = validated
	}

	// Step 3: Run fraud checks BEFORE adjudication so flags can influence the decision
	if len(lineItems) > 0 {
		firstItem := lineItems[0]

		isFrequent, _ := s.fraudSvc.CheckFrequency(ctx, memberID, providerID, firstItem.ProcedureCode, created.ID)
		if isFrequent {
			flag := &entity.FraudFlag{
				ClaimID:  created.ID,
				FlagType: string(shared.FraudFlagFrequency),
				Severity: string(shared.FraudSeverityMedium),
				Details:  fmt.Sprintf("High frequency of procedure %s for member", firstItem.ProcedureCode),
			}
			if flagErr := s.fraudSvc.FlagClaim(ctx, flag); flagErr != nil {
				log.Printf("Failed to flag claim for frequency: %v", flagErr)
			}
		}

		exceedsThreshold, _ := s.fraudSvc.CheckAmountThreshold(ctx, providerID, firstItem.ProcedureCode, totalAmount)
		if exceedsThreshold {
			flag := &entity.FraudFlag{
				ClaimID:  created.ID,
				FlagType: string(shared.FraudFlagAmountThreshold),
				Severity: string(shared.FraudSeverityHigh),
				Details:  fmt.Sprintf("Claim amount %d exceeds threshold for procedure %s", totalAmount, firstItem.ProcedureCode),
			}
			if flagErr := s.fraudSvc.FlagClaim(ctx, flag); flagErr != nil {
				log.Printf("Failed to flag claim for amount threshold: %v", flagErr)
			}
		}

		isExpired, _ := s.fraudSvc.CheckExpiredContract(ctx, providerID, req.ServiceDate)
		if isExpired {
			flag := &entity.FraudFlag{
				ClaimID:  created.ID,
				FlagType: string(shared.FraudFlagExpiredContract),
				Severity: string(shared.FraudSeverityHigh),
				Details:  "Provider has no valid contract covering the service date",
			}
			s.fraudSvc.FlagClaim(ctx, flag)
		}

		isSuspended, _ := s.fraudSvc.CheckSuspendedProvider(ctx, providerID)
		if isSuspended {
			flag := &entity.FraudFlag{
				ClaimID:  created.ID,
				FlagType: string(shared.FraudFlagSuspendedProvider),
				Severity: string(shared.FraudSeverityCritical),
				Details:  "Provider is suspended",
			}
			s.fraudSvc.FlagClaim(ctx, flag)
		}

		isOvercharge, _ := s.fraudSvc.CheckRateCardOvercharge(ctx, providerID, firstItem.ProcedureCode, firstItem.UnitPrice)
		if isOvercharge {
			flag := &entity.FraudFlag{
				ClaimID:  created.ID,
				FlagType: string(shared.FraudFlagRateCardOvercharge),
				Severity: string(shared.FraudSeverityMedium),
				Details:  fmt.Sprintf("Unit price exceeds rate card for procedure %s", firstItem.ProcedureCode),
			}
			s.fraudSvc.FlagClaim(ctx, flag)
		}

		isRepeat, _ := s.fraudSvc.CheckRepeatVisit(ctx, memberID, providerID, firstItem.ProcedureCode, req.ServiceDate, created.ID)
		if isRepeat {
			flag := &entity.FraudFlag{
				ClaimID:  created.ID,
				FlagType: string(shared.FraudFlagRepeatVisit),
				Severity: string(shared.FraudSeverityLow),
				Details:  fmt.Sprintf("Repeat visit detected for procedure %s", firstItem.ProcedureCode),
			}
			s.fraudSvc.FlagClaim(ctx, flag)
		}
	}

	// Step 3b: Check for critical fraud flags that should override adjudication
	hasCriticalFraud := false
	fraudFlags, _ := s.fraudFlagRepo.ListByClaim(ctx, created.ID)
	for _, f := range fraudFlags {
		if f.Severity == string(shared.FraudSeverityHigh) || f.Severity == string(shared.FraudSeverityCritical) {
			hasCriticalFraud = true
			break
		}
	}

	// Step 4: Adjudicate claim
	result, adjErr := s.adjudicatorSvc.Adjudicate(ctx, created, lineItems)
	if adjErr != nil {
		log.Printf("Adjudication error: %v", adjErr)
		s.logAudit(ctx, createdBy, string(shared.AuditEntityTypeClaim), created.ID, string(shared.AuditActionCreate))
		return schema.NewServiceResponse(claimsSchema.ToClaimResponse(created), http.StatusCreated, "Claim submitted, adjudication failed")
	}

	// Step 4b: Override to MANUAL_REVIEW if critical fraud detected on an approved claim
	if hasCriticalFraud && result.Decision == string(shared.AdjudicationDecisionApprove) {
		result.Decision = string(shared.AdjudicationDecisionManualReview)
		result.Reasons = append(result.Reasons, entity.RuleResult{
			Category: string(shared.RuleCategoryFraud),
			Rule:     "critical_fraud_override",
			Result:   string(shared.RuleResultFlag),
			Details:  "Adjudication overridden to MANUAL_REVIEW due to HIGH/CRITICAL fraud flags",
		})
	}

	// Step 5: Store adjudication decision
	reasonsJSON, _ := json.Marshal(result.Reasons)
	decision := &entity.AdjudicationDecision{
		ClaimID:              created.ID,
		Decision:             result.Decision,
		PayableAmount:        result.PayableAmount,
		MemberResponsibility: result.MemberResponsibility,
		Reasons:              reasonsJSON,
		RuleResults:          reasonsJSON,
		AdjudicatedAt:        time.Now(),
	}
	_, decisionErr := s.adjRepo.Create(ctx, decision)
	if decisionErr != nil {
		log.Printf("Failed to store adjudication decision: %v", decisionErr)
	}

	// Step 6: Map decision to claim status
	var newStatus string
	switch result.Decision {
	case string(shared.AdjudicationDecisionApprove):
		newStatus = string(shared.ClaimStatusAdjudicated)
	case string(shared.AdjudicationDecisionReject):
		newStatus = string(shared.ClaimStatusRejected)
	case string(shared.AdjudicationDecisionManualReview):
		newStatus = string(shared.ClaimStatusManualReview)
	default:
		newStatus = string(shared.ClaimStatusAdjudicated)
	}

	// Step 7: Update claim amounts
	coPayAmount := created.TotalAmount - result.PayableAmount - result.MemberResponsibility
	if coPayAmount < 0 {
		coPayAmount = 0
	}
	_, amtErr := s.claimRepo.UpdateAmounts(ctx, created.ID, result.PayableAmount, coPayAmount, result.MemberResponsibility)
	if amtErr != nil {
		log.Printf("Failed to update claim amounts: %v", amtErr)
	}

	// Step 8: Update claim status
	if newStatus == string(shared.ClaimStatusRejected) {
		rejReason := "Adjudication rejected"
		if len(result.Reasons) > 0 {
			var failReasons []string
			for _, r := range result.Reasons {
				if r.Result == string(shared.RuleResultFail) {
					failReasons = append(failReasons, r.Details)
				}
			}
			if len(failReasons) > 0 {
				rejReason = strings.Join(failReasons, "; ")
			}
		}
		s.claimRepo.Reject(ctx, created.ID, rejReason)
	} else {
		s.claimRepo.UpdateStatus(ctx, created.ID, newStatus)
	}

	// Step 9: Evaluate escalation rules
	if s.escalationRuleRepo != nil && newStatus != string(shared.ClaimStatusRejected) {
		s.evaluateEscalationRules(ctx, created, fraudFlags, newStatus)
	}

	// Step 10: Publish ClaimSubmittedEvent
	go func() {
		if s.queueManager != nil {
			event := map[string]interface{}{
				"event":        "ClaimSubmitted",
				"claim_id":     created.ID.String(),
				"claim_number": created.ClaimNumber,
				"status":       newStatus,
				"amount":       created.TotalAmount,
			}
			eventJSON, _ := json.Marshal(event)
			if err := s.queueManager.Publish(context.Background(), queue.TopicClaimProcessing, eventJSON); err != nil {
				log.Printf("Failed to publish ClaimSubmittedEvent: %v", err)
			}
		}
	}()

	// Step 11: Update PreAuth status to CLAIMED if preauth_id was provided
	if claim.PreAuthID != uuid.Nil && s.preauthRepo != nil {
		s.preauthRepo.UpdateStatus(ctx, claim.PreAuthID, string(shared.PreAuthStatusClaimed))
	}

	// Re-fetch claim to get updated amounts/status
	finalClaim, err := s.claimRepo.GetByID(ctx, created.ID)
	if err != nil {
		log.Printf("Failed to re-fetch claim: %v", err)
		finalClaim = created
	}

	s.logAudit(ctx, createdBy, string(shared.AuditEntityTypeClaim), created.ID, string(shared.AuditActionCreate))

	return schema.NewServiceResponse(claimsSchema.ToClaimResponse(finalClaim), http.StatusCreated, "Claim submitted and processed")
}

func (s *claimServiceImpl) GetClaim(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse] {
	claim, err := s.claimRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusNotFound, "Claim not found", err)
	}

	response := claimsSchema.ToClaimResponse(claim)

	// Load line items
	lineItems, _ := s.lineItemRepo.ListByClaim(ctx, id)
	if lineItems != nil {
		liResponses := make([]claimsSchema.LineItemResponse, len(lineItems))
		for i, li := range lineItems {
			liResponses[i] = claimsSchema.LineItemResponse{
				ID: li.ID, ProcedureCode: li.ProcedureCode, ProcedureName: li.ProcedureName,
				DiagnosisCode: li.DiagnosisCode, Quantity: li.Quantity, UnitPrice: li.UnitPrice,
				TotalPrice: li.TotalPrice, ApprovedAmount: li.ApprovedAmount,
			}
		}
		response.LineItems = liResponses
	}

	// Load adjudication decision
	decision, _ := s.adjRepo.GetByClaimID(ctx, id)
	if decision != nil {
		response.Decision = &claimsSchema.AdjudicationResponse{
			Decision:             decision.Decision,
			PayableAmount:        decision.PayableAmount,
			MemberResponsibility: decision.MemberResponsibility,
			Reasons:              decision.Reasons,
			RuleResults:          decision.RuleResults,
			AdjudicatedAt:        decision.AdjudicatedAt,
		}
	}

	// Load fraud flags
	flags, _ := s.fraudFlagRepo.ListByClaim(ctx, id)
	if flags != nil {
		flagResponses := make([]claimsSchema.FraudFlagResponse, len(flags))
		for i, f := range flags {
			flagResponses[i] = claimsSchema.FraudFlagResponse{
				ID: f.ID, FlagType: f.FlagType, Severity: f.Severity,
				Details: f.Details, Resolved: f.Resolved,
			}
		}
		response.FraudFlags = flagResponses
	}

	return schema.NewServiceResponse(response, http.StatusOK, "Claim retrieved")
}

func (s *claimServiceImpl) ListClaims(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]claimsSchema.ClaimResponse] {
	offset := (page - 1) * pageSize
	claims, err := s.claimRepo.List(ctx, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]claimsSchema.ClaimResponse](http.StatusInternalServerError, "Failed to list claims", err)
	}

	responses := make([]claimsSchema.ClaimResponse, len(claims))
	for i, c := range claims {
		responses[i] = claimsSchema.ToClaimResponse(c)
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Claims retrieved")
}

func (s *claimServiceImpl) ListClaimsByStatus(ctx context.Context, status string, page, pageSize int) *schema.ServiceResponse[[]claimsSchema.ClaimResponse] {
	offset := (page - 1) * pageSize
	claims, err := s.claimRepo.ListByStatus(ctx, status, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]claimsSchema.ClaimResponse](http.StatusInternalServerError, "Failed to list claims by status", err)
	}

	responses := make([]claimsSchema.ClaimResponse, len(claims))
	for i, c := range claims {
		responses[i] = claimsSchema.ToClaimResponse(c)
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Claims retrieved")
}

func (s *claimServiceImpl) ListClaimsByProvider(ctx context.Context, providerID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]claimsSchema.ClaimResponse] {
	offset := (page - 1) * pageSize
	claims, err := s.claimRepo.ListByProvider(ctx, providerID, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]claimsSchema.ClaimResponse](http.StatusInternalServerError, "Failed to list claims by provider", err)
	}

	responses := make([]claimsSchema.ClaimResponse, len(claims))
	for i, c := range claims {
		responses[i] = claimsSchema.ToClaimResponse(c)
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Claims retrieved")
}

func (s *claimServiceImpl) ApproveClaim(ctx context.Context, id uuid.UUID, approvedBy uuid.UUID, approverRole string) *schema.ServiceResponse[claimsSchema.ClaimResponse] {
	// Fetch claim and validate status
	claim, err := s.claimRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusNotFound, "Claim not found", err)
	}

	if claim.Status != string(shared.ClaimStatusAdjudicated) && claim.Status != string(shared.ClaimStatusManualReview) && claim.Status != string(shared.ClaimStatusEscalated) {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Cannot approve claim in %s status; must be ADJUDICATED, MANUAL_REVIEW, or ESCALATED", claim.Status),
			nil,
		)
	}

	// Check approval limits if approverRole is provided
	if s.approvalLimitRepo != nil && approverRole != "" {
		limit, limErr := s.approvalLimitRepo.GetByRole(ctx, approverRole)
		if limErr == nil && limit != nil && limit.MaxClaimAmount > 0 {
			if claim.TotalAmount > limit.MaxClaimAmount {
				return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](
					http.StatusForbidden,
					fmt.Sprintf("Claim amount %d exceeds your approval limit of %d; escalate to %s", claim.TotalAmount, limit.MaxClaimAmount, limit.EscalationRole),
					nil,
				)
			}
		}
	}

	// Fetch adjudication decision and use its amounts
	decision, err := s.adjRepo.GetByClaimID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusInternalServerError, "Adjudication decision not found", err)
	}

	coPayAmount := claim.TotalAmount - decision.PayableAmount - decision.MemberResponsibility
	if coPayAmount < 0 {
		coPayAmount = 0
	}

	_, amtErr := s.claimRepo.UpdateAmounts(ctx, id, decision.PayableAmount, coPayAmount, decision.MemberResponsibility)
	if amtErr != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusInternalServerError, "Failed to update claim amounts", amtErr)
	}

	updated, err := s.claimRepo.UpdateStatus(ctx, id, string(shared.ClaimStatusApproved))
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusInternalServerError, "Failed to approve claim", err)
	}

	s.recordTimeline(ctx, id, claim.Status, string(shared.ClaimStatusApproved), "Claim Approved", "", approvedBy)
	s.logAudit(ctx, approvedBy, string(shared.AuditEntityTypeClaim), id, string(shared.AuditActionStateChange))
	s.notifyClaim(ctx, claim.MemberID, claim.ClaimNumber, "APPROVED")

	go func() {
		if s.queueManager != nil {
			event := map[string]interface{}{
				"event":    "ClaimApproved",
				"claim_id": id.String(),
				"amount":   claim.TotalAmount,
			}
			eventJSON, _ := json.Marshal(event)
			if err := s.queueManager.Publish(context.Background(), queue.TopicClaimProcessing, eventJSON); err != nil {
				log.Printf("Failed to publish ClaimApprovedEvent: %v", err)
			}
		}
	}()

	return schema.NewServiceResponse(claimsSchema.ToClaimResponse(updated), http.StatusOK, "Claim approved")
}

func (s *claimServiceImpl) RejectClaim(ctx context.Context, id uuid.UUID, reason string, rejectedBy uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse] {
	// Fetch claim and validate status
	claim, err := s.claimRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusNotFound, "Claim not found", err)
	}

	allowedStatuses := map[string]bool{
		string(shared.ClaimStatusReceived):     true,
		string(shared.ClaimStatusValidated):    true,
		string(shared.ClaimStatusAdjudicated):  true,
		string(shared.ClaimStatusManualReview): true,
	}
	if !allowedStatuses[claim.Status] {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Cannot reject claim in %s status", claim.Status),
			nil,
		)
	}

	updated, err := s.claimRepo.Reject(ctx, id, reason)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusInternalServerError, "Failed to reject claim", err)
	}

	s.recordTimeline(ctx, id, claim.Status, string(shared.ClaimStatusRejected), "Claim Rejected", reason, rejectedBy)
	s.logAudit(ctx, rejectedBy, string(shared.AuditEntityTypeClaim), id, string(shared.AuditActionStateChange))
	s.notifyClaim(ctx, claim.MemberID, claim.ClaimNumber, "REJECTED")

	return schema.NewServiceResponse(claimsSchema.ToClaimResponse(updated), http.StatusOK, "Claim rejected")
}

func (s *claimServiceImpl) GetTotalCount(ctx context.Context) *schema.ServiceResponse[int64] {
	count, err := s.claimRepo.Count(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to get count", err)
	}
	return schema.NewServiceResponse(count, http.StatusOK, "Count retrieved")
}

func (s *claimServiceImpl) VetClaim(ctx context.Context, id uuid.UUID, req claimsSchema.VetClaimRequest, vettedBy uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse] {
	claim, err := s.claimRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusNotFound, "Claim not found", err)
	}

	if claim.Status != string(shared.ClaimStatusAdjudicated) && claim.Status != string(shared.ClaimStatusApproved) {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Cannot vet claim in %s status; must be ADJUDICATED or APPROVED", claim.Status),
			nil,
		)
	}

	// Claim-type-specific vetting rules
	switch shared.ClaimType(claim.ClaimType) {
	case shared.ClaimTypeDirect:
		// Direct (provider-submitted): if inpatient (has admission date), require pre-authorization
		if claim.AdmissionDate != nil && claim.PreAuthID == uuid.Nil {
			return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](
				http.StatusBadRequest,
				"Inpatient direct claims require pre-authorization reference",
				nil,
			)
		}
	case shared.ClaimTypeReimbursement:
		// Reimbursement: verify vetted amount doesn't exceed total claimed
		if req.VettedAmount > claim.TotalAmount {
			return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](
				http.StatusBadRequest,
				"Vetted amount cannot exceed total claimed amount for reimbursement claims",
				nil,
			)
		}
	case shared.ClaimTypeException:
		// Exception claims: must be vetted by manager role (enforced by RBAC middleware)
		// Verify amount is reasonable (within 150% of approved)
		if claim.ApprovedAmount > 0 && req.VettedAmount > claim.ApprovedAmount*3/2 {
			return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](
				http.StatusBadRequest,
				"Exception claim vetted amount exceeds 150% of approved amount — requires manual override",
				nil,
			)
		}
	}

	// Determine status based on vetted vs approved amount
	vetStatus := string(shared.ClaimStatusVetted)
	if claim.ApprovedAmount > 0 && req.VettedAmount < claim.ApprovedAmount {
		vetStatus = string(shared.ClaimStatusPartiallyVetted)
	}

	updated, err := s.claimRepo.VetClaim(ctx, id, req.VettedAmount, vettedBy, vetStatus)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusInternalServerError, "Failed to vet claim", err)
	}

	s.recordTimeline(ctx, id, claim.Status, vetStatus, "Claim Vetted", "", vettedBy)
	s.logAudit(ctx, vettedBy, string(shared.AuditEntityTypeClaim), id, string(shared.AuditActionStateChange))
	return schema.NewServiceResponse(claimsSchema.ToClaimResponse(updated), http.StatusOK, "Claim vetted")
}

func (s *claimServiceImpl) MarkReadyForPayment(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse] {
	claim, err := s.claimRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusNotFound, "Claim not found", err)
	}

	if claim.Status != string(shared.ClaimStatusVetted) && claim.Status != string(shared.ClaimStatusPartiallyVetted) {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Cannot mark claim as ready for payment in %s status; must be VETTED or PARTIALLY_VETTED", claim.Status),
			nil,
		)
	}

	updated, err := s.claimRepo.MarkReadyForPayment(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusInternalServerError, "Failed to mark claim ready for payment", err)
	}

	s.recordTimeline(ctx, id, claim.Status, string(shared.ClaimStatusReadyForPayment), "Marked Ready for Payment", "", userID)
	s.logAudit(ctx, userID, string(shared.AuditEntityTypeClaim), id, string(shared.AuditActionStateChange))
	return schema.NewServiceResponse(claimsSchema.ToClaimResponse(updated), http.StatusOK, "Claim marked ready for payment")
}

func (s *claimServiceImpl) MarkPaid(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse] {
	claim, err := s.claimRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusNotFound, "Claim not found", err)
	}

	if claim.Status != string(shared.ClaimStatusReadyForPayment) {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Cannot mark claim as paid in %s status; must be READY_FOR_PAYMENT", claim.Status),
			nil,
		)
	}

	updated, err := s.claimRepo.UpdateStatus(ctx, id, string(shared.ClaimStatusPaid))
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusInternalServerError, "Failed to mark claim as paid", err)
	}

	s.recordTimeline(ctx, id, claim.Status, string(shared.ClaimStatusPaid), "Claim Paid", "", userID)
	s.logAudit(ctx, userID, string(shared.AuditEntityTypeClaim), id, string(shared.AuditActionStateChange))
	s.notifyClaim(ctx, claim.MemberID, claim.ClaimNumber, "PAID")
	return schema.NewServiceResponse(claimsSchema.ToClaimResponse(updated), http.StatusOK, "Claim marked as paid")
}

func (s *claimServiceImpl) MarkPartPaid(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse] {
	claim, err := s.claimRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusNotFound, "Claim not found", err)
	}

	if claim.Status != string(shared.ClaimStatusReadyForPayment) {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Cannot mark claim as part paid in %s status; must be READY_FOR_PAYMENT", claim.Status),
			nil,
		)
	}

	updated, err := s.claimRepo.UpdateStatus(ctx, id, string(shared.ClaimStatusPartPaid))
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusInternalServerError, "Failed to mark claim as part paid", err)
	}

	s.recordTimeline(ctx, id, claim.Status, string(shared.ClaimStatusPartPaid), "Claim Partially Paid", "", userID)
	s.logAudit(ctx, userID, string(shared.AuditEntityTypeClaim), id, string(shared.AuditActionStateChange))
	return schema.NewServiceResponse(claimsSchema.ToClaimResponse(updated), http.StatusOK, "Claim marked as part paid")
}

func (s *claimServiceImpl) BulkSubmitClaims(ctx context.Context, req claimsSchema.BulkSubmitClaimsRequest, createdBy uuid.UUID) *schema.ServiceResponse[[]claimsSchema.ClaimResponse] {
	var responses []claimsSchema.ClaimResponse
	for _, claimReq := range req.Claims {
		resp := s.SubmitClaim(ctx, claimReq, createdBy)
		if resp.Error != nil {
			log.Printf("Failed to submit claim in bulk: %v", resp.Error)
			continue
		}
		responses = append(responses, resp.Data)
	}
	return schema.NewServiceResponse(responses, http.StatusCreated, fmt.Sprintf("%d claims submitted", len(responses)))
}

func (s *claimServiceImpl) ListSLABreached(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]claimsSchema.ClaimResponse] {
	offset := (page - 1) * pageSize
	claims, err := s.claimRepo.ListSLABreached(ctx, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]claimsSchema.ClaimResponse](http.StatusInternalServerError, "Failed to list SLA breached claims", err)
	}

	responses := make([]claimsSchema.ClaimResponse, len(claims))
	for i, c := range claims {
		responses[i] = claimsSchema.ToClaimResponse(c)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "SLA breached claims retrieved")
}

func (s *claimServiceImpl) UploadClaimDocument(ctx context.Context, claimID uuid.UUID, fileName, fileType string, fileSize int64, s3Key string, uploadedBy uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimDocumentResponse] {
	// Verify claim exists
	_, err := s.claimRepo.GetByID(ctx, claimID)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimDocumentResponse](http.StatusNotFound, "Claim not found", err)
	}

	doc := &entity.ClaimDocument{
		ClaimID:    claimID,
		FileName:   fileName,
		FileType:   fileType,
		FileSize:   fileSize,
		S3Key:      s3Key,
		UploadedBy: uploadedBy,
	}

	created, err := s.claimDocRepo.Create(ctx, doc)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimDocumentResponse](http.StatusInternalServerError, "Failed to upload document", err)
	}

	s.logAudit(ctx, uploadedBy, string(shared.AuditEntityTypeClaimDocument), created.ID, string(shared.AuditActionCreate))
	return schema.NewServiceResponse(claimsSchema.ToClaimDocumentResponse(created), http.StatusCreated, "Document uploaded")
}

func (s *claimServiceImpl) ListClaimDocuments(ctx context.Context, claimID uuid.UUID) *schema.ServiceResponse[[]claimsSchema.ClaimDocumentResponse] {
	docs, err := s.claimDocRepo.ListByClaim(ctx, claimID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]claimsSchema.ClaimDocumentResponse](http.StatusInternalServerError, "Failed to list documents", err)
	}

	responses := make([]claimsSchema.ClaimDocumentResponse, len(docs))
	for i, d := range docs {
		responses[i] = claimsSchema.ToClaimDocumentResponse(d)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Documents retrieved")
}

func (s *claimServiceImpl) DeleteClaimDocument(ctx context.Context, docID uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimDocumentResponse] {
	deleted, err := s.claimDocRepo.SoftDelete(ctx, docID)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimDocumentResponse](http.StatusInternalServerError, "Failed to delete document", err)
	}
	return schema.NewServiceResponse(claimsSchema.ToClaimDocumentResponse(deleted), http.StatusOK, "Document deleted")
}

func (s *claimServiceImpl) ImportClaimsCSV(ctx context.Context, csvData []byte, createdBy uuid.UUID) *schema.ServiceResponse[claimsSchema.BulkClaimResultResponse] {
	reader := csv.NewReader(bytes.NewReader(csvData))

	header, err := reader.Read()
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.BulkClaimResultResponse](http.StatusBadRequest, "Failed to read CSV header", err)
	}

	colIndex := make(map[string]int)
	for i, col := range header {
		colIndex[strings.TrimSpace(strings.ToLower(col))] = i
	}

	requiredCols := []string{"policy_id", "member_id", "provider_id", "service_date", "procedure_code", "procedure_name", "quantity", "unit_price"}
	for _, col := range requiredCols {
		if _, ok := colIndex[col]; !ok {
			return schema.NewServiceErrorResponse[claimsSchema.BulkClaimResultResponse](
				http.StatusBadRequest,
				fmt.Sprintf("Missing required CSV column: %s", col),
				nil,
			)
		}
	}

	result := claimsSchema.BulkClaimResultResponse{}
	lineNum := 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		lineNum++
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Line %d: CSV parse error: %v", lineNum, err))
			continue
		}

		getField := func(field string) string {
			if idx, ok := colIndex[field]; ok && idx < len(record) {
				return strings.TrimSpace(record[idx])
			}
			return ""
		}

		serviceDate, err := utils.ParseFlexibleDate(getField("service_date"))
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Line %d: invalid service_date", lineNum))
			continue
		}

		qty, _ := strconv.Atoi(getField("quantity"))
		if qty <= 0 {
			qty = 1
		}
		unitPrice, _ := strconv.ParseInt(getField("unit_price"), 10, 64)

		claimType := getField("claim_type")
		if claimType == "" {
			claimType = string(shared.ClaimTypeDirect)
		}

		var diagCodes []string
		if dc := getField("diagnosis_code"); dc != "" {
			diagCodes = strings.Split(dc, ";")
		} else {
			diagCodes = []string{"UNSPECIFIED"}
		}

		req := claimsSchema.SubmitClaimRequest{
			PolicyID:       getField("policy_id"),
			MemberID:       getField("member_id"),
			ProviderID:     getField("provider_id"),
			PreAuthID:      getField("preauth_id"),
			DiagnosisCodes: diagCodes,
			ServiceDate:    serviceDate,
			Notes:          getField("notes"),
			ClaimType:      claimType,
			LineItems: []claimsSchema.LineItemRequest{
				{
					ProcedureCode: getField("procedure_code"),
					ProcedureName: getField("procedure_name"),
					DiagnosisCode: getField("diagnosis_code"),
					Quantity:      qty,
					UnitPrice:     unitPrice,
				},
			},
		}

		resp := s.SubmitClaim(ctx, req, createdBy)
		if resp.Error != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Line %d: %s", lineNum, resp.Message))
		} else {
			result.Succeeded++
			result.Claims = append(result.Claims, resp.Data)
		}
	}

	if result.Succeeded == 0 && result.Failed == 0 {
		return schema.NewServiceErrorResponse[claimsSchema.BulkClaimResultResponse](http.StatusBadRequest, "CSV contains no data rows", nil)
	}

	return schema.NewServiceResponse(result, http.StatusOK, fmt.Sprintf("CSV import: %d succeeded, %d failed", result.Succeeded, result.Failed))
}

func (s *claimServiceImpl) notifyClaim(ctx context.Context, memberID uuid.UUID, claimNumber, newStatus string) {
	if s.notifSvc == nil {
		return
	}
	subject := fmt.Sprintf("Claim %s — %s", claimNumber, newStatus)
	body := fmt.Sprintf("Your claim %s has been updated to status: %s", claimNumber, newStatus)
	resp := s.notifSvc.Send(ctx, memberID, string(shared.NotificationChannelInApp), string(shared.NotificationTypeClaim), subject, body)
	if resp.Error != nil {
		log.Printf("Failed to send claim notification: %v", resp.Error)
	}
}

func (s *claimServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}

func (s *claimServiceImpl) evaluateEscalationRules(ctx context.Context, claim *entity.Claim, fraudFlags []*entity.FraudFlag, currentStatus string) {
	rules, err := s.escalationRuleRepo.ListActive(ctx)
	if err != nil || len(rules) == 0 {
		return
	}

	for _, rule := range rules {
		shouldEscalate := false

		switch shared.EscalationConditionType(rule.ConditionType) {
		case shared.EscalationConditionAmountExceeds:
			shouldEscalate = claim.TotalAmount > rule.ThresholdAmount

		case shared.EscalationConditionFraudFlag:
			for _, f := range fraudFlags {
				if f.Severity == string(shared.FraudSeverityHigh) || f.Severity == string(shared.FraudSeverityCritical) {
					shouldEscalate = true
					break
				}
			}

		case shared.EscalationConditionManualReview:
			shouldEscalate = currentStatus == string(shared.ClaimStatusManualReview)
		}

		if shouldEscalate {
			s.claimRepo.UpdateStatus(ctx, claim.ID, string(shared.ClaimStatusEscalated))
			s.claimRepo.SetEscalatedTo(ctx, claim.ID, rule.EscalationRole)
			log.Printf("Claim %s escalated to %s (rule: %s)", claim.ClaimNumber, rule.EscalationRole, rule.Name)
			return
		}
	}
}

func (s *claimServiceImpl) GetClaimTimeline(ctx context.Context, claimID uuid.UUID) *schema.ServiceResponse[[]claimsSchema.ClaimTimelineEntry] {
	entries, err := s.claimRepo.ListTimeline(ctx, claimID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]claimsSchema.ClaimTimelineEntry](http.StatusInternalServerError, "Failed to get claim timeline", err)
	}

	result := make([]claimsSchema.ClaimTimelineEntry, len(entries))
	for i, e := range entries {
		result[i] = claimsSchema.ClaimTimelineEntry{
			ID:              e.ID,
			Action:          e.Action,
			FromStatus:      e.FromStatus,
			ToStatus:        e.ToStatus,
			PerformedBy:     e.PerformedBy,
			PerformedByName: e.PerformedByName,
			Notes:           e.Notes,
			CreatedAt:       e.CreatedAt,
		}
	}

	return schema.NewServiceResponse(result, http.StatusOK, "Claim timeline retrieved")
}

func (s *claimServiceImpl) ListClaimsFiltered(ctx context.Context, status string, dateFrom, dateTo *time.Time, search string, page, pageSize int) *schema.ServiceResponse[[]claimsSchema.ClaimResponse] {
	offset := (page - 1) * pageSize
	claims, err := s.claimRepo.ListFiltered(ctx, status, dateFrom, dateTo, search, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]claimsSchema.ClaimResponse](http.StatusInternalServerError, "Failed to list claims", err)
	}

	responses := make([]claimsSchema.ClaimResponse, len(claims))
	for i, c := range claims {
		responses[i] = claimsSchema.ToClaimResponse(c)
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Claims retrieved")
}

func (s *claimServiceImpl) CountClaimsFiltered(ctx context.Context, status string, dateFrom, dateTo *time.Time, search string) *schema.ServiceResponse[int64] {
	count, err := s.claimRepo.CountFiltered(ctx, status, dateFrom, dateTo, search)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to count claims", err)
	}
	return schema.NewServiceResponse(count, http.StatusOK, "Claims counted")
}

func (s *claimServiceImpl) recordTimeline(ctx context.Context, claimID uuid.UUID, fromStatus, toStatus, action, notes string, performedBy uuid.UUID) {
	if err := s.claimRepo.CreateStatusHistory(ctx, claimID, fromStatus, toStatus, action, notes, performedBy); err != nil {
		log.Printf("Failed to record timeline for claim %s: %v", claimID, err)
	}
}
