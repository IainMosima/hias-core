package policy

import (
	"context"
	"log"
	"net/http"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/domains/policy/repository"
	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/bitbiz/hias-core/domains/policy/service"
	"github.com/bitbiz/hias-core/shared"
	"github.com/google/uuid"
)

type underwritingFlagServiceImpl struct {
	flagRepo repository.UnderwritingFlagRepository
	auditSvc auditService.AuditService
}

func NewUnderwritingFlagService(
	flagRepo repository.UnderwritingFlagRepository,
	auditSvc auditService.AuditService,
) service.UnderwritingFlagService {
	return &underwritingFlagServiceImpl{
		flagRepo: flagRepo,
		auditSvc: auditSvc,
	}
}

func (s *underwritingFlagServiceImpl) ListByPolicy(ctx context.Context, policyID uuid.UUID) *schema.ServiceResponse[[]policySchema.UnderwritingFlagResponse] {
	flags, err := s.flagRepo.ListByPolicy(ctx, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]policySchema.UnderwritingFlagResponse](http.StatusInternalServerError, "Failed to list flags", err)
	}
	responses := make([]policySchema.UnderwritingFlagResponse, len(flags))
	for i, f := range flags {
		responses[i] = policySchema.ToUnderwritingFlagResponse(f)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Flags retrieved")
}

func (s *underwritingFlagServiceImpl) ListByMember(ctx context.Context, memberID uuid.UUID) *schema.ServiceResponse[[]policySchema.UnderwritingFlagResponse] {
	flags, err := s.flagRepo.ListByMember(ctx, memberID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]policySchema.UnderwritingFlagResponse](http.StatusInternalServerError, "Failed to list flags", err)
	}
	responses := make([]policySchema.UnderwritingFlagResponse, len(flags))
	for i, f := range flags {
		responses[i] = policySchema.ToUnderwritingFlagResponse(f)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Flags retrieved")
}

func (s *underwritingFlagServiceImpl) GetFlag(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[policySchema.UnderwritingFlagResponse] {
	flag, err := s.flagRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.UnderwritingFlagResponse](http.StatusNotFound, "Flag not found", err)
	}
	return schema.NewServiceResponse(policySchema.ToUnderwritingFlagResponse(flag), http.StatusOK, "Flag retrieved")
}

func (s *underwritingFlagServiceImpl) ListOpen(ctx context.Context, limit, offset int32) *schema.ServiceResponse[[]policySchema.UnderwritingFlagResponse] {
	flags, err := s.flagRepo.ListOpen(ctx, limit, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]policySchema.UnderwritingFlagResponse](http.StatusInternalServerError, "Failed to list open flags", err)
	}
	responses := make([]policySchema.UnderwritingFlagResponse, len(flags))
	for i, f := range flags {
		responses[i] = policySchema.ToUnderwritingFlagResponse(f)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Open flags retrieved")
}

func (s *underwritingFlagServiceImpl) CountOpen(ctx context.Context) *schema.ServiceResponse[int64] {
	count, err := s.flagRepo.CountOpen(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to count open flags", err)
	}
	return schema.NewServiceResponse(count, http.StatusOK, "Open flag count retrieved")
}

func (s *underwritingFlagServiceImpl) ResolveFlag(ctx context.Context, id uuid.UUID, req policySchema.ResolveFlagRequest, resolvedBy uuid.UUID) *schema.ServiceResponse[policySchema.UnderwritingFlagResponse] {
	flag, err := s.flagRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.UnderwritingFlagResponse](http.StatusNotFound, "Flag not found", err)
	}
	if flag.Status != string(shared.UnderwritingFlagStatusOpen) {
		return schema.NewServiceErrorResponse[policySchema.UnderwritingFlagResponse](http.StatusBadRequest, "Only OPEN flags can be resolved", nil)
	}

	resolved, err := s.flagRepo.Resolve(ctx, id, resolvedBy, req.Resolution)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.UnderwritingFlagResponse](http.StatusInternalServerError, "Failed to resolve flag", err)
	}

	s.logAudit(ctx, resolvedBy, string(shared.AuditEntityTypeUnderwritingFlag), id, string(shared.AuditActionStateChange))

	return schema.NewServiceResponse(policySchema.ToUnderwritingFlagResponse(resolved), http.StatusOK, "Flag resolved")
}

func (s *underwritingFlagServiceImpl) OverrideFlag(ctx context.Context, id uuid.UUID, req policySchema.OverrideFlagRequest, overriddenBy uuid.UUID) *schema.ServiceResponse[policySchema.UnderwritingFlagResponse] {
	flag, err := s.flagRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.UnderwritingFlagResponse](http.StatusNotFound, "Flag not found", err)
	}
	if flag.Status != string(shared.UnderwritingFlagStatusOpen) {
		return schema.NewServiceErrorResponse[policySchema.UnderwritingFlagResponse](http.StatusBadRequest, "Only OPEN flags can be overridden", nil)
	}

	overridden, err := s.flagRepo.Override(ctx, id, overriddenBy, req.Reason)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.UnderwritingFlagResponse](http.StatusInternalServerError, "Failed to override flag", err)
	}

	s.logAudit(ctx, overriddenBy, string(shared.AuditEntityTypeUnderwritingFlag), id, string(shared.AuditActionStateChange))

	return schema.NewServiceResponse(policySchema.ToUnderwritingFlagResponse(overridden), http.StatusOK, "Flag overridden")
}

func (s *underwritingFlagServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}
