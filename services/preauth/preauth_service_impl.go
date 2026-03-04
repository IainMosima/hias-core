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
	"github.com/bitbiz/hias-core/shared"
	"github.com/google/uuid"
)

type preauthServiceImpl struct {
	preauthRepo repository.PreAuthRepository
}

func NewPreAuthService(preauthRepo repository.PreAuthRepository) service.PreAuthService {
	return &preauthServiceImpl{preauthRepo: preauthRepo}
}

func (s *preauthServiceImpl) SubmitPreAuth(ctx context.Context, req preauthSchema.SubmitPreAuthRequest, createdBy uuid.UUID) *schema.ServiceResponse[preauthSchema.PreAuthResponse] {
	policyID, _ := uuid.Parse(req.PolicyID)
	memberID, _ := uuid.Parse(req.MemberID)
	providerID, _ := uuid.Parse(req.ProviderID)

	procJSON, _ := json.Marshal(req.ProcedureCodes)
	diagJSON, _ := json.Marshal(req.DiagnosisCodes)

	preauth := &entity.PreAuthorization{
		PolicyID:       policyID,
		MemberID:       memberID,
		ProviderID:     providerID,
		ProcedureCodes: procJSON,
		DiagnosisCodes: diagJSON,
		EstimatedCost:  req.EstimatedCost,
		Status:         string(shared.PreAuthStatusSubmitted),
		Notes:          req.Notes,
		CreatedBy:      createdBy,
	}

	created, err := s.preauthRepo.Create(ctx, preauth)
	if err != nil {
		return schema.NewServiceErrorResponse[preauthSchema.PreAuthResponse](http.StatusInternalServerError, "Failed to submit pre-auth", err)
	}

	return schema.NewServiceResponse(preauthSchema.ToPreAuthResponse(created), http.StatusCreated, "Pre-authorization submitted")
}

func (s *preauthServiceImpl) GetPreAuth(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[preauthSchema.PreAuthResponse] {
	preauth, err := s.preauthRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[preauthSchema.PreAuthResponse](http.StatusNotFound, "Pre-auth not found", err)
	}
	return schema.NewServiceResponse(preauthSchema.ToPreAuthResponse(preauth), http.StatusOK, "Pre-auth retrieved")
}

func (s *preauthServiceImpl) ListPreAuths(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]preauthSchema.PreAuthResponse] {
	offset := (page - 1) * pageSize
	preauths, err := s.preauthRepo.List(ctx, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]preauthSchema.PreAuthResponse](http.StatusInternalServerError, "Failed to list pre-auths", err)
	}

	responses := make([]preauthSchema.PreAuthResponse, len(preauths))
	for i, p := range preauths {
		responses[i] = preauthSchema.ToPreAuthResponse(p)
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Pre-auths retrieved")
}

func (s *preauthServiceImpl) ReviewPreAuth(ctx context.Context, id uuid.UUID, req preauthSchema.ReviewPreAuthRequest, reviewedBy uuid.UUID) *schema.ServiceResponse[preauthSchema.PreAuthResponse] {
	switch req.Decision {
	case string(shared.PreAuthStatusApproved):
		return s.ApprovePreAuth(ctx, id, reviewedBy)
	case string(shared.PreAuthStatusDenied):
		return s.DenyPreAuth(ctx, id, req.DenialReason, reviewedBy)
	case string(shared.PreAuthStatusInfoRequested):
		updated, err := s.preauthRepo.UpdateStatus(ctx, id, string(shared.PreAuthStatusInfoRequested))
		if err != nil {
			return schema.NewServiceErrorResponse[preauthSchema.PreAuthResponse](http.StatusInternalServerError, "Failed to update status", err)
		}
		return schema.NewServiceResponse(preauthSchema.ToPreAuthResponse(updated), http.StatusOK, "Info requested")
	default:
		return schema.NewServiceErrorResponse[preauthSchema.PreAuthResponse](http.StatusBadRequest, "Invalid decision", nil)
	}
}

func (s *preauthServiceImpl) ApprovePreAuth(ctx context.Context, id uuid.UUID, reviewedBy uuid.UUID) *schema.ServiceResponse[preauthSchema.PreAuthResponse] {
	preauth, err := s.preauthRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[preauthSchema.PreAuthResponse](http.StatusNotFound, "Pre-auth not found", err)
	}

	now := time.Now()
	validEnd := now.AddDate(0, 0, shared.PreAuthValidityDays)
	authCode := fmt.Sprintf("AUTH-%d-%06d", time.Now().Year(), time.Now().UnixNano()%1000000)

	preauth.Status = string(shared.PreAuthStatusApproved)
	preauth.AuthCode = authCode
	preauth.ApprovedAmount = preauth.EstimatedCost
	preauth.ValidityStart = &now
	preauth.ValidityEnd = &validEnd
	preauth.ReviewedBy = reviewedBy
	preauth.ReviewedAt = &now

	approved, err := s.preauthRepo.Approve(ctx, preauth)
	if err != nil {
		return schema.NewServiceErrorResponse[preauthSchema.PreAuthResponse](http.StatusInternalServerError, "Failed to approve pre-auth", err)
	}

	return schema.NewServiceResponse(preauthSchema.ToPreAuthResponse(approved), http.StatusOK, "Pre-auth approved")
}

func (s *preauthServiceImpl) DenyPreAuth(ctx context.Context, id uuid.UUID, reason string, reviewedBy uuid.UUID) *schema.ServiceResponse[preauthSchema.PreAuthResponse] {
	denied, err := s.preauthRepo.Deny(ctx, id, reason, reviewedBy)
	if err != nil {
		return schema.NewServiceErrorResponse[preauthSchema.PreAuthResponse](http.StatusInternalServerError, "Failed to deny pre-auth", err)
	}
	return schema.NewServiceResponse(preauthSchema.ToPreAuthResponse(denied), http.StatusOK, "Pre-auth denied")
}

func (s *preauthServiceImpl) ExpirePreAuths(ctx context.Context) *schema.ServiceResponse[int] {
	expiring, err := s.preauthRepo.GetExpiring(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int](http.StatusInternalServerError, "Failed to get expiring pre-auths", err)
	}

	expired := 0
	for _, p := range expiring {
		_, updateErr := s.preauthRepo.UpdateStatus(ctx, p.ID, string(shared.PreAuthStatusExpired))
		if updateErr == nil {
			expired++
		}
	}

	return schema.NewServiceResponse(expired, http.StatusOK, fmt.Sprintf("%d pre-auths expired", expired))
}

func (s *preauthServiceImpl) GetTotalCount(ctx context.Context) *schema.ServiceResponse[int64] {
	count, err := s.preauthRepo.Count(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to get count", err)
	}
	return schema.NewServiceResponse(count, http.StatusOK, "Count retrieved")
}
