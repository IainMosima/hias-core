package preauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/domains/preauth/entity"
	"github.com/bitbiz/hias-core/domains/preauth/repository"
	preauthSchema "github.com/bitbiz/hias-core/domains/preauth/schema"
	"github.com/bitbiz/hias-core/domains/preauth/service"
	policyDomainRepo "github.com/bitbiz/hias-core/domains/policy/repository"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/google/uuid"
)

type preAuthServiceImpl struct {
	preauthRepo repository.PreAuthRepository
	policyRepo  policyDomainRepo.PolicyRepository
	memberRepo  policyDomainRepo.MemberRepository
}

func NewPreAuthService(
	preauthRepo repository.PreAuthRepository,
	policyRepo policyDomainRepo.PolicyRepository,
	memberRepo policyDomainRepo.MemberRepository,
) service.PreAuthService {
	return &preAuthServiceImpl{
		preauthRepo: preauthRepo,
		policyRepo:  policyRepo,
		memberRepo:  memberRepo,
	}
}

// SubmitPreAuth creates a new pre-authorization request with SUBMITTED status.
// State machine: SUBMITTED -> UNDER_REVIEW -> APPROVED|DENIED|INFO_REQUESTED -> EXPIRED|CLAIMED
func (s *preAuthServiceImpl) SubmitPreAuth(ctx context.Context, req preauthSchema.SubmitPreAuthRequest, createdBy uuid.UUID) *schema.ServiceResponse[preauthSchema.PreAuthResponse] {
	policyID, err := uuid.Parse(req.PolicyID)
	if err != nil {
		return schema.NewServiceErrorResponse[preauthSchema.PreAuthResponse](http.StatusBadRequest, "Invalid policy ID", err)
	}
	memberID, err := uuid.Parse(req.MemberID)
	if err != nil {
		return schema.NewServiceErrorResponse[preauthSchema.PreAuthResponse](http.StatusBadRequest, "Invalid member ID", err)
	}
	providerID, err := uuid.Parse(req.ProviderID)
	if err != nil {
		return schema.NewServiceErrorResponse[preauthSchema.PreAuthResponse](http.StatusBadRequest, "Invalid provider ID", err)
	}

	// Verify policy exists and is active
	policy, err := s.policyRepo.GetByID(ctx, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[preauthSchema.PreAuthResponse](http.StatusNotFound, "Policy not found", err)
	}
	if policy.Status != string(shared.PolicyStatusActive) {
		return schema.NewServiceErrorResponse[preauthSchema.PreAuthResponse](
			http.StatusConflict,
			fmt.Sprintf("Policy is not active (status: %s)", policy.Status),
			fmt.Errorf("policy %s is not active", policyID),
		)
	}

	// Verify member exists
	_, err = s.memberRepo.GetByID(ctx, memberID)
	if err != nil {
		return schema.NewServiceErrorResponse[preauthSchema.PreAuthResponse](http.StatusNotFound, "Member not found", err)
	}

	// Marshal codes
	procCodes, err := json.Marshal(req.ProcedureCodes)
	if err != nil {
		return schema.NewServiceErrorResponse[preauthSchema.PreAuthResponse](http.StatusInternalServerError, "Failed to encode procedure codes", err)
	}
	diagCodes, err := json.Marshal(req.DiagnosisCodes)
	if err != nil {
		return schema.NewServiceErrorResponse[preauthSchema.PreAuthResponse](http.StatusInternalServerError, "Failed to encode diagnosis codes", err)
	}

	preauth := &entity.PreAuthorization{
		PolicyID:       policyID,
		MemberID:       memberID,
		ProviderID:     providerID,
		ProcedureCodes: procCodes,
		DiagnosisCodes: diagCodes,
		EstimatedCost:  req.EstimatedCost,
		Status:         string(shared.PreAuthStatusSubmitted),
		Notes:          req.Notes,
		CreatedBy:      createdBy,
	}

	preauth, err = s.preauthRepo.Create(ctx, preauth)
	if err != nil {
		utils.LogError("Failed to create pre-auth: %v", err)
		return schema.NewServiceErrorResponse[preauthSchema.PreAuthResponse](http.StatusInternalServerError, "Failed to create pre-authorization", err)
	}

	utils.LogInfo("Pre-auth %s created for policy %s, member %s", preauth.ID, policyID, memberID)

	return schema.NewServiceResponse(preauthSchema.ToPreAuthResponse(preauth), http.StatusCreated, "Pre-authorization submitted")
}

// ReviewPreAuth handles the review decision for a pre-authorization.
// Depending on the decision, it delegates to ApprovePreAuth, DenyPreAuth, or
// moves the pre-auth to INFO_REQUESTED status.
func (s *preAuthServiceImpl) ReviewPreAuth(ctx context.Context, id uuid.UUID, req preauthSchema.ReviewPreAuthRequest, reviewedBy uuid.UUID) *schema.ServiceResponse[preauthSchema.PreAuthResponse] {
	preauth, err := s.preauthRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[preauthSchema.PreAuthResponse](http.StatusNotFound, "Pre-authorization not found", err)
	}

	// Can review from SUBMITTED or UNDER_REVIEW
	if preauth.Status != string(shared.PreAuthStatusSubmitted) && preauth.Status != string(shared.PreAuthStatusUnderReview) {
		return schema.NewServiceErrorResponse[preauthSchema.PreAuthResponse](
			http.StatusConflict,
			fmt.Sprintf("Cannot review pre-auth in status %s", preauth.Status),
			fmt.Errorf("invalid state transition from %s", preauth.Status),
		)
	}

	// Move to UNDER_REVIEW if currently SUBMITTED
	if preauth.Status == string(shared.PreAuthStatusSubmitted) {
		preauth, err = s.preauthRepo.UpdateStatus(ctx, id, string(shared.PreAuthStatusUnderReview))
		if err != nil {
			return schema.NewServiceErrorResponse[preauthSchema.PreAuthResponse](http.StatusInternalServerError, "Failed to update status", err)
		}
	}

	switch req.Decision {
	case "APPROVED":
		return s.ApprovePreAuth(ctx, id, reviewedBy)
	case "DENIED":
		return s.DenyPreAuth(ctx, id, req.DenialReason, reviewedBy)
	case "INFO_REQUESTED":
		preauth, err = s.preauthRepo.UpdateStatus(ctx, id, string(shared.PreAuthStatusInfoRequested))
		if err != nil {
			return schema.NewServiceErrorResponse[preauthSchema.PreAuthResponse](http.StatusInternalServerError, "Failed to request info", err)
		}
		utils.LogInfo("Pre-auth %s: additional information requested by %s", id, reviewedBy)
		return schema.NewServiceResponse(preauthSchema.ToPreAuthResponse(preauth), http.StatusOK, "Additional information requested")
	default:
		return schema.NewServiceErrorResponse[preauthSchema.PreAuthResponse](http.StatusBadRequest, "Invalid decision", fmt.Errorf("unknown decision: %s", req.Decision))
	}
}

// ApprovePreAuth approves a pre-authorization, generates an auth code,
// and sets validity dates.
// State machine: SUBMITTED|UNDER_REVIEW -> APPROVED
func (s *preAuthServiceImpl) ApprovePreAuth(ctx context.Context, id uuid.UUID, reviewedBy uuid.UUID) *schema.ServiceResponse[preauthSchema.PreAuthResponse] {
	preauth, err := s.preauthRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[preauthSchema.PreAuthResponse](http.StatusNotFound, "Pre-authorization not found", err)
	}

	// Can approve from SUBMITTED, UNDER_REVIEW, or INFO_REQUESTED
	validStatuses := map[string]bool{
		string(shared.PreAuthStatusSubmitted):     true,
		string(shared.PreAuthStatusUnderReview):   true,
		string(shared.PreAuthStatusInfoRequested): true,
	}
	if !validStatuses[preauth.Status] {
		return schema.NewServiceErrorResponse[preauthSchema.PreAuthResponse](
			http.StatusConflict,
			fmt.Sprintf("Cannot approve pre-auth in status %s", preauth.Status),
			fmt.Errorf("invalid state transition from %s to APPROVED", preauth.Status),
		)
	}

	// Generate auth code and set validity
	authCode := generateAuthCode()
	now := time.Now()
	validityEnd := now.AddDate(0, 0, 30) // 30-day validity by default

	preauth.AuthCode = authCode
	preauth.Status = string(shared.PreAuthStatusApproved)
	preauth.ApprovedAmount = preauth.EstimatedCost
	preauth.ValidityStart = &now
	preauth.ValidityEnd = &validityEnd
	preauth.ReviewedBy = reviewedBy
	preauth.ReviewedAt = &now

	preauth, err = s.preauthRepo.Approve(ctx, preauth)
	if err != nil {
		return schema.NewServiceErrorResponse[preauthSchema.PreAuthResponse](http.StatusInternalServerError, "Failed to approve pre-authorization", err)
	}

	utils.LogInfo("Pre-auth %s approved by %s with auth_code=%s, valid until %s",
		id, reviewedBy, authCode, validityEnd.Format("2006-01-02"))

	return schema.NewServiceResponse(preauthSchema.ToPreAuthResponse(preauth), http.StatusOK, "Pre-authorization approved")
}

// DenyPreAuth denies a pre-authorization with a reason.
// State machine: SUBMITTED|UNDER_REVIEW|INFO_REQUESTED -> DENIED
func (s *preAuthServiceImpl) DenyPreAuth(ctx context.Context, id uuid.UUID, reason string, reviewedBy uuid.UUID) *schema.ServiceResponse[preauthSchema.PreAuthResponse] {
	preauth, err := s.preauthRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[preauthSchema.PreAuthResponse](http.StatusNotFound, "Pre-authorization not found", err)
	}

	validStatuses := map[string]bool{
		string(shared.PreAuthStatusSubmitted):     true,
		string(shared.PreAuthStatusUnderReview):   true,
		string(shared.PreAuthStatusInfoRequested): true,
	}
	if !validStatuses[preauth.Status] {
		return schema.NewServiceErrorResponse[preauthSchema.PreAuthResponse](
			http.StatusConflict,
			fmt.Sprintf("Cannot deny pre-auth in status %s", preauth.Status),
			fmt.Errorf("invalid state transition from %s to DENIED", preauth.Status),
		)
	}

	preauth, err = s.preauthRepo.Deny(ctx, id, reason, reviewedBy)
	if err != nil {
		return schema.NewServiceErrorResponse[preauthSchema.PreAuthResponse](http.StatusInternalServerError, "Failed to deny pre-authorization", err)
	}

	utils.LogInfo("Pre-auth %s denied by %s: %s", id, reviewedBy, reason)

	return schema.NewServiceResponse(preauthSchema.ToPreAuthResponse(preauth), http.StatusOK, "Pre-authorization denied")
}

// ExpirePreAuths finds all approved pre-authorizations that have passed their
// validity end date and marks them as EXPIRED. This is intended to be run
// by a scheduled job.
func (s *preAuthServiceImpl) ExpirePreAuths(ctx context.Context) *schema.ServiceResponse[int] {
	expiring, err := s.preauthRepo.GetExpiring(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int](http.StatusInternalServerError, "Failed to get expiring pre-auths", err)
	}

	expiredCount := 0
	for _, pa := range expiring {
		_, err := s.preauthRepo.UpdateStatus(ctx, pa.ID, string(shared.PreAuthStatusExpired))
		if err != nil {
			utils.LogError("Failed to expire pre-auth %s: %v", pa.ID, err)
			continue
		}
		expiredCount++
	}

	utils.LogInfo("Expired %d pre-authorizations out of %d eligible", expiredCount, len(expiring))

	return schema.NewServiceResponse(expiredCount, http.StatusOK, fmt.Sprintf("Expired %d pre-authorizations", expiredCount))
}

// GetPreAuth retrieves a single pre-authorization by ID.
func (s *preAuthServiceImpl) GetPreAuth(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[preauthSchema.PreAuthResponse] {
	preauth, err := s.preauthRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[preauthSchema.PreAuthResponse](http.StatusNotFound, "Pre-authorization not found", err)
	}

	return schema.NewServiceResponse(preauthSchema.ToPreAuthResponse(preauth), http.StatusOK, "Pre-authorization retrieved")
}

// ListPreAuths returns a paginated list of pre-authorizations.
func (s *preAuthServiceImpl) ListPreAuths(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]preauthSchema.PreAuthResponse] {
	offset := (page - 1) * pageSize
	preauths, err := s.preauthRepo.List(ctx, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]preauthSchema.PreAuthResponse](http.StatusInternalServerError, "Failed to list pre-authorizations", err)
	}

	responses := make([]preauthSchema.PreAuthResponse, 0, len(preauths))
	for _, pa := range preauths {
		responses = append(responses, preauthSchema.ToPreAuthResponse(pa))
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Pre-authorizations retrieved")
}

// GetTotalCount returns the total number of pre-authorizations.
func (s *preAuthServiceImpl) GetTotalCount(ctx context.Context) *schema.ServiceResponse[int64] {
	count, err := s.preauthRepo.Count(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to get pre-auth count", err)
	}
	return schema.NewServiceResponse(count, http.StatusOK, "Pre-auth count retrieved")
}

// generateAuthCode generates a unique pre-authorization code.
// Format: PA-<8 hex chars from UUID>
func generateAuthCode() string {
	id := uuid.New()
	return fmt.Sprintf("PA-%s", id.String()[:8])
}
