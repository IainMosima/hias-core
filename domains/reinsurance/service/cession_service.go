package service

import (
	"context"

	schema "github.com/bitbiz/hias-core/domains/identity/schema"
	reinsuranceSchema "github.com/bitbiz/hias-core/domains/reinsurance/schema"
	"github.com/google/uuid"
)

type CessionService interface {
	CedePremium(ctx context.Context, req reinsuranceSchema.CedePremiumRequest, createdBy uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.CessionResponse]
	GetCession(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.CessionResponse]
	ListByTreaty(ctx context.Context, treatyID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]reinsuranceSchema.CessionResponse]
	ListByPolicy(ctx context.Context, policyID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]reinsuranceSchema.CessionResponse]
	BookCession(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.CessionResponse]
	ReverseCession(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.CessionResponse]
	AutoCedePolicyPremium(ctx context.Context, req reinsuranceSchema.AutoCedePolicyPremiumRequest, createdBy uuid.UUID) *schema.ServiceResponse[[]reinsuranceSchema.CessionResponse]
	GetCessionCount(ctx context.Context) *schema.ServiceResponse[int64]
	GetTotalCededAmount(ctx context.Context) *schema.ServiceResponse[int64]
	GetTotalGrossAmount(ctx context.Context) *schema.ServiceResponse[int64]
}
