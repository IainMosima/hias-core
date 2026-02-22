package claims

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	"github.com/bitbiz/hias-core/domains/claims/repository"
	claimsSchema "github.com/bitbiz/hias-core/domains/claims/schema"
	"github.com/bitbiz/hias-core/domains/claims/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	policyDomainRepo "github.com/bitbiz/hias-core/domains/policy/repository"
	providerRepo "github.com/bitbiz/hias-core/domains/provider/repository"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/google/uuid"
)

type claimServiceImpl struct {
	claimRepo    repository.ClaimRepository
	lineItemRepo repository.ClaimLineItemRepository
	adjRepo      repository.AdjudicationRepository
	policyRepo   policyDomainRepo.PolicyRepository
	memberRepo   policyDomainRepo.MemberRepository
	providerRepo providerRepo.ProviderRepository
	validator    service.ValidatorService
	adjudicator  service.AdjudicatorService
}

func NewClaimService(
	claimRepo repository.ClaimRepository,
	lineItemRepo repository.ClaimLineItemRepository,
	adjRepo repository.AdjudicationRepository,
	policyRepo policyDomainRepo.PolicyRepository,
	memberRepo policyDomainRepo.MemberRepository,
	providerRepo providerRepo.ProviderRepository,
	validator service.ValidatorService,
	adjudicator service.AdjudicatorService,
) service.ClaimService {
	return &claimServiceImpl{
		claimRepo:    claimRepo,
		lineItemRepo: lineItemRepo,
		adjRepo:      adjRepo,
		policyRepo:   policyRepo,
		memberRepo:   memberRepo,
		providerRepo: providerRepo,
		validator:    validator,
		adjudicator:  adjudicator,
	}
}

// SubmitClaim creates a new claim, persists line items, then runs the
// validation -> adjudication pipeline automatically.
// State machine: RECEIVED -> VALIDATED -> ADJUDICATED -> APPROVED|REJECTED|MANUAL_REVIEW
func (s *claimServiceImpl) SubmitClaim(ctx context.Context, req claimsSchema.SubmitClaimRequest, createdBy uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse] {
	policyID, err := uuid.Parse(req.PolicyID)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusBadRequest, "Invalid policy ID", err)
	}
	memberID, err := uuid.Parse(req.MemberID)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusBadRequest, "Invalid member ID", err)
	}
	providerID, err := uuid.Parse(req.ProviderID)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusBadRequest, "Invalid provider ID", err)
	}

	var preAuthID uuid.UUID
	if req.PreAuthID != "" {
		preAuthID, err = uuid.Parse(req.PreAuthID)
		if err != nil {
			return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusBadRequest, "Invalid preauth ID", err)
		}
	}

	// Calculate total amount from line items
	var totalAmount int64
	for _, li := range req.LineItems {
		totalAmount += li.UnitPrice * int64(li.Quantity)
	}

	// Marshal diagnosis codes
	diagCodes, err := json.Marshal(req.DiagnosisCodes)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusInternalServerError, "Failed to encode diagnosis codes", err)
	}

	// Generate claim number
	claimNumber := utils.GenerateClaimNumber()

	claim := &entity.Claim{
		ClaimNumber:    claimNumber,
		PolicyID:       policyID,
		MemberID:       memberID,
		ProviderID:     providerID,
		PreAuthID:      preAuthID,
		Status:         string(shared.ClaimStatusReceived),
		TotalAmount:    totalAmount,
		DiagnosisCodes: diagCodes,
		ServiceDate:    req.ServiceDate,
		AdmissionDate:  req.AdmissionDate,
		DischargeDate:  req.DischargeDate,
		Notes:          req.Notes,
		CreatedBy:      createdBy,
	}

	// Create the claim
	claim, err = s.claimRepo.Create(ctx, claim)
	if err != nil {
		utils.LogError("Failed to create claim: %v", err)
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusInternalServerError, "Failed to create claim", err)
	}

	// Create line items
	lineItems := make([]*entity.ClaimLineItem, 0, len(req.LineItems))
	for _, li := range req.LineItems {
		item := &entity.ClaimLineItem{
			ClaimID:       claim.ID,
			ProcedureCode: li.ProcedureCode,
			ProcedureName: li.ProcedureName,
			DiagnosisCode: li.DiagnosisCode,
			Quantity:      li.Quantity,
			UnitPrice:     li.UnitPrice,
			TotalPrice:    li.UnitPrice * int64(li.Quantity),
		}
		created, err := s.lineItemRepo.Create(ctx, item)
		if err != nil {
			utils.LogError("Failed to create line item for claim %s: %v", claimNumber, err)
			return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusInternalServerError, "Failed to create claim line item", err)
		}
		lineItems = append(lineItems, created)
	}

	utils.LogInfo("Claim %s created with %d line items, total: %d cents", claimNumber, len(lineItems), totalAmount)

	// --- Auto-run validation -> adjudication pipeline ---

	// Step 1: Validate
	valid, validationErrors, err := s.validator.ValidateClaim(ctx, claim, lineItems)
	if err != nil {
		utils.LogError("Validation error for claim %s: %v", claimNumber, err)
		// Claim stays in RECEIVED status; don't fail the submission
		return schema.NewServiceResponse(s.buildClaimResponse(claim, lineItems, nil), http.StatusCreated, "Claim submitted but validation encountered an error")
	}

	if !valid {
		// Reject the claim with validation errors
		reason := fmt.Sprintf("Validation failed: %v", validationErrors)
		claim, _ = s.claimRepo.Reject(ctx, claim.ID, reason)
		utils.LogInfo("Claim %s rejected during validation: %v", claimNumber, validationErrors)
		return schema.NewServiceResponse(s.buildClaimResponse(claim, lineItems, nil), http.StatusCreated, "Claim submitted but rejected during validation")
	}

	// Move to VALIDATED
	claim, err = s.claimRepo.UpdateStatus(ctx, claim.ID, string(shared.ClaimStatusValidated))
	if err != nil {
		utils.LogError("Failed to update claim %s to VALIDATED: %v", claimNumber, err)
		return schema.NewServiceResponse(s.buildClaimResponse(claim, lineItems, nil), http.StatusCreated, "Claim submitted and validated but status update failed")
	}

	utils.LogInfo("Claim %s validated successfully", claimNumber)

	// Step 2: Adjudicate
	result, err := s.adjudicator.Adjudicate(ctx, claim, lineItems)
	if err != nil {
		utils.LogError("Adjudication error for claim %s: %v", claimNumber, err)
		return schema.NewServiceResponse(s.buildClaimResponse(claim, lineItems, nil), http.StatusCreated, "Claim submitted and validated but adjudication encountered an error")
	}

	// Persist adjudication decision
	ruleResultsJSON, _ := json.Marshal(result.Reasons)
	reasonsJSON, _ := json.Marshal(result.Reasons)

	decision := &entity.AdjudicationDecision{
		ClaimID:              claim.ID,
		Decision:             result.Decision,
		PayableAmount:        result.PayableAmount,
		MemberResponsibility: result.MemberResponsibility,
		Reasons:              reasonsJSON,
		RuleResults:          ruleResultsJSON,
	}
	_, err = s.adjRepo.Create(ctx, decision)
	if err != nil {
		utils.LogError("Failed to persist adjudication decision for claim %s: %v", claimNumber, err)
	}

	// Update claim amounts
	claim, err = s.claimRepo.UpdateAmounts(ctx, claim.ID, result.PayableAmount, result.MemberResponsibility-result.PayableAmount, result.MemberResponsibility)
	if err != nil {
		utils.LogError("Failed to update amounts for claim %s: %v", claimNumber, err)
	}

	// Apply final status based on decision
	var finalStatus string
	switch result.Decision {
	case string(shared.AdjudicationDecisionApprove):
		finalStatus = string(shared.ClaimStatusApproved)
	case string(shared.AdjudicationDecisionReject):
		finalStatus = string(shared.ClaimStatusRejected)
	case string(shared.AdjudicationDecisionManualReview):
		finalStatus = string(shared.ClaimStatusManualReview)
	default:
		finalStatus = string(shared.ClaimStatusAdjudicated)
	}

	claim, err = s.claimRepo.UpdateStatus(ctx, claim.ID, finalStatus)
	if err != nil {
		utils.LogError("Failed to update claim %s to %s: %v", claimNumber, finalStatus, err)
	}

	utils.LogInfo("Claim %s adjudicated: decision=%s, payable=%d, member_resp=%d", claimNumber, result.Decision, result.PayableAmount, result.MemberResponsibility)

	return schema.NewServiceResponse(s.buildClaimResponse(claim, lineItems, decision), http.StatusCreated, fmt.Sprintf("Claim submitted and adjudicated: %s", result.Decision))
}

// ApproveClaim manually approves a claim that is in MANUAL_REVIEW status.
func (s *claimServiceImpl) ApproveClaim(ctx context.Context, id uuid.UUID, approvedBy uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse] {
	claim, err := s.claimRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusNotFound, "Claim not found", err)
	}

	// State machine: only MANUAL_REVIEW -> APPROVED
	if claim.Status != string(shared.ClaimStatusManualReview) {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](
			http.StatusConflict,
			fmt.Sprintf("Cannot approve claim in status %s; must be in MANUAL_REVIEW", claim.Status),
			fmt.Errorf("invalid state transition from %s to APPROVED", claim.Status),
		)
	}

	claim, err = s.claimRepo.UpdateStatus(ctx, id, string(shared.ClaimStatusApproved))
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusInternalServerError, "Failed to approve claim", err)
	}

	utils.LogInfo("Claim %s manually approved by %s", claim.ClaimNumber, approvedBy.String())

	lineItems, _ := s.lineItemRepo.ListByClaim(ctx, id)
	decision, _ := s.adjRepo.GetByClaimID(ctx, id)

	return schema.NewServiceResponse(s.buildClaimResponse(claim, lineItems, decision), http.StatusOK, "Claim approved")
}

// RejectClaim manually rejects a claim that is in MANUAL_REVIEW status.
func (s *claimServiceImpl) RejectClaim(ctx context.Context, id uuid.UUID, reason string, rejectedBy uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse] {
	claim, err := s.claimRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusNotFound, "Claim not found", err)
	}

	// State machine: only MANUAL_REVIEW -> REJECTED
	if claim.Status != string(shared.ClaimStatusManualReview) {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](
			http.StatusConflict,
			fmt.Sprintf("Cannot reject claim in status %s; must be in MANUAL_REVIEW", claim.Status),
			fmt.Errorf("invalid state transition from %s to REJECTED", claim.Status),
		)
	}

	claim, err = s.claimRepo.Reject(ctx, id, reason)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusInternalServerError, "Failed to reject claim", err)
	}

	utils.LogInfo("Claim %s manually rejected by %s: %s", claim.ClaimNumber, rejectedBy.String(), reason)

	lineItems, _ := s.lineItemRepo.ListByClaim(ctx, id)
	decision, _ := s.adjRepo.GetByClaimID(ctx, id)

	return schema.NewServiceResponse(s.buildClaimResponse(claim, lineItems, decision), http.StatusOK, "Claim rejected")
}

// GetClaim retrieves a single claim by ID with all related data.
func (s *claimServiceImpl) GetClaim(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse] {
	claim, err := s.claimRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](http.StatusNotFound, "Claim not found", err)
	}

	lineItems, _ := s.lineItemRepo.ListByClaim(ctx, id)
	decision, _ := s.adjRepo.GetByClaimID(ctx, id)

	return schema.NewServiceResponse(s.buildClaimResponse(claim, lineItems, decision), http.StatusOK, "Claim retrieved")
}

// ListClaims returns a paginated list of claims.
func (s *claimServiceImpl) ListClaims(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]claimsSchema.ClaimResponse] {
	offset := (page - 1) * pageSize
	claims, err := s.claimRepo.List(ctx, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]claimsSchema.ClaimResponse](http.StatusInternalServerError, "Failed to list claims", err)
	}

	responses := make([]claimsSchema.ClaimResponse, 0, len(claims))
	for _, c := range claims {
		responses = append(responses, claimsSchema.ToClaimResponse(c))
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Claims retrieved")
}

// ListClaimsByStatus returns claims filtered by status.
func (s *claimServiceImpl) ListClaimsByStatus(ctx context.Context, status string, page, pageSize int) *schema.ServiceResponse[[]claimsSchema.ClaimResponse] {
	offset := (page - 1) * pageSize
	claims, err := s.claimRepo.ListByStatus(ctx, status, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]claimsSchema.ClaimResponse](http.StatusInternalServerError, "Failed to list claims by status", err)
	}

	responses := make([]claimsSchema.ClaimResponse, 0, len(claims))
	for _, c := range claims {
		responses = append(responses, claimsSchema.ToClaimResponse(c))
	}

	return schema.NewServiceResponse(responses, http.StatusOK, fmt.Sprintf("Claims with status %s retrieved", status))
}

// ListClaimsByProvider returns claims filtered by provider.
func (s *claimServiceImpl) ListClaimsByProvider(ctx context.Context, providerID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]claimsSchema.ClaimResponse] {
	offset := (page - 1) * pageSize
	claims, err := s.claimRepo.ListByProvider(ctx, providerID, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]claimsSchema.ClaimResponse](http.StatusInternalServerError, "Failed to list claims by provider", err)
	}

	responses := make([]claimsSchema.ClaimResponse, 0, len(claims))
	for _, c := range claims {
		responses = append(responses, claimsSchema.ToClaimResponse(c))
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Provider claims retrieved")
}

// GetTotalCount returns the total number of claims.
func (s *claimServiceImpl) GetTotalCount(ctx context.Context) *schema.ServiceResponse[int64] {
	count, err := s.claimRepo.Count(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to get claim count", err)
	}
	return schema.NewServiceResponse(count, http.StatusOK, "Claim count retrieved")
}

// buildClaimResponse constructs a full ClaimResponse with line items and adjudication decision.
func (s *claimServiceImpl) buildClaimResponse(claim *entity.Claim, lineItems []*entity.ClaimLineItem, decision *entity.AdjudicationDecision) claimsSchema.ClaimResponse {
	resp := claimsSchema.ToClaimResponse(claim)

	// Attach line items
	if lineItems != nil {
		liResponses := make([]claimsSchema.LineItemResponse, 0, len(lineItems))
		for _, li := range lineItems {
			liResponses = append(liResponses, claimsSchema.LineItemResponse{
				ID:             li.ID,
				ProcedureCode:  li.ProcedureCode,
				ProcedureName:  li.ProcedureName,
				DiagnosisCode:  li.DiagnosisCode,
				Quantity:       li.Quantity,
				UnitPrice:      li.UnitPrice,
				TotalPrice:     li.TotalPrice,
				ApprovedAmount: li.ApprovedAmount,
			})
		}
		resp.LineItems = liResponses
	}

	// Attach adjudication decision
	if decision != nil {
		resp.Decision = &claimsSchema.AdjudicationResponse{
			Decision:             decision.Decision,
			PayableAmount:        decision.PayableAmount,
			MemberResponsibility: decision.MemberResponsibility,
			Reasons:              decision.Reasons,
			RuleResults:          decision.RuleResults,
			AdjudicatedAt:        decision.AdjudicatedAt,
		}
	}

	return resp
}
