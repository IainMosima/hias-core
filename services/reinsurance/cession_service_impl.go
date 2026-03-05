package reinsurance

import (
	"context"
	"fmt"
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

type cessionServiceImpl struct {
	cessionRepo     repository.CessionRepository
	treatyRepo      repository.TreatyRepository
	participantRepo repository.TreatyParticipantRepository
	auditSvc        auditService.AuditService
}

func NewCessionService(
	cessionRepo repository.CessionRepository,
	treatyRepo repository.TreatyRepository,
	participantRepo repository.TreatyParticipantRepository,
	auditSvc auditService.AuditService,
) service.CessionService {
	return &cessionServiceImpl{
		cessionRepo:     cessionRepo,
		treatyRepo:      treatyRepo,
		participantRepo: participantRepo,
		auditSvc:        auditSvc,
	}
}

func (s *cessionServiceImpl) CedePremium(ctx context.Context, req reinsuranceSchema.CedePremiumRequest, createdBy uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.CessionResponse] {
	treatyID, err := uuid.Parse(req.TreatyID)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.CessionResponse](http.StatusBadRequest, "Invalid treaty ID", err)
	}

	policyID, err := uuid.Parse(req.PolicyID)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.CessionResponse](http.StatusBadRequest, "Invalid policy ID", err)
	}

	treaty, err := s.treatyRepo.GetByID(ctx, treatyID)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.CessionResponse](http.StatusNotFound, "Treaty not found", err)
	}

	if treaty.Status != string(shared.TreatyStatusActive) {
		return schema.NewServiceErrorResponse[reinsuranceSchema.CessionResponse](http.StatusBadRequest, "Treaty is not active", fmt.Errorf("treaty status: %s", treaty.Status))
	}

	// Get participants and calculate total share
	participants, err := s.participantRepo.ListByTreaty(ctx, treatyID)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.CessionResponse](http.StatusInternalServerError, "Failed to fetch treaty participants", err)
	}

	if len(participants) == 0 {
		return schema.NewServiceErrorResponse[reinsuranceSchema.CessionResponse](http.StatusBadRequest, "Treaty has no participants", fmt.Errorf("no participants for treaty %s", treatyID))
	}

	var totalShare float64
	var totalCommissionRate float64
	for _, p := range participants {
		totalShare += p.SharePercentage
		totalCommissionRate += p.CommissionRate
	}
	avgCommissionRate := totalCommissionRate / float64(len(participants))

	grossAmount := req.Amount
	cededAmount := int64(float64(grossAmount) * totalShare / 100)
	retainedAmount := grossAmount - cededAmount
	commissionAmount := int64(float64(cededAmount) * avgCommissionRate / 100)

	// Enforce retention limit cap
	if treaty.RetentionLimit > 0 && retainedAmount < treaty.RetentionLimit {
		retainedAmount = treaty.RetentionLimit
		cededAmount = grossAmount - retainedAmount
		commissionAmount = int64(float64(cededAmount) * avgCommissionRate / 100)
	}

	cession := &entity.Cession{
		CessionNumber:    utils.GenerateCessionNumber(),
		TreatyID:         treatyID,
		PolicyID:         policyID,
		CessionType:      string(shared.CessionTypePremium),
		GrossAmount:      grossAmount,
		CededAmount:      cededAmount,
		RetainedAmount:   retainedAmount,
		CommissionAmount: commissionAmount,
		SharePercentage:  totalShare,
		Status:           string(shared.CessionStatusPending),
		CreatedBy:        createdBy,
	}

	created, err := s.cessionRepo.Create(ctx, cession)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.CessionResponse](http.StatusInternalServerError, "Failed to create cession", err)
	}

	logAudit(ctx, s.auditSvc, createdBy, string(shared.AuditEntityTypeCession), created.ID, string(shared.AuditActionCreate))

	resp := reinsuranceSchema.ToCessionResponse(created)
	return schema.NewServiceResponse(resp, http.StatusCreated, "Premium ceded successfully")
}

func (s *cessionServiceImpl) AutoCedePolicyPremium(ctx context.Context, req reinsuranceSchema.AutoCedePolicyPremiumRequest, createdBy uuid.UUID) *schema.ServiceResponse[[]reinsuranceSchema.CessionResponse] {
	policyID, err := uuid.Parse(req.PolicyID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]reinsuranceSchema.CessionResponse](http.StatusBadRequest, "Invalid policy ID", err)
	}

	// Find all ACTIVE Quota Share treaties
	activeTreaties, err := s.treatyRepo.ListByStatus(ctx, string(shared.TreatyStatusActive), 1000, 0)
	if err != nil {
		return schema.NewServiceErrorResponse[[]reinsuranceSchema.CessionResponse](http.StatusInternalServerError, "Failed to list active treaties", err)
	}

	var quotaShareTreaties []*entity.Treaty
	for _, t := range activeTreaties {
		if t.TreatyType == string(shared.TreatyTypeQuotaShare) {
			quotaShareTreaties = append(quotaShareTreaties, t)
		}
	}

	if len(quotaShareTreaties) == 0 {
		return schema.NewServiceErrorResponse[[]reinsuranceSchema.CessionResponse](http.StatusNotFound, "No active Quota Share treaties found", fmt.Errorf("no active quota share treaties"))
	}

	var responses []reinsuranceSchema.CessionResponse

	for _, treaty := range quotaShareTreaties {
		participants, err := s.participantRepo.ListByTreaty(ctx, treaty.ID)
		if err != nil || len(participants) == 0 {
			continue
		}

		var totalShare float64
		var totalCommissionRate float64
		for _, p := range participants {
			totalShare += p.SharePercentage
			totalCommissionRate += p.CommissionRate
		}
		avgCommissionRate := totalCommissionRate / float64(len(participants))

		grossAmount := req.Amount
		cededAmount := int64(float64(grossAmount) * totalShare / 100)
		retainedAmount := grossAmount - cededAmount
		commissionAmount := int64(float64(cededAmount) * avgCommissionRate / 100)

		// Enforce retention limit cap
		if treaty.RetentionLimit > 0 && retainedAmount < treaty.RetentionLimit {
			retainedAmount = treaty.RetentionLimit
			cededAmount = grossAmount - retainedAmount
			commissionAmount = int64(float64(cededAmount) * avgCommissionRate / 100)
		}

		cession := &entity.Cession{
			CessionNumber:    utils.GenerateCessionNumber(),
			TreatyID:         treaty.ID,
			PolicyID:         policyID,
			CessionType:      string(shared.CessionTypePremium),
			GrossAmount:      grossAmount,
			CededAmount:      cededAmount,
			RetainedAmount:   retainedAmount,
			CommissionAmount: commissionAmount,
			SharePercentage:  totalShare,
			Status:           string(shared.CessionStatusPending),
			CreatedBy:        createdBy,
		}

		created, err := s.cessionRepo.Create(ctx, cession)
		if err != nil {
			continue
		}

		logAudit(ctx, s.auditSvc, createdBy, string(shared.AuditEntityTypeCession), created.ID, string(shared.AuditActionCreate))
		responses = append(responses, reinsuranceSchema.ToCessionResponse(created))
	}

	if len(responses) == 0 {
		return schema.NewServiceErrorResponse[[]reinsuranceSchema.CessionResponse](http.StatusInternalServerError, "Failed to create any cessions", fmt.Errorf("no cessions created"))
	}

	return schema.NewServiceResponse(responses, http.StatusCreated, fmt.Sprintf("Auto-ceded premium across %d treaties", len(responses)))
}

func (s *cessionServiceImpl) BookCession(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.CessionResponse] {
	existing, err := s.cessionRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.CessionResponse](http.StatusNotFound, "Cession not found", err)
	}

	if existing.Status != string(shared.CessionStatusPending) {
		return schema.NewServiceErrorResponse[reinsuranceSchema.CessionResponse](http.StatusBadRequest, "Only PENDING cessions can be booked", fmt.Errorf("invalid status: %s", existing.Status))
	}

	updated, err := s.cessionRepo.UpdateStatus(ctx, id, string(shared.CessionStatusBooked))
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.CessionResponse](http.StatusInternalServerError, "Failed to book cession", err)
	}

	logAudit(ctx, s.auditSvc, userID, string(shared.AuditEntityTypeCession), id, string(shared.AuditActionStateChange))

	resp := reinsuranceSchema.ToCessionResponse(updated)
	return schema.NewServiceResponse(resp, http.StatusOK, "Cession booked successfully")
}

func (s *cessionServiceImpl) ReverseCession(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.CessionResponse] {
	existing, err := s.cessionRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.CessionResponse](http.StatusNotFound, "Cession not found", err)
	}

	if existing.Status != string(shared.CessionStatusBooked) {
		return schema.NewServiceErrorResponse[reinsuranceSchema.CessionResponse](http.StatusBadRequest, "Only BOOKED cessions can be reversed", fmt.Errorf("invalid status: %s", existing.Status))
	}

	updated, err := s.cessionRepo.UpdateStatus(ctx, id, string(shared.CessionStatusReversed))
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.CessionResponse](http.StatusInternalServerError, "Failed to reverse cession", err)
	}

	logAudit(ctx, s.auditSvc, userID, string(shared.AuditEntityTypeCession), id, string(shared.AuditActionStateChange))

	resp := reinsuranceSchema.ToCessionResponse(updated)
	return schema.NewServiceResponse(resp, http.StatusOK, "Cession reversed successfully")
}

func (s *cessionServiceImpl) GetCession(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.CessionResponse] {
	cession, err := s.cessionRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.CessionResponse](http.StatusNotFound, "Cession not found", err)
	}

	resp := reinsuranceSchema.ToCessionResponse(cession)
	return schema.NewServiceResponse(resp, http.StatusOK, "Cession retrieved successfully")
}

func (s *cessionServiceImpl) ListByTreaty(ctx context.Context, treatyID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]reinsuranceSchema.CessionResponse] {
	offset := (page - 1) * pageSize
	cessions, err := s.cessionRepo.ListByTreaty(ctx, treatyID, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]reinsuranceSchema.CessionResponse](http.StatusInternalServerError, "Failed to list cessions by treaty", err)
	}

	responses := make([]reinsuranceSchema.CessionResponse, len(cessions))
	for i, c := range cessions {
		responses[i] = reinsuranceSchema.ToCessionResponse(c)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Cessions retrieved successfully")
}

func (s *cessionServiceImpl) ListByPolicy(ctx context.Context, policyID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]reinsuranceSchema.CessionResponse] {
	offset := (page - 1) * pageSize
	cessions, err := s.cessionRepo.ListByPolicy(ctx, policyID, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]reinsuranceSchema.CessionResponse](http.StatusInternalServerError, "Failed to list cessions by policy", err)
	}

	responses := make([]reinsuranceSchema.CessionResponse, len(cessions))
	for i, c := range cessions {
		responses[i] = reinsuranceSchema.ToCessionResponse(c)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Cessions retrieved successfully")
}

func (s *cessionServiceImpl) GetCessionCount(ctx context.Context) *schema.ServiceResponse[int64] {
	count, err := s.cessionRepo.Count(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to count cessions", err)
	}
	return schema.NewServiceResponse(count, http.StatusOK, "Cession count retrieved")
}

func (s *cessionServiceImpl) GetTotalCededAmount(ctx context.Context) *schema.ServiceResponse[int64] {
	total, err := s.cessionRepo.GetTotalCededAmountAll(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to get total ceded amount", err)
	}
	return schema.NewServiceResponse(total, http.StatusOK, "Total ceded amount retrieved")
}

func (s *cessionServiceImpl) GetTotalGrossAmount(ctx context.Context) *schema.ServiceResponse[int64] {
	total, err := s.cessionRepo.GetTotalGrossAmountAll(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to get total gross amount", err)
	}
	return schema.NewServiceResponse(total, http.StatusOK, "Total gross amount retrieved")
}
