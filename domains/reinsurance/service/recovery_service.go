package service

import (
	"context"

	schema "github.com/bitbiz/hias-core/domains/identity/schema"
	reinsuranceSchema "github.com/bitbiz/hias-core/domains/reinsurance/schema"
	"github.com/google/uuid"
)

type RecoveryService interface {
	CreateRecovery(ctx context.Context, req reinsuranceSchema.CreateRecoveryRequest, createdBy uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.RecoveryResponse]
	GetRecovery(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.RecoveryDetailResponse]
	ListByClaim(ctx context.Context, claimID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]reinsuranceSchema.RecoveryResponse]
	ListByTreaty(ctx context.Context, treatyID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]reinsuranceSchema.RecoveryResponse]
	ListOutstanding(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]reinsuranceSchema.RecoveryResponse]
	ApplyRecoveryForClaim(ctx context.Context, claimID uuid.UUID, req reinsuranceSchema.ApplyRecoveryForClaimRequest, createdBy uuid.UUID) *schema.ServiceResponse[[]reinsuranceSchema.RecoveryResponse]

	// Workflow transitions
	AcknowledgeRecovery(ctx context.Context, id uuid.UUID, req reinsuranceSchema.RecoveryWorkflowRequest, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.RecoveryResponse]
	RequestInfo(ctx context.Context, id uuid.UUID, req reinsuranceSchema.RecoveryWorkflowRequest, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.RecoveryResponse]
	ApproveRecovery(ctx context.Context, id uuid.UUID, req reinsuranceSchema.RecoveryWorkflowRequest, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.RecoveryResponse]
	RecordPayment(ctx context.Context, id uuid.UUID, req reinsuranceSchema.RecordPaymentRequest, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.RecoveryResponse]
	WriteOffRecovery(ctx context.Context, id uuid.UUID, req reinsuranceSchema.RecoveryWorkflowRequest, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.RecoveryResponse]

	GetWorkflowEvents(ctx context.Context, recoveryID uuid.UUID) *schema.ServiceResponse[[]reinsuranceSchema.RecoveryWorkflowEventResponse]
	GetRecoveryCount(ctx context.Context) *schema.ServiceResponse[int64]
	GetAgedAnalysis(ctx context.Context) *schema.ServiceResponse[[]reinsuranceSchema.AgedRecoveryBucketResponse]
	GetTotalRecoverableAmount(ctx context.Context) *schema.ServiceResponse[int64]
	GetTotalRecoveredAmount(ctx context.Context) *schema.ServiceResponse[int64]
}
