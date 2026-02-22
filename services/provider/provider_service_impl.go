package provider

import (
	"context"
	"fmt"
	"net/http"

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
}

func NewProviderService(
	providerRepo repository.ProviderRepository,
	contractRepo repository.ContractRepository,
	rateCardRepo repository.RateCardRepository,
) service.ProviderService {
	return &providerServiceImpl{
		providerRepo: providerRepo,
		contractRepo: contractRepo,
		rateCardRepo: rateCardRepo,
	}
}

func (s *providerServiceImpl) RegisterProvider(ctx context.Context, req providerSchema.RegisterProviderRequest, createdBy uuid.UUID) *schema.ServiceResponse[providerSchema.ProviderResponse] {
	// Check if a provider with this license number already exists
	existing, _ := s.providerRepo.GetByLicense(ctx, req.LicenseNumber)
	if existing != nil {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](
			http.StatusConflict,
			fmt.Sprintf("Provider with license number %s already exists", req.LicenseNumber),
			nil,
		)
	}

	provider, err := s.providerRepo.Create(ctx, &entity.Provider{
		Name:          req.Name,
		Type:          req.Type,
		LicenseNumber: req.LicenseNumber,
		Status:        string(shared.ProviderStatusPending),
		County:        req.County,
		Address:       req.Address,
		Phone:         req.Phone,
		Email:         req.Email,
		ContactPerson: req.ContactPerson,
		CreatedBy:     createdBy,
	})
	if err != nil {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](http.StatusInternalServerError, "Failed to register provider", err)
	}

	return schema.NewServiceResponse(providerSchema.ToProviderResponse(provider), http.StatusCreated, "Provider registered successfully")
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
	provider, err := s.providerRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](http.StatusNotFound, "Provider not found", err)
	}

	// State machine: only PENDING → CREDENTIALING
	if provider.Status != string(shared.ProviderStatusPending) {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](
			http.StatusConflict,
			fmt.Sprintf("Cannot credential provider: current status is %s, expected PENDING", provider.Status),
			nil,
		)
	}

	updated, err := s.providerRepo.UpdateStatus(ctx, id, string(shared.ProviderStatusCredentialing))
	if err != nil {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](http.StatusInternalServerError, "Failed to credential provider", err)
	}

	return schema.NewServiceResponse(providerSchema.ToProviderResponse(updated), http.StatusOK, "Provider credentialing started")
}

func (s *providerServiceImpl) ActivateProvider(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[providerSchema.ProviderResponse] {
	provider, err := s.providerRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](http.StatusNotFound, "Provider not found", err)
	}

	// State machine: only CREDENTIALING → ACTIVE
	if provider.Status != string(shared.ProviderStatusCredentialing) {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](
			http.StatusConflict,
			fmt.Sprintf("Cannot activate provider: current status is %s, expected CREDENTIALING", provider.Status),
			nil,
		)
	}

	updated, err := s.providerRepo.UpdateStatus(ctx, id, string(shared.ProviderStatusActive))
	if err != nil {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](http.StatusInternalServerError, "Failed to activate provider", err)
	}

	return schema.NewServiceResponse(providerSchema.ToProviderResponse(updated), http.StatusOK, "Provider activated")
}

func (s *providerServiceImpl) SuspendProvider(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[providerSchema.ProviderResponse] {
	provider, err := s.providerRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](http.StatusNotFound, "Provider not found", err)
	}

	// State machine: only ACTIVE → SUSPENDED
	if provider.Status != string(shared.ProviderStatusActive) {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](
			http.StatusConflict,
			fmt.Sprintf("Cannot suspend provider: current status is %s, expected ACTIVE", provider.Status),
			nil,
		)
	}

	updated, err := s.providerRepo.UpdateStatus(ctx, id, string(shared.ProviderStatusSuspended))
	if err != nil {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](http.StatusInternalServerError, "Failed to suspend provider", err)
	}

	return schema.NewServiceResponse(providerSchema.ToProviderResponse(updated), http.StatusOK, "Provider suspended")
}

func (s *providerServiceImpl) TerminateProvider(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[providerSchema.ProviderResponse] {
	provider, err := s.providerRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](http.StatusNotFound, "Provider not found", err)
	}

	// State machine: ACTIVE or SUSPENDED → TERMINATED
	if provider.Status != string(shared.ProviderStatusActive) && provider.Status != string(shared.ProviderStatusSuspended) {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](
			http.StatusConflict,
			fmt.Sprintf("Cannot terminate provider: current status is %s, expected ACTIVE or SUSPENDED", provider.Status),
			nil,
		)
	}

	updated, err := s.providerRepo.UpdateStatus(ctx, id, string(shared.ProviderStatusTerminated))
	if err != nil {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](http.StatusInternalServerError, "Failed to terminate provider", err)
	}

	return schema.NewServiceResponse(providerSchema.ToProviderResponse(updated), http.StatusOK, "Provider terminated")
}

func (s *providerServiceImpl) UpdateProvider(ctx context.Context, id uuid.UUID, req providerSchema.UpdateProviderRequest) *schema.ServiceResponse[providerSchema.ProviderResponse] {
	existing, err := s.providerRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](http.StatusNotFound, "Provider not found", err)
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.County != nil {
		existing.County = *req.County
	}
	if req.Address != nil {
		existing.Address = *req.Address
	}
	if req.Phone != nil {
		existing.Phone = *req.Phone
	}
	if req.Email != nil {
		existing.Email = *req.Email
	}
	if req.ContactPerson != nil {
		existing.ContactPerson = *req.ContactPerson
	}

	updated, err := s.providerRepo.Update(ctx, existing)
	if err != nil {
		return schema.NewServiceErrorResponse[providerSchema.ProviderResponse](http.StatusInternalServerError, "Failed to update provider", err)
	}

	return schema.NewServiceResponse(providerSchema.ToProviderResponse(updated), http.StatusOK, "Provider updated")
}

func (s *providerServiceImpl) GetTotalCount(ctx context.Context) *schema.ServiceResponse[int64] {
	count, err := s.providerRepo.Count(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to count providers", err)
	}
	return schema.NewServiceResponse(count, http.StatusOK, "Count retrieved")
}
