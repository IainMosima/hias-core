package reinsurance

import (
	"context"
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

type bordereauServiceImpl struct {
	bordereauRepo     repository.BordereauRepository
	bordereauItemRepo repository.BordereauItemRepository
	cessionRepo       repository.CessionRepository
	recoveryRepo      repository.RecoveryRepository
	auditSvc          auditService.AuditService
}

func NewBordereauService(
	bordereauRepo repository.BordereauRepository,
	bordereauItemRepo repository.BordereauItemRepository,
	cessionRepo repository.CessionRepository,
	recoveryRepo repository.RecoveryRepository,
	auditSvc auditService.AuditService,
) service.BordereauService {
	return &bordereauServiceImpl{
		bordereauRepo:     bordereauRepo,
		bordereauItemRepo: bordereauItemRepo,
		cessionRepo:       cessionRepo,
		recoveryRepo:      recoveryRepo,
		auditSvc:          auditSvc,
	}
}

func (s *bordereauServiceImpl) GeneratePremiumBordereau(ctx context.Context, req reinsuranceSchema.GenerateBordereauRequest, createdBy uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.BordereauResponse] {
	treatyID, _ := uuid.Parse(req.TreatyID)

	cessions, err := s.cessionRepo.ListBookedByTreatyAndPeriod(ctx, treatyID, req.PeriodStart, req.PeriodEnd)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.BordereauResponse](http.StatusInternalServerError, "Failed to fetch cessions", err)
	}

	bordereau := &entity.Bordereau{
		BordereauNumber: utils.GenerateBordereauNumber(),
		TreatyID:        treatyID,
		BordereauType:   string(shared.BordereauTypePremium),
		PeriodStart:     req.PeriodStart,
		PeriodEnd:       req.PeriodEnd,
		Status:          string(shared.BordereauStatusDraft),
		CreatedBy:       createdBy,
	}

	created, err := s.bordereauRepo.Create(ctx, bordereau)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.BordereauResponse](http.StatusInternalServerError, "Failed to create bordereau", err)
	}

	var totalGross, totalCeded, totalCommission int64
	for _, c := range cessions {
		item := &entity.BordereauItem{
			BordereauID:      created.ID,
			CessionID:        c.ID,
			GrossAmount:      c.GrossAmount,
			CededAmount:      c.CededAmount,
			CommissionAmount: c.CommissionAmount,
		}
		s.bordereauItemRepo.Create(ctx, item)
		totalGross += c.GrossAmount
		totalCeded += c.CededAmount
		totalCommission += c.CommissionAmount
	}

	created.TotalGross = totalGross
	created.TotalCeded = totalCeded
	created.TotalCommission = totalCommission
	created.ItemCount = len(cessions)

	updated, err := s.bordereauRepo.Update(ctx, created)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.BordereauResponse](http.StatusInternalServerError, "Failed to update bordereau totals", err)
	}

	logAudit(ctx, s.auditSvc, createdBy, string(shared.AuditEntityTypeBordereau), updated.ID, string(shared.AuditActionCreate))

	resp := reinsuranceSchema.ToBordereauResponse(updated)
	return schema.NewServiceResponse(resp, http.StatusCreated, "Premium bordereau generated successfully")
}

func (s *bordereauServiceImpl) GenerateClaimBordereau(ctx context.Context, req reinsuranceSchema.GenerateBordereauRequest, createdBy uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.BordereauResponse] {
	treatyID, _ := uuid.Parse(req.TreatyID)

	recoveries, err := s.recoveryRepo.ListByTreaty(ctx, treatyID, 10000, 0)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.BordereauResponse](http.StatusInternalServerError, "Failed to fetch recoveries", err)
	}

	bordereau := &entity.Bordereau{
		BordereauNumber: utils.GenerateBordereauNumber(),
		TreatyID:        treatyID,
		BordereauType:   string(shared.BordereauTypeClaim),
		PeriodStart:     req.PeriodStart,
		PeriodEnd:       req.PeriodEnd,
		Status:          string(shared.BordereauStatusDraft),
		CreatedBy:       createdBy,
	}

	created, err := s.bordereauRepo.Create(ctx, bordereau)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.BordereauResponse](http.StatusInternalServerError, "Failed to create bordereau", err)
	}

	var totalGross, totalCeded int64
	var itemCount int
	for _, r := range recoveries {
		if r.CreatedAt.Before(req.PeriodStart) || r.CreatedAt.After(req.PeriodEnd) {
			continue
		}
		item := &entity.BordereauItem{
			BordereauID:      created.ID,
			RecoveryID:       r.ID,
			GrossAmount:      r.GrossClaimAmount,
			CededAmount:      r.RecoverableAmount,
			CommissionAmount: 0,
		}
		s.bordereauItemRepo.Create(ctx, item)
		totalGross += r.GrossClaimAmount
		totalCeded += r.RecoverableAmount
		itemCount++
	}

	created.TotalGross = totalGross
	created.TotalCeded = totalCeded
	created.ItemCount = itemCount

	updated, err := s.bordereauRepo.Update(ctx, created)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.BordereauResponse](http.StatusInternalServerError, "Failed to update bordereau totals", err)
	}

	logAudit(ctx, s.auditSvc, createdBy, string(shared.AuditEntityTypeBordereau), updated.ID, string(shared.AuditActionCreate))

	resp := reinsuranceSchema.ToBordereauResponse(updated)
	return schema.NewServiceResponse(resp, http.StatusCreated, "Claim bordereau generated successfully")
}

func (s *bordereauServiceImpl) GetBordereau(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.BordereauDetailResponse] {
	bordereau, err := s.bordereauRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.BordereauDetailResponse](http.StatusNotFound, "Bordereau not found", err)
	}

	items, _ := s.bordereauItemRepo.ListByBordereau(ctx, id)
	itemResponses := make([]reinsuranceSchema.BordereauItemResponse, len(items))
	for i, item := range items {
		itemResponses[i] = reinsuranceSchema.ToBordereauItemResponse(item)
	}

	resp := reinsuranceSchema.BordereauDetailResponse{
		BordereauResponse: reinsuranceSchema.ToBordereauResponse(bordereau),
		Items:             itemResponses,
	}
	return schema.NewServiceResponse(resp, http.StatusOK, "Bordereau retrieved successfully")
}

func (s *bordereauServiceImpl) ListByTreaty(ctx context.Context, treatyID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]reinsuranceSchema.BordereauResponse] {
	offset := (page - 1) * pageSize
	bordereaux, err := s.bordereauRepo.ListByTreaty(ctx, treatyID, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]reinsuranceSchema.BordereauResponse](http.StatusInternalServerError, "Failed to list bordereaux", err)
	}

	responses := make([]reinsuranceSchema.BordereauResponse, len(bordereaux))
	for i, b := range bordereaux {
		responses[i] = reinsuranceSchema.ToBordereauResponse(b)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Bordereaux retrieved successfully")
}

func (s *bordereauServiceImpl) FinalizeBordereau(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.BordereauResponse] {
	existing, err := s.bordereauRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.BordereauResponse](http.StatusNotFound, "Bordereau not found", err)
	}
	if existing.Status != string(shared.BordereauStatusDraft) {
		return schema.NewServiceErrorResponse[reinsuranceSchema.BordereauResponse](http.StatusBadRequest, "Only DRAFT bordereaux can be finalized", nil)
	}

	updated, err := s.bordereauRepo.UpdateStatus(ctx, id, string(shared.BordereauStatusFinalized))
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.BordereauResponse](http.StatusInternalServerError, "Failed to finalize bordereau", err)
	}

	logAudit(ctx, s.auditSvc, userID, string(shared.AuditEntityTypeBordereau), id, string(shared.AuditActionStateChange))

	resp := reinsuranceSchema.ToBordereauResponse(updated)
	return schema.NewServiceResponse(resp, http.StatusOK, "Bordereau finalized successfully")
}

func (s *bordereauServiceImpl) MarkSent(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.BordereauResponse] {
	existing, err := s.bordereauRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.BordereauResponse](http.StatusNotFound, "Bordereau not found", err)
	}
	if existing.Status != string(shared.BordereauStatusFinalized) {
		return schema.NewServiceErrorResponse[reinsuranceSchema.BordereauResponse](http.StatusBadRequest, "Only FINALIZED bordereaux can be marked as sent", nil)
	}

	updated, err := s.bordereauRepo.UpdateStatus(ctx, id, string(shared.BordereauStatusSent))
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.BordereauResponse](http.StatusInternalServerError, "Failed to mark bordereau as sent", err)
	}

	logAudit(ctx, s.auditSvc, userID, string(shared.AuditEntityTypeBordereau), id, string(shared.AuditActionStateChange))

	resp := reinsuranceSchema.ToBordereauResponse(updated)
	return schema.NewServiceResponse(resp, http.StatusOK, "Bordereau marked as sent")
}

func (s *bordereauServiceImpl) ListItems(ctx context.Context, bordereauID uuid.UUID) *schema.ServiceResponse[[]reinsuranceSchema.BordereauItemResponse] {
	items, err := s.bordereauItemRepo.ListByBordereau(ctx, bordereauID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]reinsuranceSchema.BordereauItemResponse](http.StatusInternalServerError, "Failed to list bordereau items", err)
	}

	responses := make([]reinsuranceSchema.BordereauItemResponse, len(items))
	for i, item := range items {
		responses[i] = reinsuranceSchema.ToBordereauItemResponse(item)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Bordereau items retrieved successfully")
}
