package service

import (
	"context"

	schema "github.com/bitbiz/hias-core/domains/identity/schema"
	reinsuranceSchema "github.com/bitbiz/hias-core/domains/reinsurance/schema"
	"github.com/google/uuid"
)

type ReinsurerStatementService interface {
	GenerateStatement(ctx context.Context, req reinsuranceSchema.GenerateStatementRequest, createdBy uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.ReinsurerStatementResponse]
	GetStatement(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.ReinsurerStatementResponse]
	ListByTreaty(ctx context.Context, treatyID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]reinsuranceSchema.ReinsurerStatementResponse]
	IssueStatement(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.ReinsurerStatementResponse]
	AcknowledgeStatement(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.ReinsurerStatementResponse]
	SettleStatement(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.ReinsurerStatementResponse]
	CalculateProfitCommission(ctx context.Context, req reinsuranceSchema.CalculateProfitCommissionRequest) *schema.ServiceResponse[reinsuranceSchema.ProfitCommissionCalculationResponse]
}
