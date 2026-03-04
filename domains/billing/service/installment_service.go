package service

import (
	"context"

	"github.com/bitbiz/hias-core/domains/billing/schema"
	identitySchema "github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/google/uuid"
)

type InstallmentService interface {
	CreateSchedule(ctx context.Context, req schema.CreateInstallmentScheduleRequest, createdBy uuid.UUID) *identitySchema.ServiceResponse[schema.InstallmentScheduleResponse]
	GetSchedulesByPolicy(ctx context.Context, policyID uuid.UUID) *identitySchema.ServiceResponse[[]schema.InstallmentScheduleResponse]
	ListInstallmentsBySchedule(ctx context.Context, scheduleID uuid.UUID) *identitySchema.ServiceResponse[[]schema.InstallmentResponse]
	MarkInstallmentPaid(ctx context.Context, installmentID uuid.UUID, invoiceID uuid.UUID) *identitySchema.ServiceResponse[schema.InstallmentResponse]
}
