package service

import (
	"context"

	claimsSchema "github.com/bitbiz/hias-core/domains/claims/schema"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/google/uuid"
)

type CaseService interface {
	CreateCase(ctx context.Context, req claimsSchema.CreateCaseRequest, createdBy uuid.UUID) *schema.ServiceResponse[claimsSchema.CaseRecordResponse]
	GetCase(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[claimsSchema.CaseRecordResponse]
	ListByPolicy(ctx context.Context, policyID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]claimsSchema.CaseRecordResponse]
	ListByMember(ctx context.Context, memberID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]claimsSchema.CaseRecordResponse]
	ListByProvider(ctx context.Context, providerID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]claimsSchema.CaseRecordResponse]
	ListByStatus(ctx context.Context, status string, page, pageSize int) *schema.ServiceResponse[[]claimsSchema.CaseRecordResponse]
	AdmitCase(ctx context.Context, id uuid.UUID, req claimsSchema.AdmitCaseRequest, userID uuid.UUID) *schema.ServiceResponse[claimsSchema.CaseRecordResponse]
	UpdateCase(ctx context.Context, id uuid.UUID, req claimsSchema.UpdateCaseRequest, userID uuid.UUID) *schema.ServiceResponse[claimsSchema.CaseRecordResponse]
	DischargeCase(ctx context.Context, id uuid.UUID, req claimsSchema.DischargeCaseRequest, userID uuid.UUID) *schema.ServiceResponse[claimsSchema.CaseRecordResponse]
	StartTreatment(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[claimsSchema.CaseRecordResponse]
	CloseCase(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[claimsSchema.CaseRecordResponse]
	CountByStatus(ctx context.Context, status string) *schema.ServiceResponse[int64]
}
