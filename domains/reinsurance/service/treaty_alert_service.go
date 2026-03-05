package service

import (
	"context"

	schema "github.com/bitbiz/hias-core/domains/identity/schema"
	reinsuranceSchema "github.com/bitbiz/hias-core/domains/reinsurance/schema"
	"github.com/google/uuid"
)

type TreatyAlertService interface {
	CheckTreatyLimits(ctx context.Context, treatyID uuid.UUID) *schema.ServiceResponse[[]reinsuranceSchema.TreatyAlertResponse]
	CheckCatastropheThresholds(ctx context.Context, treatyID uuid.UUID) *schema.ServiceResponse[[]reinsuranceSchema.TreatyAlertResponse]
	ListAlerts(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]reinsuranceSchema.TreatyAlertResponse]
	ListByTreaty(ctx context.Context, treatyID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]reinsuranceSchema.TreatyAlertResponse]
	ListUnacknowledged(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]reinsuranceSchema.TreatyAlertResponse]
	AcknowledgeAlert(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.TreatyAlertResponse]
	CountUnacknowledged(ctx context.Context) *schema.ServiceResponse[int64]
	CheckExpiryWarnings(ctx context.Context) *schema.ServiceResponse[[]reinsuranceSchema.TreatyAlertResponse]
}
