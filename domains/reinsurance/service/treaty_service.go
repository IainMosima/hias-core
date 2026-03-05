package service

import (
	"context"

	schema "github.com/bitbiz/hias-core/domains/identity/schema"
	reinsuranceSchema "github.com/bitbiz/hias-core/domains/reinsurance/schema"
	"github.com/google/uuid"
)

type TreatyService interface {
	CreateTreaty(ctx context.Context, req reinsuranceSchema.CreateTreatyRequest, createdBy uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.TreatyResponse]
	GetTreaty(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.TreatyDetailResponse]
	ListTreaties(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]reinsuranceSchema.TreatyResponse]
	ListTreatiesByStatus(ctx context.Context, status string, page, pageSize int) *schema.ServiceResponse[[]reinsuranceSchema.TreatyResponse]
	ListTreatiesByType(ctx context.Context, treatyType string, page, pageSize int) *schema.ServiceResponse[[]reinsuranceSchema.TreatyResponse]
	UpdateTreaty(ctx context.Context, id uuid.UUID, req reinsuranceSchema.UpdateTreatyRequest, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.TreatyResponse]
	ActivateTreaty(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.TreatyResponse]
	TerminateTreaty(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.TreatyResponse]
	ExpireOverdue(ctx context.Context) *schema.ServiceResponse[int64]
	GetTreatyCount(ctx context.Context) *schema.ServiceResponse[int64]

	// Participants
	AddParticipant(ctx context.Context, treatyID uuid.UUID, req reinsuranceSchema.AddParticipantRequest) *schema.ServiceResponse[reinsuranceSchema.TreatyParticipantResponse]
	ListParticipants(ctx context.Context, treatyID uuid.UUID) *schema.ServiceResponse[[]reinsuranceSchema.TreatyParticipantResponse]
	UpdateParticipant(ctx context.Context, id uuid.UUID, req reinsuranceSchema.UpdateParticipantRequest) *schema.ServiceResponse[reinsuranceSchema.TreatyParticipantResponse]
	RemoveParticipant(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[bool]

	// Layers
	AddLayer(ctx context.Context, treatyID uuid.UUID, req reinsuranceSchema.AddLayerRequest) *schema.ServiceResponse[reinsuranceSchema.TreatyLayerResponse]
	ListLayers(ctx context.Context, treatyID uuid.UUID) *schema.ServiceResponse[[]reinsuranceSchema.TreatyLayerResponse]
	UpdateLayer(ctx context.Context, id uuid.UUID, req reinsuranceSchema.UpdateLayerRequest) *schema.ServiceResponse[reinsuranceSchema.TreatyLayerResponse]
	RemoveLayer(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[bool]

	// Profit Commission Rules
	AddProfitCommissionRule(ctx context.Context, treatyID uuid.UUID, req reinsuranceSchema.AddProfitCommissionRuleRequest) *schema.ServiceResponse[reinsuranceSchema.ProfitCommissionResponse]
	ListProfitCommissionRules(ctx context.Context, treatyID uuid.UUID) *schema.ServiceResponse[[]reinsuranceSchema.ProfitCommissionResponse]
	RemoveProfitCommissionRule(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[bool]
}
