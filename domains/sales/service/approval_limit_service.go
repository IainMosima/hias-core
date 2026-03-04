package service

import (
	"context"

	"github.com/bitbiz/hias-core/domains/identity/schema"
	salesSchema "github.com/bitbiz/hias-core/domains/sales/schema"
	"github.com/google/uuid"
)

type ApprovalLimitService interface {
	GetLimits(ctx context.Context) *schema.ServiceResponse[[]salesSchema.ApprovalLimitResponse]
	GetLimitByRole(ctx context.Context, roleName string) *schema.ServiceResponse[salesSchema.ApprovalLimitResponse]
	CreateLimit(ctx context.Context, req salesSchema.CreateApprovalLimitRequest, createdBy uuid.UUID) *schema.ServiceResponse[salesSchema.ApprovalLimitResponse]
	UpdateLimit(ctx context.Context, id uuid.UUID, req salesSchema.UpdateApprovalLimitRequest, updatedBy uuid.UUID) *schema.ServiceResponse[salesSchema.ApprovalLimitResponse]
}
