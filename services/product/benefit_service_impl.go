package product

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/domains/product/entity"
	"github.com/bitbiz/hias-core/domains/product/repository"
	productSchema "github.com/bitbiz/hias-core/domains/product/schema"
	"github.com/bitbiz/hias-core/domains/product/service"
	"github.com/google/uuid"
)

type benefitServiceImpl struct {
	benefitRepo   repository.BenefitRepository
	planRepo      repository.PlanRepository
	exclusionRepo repository.ExclusionRepository
}

func NewBenefitService(
	benefitRepo repository.BenefitRepository,
	planRepo repository.PlanRepository,
	exclusionRepo repository.ExclusionRepository,
) service.BenefitService {
	return &benefitServiceImpl{
		benefitRepo:   benefitRepo,
		planRepo:      planRepo,
		exclusionRepo: exclusionRepo,
	}
}

func (s *benefitServiceImpl) CreateBenefit(ctx context.Context, planID uuid.UUID, req productSchema.CreateBenefitRequest) *schema.ServiceResponse[productSchema.BenefitResponse] {
	// Verify the plan exists
	_, err := s.planRepo.GetByID(ctx, planID)
	if err != nil {
		return schema.NewServiceErrorResponse[productSchema.BenefitResponse](http.StatusNotFound, "Plan not found", err)
	}

	benefit, err := s.benefitRepo.Create(ctx, &entity.Benefit{
		PlanID:            planID,
		Name:              req.Name,
		Category:          req.Category,
		AnnualLimit:       req.AnnualLimit,
		CoPayType:         req.CoPayType,
		CoPayValue:        req.CoPayValue,
		WaitingPeriodDays: req.WaitingPeriodDays,
	})
	if err != nil {
		return schema.NewServiceErrorResponse[productSchema.BenefitResponse](http.StatusInternalServerError, "Failed to create benefit", err)
	}

	return schema.NewServiceResponse(productSchema.ToBenefitResponse(benefit), http.StatusCreated, "Benefit created successfully")
}

func (s *benefitServiceImpl) ListBenefitsByPlan(ctx context.Context, planID uuid.UUID) *schema.ServiceResponse[[]productSchema.BenefitResponse] {
	// Verify the plan exists
	_, err := s.planRepo.GetByID(ctx, planID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]productSchema.BenefitResponse](http.StatusNotFound, "Plan not found", err)
	}

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
	// Verify the plan exists
	_, err := s.planRepo.GetByID(ctx, planID)
	if err != nil {
		return schema.NewServiceErrorResponse[bool](http.StatusNotFound, "Plan not found", err)
	}

	// Check if the procedure code is excluded by any exclusion on this plan
	exclusions, err := s.exclusionRepo.ListByPlan(ctx, planID)
	if err != nil {
		return schema.NewServiceErrorResponse[bool](http.StatusInternalServerError, "Failed to check exclusions", err)
	}

	for _, exclusion := range exclusions {
		if exclusion.ICDCodes != nil && len(exclusion.ICDCodes) > 0 {
			var icdCodes []string
			if err := json.Unmarshal(exclusion.ICDCodes, &icdCodes); err == nil {
				for _, code := range icdCodes {
					if code == procedureCode {
						return schema.NewServiceResponse(false, http.StatusOK, "Procedure is excluded from coverage")
					}
				}
			}
		}
	}

	// Check if the plan has any benefits at all (plan must have benefits to provide coverage)
	benefits, err := s.benefitRepo.ListByPlan(ctx, planID)
	if err != nil {
		return schema.NewServiceErrorResponse[bool](http.StatusInternalServerError, "Failed to check benefits", err)
	}

	if len(benefits) == 0 {
		return schema.NewServiceResponse(false, http.StatusOK, "Plan has no benefits configured")
	}

	return schema.NewServiceResponse(true, http.StatusOK, "Procedure is covered")
}

func (s *benefitServiceImpl) CalculateCoPay(ctx context.Context, benefitID uuid.UUID, amount int64) *schema.ServiceResponse[int64] {
	benefit, err := s.benefitRepo.GetByID(ctx, benefitID)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusNotFound, "Benefit not found", err)
	}

	var coPay int64
	switch benefit.CoPayType {
	case "percentage":
		// CoPayValue is stored as a percentage (e.g., 20 for 20%)
		coPay = (amount * benefit.CoPayValue) / 100
	case "fixed":
		// CoPayValue is stored as a fixed amount in cents
		coPay = benefit.CoPayValue
		if coPay > amount {
			coPay = amount
		}
	default:
		return schema.NewServiceErrorResponse[int64](http.StatusBadRequest, fmt.Sprintf("Unknown co-pay type: %s", benefit.CoPayType), nil)
	}

	return schema.NewServiceResponse(coPay, http.StatusOK, "Co-pay calculated")
}
