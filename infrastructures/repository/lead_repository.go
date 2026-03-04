package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/sales/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/sales/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type leadRepository struct {
	store db.Store
}

func NewLeadRepository(store db.Store) domainRepo.LeadRepository {
	return &leadRepository{store: store}
}

func (r *leadRepository) Create(ctx context.Context, lead *entity.Lead) (*entity.Lead, error) {
	dbLead, err := r.store.CreateLead(ctx, db.CreateLeadParams{
		LeadNumber:         lead.LeadNumber,
		ContactName:        lead.ContactName,
		ContactEmail:       stringToPgtypeText(lead.ContactEmail),
		ContactPhone:       stringToPgtypeText(lead.ContactPhone),
		CompanyName:        stringToPgtypeText(lead.CompanyName),
		Source:             lead.Source,
		Segment:            lead.Segment,
		PlanType:           lead.PlanType,
		EstimatedMembers:   int32(lead.EstimatedMembers),
		ExpectedPremium:    lead.ExpectedPremium,
		ClosureProbability: int32(lead.ClosureProbability),
		Currency:           lead.Currency,
		Status:             lead.Status,
		AssignedTo:         uuidToPgtype(lead.AssignedTo),
		NextFollowUpDate:   timePtrToPgtypeTimestamptz(lead.NextFollowUpDate),
		Notes:              stringToPgtypeText(lead.Notes),
		CreatedBy:          uuidToPgtype(lead.CreatedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create lead: %w", err)
	}
	return sqlcLeadToDomain(dbLead), nil
}

func (r *leadRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Lead, error) {
	dbLead, err := r.store.GetLeadByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get lead by ID: %w", err)
	}
	return sqlcLeadToDomain(dbLead), nil
}

func (r *leadRepository) GetByNumber(ctx context.Context, number string) (*entity.Lead, error) {
	dbLead, err := r.store.GetLeadByNumber(ctx, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get lead by number: %w", err)
	}
	return sqlcLeadToDomain(dbLead), nil
}

func (r *leadRepository) List(ctx context.Context, limit, offset int) ([]*entity.Lead, error) {
	dbLeads, err := r.store.ListLeads(ctx, db.ListLeadsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list leads: %w", err)
	}
	leads := make([]*entity.Lead, len(dbLeads))
	for i, l := range dbLeads {
		leads[i] = sqlcLeadToDomain(l)
	}
	return leads, nil
}

func (r *leadRepository) ListByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.Lead, error) {
	dbLeads, err := r.store.ListLeadsByStatus(ctx, db.ListLeadsByStatusParams{
		Status: status,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list leads by status: %w", err)
	}
	leads := make([]*entity.Lead, len(dbLeads))
	for i, l := range dbLeads {
		leads[i] = sqlcLeadToDomain(l)
	}
	return leads, nil
}

func (r *leadRepository) ListByAssignedTo(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Lead, error) {
	dbLeads, err := r.store.ListLeadsByAssignedTo(ctx, db.ListLeadsByAssignedToParams{
		AssignedTo: uuidToPgtype(userID),
		Limit:      int32(limit),
		Offset:     int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list leads by assigned to: %w", err)
	}
	leads := make([]*entity.Lead, len(dbLeads))
	for i, l := range dbLeads {
		leads[i] = sqlcLeadToDomain(l)
	}
	return leads, nil
}

func (r *leadRepository) ListDueFollowUps(ctx context.Context, limit, offset int) ([]*entity.Lead, error) {
	dbLeads, err := r.store.ListDueFollowUps(ctx, db.ListDueFollowUpsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list due follow ups: %w", err)
	}
	leads := make([]*entity.Lead, len(dbLeads))
	for i, l := range dbLeads {
		leads[i] = sqlcLeadToDomain(l)
	}
	return leads, nil
}

func (r *leadRepository) Update(ctx context.Context, lead *entity.Lead) (*entity.Lead, error) {
	dbLead, err := r.store.UpdateLead(ctx, db.UpdateLeadParams{
		ID:                 lead.ID,
		Column2:            lead.ContactName,
		Column3:            lead.ContactEmail,
		Column4:            lead.ContactPhone,
		Column5:            lead.CompanyName,
		Column6:            lead.Source,
		Column7:            lead.Segment,
		Column8:            lead.PlanType,
		EstimatedMembers:   int32(lead.EstimatedMembers),
		ExpectedPremium:    lead.ExpectedPremium,
		ClosureProbability: int32(lead.ClosureProbability),
		AssignedTo:         uuidToPgtype(lead.AssignedTo),
		NextFollowUpDate:   timePtrToPgtypeTimestamptz(lead.NextFollowUpDate),
		Column14:           lead.Notes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update lead: %w", err)
	}
	return sqlcLeadToDomain(dbLead), nil
}

func (r *leadRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Lead, error) {
	dbLead, err := r.store.UpdateLeadStatus(ctx, db.UpdateLeadStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update lead status: %w", err)
	}
	return sqlcLeadToDomain(dbLead), nil
}

func (r *leadRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.store.CountLeads(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count leads: %w", err)
	}
	return count, nil
}

func (r *leadRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	count, err := r.store.CountLeadsByStatus(ctx, status)
	if err != nil {
		return 0, fmt.Errorf("failed to count leads by status: %w", err)
	}
	return count, nil
}

func sqlcLeadToDomain(l db.Lead) *entity.Lead {
	return &entity.Lead{
		ID:                 l.ID,
		LeadNumber:         l.LeadNumber,
		ContactName:        l.ContactName,
		ContactEmail:       l.ContactEmail.String,
		ContactPhone:       l.ContactPhone.String,
		CompanyName:        l.CompanyName.String,
		Source:             l.Source,
		Segment:            l.Segment,
		PlanType:           l.PlanType,
		EstimatedMembers:   int(l.EstimatedMembers),
		ExpectedPremium:    l.ExpectedPremium,
		ClosureProbability: int(l.ClosureProbability),
		Currency:           l.Currency,
		Status:             l.Status,
		AssignedTo:         pgtypeToUUID(l.AssignedTo),
		NextFollowUpDate:   pgtypeTimestamptzToTimePtr(l.NextFollowUpDate),
		Notes:              l.Notes.String,
		CreatedBy:          pgtypeToUUID(l.CreatedBy),
		CreatedAt:          l.CreatedAt,
		UpdatedAt:          l.UpdatedAt,
	}
}
