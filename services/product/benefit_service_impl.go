package product

import (
	"context"
	"log"
	"net/http"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/domains/product/entity"
	"github.com/bitbiz/hias-core/domains/product/repository"
	productSchema "github.com/bitbiz/hias-core/domains/product/schema"
	"github.com/bitbiz/hias-core/domains/product/service"
	"github.com/bitbiz/hias-core/shared"
	"github.com/google/uuid"
)

type benefitServiceImpl struct {
	benefitRepo repository.BenefitRepository
	auditSvc    auditService.AuditService
}

func NewBenefitService(benefitRepo repository.BenefitRepository, auditSvc auditService.AuditService) service.BenefitService {
	return &benefitServiceImpl{benefitRepo: benefitRepo, auditSvc: auditSvc}
}

func (s *benefitServiceImpl) CreateBenefit(ctx context.Context, planID uuid.UUID, req productSchema.CreateBenefitRequest) *schema.ServiceResponse[productSchema.BenefitResponse] {
	subLimitType := req.SubLimitType
	if subLimitType == "" {
		subLimitType = string(shared.SubLimitTypeNone)
	}
	maxAge := req.MaxAge
	if maxAge == 0 {
		maxAge = shared.DefaultMaxAge
	}
	waitingPeriodType := req.WaitingPeriodType
	if waitingPeriodType == "" {
		waitingPeriodType = string(shared.WaitingPeriodTypeGeneral)
	}

	benefit := &entity.Benefit{
		PlanID:            planID,
		Name:              req.Name,
		Category:          req.Category,
		AnnualLimit:       req.AnnualLimit,
		CoPayType:         req.CoPayType,
		CoPayValue:        req.CoPayValue,
		WaitingPeriodDays: req.WaitingPeriodDays,
		SubLimitType:      subLimitType,
		SubLimitValue:     req.SubLimitValue,
		MinAge:            req.MinAge,
		MaxAge:            maxAge,
		WaitingPeriodType: waitingPeriodType,
		DeductibleAmount:  req.DeductibleAmount,
	}

	created, err := s.benefitRepo.Create(ctx, benefit)
	if err != nil {
		return schema.NewServiceErrorResponse[productSchema.BenefitResponse](http.StatusInternalServerError, "Failed to create benefit", err)
	}

	s.logAudit(ctx, uuid.Nil, string(shared.AuditEntityTypeBenefit), created.ID, string(shared.AuditActionCreate))

	return schema.NewServiceResponse(productSchema.ToBenefitResponse(created), http.StatusCreated, "Benefit created")
}

func (s *benefitServiceImpl) ListBenefitsByPlan(ctx context.Context, planID uuid.UUID) *schema.ServiceResponse[[]productSchema.BenefitResponse] {
	benefits, err := s.benefitRepo.ListByPlan(ctx, planID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]productSchema.BenefitResponse](http.StatusInternalServerError, "Failed to list benefits", err)
	}

	responses := make([]productSchema.BenefitResponse, len(benefits))
	for i, b := range benefits {
		responses[i] = productSchema.ToBenefitResponse(b)
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Benefits retrieved")
}

func (s *benefitServiceImpl) CheckCoverage(ctx context.Context, planID uuid.UUID, procedureCode string) *schema.ServiceResponse[bool] {
	benefits, err := s.benefitRepo.ListByPlan(ctx, planID)
	if err != nil {
		return schema.NewServiceErrorResponse[bool](http.StatusInternalServerError, "Failed to check coverage", err)
	}
	return schema.NewServiceResponse(len(benefits) > 0, http.StatusOK, "Coverage checked")
}

func (s *benefitServiceImpl) CalculateCoPay(ctx context.Context, benefitID uuid.UUID, amount int64) *schema.ServiceResponse[int64] {
	benefit, err := s.benefitRepo.GetByID(ctx, benefitID)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusNotFound, "Benefit not found", err)
	}

	var copay int64
	switch benefit.CoPayType {
	case string(shared.CoPayTypePercentage):
		copay = amount * benefit.CoPayValue / 100
	case string(shared.CoPayTypeFixed):
		copay = benefit.CoPayValue
	}

	return schema.NewServiceResponse(copay, http.StatusOK, "Co-pay calculated")
}

func (s *benefitServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}
