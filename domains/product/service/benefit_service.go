package service

import (
	"context"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	productSchema "github.com/bitbiz/hias-core/domains/product/schema"
	"github.com/google/uuid"
)

type BenefitService interface {
	CreateBenefit(ctx context.Context, planID uuid.UUID, req productSchema.CreateBenefitRequest) *schema.ServiceResponse[productSchema.BenefitResponse]
	GetBenefit(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[productSchema.BenefitResponse]
	ListBenefitsByPlan(ctx context.Context, planID uuid.UUID) *schema.ServiceResponse[[]productSchema.BenefitResponse]
	UpdateBenefit(ctx context.Context, id uuid.UUID, req productSchema.UpdateBenefitRequest) *schema.ServiceResponse[productSchema.BenefitResponse]
	DeleteBenefit(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[string]
	CheckCoverage(ctx context.Context, planID uuid.UUID, procedureCode string) *schema.ServiceResponse[bool]
	CalculateCoPay(ctx context.Context, benefitID uuid.UUID, amount int64) *schema.ServiceResponse[int64]
	CreateSubBenefit(ctx context.Context, parentID uuid.UUID, req productSchema.CreateBenefitRequest) *schema.ServiceResponse[productSchema.BenefitResponse]
	ListSubBenefits(ctx context.Context, parentID uuid.UUID) *schema.ServiceResponse[[]productSchema.BenefitResponse]
}
