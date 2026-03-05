package reinsurance

import (
	"context"
	"fmt"
	"log"
	"net/http"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	schema "github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/domains/reinsurance/entity"
	"github.com/bitbiz/hias-core/domains/reinsurance/repository"
	reinsuranceSchema "github.com/bitbiz/hias-core/domains/reinsurance/schema"
	"github.com/bitbiz/hias-core/domains/reinsurance/service"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/google/uuid"
)

func logAudit(ctx context.Context, auditSvc auditService.AuditService, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if auditSvc != nil {
		resp := auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}

type treatyServiceImpl struct {
	treatyRepo      repository.TreatyRepository
	participantRepo repository.TreatyParticipantRepository
	layerRepo       repository.TreatyLayerRepository
	profitCommRepo  repository.ProfitCommissionRepository
	auditSvc        auditService.AuditService
}

func NewTreatyService(
	treatyRepo repository.TreatyRepository,
	participantRepo repository.TreatyParticipantRepository,
	layerRepo repository.TreatyLayerRepository,
	profitCommRepo repository.ProfitCommissionRepository,
	auditSvc auditService.AuditService,
) service.TreatyService {
	return &treatyServiceImpl{
		treatyRepo:      treatyRepo,
		participantRepo: participantRepo,
		layerRepo:       layerRepo,
		profitCommRepo:  profitCommRepo,
		auditSvc:        auditSvc,
	}
}

func (s *treatyServiceImpl) CreateTreaty(ctx context.Context, req reinsuranceSchema.CreateTreatyRequest, createdBy uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.TreatyResponse] {
	currency := req.Currency
	if currency == "" {
		currency = string(shared.CurrencyKES)
	}

	treaty := &entity.Treaty{
		TreatyNumber:   utils.GenerateTreatyNumber(),
		Name:           req.Name,
		TreatyType:     req.TreatyType,
		Status:         string(shared.TreatyStatusDraft),
		EffectiveDate:  req.EffectiveDate,
		ExpiryDate:     req.ExpiryDate,
		RetentionLimit: req.RetentionLimit,
		Currency:       currency,
		Notes:          req.Notes,
		CreatedBy:      createdBy,
	}

	created, err := s.treatyRepo.Create(ctx, treaty)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.TreatyResponse](http.StatusInternalServerError, "Failed to create treaty", err)
	}

	logAudit(ctx, s.auditSvc, createdBy, string(shared.AuditEntityTypeTreaty), created.ID, string(shared.AuditActionCreate))

	resp := reinsuranceSchema.ToTreatyResponse(created)
	return schema.NewServiceResponse(resp, http.StatusCreated, "Treaty created successfully")
}

func (s *treatyServiceImpl) GetTreaty(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.TreatyDetailResponse] {
	treaty, err := s.treatyRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.TreatyDetailResponse](http.StatusNotFound, "Treaty not found", err)
	}

	participants, _ := s.participantRepo.ListByTreaty(ctx, id)
	layers, _ := s.layerRepo.ListByTreaty(ctx, id)
	rules, _ := s.profitCommRepo.ListByTreaty(ctx, id)

	participantResponses := make([]reinsuranceSchema.TreatyParticipantResponse, len(participants))
	for i, p := range participants {
		participantResponses[i] = reinsuranceSchema.ToTreatyParticipantResponse(p)
	}

	layerResponses := make([]reinsuranceSchema.TreatyLayerResponse, len(layers))
	for i, l := range layers {
		layerResponses[i] = reinsuranceSchema.ToTreatyLayerResponse(l)
	}

	ruleResponses := make([]reinsuranceSchema.ProfitCommissionResponse, len(rules))
	for i, r := range rules {
		ruleResponses[i] = reinsuranceSchema.ToProfitCommissionResponse(r)
	}

	resp := reinsuranceSchema.TreatyDetailResponse{
		TreatyResponse:  reinsuranceSchema.ToTreatyResponse(treaty),
		Participants:    participantResponses,
		Layers:          layerResponses,
		ProfitCommRules: ruleResponses,
	}
	return schema.NewServiceResponse(resp, http.StatusOK, "Treaty retrieved successfully")
}

func (s *treatyServiceImpl) ListTreaties(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]reinsuranceSchema.TreatyResponse] {
	offset := (page - 1) * pageSize
	treaties, err := s.treatyRepo.List(ctx, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]reinsuranceSchema.TreatyResponse](http.StatusInternalServerError, "Failed to list treaties", err)
	}

	responses := make([]reinsuranceSchema.TreatyResponse, len(treaties))
	for i, t := range treaties {
		responses[i] = reinsuranceSchema.ToTreatyResponse(t)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Treaties retrieved successfully")
}

func (s *treatyServiceImpl) ListTreatiesByStatus(ctx context.Context, status string, page, pageSize int) *schema.ServiceResponse[[]reinsuranceSchema.TreatyResponse] {
	offset := (page - 1) * pageSize
	treaties, err := s.treatyRepo.ListByStatus(ctx, status, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]reinsuranceSchema.TreatyResponse](http.StatusInternalServerError, "Failed to list treaties by status", err)
	}

	responses := make([]reinsuranceSchema.TreatyResponse, len(treaties))
	for i, t := range treaties {
		responses[i] = reinsuranceSchema.ToTreatyResponse(t)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Treaties retrieved successfully")
}

func (s *treatyServiceImpl) ListTreatiesByType(ctx context.Context, treatyType string, page, pageSize int) *schema.ServiceResponse[[]reinsuranceSchema.TreatyResponse] {
	offset := (page - 1) * pageSize
	treaties, err := s.treatyRepo.ListByType(ctx, treatyType, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]reinsuranceSchema.TreatyResponse](http.StatusInternalServerError, "Failed to list treaties by type", err)
	}

	responses := make([]reinsuranceSchema.TreatyResponse, len(treaties))
	for i, t := range treaties {
		responses[i] = reinsuranceSchema.ToTreatyResponse(t)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Treaties retrieved successfully")
}

func (s *treatyServiceImpl) UpdateTreaty(ctx context.Context, id uuid.UUID, req reinsuranceSchema.UpdateTreatyRequest, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.TreatyResponse] {
	existing, err := s.treatyRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.TreatyResponse](http.StatusNotFound, "Treaty not found", err)
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if !req.EffectiveDate.IsZero() {
		existing.EffectiveDate = req.EffectiveDate
	}
	if !req.ExpiryDate.IsZero() {
		existing.ExpiryDate = req.ExpiryDate
	}
	if req.RetentionLimit > 0 {
		existing.RetentionLimit = req.RetentionLimit
	}
	if req.Currency != "" {
		existing.Currency = req.Currency
	}
	if req.Notes != "" {
		existing.Notes = req.Notes
	}

	updated, err := s.treatyRepo.Update(ctx, existing)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.TreatyResponse](http.StatusInternalServerError, "Failed to update treaty", err)
	}

	logAudit(ctx, s.auditSvc, userID, string(shared.AuditEntityTypeTreaty), id, string(shared.AuditActionUpdate))

	resp := reinsuranceSchema.ToTreatyResponse(updated)
	return schema.NewServiceResponse(resp, http.StatusOK, "Treaty updated successfully")
}

func (s *treatyServiceImpl) ActivateTreaty(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.TreatyResponse] {
	treaty, err := s.treatyRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.TreatyResponse](http.StatusNotFound, "Treaty not found", err)
	}

	if treaty.Status != string(shared.TreatyStatusDraft) {
		return schema.NewServiceErrorResponse[reinsuranceSchema.TreatyResponse](http.StatusBadRequest, "Only DRAFT treaties can be activated", fmt.Errorf("invalid status: %s", treaty.Status))
	}

	// Validate participants exist and total share <= 100
	totalShare, err := s.participantRepo.GetTotalShareByTreaty(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.TreatyResponse](http.StatusInternalServerError, "Failed to validate participants", err)
	}
	if totalShare <= 0 {
		return schema.NewServiceErrorResponse[reinsuranceSchema.TreatyResponse](http.StatusBadRequest, "Treaty must have at least one participant", fmt.Errorf("no participants"))
	}
	if totalShare > 100 {
		return schema.NewServiceErrorResponse[reinsuranceSchema.TreatyResponse](http.StatusBadRequest, "Total participant share exceeds 100%", fmt.Errorf("total share: %.2f", totalShare))
	}

	updated, err := s.treatyRepo.UpdateStatus(ctx, id, string(shared.TreatyStatusActive))
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.TreatyResponse](http.StatusInternalServerError, "Failed to activate treaty", err)
	}

	logAudit(ctx, s.auditSvc, userID, string(shared.AuditEntityTypeTreaty), id, string(shared.AuditActionStateChange))

	resp := reinsuranceSchema.ToTreatyResponse(updated)
	return schema.NewServiceResponse(resp, http.StatusOK, "Treaty activated successfully")
}

func (s *treatyServiceImpl) TerminateTreaty(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.TreatyResponse] {
	treaty, err := s.treatyRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.TreatyResponse](http.StatusNotFound, "Treaty not found", err)
	}

	if treaty.Status != string(shared.TreatyStatusActive) {
		return schema.NewServiceErrorResponse[reinsuranceSchema.TreatyResponse](http.StatusBadRequest, "Only ACTIVE treaties can be terminated", fmt.Errorf("invalid status: %s", treaty.Status))
	}

	updated, err := s.treatyRepo.UpdateStatus(ctx, id, string(shared.TreatyStatusTerminated))
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.TreatyResponse](http.StatusInternalServerError, "Failed to terminate treaty", err)
	}

	logAudit(ctx, s.auditSvc, userID, string(shared.AuditEntityTypeTreaty), id, string(shared.AuditActionStateChange))

	resp := reinsuranceSchema.ToTreatyResponse(updated)
	return schema.NewServiceResponse(resp, http.StatusOK, "Treaty terminated successfully")
}

func (s *treatyServiceImpl) ExpireOverdue(ctx context.Context) *schema.ServiceResponse[int64] {
	expired, err := s.treatyRepo.ListByStatus(ctx, string(shared.TreatyStatusActive), 1000, 0)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to list active treaties", err)
	}

	var count int64
	for _, t := range expired {
		if t.ExpiryDate.Before(t.CreatedAt) || t.Status == string(shared.TreatyStatusActive) {
			// Use ListExpiring approach - just expire treaties past expiry_date
			_, err := s.treatyRepo.UpdateStatus(ctx, t.ID, string(shared.TreatyStatusExpired))
			if err == nil {
				count++
			}
		}
	}
	return schema.NewServiceResponse(count, http.StatusOK, fmt.Sprintf("Expired %d treaties", count))
}

func (s *treatyServiceImpl) GetTreatyCount(ctx context.Context) *schema.ServiceResponse[int64] {
	count, err := s.treatyRepo.Count(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to count treaties", err)
	}
	return schema.NewServiceResponse(count, http.StatusOK, "Treaty count retrieved")
}

// Participants

func (s *treatyServiceImpl) AddParticipant(ctx context.Context, treatyID uuid.UUID, req reinsuranceSchema.AddParticipantRequest) *schema.ServiceResponse[reinsuranceSchema.TreatyParticipantResponse] {
	// Validate total share won't exceed 100%
	currentShare, _ := s.participantRepo.GetTotalShareByTreaty(ctx, treatyID)
	if currentShare+req.SharePercentage > 100 {
		return schema.NewServiceErrorResponse[reinsuranceSchema.TreatyParticipantResponse](http.StatusBadRequest, fmt.Sprintf("Adding %.2f%% would exceed 100%% (current: %.2f%%)", req.SharePercentage, currentShare), fmt.Errorf("share exceeded"))
	}

	participant := &entity.TreatyParticipant{
		TreatyID:        treatyID,
		ReinsurerName:   req.ReinsurerName,
		SharePercentage: req.SharePercentage,
		CommissionRate:  req.CommissionRate,
		IsLead:          req.IsLead,
	}

	created, err := s.participantRepo.Create(ctx, participant)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.TreatyParticipantResponse](http.StatusInternalServerError, "Failed to add participant", err)
	}

	resp := reinsuranceSchema.ToTreatyParticipantResponse(created)
	return schema.NewServiceResponse(resp, http.StatusCreated, "Participant added successfully")
}

func (s *treatyServiceImpl) ListParticipants(ctx context.Context, treatyID uuid.UUID) *schema.ServiceResponse[[]reinsuranceSchema.TreatyParticipantResponse] {
	participants, err := s.participantRepo.ListByTreaty(ctx, treatyID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]reinsuranceSchema.TreatyParticipantResponse](http.StatusInternalServerError, "Failed to list participants", err)
	}

	responses := make([]reinsuranceSchema.TreatyParticipantResponse, len(participants))
	for i, p := range participants {
		responses[i] = reinsuranceSchema.ToTreatyParticipantResponse(p)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Participants retrieved successfully")
}

func (s *treatyServiceImpl) UpdateParticipant(ctx context.Context, id uuid.UUID, req reinsuranceSchema.UpdateParticipantRequest) *schema.ServiceResponse[reinsuranceSchema.TreatyParticipantResponse] {
	existing, err := s.participantRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.TreatyParticipantResponse](http.StatusNotFound, "Participant not found", err)
	}

	if req.ReinsurerName != "" {
		existing.ReinsurerName = req.ReinsurerName
	}
	if req.SharePercentage > 0 {
		existing.SharePercentage = req.SharePercentage
	}
	if req.CommissionRate > 0 {
		existing.CommissionRate = req.CommissionRate
	}
	existing.IsLead = req.IsLead

	updated, err := s.participantRepo.Update(ctx, existing)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.TreatyParticipantResponse](http.StatusInternalServerError, "Failed to update participant", err)
	}

	resp := reinsuranceSchema.ToTreatyParticipantResponse(updated)
	return schema.NewServiceResponse(resp, http.StatusOK, "Participant updated successfully")
}

func (s *treatyServiceImpl) RemoveParticipant(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[bool] {
	err := s.participantRepo.Delete(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[bool](http.StatusInternalServerError, "Failed to remove participant", err)
	}
	return schema.NewServiceResponse(true, http.StatusOK, "Participant removed successfully")
}

// Layers

func (s *treatyServiceImpl) AddLayer(ctx context.Context, treatyID uuid.UUID, req reinsuranceSchema.AddLayerRequest) *schema.ServiceResponse[reinsuranceSchema.TreatyLayerResponse] {
	layer := &entity.TreatyLayer{
		TreatyID:         treatyID,
		LayerNumber:      req.LayerNumber,
		AttachmentPoint:  req.AttachmentPoint,
		LayerLimit:       req.LayerLimit,
		DeductibleAmount: req.DeductibleAmount,
		PremiumRate:      req.PremiumRate,
		AggregateLimit:   req.AggregateLimit,
	}

	created, err := s.layerRepo.Create(ctx, layer)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.TreatyLayerResponse](http.StatusInternalServerError, "Failed to add layer", err)
	}

	resp := reinsuranceSchema.ToTreatyLayerResponse(created)
	return schema.NewServiceResponse(resp, http.StatusCreated, "Layer added successfully")
}

func (s *treatyServiceImpl) ListLayers(ctx context.Context, treatyID uuid.UUID) *schema.ServiceResponse[[]reinsuranceSchema.TreatyLayerResponse] {
	layers, err := s.layerRepo.ListByTreaty(ctx, treatyID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]reinsuranceSchema.TreatyLayerResponse](http.StatusInternalServerError, "Failed to list layers", err)
	}

	responses := make([]reinsuranceSchema.TreatyLayerResponse, len(layers))
	for i, l := range layers {
		responses[i] = reinsuranceSchema.ToTreatyLayerResponse(l)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Layers retrieved successfully")
}

func (s *treatyServiceImpl) UpdateLayer(ctx context.Context, id uuid.UUID, req reinsuranceSchema.UpdateLayerRequest) *schema.ServiceResponse[reinsuranceSchema.TreatyLayerResponse] {
	existing, err := s.layerRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.TreatyLayerResponse](http.StatusNotFound, "Layer not found", err)
	}

	if req.AttachmentPoint > 0 {
		existing.AttachmentPoint = req.AttachmentPoint
	}
	if req.LayerLimit > 0 {
		existing.LayerLimit = req.LayerLimit
	}
	if req.DeductibleAmount >= 0 {
		existing.DeductibleAmount = req.DeductibleAmount
	}
	if req.PremiumRate > 0 {
		existing.PremiumRate = req.PremiumRate
	}
	if req.AggregateLimit != nil {
		existing.AggregateLimit = req.AggregateLimit
	}

	updated, err := s.layerRepo.Update(ctx, existing)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.TreatyLayerResponse](http.StatusInternalServerError, "Failed to update layer", err)
	}

	resp := reinsuranceSchema.ToTreatyLayerResponse(updated)
	return schema.NewServiceResponse(resp, http.StatusOK, "Layer updated successfully")
}

func (s *treatyServiceImpl) RemoveLayer(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[bool] {
	err := s.layerRepo.Delete(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[bool](http.StatusInternalServerError, "Failed to remove layer", err)
	}
	return schema.NewServiceResponse(true, http.StatusOK, "Layer removed successfully")
}

// Profit Commission Rules

func (s *treatyServiceImpl) AddProfitCommissionRule(ctx context.Context, treatyID uuid.UUID, req reinsuranceSchema.AddProfitCommissionRuleRequest) *schema.ServiceResponse[reinsuranceSchema.ProfitCommissionResponse] {
	pc := &entity.ProfitCommission{
		TreatyID:          treatyID,
		CommissionType:    req.CommissionType,
		LossRatioFrom:     req.LossRatioFrom,
		LossRatioTo:       req.LossRatioTo,
		CommissionRate:    req.CommissionRate,
		CarryForwardYears: req.CarryForwardYears,
	}

	created, err := s.profitCommRepo.Create(ctx, pc)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.ProfitCommissionResponse](http.StatusInternalServerError, "Failed to add profit commission rule", err)
	}

	resp := reinsuranceSchema.ToProfitCommissionResponse(created)
	return schema.NewServiceResponse(resp, http.StatusCreated, "Profit commission rule added successfully")
}

func (s *treatyServiceImpl) ListProfitCommissionRules(ctx context.Context, treatyID uuid.UUID) *schema.ServiceResponse[[]reinsuranceSchema.ProfitCommissionResponse] {
	rules, err := s.profitCommRepo.ListByTreaty(ctx, treatyID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]reinsuranceSchema.ProfitCommissionResponse](http.StatusInternalServerError, "Failed to list profit commission rules", err)
	}

	responses := make([]reinsuranceSchema.ProfitCommissionResponse, len(rules))
	for i, r := range rules {
		responses[i] = reinsuranceSchema.ToProfitCommissionResponse(r)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Profit commission rules retrieved successfully")
}

func (s *treatyServiceImpl) RemoveProfitCommissionRule(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[bool] {
	err := s.profitCommRepo.Delete(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[bool](http.StatusInternalServerError, "Failed to remove profit commission rule", err)
	}
	return schema.NewServiceResponse(true, http.StatusOK, "Profit commission rule removed successfully")
}
