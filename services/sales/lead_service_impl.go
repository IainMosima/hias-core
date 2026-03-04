package sales

import (
	"context"
	"fmt"
	"log"
	"net/http"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/domains/sales/entity"
	salesRepo "github.com/bitbiz/hias-core/domains/sales/repository"
	salesSchema "github.com/bitbiz/hias-core/domains/sales/schema"
	"github.com/bitbiz/hias-core/domains/sales/service"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/google/uuid"
)

type leadServiceImpl struct {
	leadRepo     salesRepo.LeadRepository
	activityRepo salesRepo.LeadActivityRepository
	auditSvc     auditService.AuditService
}

func NewLeadService(
	leadRepo salesRepo.LeadRepository,
	activityRepo salesRepo.LeadActivityRepository,
	auditSvc auditService.AuditService,
) service.LeadService {
	return &leadServiceImpl{
		leadRepo:     leadRepo,
		activityRepo: activityRepo,
		auditSvc:     auditSvc,
	}
}

func (s *leadServiceImpl) CreateLead(ctx context.Context, req salesSchema.CreateLeadRequest, createdBy uuid.UUID) *schema.ServiceResponse[salesSchema.LeadResponse] {
	lead := &entity.Lead{
		LeadNumber:         utils.GenerateLeadNumber(),
		ContactName:        req.ContactName,
		ContactEmail:       req.ContactEmail,
		ContactPhone:       req.ContactPhone,
		CompanyName:        req.CompanyName,
		Source:             req.Source,
		Segment:            req.Segment,
		PlanType:           req.PlanType,
		EstimatedMembers:   req.EstimatedMembers,
		ExpectedPremium:    req.ExpectedPremium,
		Currency:           string(shared.CurrencyKES),
		Status:             string(shared.LeadStatusNew),
		AssignedTo:         createdBy,
		ClosureProbability: req.ClosureProbability,
		NextFollowUpDate:   req.NextFollowUpDate,
		Notes:              req.Notes,
		CreatedBy:          createdBy,
	}

	if lead.EstimatedMembers == 0 {
		lead.EstimatedMembers = 1
	}

	created, err := s.leadRepo.Create(ctx, lead)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.LeadResponse](http.StatusInternalServerError, "Failed to create lead", err)
	}

	s.logAudit(ctx, createdBy, string(shared.AuditEntityTypeLead), created.ID, string(shared.AuditActionCreate))

	return schema.NewServiceResponse(salesSchema.ToLeadResponse(created), http.StatusCreated, "Lead created")
}

func (s *leadServiceImpl) GetLead(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[salesSchema.LeadResponse] {
	lead, err := s.leadRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.LeadResponse](http.StatusNotFound, "Lead not found", err)
	}
	return schema.NewServiceResponse(salesSchema.ToLeadResponse(lead), http.StatusOK, "Lead retrieved")
}

func (s *leadServiceImpl) ListLeads(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]salesSchema.LeadResponse] {
	offset := (page - 1) * pageSize
	leads, err := s.leadRepo.List(ctx, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]salesSchema.LeadResponse](http.StatusInternalServerError, "Failed to list leads", err)
	}
	responses := make([]salesSchema.LeadResponse, len(leads))
	for i, l := range leads {
		responses[i] = salesSchema.ToLeadResponse(l)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Leads retrieved")
}

func (s *leadServiceImpl) ListLeadsByStatus(ctx context.Context, status string, page, pageSize int) *schema.ServiceResponse[[]salesSchema.LeadResponse] {
	offset := (page - 1) * pageSize
	leads, err := s.leadRepo.ListByStatus(ctx, status, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]salesSchema.LeadResponse](http.StatusInternalServerError, "Failed to list leads by status", err)
	}
	responses := make([]salesSchema.LeadResponse, len(leads))
	for i, l := range leads {
		responses[i] = salesSchema.ToLeadResponse(l)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Leads retrieved")
}

func (s *leadServiceImpl) ListMyLeads(ctx context.Context, userID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]salesSchema.LeadResponse] {
	offset := (page - 1) * pageSize
	leads, err := s.leadRepo.ListByAssignedTo(ctx, userID, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]salesSchema.LeadResponse](http.StatusInternalServerError, "Failed to list my leads", err)
	}
	responses := make([]salesSchema.LeadResponse, len(leads))
	for i, l := range leads {
		responses[i] = salesSchema.ToLeadResponse(l)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Leads retrieved")
}

func (s *leadServiceImpl) UpdateLead(ctx context.Context, id uuid.UUID, req salesSchema.UpdateLeadRequest, updatedBy uuid.UUID) *schema.ServiceResponse[salesSchema.LeadResponse] {
	existing, err := s.leadRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.LeadResponse](http.StatusNotFound, "Lead not found", err)
	}

	// Apply updates
	if req.ContactName != "" {
		existing.ContactName = req.ContactName
	}
	if req.ContactEmail != "" {
		existing.ContactEmail = req.ContactEmail
	}
	if req.ContactPhone != "" {
		existing.ContactPhone = req.ContactPhone
	}
	if req.CompanyName != "" {
		existing.CompanyName = req.CompanyName
	}
	if req.Source != "" {
		existing.Source = req.Source
	}
	if req.Segment != "" {
		existing.Segment = req.Segment
	}
	if req.PlanType != "" {
		existing.PlanType = req.PlanType
	}
	if req.EstimatedMembers != nil {
		existing.EstimatedMembers = *req.EstimatedMembers
	}
	if req.ExpectedPremium != nil {
		existing.ExpectedPremium = *req.ExpectedPremium
	}
	if req.ClosureProbability != nil {
		existing.ClosureProbability = *req.ClosureProbability
	}
	if req.AssignedTo != "" {
		assignedTo, _ := uuid.Parse(req.AssignedTo)
		existing.AssignedTo = assignedTo
	}
	if req.NextFollowUpDate != nil {
		existing.NextFollowUpDate = req.NextFollowUpDate
	}
	if req.Notes != "" {
		existing.Notes = req.Notes
	}

	updated, err := s.leadRepo.Update(ctx, existing)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.LeadResponse](http.StatusInternalServerError, "Failed to update lead", err)
	}

	s.logAudit(ctx, updatedBy, string(shared.AuditEntityTypeLead), id, string(shared.AuditActionUpdate))

	return schema.NewServiceResponse(salesSchema.ToLeadResponse(updated), http.StatusOK, "Lead updated")
}

func (s *leadServiceImpl) UpdateLeadStatus(ctx context.Context, id uuid.UUID, req salesSchema.UpdateLeadStatusRequest, updatedBy uuid.UUID) *schema.ServiceResponse[salesSchema.LeadResponse] {
	existing, err := s.leadRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.LeadResponse](http.StatusNotFound, "Lead not found", err)
	}

	// Validate transition: can't go from terminal states back to NEW
	terminalStatuses := map[string]bool{
		string(shared.LeadStatusWon):  true,
		string(shared.LeadStatusLost): true,
	}
	if terminalStatuses[existing.Status] && req.Status == string(shared.LeadStatusNew) {
		return schema.NewServiceErrorResponse[salesSchema.LeadResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Cannot transition from %s to %s", existing.Status, req.Status),
			nil,
		)
	}

	updated, err := s.leadRepo.UpdateStatus(ctx, id, req.Status)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.LeadResponse](http.StatusInternalServerError, "Failed to update lead status", err)
	}

	// Auto-log status change activity
	activity := &entity.LeadActivity{
		LeadID:       id,
		ActivityType: string(shared.LeadActivityTypeStatusChange),
		Description:  fmt.Sprintf("Status changed from %s to %s", existing.Status, req.Status),
		CreatedBy:    updatedBy,
	}
	if _, actErr := s.activityRepo.Create(ctx, activity); actErr != nil {
		log.Printf("Failed to log status change activity: %v", actErr)
	}

	s.logAudit(ctx, updatedBy, string(shared.AuditEntityTypeLead), id, string(shared.AuditActionStateChange))

	return schema.NewServiceResponse(salesSchema.ToLeadResponse(updated), http.StatusOK, "Lead status updated")
}

func (s *leadServiceImpl) AddActivity(ctx context.Context, leadID uuid.UUID, req salesSchema.CreateLeadActivityRequest, createdBy uuid.UUID) *schema.ServiceResponse[salesSchema.LeadActivityResponse] {
	// Verify lead exists
	_, err := s.leadRepo.GetByID(ctx, leadID)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.LeadActivityResponse](http.StatusNotFound, "Lead not found", err)
	}

	activity := &entity.LeadActivity{
		LeadID:       leadID,
		ActivityType: req.ActivityType,
		Description:  req.Description,
		ScheduledAt:  req.ScheduledAt,
		CompletedAt:  req.CompletedAt,
		CreatedBy:    createdBy,
	}

	created, err := s.activityRepo.Create(ctx, activity)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.LeadActivityResponse](http.StatusInternalServerError, "Failed to create activity", err)
	}

	// If follow_up with scheduled_at, update lead's next_follow_up_date
	if req.ActivityType == string(shared.LeadActivityTypeFollowUp) && req.ScheduledAt != nil {
		lead := &entity.Lead{
			ID:               leadID,
			NextFollowUpDate: req.ScheduledAt,
		}
		if _, updateErr := s.leadRepo.Update(ctx, lead); updateErr != nil {
			log.Printf("Failed to update lead follow-up date: %v", updateErr)
		}
	}

	return schema.NewServiceResponse(salesSchema.ToLeadActivityResponse(created), http.StatusCreated, "Activity created")
}

func (s *leadServiceImpl) ListActivities(ctx context.Context, leadID uuid.UUID) *schema.ServiceResponse[[]salesSchema.LeadActivityResponse] {
	activities, err := s.activityRepo.ListByLead(ctx, leadID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]salesSchema.LeadActivityResponse](http.StatusInternalServerError, "Failed to list activities", err)
	}
	responses := make([]salesSchema.LeadActivityResponse, len(activities))
	for i, a := range activities {
		responses[i] = salesSchema.ToLeadActivityResponse(a)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Activities retrieved")
}

func (s *leadServiceImpl) GetDueFollowUps(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]salesSchema.LeadResponse] {
	offset := (page - 1) * pageSize
	leads, err := s.leadRepo.ListDueFollowUps(ctx, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]salesSchema.LeadResponse](http.StatusInternalServerError, "Failed to get due follow-ups", err)
	}
	responses := make([]salesSchema.LeadResponse, len(leads))
	for i, l := range leads {
		responses[i] = salesSchema.ToLeadResponse(l)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Due follow-ups retrieved")
}

func (s *leadServiceImpl) GetTotalCount(ctx context.Context) *schema.ServiceResponse[int64] {
	count, err := s.leadRepo.Count(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to get count", err)
	}
	return schema.NewServiceResponse(count, http.StatusOK, "Count retrieved")
}

func (s *leadServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}
