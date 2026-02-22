package service

import (
	"context"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	productSchema "github.com/bitbiz/hias-core/domains/product/schema"
	"github.com/google/uuid"
)

type BenefitService interface {
	CreateBenefit(ctx context.Context, planID uuid.UUID, req productSchema.CreateBenefitRequest) *schema.ServiceResponse[productSchema.BenefitResponse]
	ListBenefitsByPlan(ctx context.Context, planID uuid.UUID) *schema.ServiceResponse[[]productSchema.BenefitResponse]
	CheckCoverage(ctx context.Context, planID uuid.UUID, procedureCode string) *schema.ServiceResponse[bool]
	CalculateCoPay(ctx context.Context, benefitID uuid.UUID, amount int64) *schema.ServiceResponse[int64]
}
