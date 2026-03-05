package service

import (
	"context"

	schema "github.com/bitbiz/hias-core/domains/identity/schema"
	reinsuranceSchema "github.com/bitbiz/hias-core/domains/reinsurance/schema"
	"github.com/google/uuid"
)

type BordereauService interface {
	GeneratePremiumBordereau(ctx context.Context, req reinsuranceSchema.GenerateBordereauRequest, createdBy uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.BordereauResponse]
	GenerateClaimBordereau(ctx context.Context, req reinsuranceSchema.GenerateBordereauRequest, createdBy uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.BordereauResponse]
	GetBordereau(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.BordereauDetailResponse]
	ListByTreaty(ctx context.Context, treatyID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]reinsuranceSchema.BordereauResponse]
	FinalizeBordereau(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.BordereauResponse]
	MarkSent(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.BordereauResponse]
	ListItems(ctx context.Context, bordereauID uuid.UUID) *schema.ServiceResponse[[]reinsuranceSchema.BordereauItemResponse]
}
