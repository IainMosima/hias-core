package claims

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	"github.com/bitbiz/hias-core/domains/claims/entity"
	claimRepo "github.com/bitbiz/hias-core/domains/claims/repository"
	claimsSchema "github.com/bitbiz/hias-core/domains/claims/schema"
	"github.com/bitbiz/hias-core/domains/claims/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/google/uuid"
)

type claimServiceImpl struct {
	claimRepo      claimRepo.ClaimRepository
	lineItemRepo   claimRepo.ClaimLineItemRepository
	adjudicatorSvc service.AdjudicatorService
	validatorSvc   service.ValidatorService
	fraudSvc       service.FraudService
	adjRepo        claimRepo.AdjudicationRepository
	fraudFlagRepo  claimRepo.FraudFlagRepository
	auditSvc       auditService.AuditService
}

func NewClaimService(
	claimRepo claimRepo.ClaimRepository,
	lineItemRepo claimRepo.ClaimLineItemRepository,
	adjudicatorSvc service.AdjudicatorService,
	validatorSvc service.ValidatorService,
	fraudSvc service.FraudService,
	adjRepo claimRepo.AdjudicationRepository,
	fraudFlagRepo claimRepo.FraudFlagRepository,
	auditSvc auditService.AuditService,
) service.ClaimService {
	return &claimServiceImpl{
		claimRepo:      claimRepo,
		lineItemRepo:   lineItemRepo,
		adjudicatorSvc: adjudicatorSvc,
		validatorSvc:   validatorSvc,
		fraudSvc:       fraudSvc,
		adjRepo:        adjRepo,
		fraudFlagRepo:  fraudFlagRepo,
		auditSvc:       auditSvc,
	}
}

func (s *claimServiceImpl) SubmitClaim(ctx context.Context, req claimsSchema.SubmitClaimRequest, createdBy uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse] {
	policyID, _ := uuid.Parse(req.PolicyID)
	memberID, _ := uuid.Parse(req.MemberID)
	providerID, _ := uuid.Parse(req.ProviderID)

	diagJSON, _ := json.Marshal(req.DiagnosisCodes)

	// Calculate total amount from line items
	var totalAmount int64
	for _, li := range req.LineItems {
		totalAmount += li.UnitPrice * int64(li.Quantity)
	}

	claim := &entity.Claim{
		ClaimNumber:    utils.GenerateClaimNumber(),
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
		CreatedBy:      createdBy,
	}

	if req.PreAuthID != "" {
		preAuthID, _ := uuid.Parse(req.PreAuthID)
		claim.PreAuthID = preAuthID
	}

	created, err := s.claimRepo.Create(ctx, claim)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusInternalServerError, "Failed to create claim", err)
	}

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

	// Step 3: Adjudicate claim
	result, adjErr := s.adjudicatorSvc.Adjudicate(ctx, created, lineItems)
	if adjErr != nil {
		log.Printf("Adjudication error: %v", adjErr)
		s.logAudit(ctx, createdBy, string(shared.AuditEntityTypeClaim), created.ID, string(shared.AuditActionCreate))
		return schema.NewServiceResponse(claimsSchema.ToClaimResponse(created), http.StatusCreated, "Claim submitted, adjudication failed")
	}

	// Step 4: Store adjudication decision
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

	// Step 5: Map decision to claim status
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

	// Step 6: Update claim amounts
	coPayAmount := created.TotalAmount - result.PayableAmount - result.MemberResponsibility
	if coPayAmount < 0 {
		coPayAmount = 0
	}
	_, amtErr := s.claimRepo.UpdateAmounts(ctx, created.ID, result.PayableAmount, coPayAmount, result.MemberResponsibility)
	if amtErr != nil {
		log.Printf("Failed to update claim amounts: %v", amtErr)
	}

	// Step 7: Update claim status
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

	// Step 8: Run additional fraud checks and store flags
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

func (s *claimServiceImpl) ApproveClaim(ctx context.Context, id uuid.UUID, approvedBy uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse] {
	// Fetch claim and validate status
	claim, err := s.claimRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusNotFound, "Claim not found", err)
	}

	if claim.Status != string(shared.ClaimStatusAdjudicated) && claim.Status != string(shared.ClaimStatusManualReview) {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Cannot approve claim in %s status; must be ADJUDICATED or MANUAL_REVIEW", claim.Status),
			nil,
		)
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

	s.logAudit(ctx, approvedBy, string(shared.AuditEntityTypeClaim), id, string(shared.AuditActionStateChange))

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

	s.logAudit(ctx, rejectedBy, string(shared.AuditEntityTypeClaim), id, string(shared.AuditActionStateChange))

	return schema.NewServiceResponse(claimsSchema.ToClaimResponse(updated), http.StatusOK, "Claim rejected")
}

func (s *claimServiceImpl) GetTotalCount(ctx context.Context) *schema.ServiceResponse[int64] {
	count, err := s.claimRepo.Count(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to get count", err)
	}
	return schema.NewServiceResponse(count, http.StatusOK, "Count retrieved")
}

func (s *claimServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}
