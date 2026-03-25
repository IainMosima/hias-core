package provider

import (
	"context"
	"fmt"
	"log"
	"net/http"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/domains/provider/entity"
	"github.com/bitbiz/hias-core/domains/provider/repository"
	providerSchema "github.com/bitbiz/hias-core/domains/provider/schema"
	"github.com/bitbiz/hias-core/domains/provider/service"
	"github.com/bitbiz/hias-core/shared"
	"github.com/google/uuid"
)

type rateCardServiceImpl struct {
	rateCardRepo repository.RateCardRepository
	auditSvc     auditService.AuditService
}

func NewRateCardService(rateCardRepo repository.RateCardRepository, auditSvc auditService.AuditService) service.RateCardService {
	return &rateCardServiceImpl{
		rateCardRepo: rateCardRepo,
		auditSvc:     auditSvc,
	}
}

func (s *rateCardServiceImpl) createSingleRateCard(ctx context.Context, providerID uuid.UUID, rc providerSchema.CreateRateCardRequest) (*entity.RateCard, error) {
	rateCard := &entity.RateCard{
		ProviderID:    providerID,
		ProcedureCode: rc.ProcedureCode,
		ProcedureName: rc.ProcedureName,
		RateAmount:    rc.RateAmount,
		EffectiveDate: rc.EffectiveDate,
		AgeFrom:       rc.AgeFrom,
		AgeTo:         rc.AgeTo,
		Gender:        rc.Gender,
		Relationship:  rc.Relationship,
	}

	created, err := s.rateCardRepo.Create(ctx, rateCard)
	if err != nil {
		return nil, fmt.Errorf("failed to create rate card for %s: %w", rc.ProcedureCode, err)
	}

	s.logAudit(ctx, uuid.Nil, string(shared.AuditEntityTypeRateCard), created.ID, "CREATE")
	return created, nil
}

func (s *rateCardServiceImpl) CreateRateCard(ctx context.Context, providerID uuid.UUID, req providerSchema.CreateRateCardRequest) *schema.ServiceResponse[providerSchema.RateCardResponse] {
	created, err := s.createSingleRateCard(ctx, providerID, req)
	if err != nil {
		return schema.NewServiceErrorResponse[providerSchema.RateCardResponse](
			http.StatusInternalServerError, err.Error(), err,
		)
	}

	resp := providerSchema.ToRateCardResponse(created)
	return schema.NewServiceResponse(resp, http.StatusCreated, "Rate card created")
}

func (s *rateCardServiceImpl) ListRateCards(ctx context.Context, providerID uuid.UUID) *schema.ServiceResponse[[]providerSchema.RateCardResponse] {
	rateCards, err := s.rateCardRepo.ListByProvider(ctx, providerID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]providerSchema.RateCardResponse](
			http.StatusInternalServerError, fmt.Sprintf("failed to list rate cards: %v", err), err,
		)
	}

	responses := make([]providerSchema.RateCardResponse, len(rateCards))
	for i, r := range rateCards {
		responses[i] = providerSchema.ToRateCardResponse(r)
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Rate cards retrieved")
}

func (s *rateCardServiceImpl) BulkCreateRateCards(ctx context.Context, providerID uuid.UUID, req providerSchema.BulkCreateRateCardRequest) *schema.ServiceResponse[[]providerSchema.RateCardResponse] {
	var responses []providerSchema.RateCardResponse
	for _, rc := range req.RateCards {
		created, err := s.createSingleRateCard(ctx, providerID, rc)
		if err != nil {
			return schema.NewServiceErrorResponse[[]providerSchema.RateCardResponse](
				http.StatusInternalServerError, err.Error(), err,
			)
		}
		responses = append(responses, providerSchema.ToRateCardResponse(created))
	}

	return schema.NewServiceResponse(responses, http.StatusCreated, fmt.Sprintf("%d rate cards created", len(responses)))
}

func (s *rateCardServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}
