package service

import (
	"context"

	"github.com/bitbiz/hias-core/domains/identity/schema"
	salesSchema "github.com/bitbiz/hias-core/domains/sales/schema"
	"github.com/google/uuid"
)

type LeadService interface {
	CreateLead(ctx context.Context, req salesSchema.CreateLeadRequest, createdBy uuid.UUID) *schema.ServiceResponse[salesSchema.LeadResponse]
	GetLead(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[salesSchema.LeadResponse]
	ListLeads(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]salesSchema.LeadResponse]
	ListLeadsByStatus(ctx context.Context, status string, page, pageSize int) *schema.ServiceResponse[[]salesSchema.LeadResponse]
	ListMyLeads(ctx context.Context, userID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]salesSchema.LeadResponse]
	UpdateLead(ctx context.Context, id uuid.UUID, req salesSchema.UpdateLeadRequest, updatedBy uuid.UUID) *schema.ServiceResponse[salesSchema.LeadResponse]
	UpdateLeadStatus(ctx context.Context, id uuid.UUID, req salesSchema.UpdateLeadStatusRequest, updatedBy uuid.UUID) *schema.ServiceResponse[salesSchema.LeadResponse]
	AddActivity(ctx context.Context, leadID uuid.UUID, req salesSchema.CreateLeadActivityRequest, createdBy uuid.UUID) *schema.ServiceResponse[salesSchema.LeadActivityResponse]
	ListActivities(ctx context.Context, leadID uuid.UUID) *schema.ServiceResponse[[]salesSchema.LeadActivityResponse]
	GetDueFollowUps(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]salesSchema.LeadResponse]
	GetTotalCount(ctx context.Context) *schema.ServiceResponse[int64]
}
