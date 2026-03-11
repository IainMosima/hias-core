package provider

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/domains/provider/entity"
	"github.com/bitbiz/hias-core/domains/provider/repository"
	providerSchema "github.com/bitbiz/hias-core/domains/provider/schema"
	"github.com/bitbiz/hias-core/domains/provider/service"
	"github.com/bitbiz/hias-core/shared"
	"github.com/google/uuid"
)

type providerServiceImpl struct {
	providerRepo repository.ProviderRepository
	contractRepo repository.ContractRepository
	rateCardRepo repository.RateCardRepository
	auditSvc     auditService.AuditService
}

func NewProviderService(
	providerRepo repository.ProviderRepository,
	contractRepo repository.ContractRepository,
	rateCardRepo repository.RateCardRepository,
	auditSvc auditService.AuditService,
) service.ProviderService {
	return &providerServiceImpl{
		providerRepo: providerRepo,
		contractRepo: contractRepo,
		rateCardRepo: rateCardRepo,
		auditSvc:     auditSvc,
	}
}

func (s *providerServiceImpl) RegisterProvider(ctx context.Context, req providerSchema.RegisterProviderRequest, createdBy uuid.UUID) *schema.ServiceResponse[providerSchema.ProviderResponse] {
	tier := req.Tier
	if tier == "" {
		tier = string(shared.ProviderTierThree)
	}

	provider := &entity.Provider{
		Name:                req.Name,
		Type:                req.Type,
		LicenseNumber:       req.LicenseNumber,
		Status:              string(shared.ProviderStatusPending),
		Tier:                tier,
		County:              req.County,
		Address:             req.Address,
		Phone:               req.Phone,
		Email:               req.Email,
		ContactPerson:       req.ContactPerson,
		AccreditationStatus: "PENDING",
		CreatedBy:           createdBy,
	}

	created, err := s.providerRepo.Create(ctx, provider)
	if err != nil {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](http.StatusInternalServerError, "Failed to register provider", err)
	}

	s.logAudit(ctx, createdBy, string(shared.AuditEntityTypeProvider), created.ID, string(shared.AuditActionCreate))

	return schema.NewServiceResponse(providerSchema.ToProviderResponse(created), http.StatusCreated, "Provider registered")
}

func (s *providerServiceImpl) GetProvider(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[providerSchema.ProviderResponse] {
	provider, err := s.providerRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](http.StatusNotFound, "Provider not found", err)
	}
	return schema.NewServiceResponse(providerSchema.ToProviderResponse(provider), http.StatusOK, "Provider retrieved")
}

func (s *providerServiceImpl) ListProviders(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]providerSchema.ProviderResponse] {
	offset := (page - 1) * pageSize
	providers, err := s.providerRepo.List(ctx, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]providerSchema.ProviderResponse](http.StatusInternalServerError, "Failed to list providers", err)
	}

	responses := make([]providerSchema.ProviderResponse, len(providers))
	for i, p := range providers {
		responses[i] = providerSchema.ToProviderResponse(p)
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Providers retrieved")
}

func (s *providerServiceImpl) CredentialProvider(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[providerSchema.ProviderResponse] {
	return s.updateProviderStatus(ctx, id, string(shared.ProviderStatusCredentialing), string(shared.ProviderStatusPending), "Provider moved to credentialing")
}

func (s *providerServiceImpl) ActivateProvider(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[providerSchema.ProviderResponse] {
	provider, err := s.providerRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](http.StatusNotFound, "Provider not found", err)
	}

	if provider.Status == string(shared.ProviderStatusActive) {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](http.StatusBadRequest, "Provider is already active", fmt.Errorf("provider is already active"))
	}
	if provider.Status == string(shared.ProviderStatusPending) {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](http.StatusBadRequest, "Provider must complete credentialing before activation", fmt.Errorf("provider must complete credentialing"))
	}

	updated, err := s.providerRepo.UpdateStatus(ctx, id, string(shared.ProviderStatusActive))
	if err != nil {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](http.StatusInternalServerError, "Failed to activate provider", err)
	}

	return schema.NewServiceResponse(providerSchema.ToProviderResponse(updated), http.StatusOK, "Provider activated")
}

func (s *providerServiceImpl) SuspendProvider(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[providerSchema.ProviderResponse] {
	return s.updateProviderStatus(ctx, id, string(shared.ProviderStatusSuspended), string(shared.ProviderStatusActive), "Provider suspended")
}

func (s *providerServiceImpl) TerminateProvider(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[providerSchema.ProviderResponse] {
	updated, err := s.providerRepo.UpdateStatus(ctx, id, string(shared.ProviderStatusTerminated))
	if err != nil {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](http.StatusInternalServerError, "Failed to terminate provider", err)
	}
	return schema.NewServiceResponse(providerSchema.ToProviderResponse(updated), http.StatusOK, "Provider terminated")
}

func (s *providerServiceImpl) UpdateProvider(ctx context.Context, id uuid.UUID, req providerSchema.UpdateProviderRequest) *schema.ServiceResponse[providerSchema.ProviderResponse] {
	provider, err := s.providerRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](http.StatusNotFound, "Provider not found", err)
	}

	if req.Name != nil {
		provider.Name = *req.Name
	}
	if req.County != nil {
		provider.County = *req.County
	}
	if req.Address != nil {
		provider.Address = *req.Address
	}
	if req.Phone != nil {
		provider.Phone = *req.Phone
	}
	if req.Email != nil {
		provider.Email = *req.Email
	}
	if req.ContactPerson != nil {
		provider.ContactPerson = *req.ContactPerson
	}

	updated, err := s.providerRepo.Update(ctx, provider)
	if err != nil {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](http.StatusInternalServerError, "Failed to update provider", err)
	}

	return schema.NewServiceResponse(providerSchema.ToProviderResponse(updated), http.StatusOK, "Provider updated")
}

func (s *providerServiceImpl) GetTotalCount(ctx context.Context) *schema.ServiceResponse[int64] {
	count, err := s.providerRepo.Count(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to get count", err)
	}
	return schema.NewServiceResponse(count, http.StatusOK, "Count retrieved")
}

func (s *providerServiceImpl) updateProviderStatus(ctx context.Context, id uuid.UUID, newStatus, requiredStatus, message string) *schema.ServiceResponse[providerSchema.ProviderResponse] {
	provider, err := s.providerRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](http.StatusNotFound, "Provider not found", err)
	}

	if provider.Status != requiredStatus {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](http.StatusBadRequest, fmt.Sprintf("Provider must be in %s status", requiredStatus), nil)
	}

	updated, err := s.providerRepo.UpdateStatus(ctx, id, newStatus)
	if err != nil {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](http.StatusInternalServerError, "Failed to update provider status", err)
	}

	return schema.NewServiceResponse(providerSchema.ToProviderResponse(updated), http.StatusOK, message)
}

func (s *providerServiceImpl) UpdateTier(ctx context.Context, id uuid.UUID, tier string, userID uuid.UUID) *schema.ServiceResponse[providerSchema.ProviderResponse] {
	updated, err := s.providerRepo.UpdateTier(ctx, id, tier)
	if err != nil {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](http.StatusInternalServerError, "Failed to update provider tier", err)
	}

	s.logAudit(ctx, userID, string(shared.AuditEntityTypeProvider), id, string(shared.AuditActionUpdate))
	return schema.NewServiceResponse(providerSchema.ToProviderResponse(updated), http.StatusOK, "Provider tier updated")
}

func (s *providerServiceImpl) ListByTier(ctx context.Context, tier string, page, pageSize int) *schema.ServiceResponse[[]providerSchema.ProviderResponse] {
	offset := (page - 1) * pageSize
	providers, err := s.providerRepo.ListByTier(ctx, tier, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]providerSchema.ProviderResponse](http.StatusInternalServerError, "Failed to list providers by tier", err)
	}

	responses := make([]providerSchema.ProviderResponse, len(providers))
	for i, p := range providers {
		responses[i] = providerSchema.ToProviderResponse(p)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Providers retrieved")
}

func (s *providerServiceImpl) UpdateAccreditation(ctx context.Context, id uuid.UUID, req providerSchema.UpdateAccreditationRequest, userID uuid.UUID) *schema.ServiceResponse[providerSchema.ProviderResponse] {
	var expiry *time.Time
	if req.AccreditationExpiry != "" {
		t, err := time.Parse("2006-01-02", req.AccreditationExpiry)
		if err != nil {
			return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](http.StatusBadRequest, "Invalid accreditation_expiry format (YYYY-MM-DD)", err)
		}
		expiry = &t
	}

	updated, err := s.providerRepo.UpdateAccreditation(ctx, id, req.AccreditationStatus, expiry, req.AccreditationBody)
	if err != nil {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](http.StatusInternalServerError, "Failed to update accreditation", err)
	}

	s.logAudit(ctx, userID, string(shared.AuditEntityTypeProvider), id, string(shared.AuditActionUpdate))
	return schema.NewServiceResponse(providerSchema.ToProviderResponse(updated), http.StatusOK, "Accreditation updated")
}

func (s *providerServiceImpl) ListByAccreditationStatus(ctx context.Context, status string, page, pageSize int) *schema.ServiceResponse[[]providerSchema.ProviderResponse] {
	offset := (page - 1) * pageSize
	providers, err := s.providerRepo.ListByAccreditationStatus(ctx, status, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]providerSchema.ProviderResponse](http.StatusInternalServerError, "Failed to list providers by accreditation status", err)
	}

	responses := make([]providerSchema.ProviderResponse, len(providers))
	for i, p := range providers {
		responses[i] = providerSchema.ToProviderResponse(p)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Providers retrieved")
}

func (s *providerServiceImpl) ListExpiringAccreditations(ctx context.Context, days, page, pageSize int) *schema.ServiceResponse[[]providerSchema.ProviderResponse] {
	offset := (page - 1) * pageSize
	providers, err := s.providerRepo.ListExpiringAccreditations(ctx, days, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]providerSchema.ProviderResponse](http.StatusInternalServerError, "Failed to list expiring accreditations", err)
	}

	responses := make([]providerSchema.ProviderResponse, len(providers))
	for i, p := range providers {
		responses[i] = providerSchema.ToProviderResponse(p)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Expiring accreditations retrieved")
}

func (s *providerServiceImpl) ListProvidersFiltered(ctx context.Context, search string, page, pageSize int) *schema.ServiceResponse[[]providerSchema.ProviderResponse] {
	offset := (page - 1) * pageSize
	providers, err := s.providerRepo.ListFiltered(ctx, search, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]providerSchema.ProviderResponse](http.StatusInternalServerError, "Failed to list providers", err)
	}

	responses := make([]providerSchema.ProviderResponse, len(providers))
	for i, p := range providers {
		responses[i] = providerSchema.ToProviderResponse(p)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Providers retrieved")
}

func (s *providerServiceImpl) CountProvidersFiltered(ctx context.Context, search string) *schema.ServiceResponse[int64] {
	count, err := s.providerRepo.CountFiltered(ctx, search)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to count providers", err)
	}
	return schema.NewServiceResponse(count, http.StatusOK, "Count retrieved")
}

func (s *providerServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}
