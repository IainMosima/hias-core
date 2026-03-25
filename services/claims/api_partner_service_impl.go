package claims

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	"github.com/bitbiz/hias-core/domains/claims/entity"
	"github.com/bitbiz/hias-core/domains/claims/repository"
	claimsSchema "github.com/bitbiz/hias-core/domains/claims/schema"
	"github.com/bitbiz/hias-core/domains/claims/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/shared"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type apiPartnerServiceImpl struct {
	partnerRepo repository.APIPartnerRepository
	auditSvc    auditService.AuditService
}

func NewAPIPartnerService(partnerRepo repository.APIPartnerRepository, auditSvc auditService.AuditService) service.APIPartnerService {
	return &apiPartnerServiceImpl{
		partnerRepo: partnerRepo,
		auditSvc:    auditSvc,
	}
}

func (s *apiPartnerServiceImpl) CreatePartner(ctx context.Context, req claimsSchema.CreateAPIPartnerRequest) *schema.ServiceResponse[claimsSchema.CreateAPIPartnerResponse] {
	apiKey, err := generateRandomHex()
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.CreateAPIPartnerResponse](
			http.StatusInternalServerError, "Failed to generate API key", err,
		)
	}

	apiSecret, err := generateRandomHex()
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.CreateAPIPartnerResponse](
			http.StatusInternalServerError, "Failed to generate API secret", err,
		)
	}

	secretHash, err := bcrypt.GenerateFromPassword([]byte(apiSecret), bcrypt.DefaultCost)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.CreateAPIPartnerResponse](
			http.StatusInternalServerError, "Failed to hash API secret", err,
		)
	}

	var providerID uuid.UUID
	if req.ProviderID != "" {
		providerID, _ = uuid.Parse(req.ProviderID)
	}

	rateLimitPerMinute := req.RateLimitPerMinute
	if rateLimitPerMinute == 0 {
		rateLimitPerMinute = 60
	}

	allowedClaimTypes := req.AllowedClaimTypes
	if len(allowedClaimTypes) == 0 {
		allowedClaimTypes = []string{"DIRECT"}
	}

	partner := &entity.APIPartner{
		Name:               req.Name,
		PartnerType:        req.PartnerType,
		APIKey:             apiKey,
		APISecretHash:      string(secretHash),
		ProviderID:         providerID,
		IsActive:           true,
		RateLimitPerMinute: rateLimitPerMinute,
		AllowedClaimTypes:  allowedClaimTypes,
		WebhookURL:         req.WebhookURL,
		ContactEmail:       req.ContactEmail,
		Metadata:           req.Metadata,
	}

	created, err := s.partnerRepo.Create(ctx, partner)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.CreateAPIPartnerResponse](
			http.StatusInternalServerError, fmt.Sprintf("failed to create API partner: %v", err), err,
		)
	}

	s.logAudit(ctx, uuid.Nil, string(shared.AuditEntityTypeAPIPartner), created.ID, "CREATE")

	resp := claimsSchema.CreateAPIPartnerResponse{
		APIPartnerResponse: claimsSchema.ToAPIPartnerResponse(created),
		APISecret:          apiSecret,
	}
	return schema.NewServiceResponse(resp, http.StatusCreated, "API partner created. Store the API secret securely — it will not be shown again.")
}

func (s *apiPartnerServiceImpl) ListPartners(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]claimsSchema.APIPartnerResponse] {
	offset := (page - 1) * pageSize

	partners, err := s.partnerRepo.List(ctx, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]claimsSchema.APIPartnerResponse](
			http.StatusInternalServerError, fmt.Sprintf("failed to list API partners: %v", err), err,
		)
	}

	responses := make([]claimsSchema.APIPartnerResponse, len(partners))
	for i, p := range partners {
		responses[i] = claimsSchema.ToAPIPartnerResponse(p)
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "API partners retrieved")
}

func (s *apiPartnerServiceImpl) GetPartner(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[claimsSchema.APIPartnerResponse] {
	partner, err := s.partnerRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.APIPartnerResponse](
			http.StatusNotFound, "API partner not found", err,
		)
	}

	resp := claimsSchema.ToAPIPartnerResponse(partner)
	return schema.NewServiceResponse(resp, http.StatusOK, "API partner retrieved")
}

func (s *apiPartnerServiceImpl) UpdatePartner(ctx context.Context, id uuid.UUID, req claimsSchema.UpdateAPIPartnerRequest) *schema.ServiceResponse[claimsSchema.APIPartnerResponse] {
	existing, err := s.partnerRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.APIPartnerResponse](
			http.StatusNotFound, "API partner not found", err,
		)
	}

	var providerID uuid.UUID
	if req.ProviderID != "" {
		providerID, _ = uuid.Parse(req.ProviderID)
	}

	existing.Name = req.Name
	existing.PartnerType = req.PartnerType
	existing.ProviderID = providerID
	existing.RateLimitPerMinute = req.RateLimitPerMinute
	existing.AllowedClaimTypes = req.AllowedClaimTypes
	existing.WebhookURL = req.WebhookURL
	existing.ContactEmail = req.ContactEmail
	existing.Metadata = req.Metadata

	updated, err := s.partnerRepo.Update(ctx, existing)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.APIPartnerResponse](
			http.StatusInternalServerError, fmt.Sprintf("failed to update API partner: %v", err), err,
		)
	}

	s.logAudit(ctx, uuid.Nil, string(shared.AuditEntityTypeAPIPartner), updated.ID, "UPDATE")

	resp := claimsSchema.ToAPIPartnerResponse(updated)
	return schema.NewServiceResponse(resp, http.StatusOK, "API partner updated")
}

func (s *apiPartnerServiceImpl) DeactivatePartner(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[string] {
	if err := s.partnerRepo.Deactivate(ctx, id); err != nil {
		return schema.NewServiceErrorResponse[string](
			http.StatusInternalServerError, fmt.Sprintf("failed to deactivate API partner: %v", err), err,
		)
	}

	s.logAudit(ctx, uuid.Nil, string(shared.AuditEntityTypeAPIPartner), id, "DEACTIVATE")

	return schema.NewServiceResponse("", http.StatusOK, "API partner deactivated")
}

func (s *apiPartnerServiceImpl) RegenerateAPIKey(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[claimsSchema.CreateAPIPartnerResponse] {
	_, err := s.partnerRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.CreateAPIPartnerResponse](
			http.StatusNotFound, "API partner not found", err,
		)
	}

	apiKey, err := generateRandomHex()
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.CreateAPIPartnerResponse](
			http.StatusInternalServerError, "Failed to generate API key", err,
		)
	}

	apiSecret, err := generateRandomHex()
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.CreateAPIPartnerResponse](
			http.StatusInternalServerError, "Failed to generate API secret", err,
		)
	}

	secretHash, err := bcrypt.GenerateFromPassword([]byte(apiSecret), bcrypt.DefaultCost)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.CreateAPIPartnerResponse](
			http.StatusInternalServerError, "Failed to hash API secret", err,
		)
	}

	updated, err := s.partnerRepo.UpdateAPIKey(ctx, id, apiKey, string(secretHash))
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.CreateAPIPartnerResponse](
			http.StatusInternalServerError, fmt.Sprintf("failed to regenerate API key: %v", err), err,
		)
	}

	s.logAudit(ctx, uuid.Nil, string(shared.AuditEntityTypeAPIPartner), id, "REGENERATE_KEY")

	resp := claimsSchema.CreateAPIPartnerResponse{
		APIPartnerResponse: claimsSchema.ToAPIPartnerResponse(updated),
		APISecret:          apiSecret,
	}
	return schema.NewServiceResponse(resp, http.StatusOK, "API key regenerated. Store the new API secret securely.")
}

func generateRandomHex() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *apiPartnerServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}
