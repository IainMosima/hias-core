package claims

import (
	"context"
	"fmt"
	"log"
	"net/http"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	"github.com/bitbiz/hias-core/domains/claims/entity"
	claimRepo "github.com/bitbiz/hias-core/domains/claims/repository"
	claimsSchema "github.com/bitbiz/hias-core/domains/claims/schema"
	"github.com/bitbiz/hias-core/domains/claims/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	preauthRepo "github.com/bitbiz/hias-core/domains/preauth/repository"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/google/uuid"
)

type caseServiceImpl struct {
	caseRepo    claimRepo.CaseRecordRepository
	preauthRepo preauthRepo.PreAuthRepository
	auditSvc    auditService.AuditService
}

func NewCaseService(
	caseRepo claimRepo.CaseRecordRepository,
	preauthRepo preauthRepo.PreAuthRepository,
	auditSvc auditService.AuditService,
) service.CaseService {
	return &caseServiceImpl{
		caseRepo:    caseRepo,
		preauthRepo: preauthRepo,
		auditSvc:    auditSvc,
	}
}

func (s *caseServiceImpl) CreateCase(ctx context.Context, req claimsSchema.CreateCaseRequest, createdBy uuid.UUID) *schema.ServiceResponse[claimsSchema.CaseRecordResponse] {
	preauthID, _ := uuid.Parse(req.PreAuthID)

	// Fetch PreAuth and validate
	preauth, err := s.preauthRepo.GetByID(ctx, preauthID)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.CaseRecordResponse](http.StatusNotFound, "Pre-authorization not found", err)
	}

	if preauth.Status != string(shared.PreAuthStatusApproved) {
		return schema.NewServiceErrorResponse[claimsSchema.CaseRecordResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Pre-authorization is %s, must be APPROVED", preauth.Status),
			nil,
		)
	}

	record := &entity.CaseRecord{
		CaseNumber:         utils.GenerateCaseNumber(),
		PreAuthID:          preauthID,
		PolicyID:           preauth.PolicyID,
		MemberID:           preauth.MemberID,
		ProviderID:         preauth.ProviderID,
		Status:             string(shared.CaseStatusScheduled),
		ExpectedDischarge:  req.ExpectedDischarge,
		Diagnosis:          req.Diagnosis,
		TreatingDoctor:     req.TreatingDoctor,
		RoomType:           req.RoomType,
		TotalEstimatedCost: req.EstimatedCost,
		Notes:              req.Notes,
		CreatedBy:          createdBy,
	}

	created, err := s.caseRepo.Create(ctx, record)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.CaseRecordResponse](http.StatusInternalServerError, "Failed to create case", err)
	}

	s.logAudit(ctx, createdBy, string(shared.AuditEntityTypeCaseRecord), created.ID, string(shared.AuditActionCreate))
	return schema.NewServiceResponse(claimsSchema.ToCaseRecordResponse(created), http.StatusCreated, "Case created")
}

func (s *caseServiceImpl) GetCase(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[claimsSchema.CaseRecordResponse] {
	record, err := s.caseRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.CaseRecordResponse](http.StatusNotFound, "Case not found", err)
	}
	return schema.NewServiceResponse(claimsSchema.ToCaseRecordResponse(record), http.StatusOK, "Case retrieved")
}

func (s *caseServiceImpl) ListByPolicy(ctx context.Context, policyID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]claimsSchema.CaseRecordResponse] {
	offset := (page - 1) * pageSize
	records, err := s.caseRepo.ListByPolicy(ctx, policyID, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]claimsSchema.CaseRecordResponse](http.StatusInternalServerError, "Failed to list cases", err)
	}
	return schema.NewServiceResponse(toCaseRecordResponses(records), http.StatusOK, "Cases retrieved")
}

func (s *caseServiceImpl) ListByMember(ctx context.Context, memberID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]claimsSchema.CaseRecordResponse] {
	offset := (page - 1) * pageSize
	records, err := s.caseRepo.ListByMember(ctx, memberID, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]claimsSchema.CaseRecordResponse](http.StatusInternalServerError, "Failed to list cases", err)
	}
	return schema.NewServiceResponse(toCaseRecordResponses(records), http.StatusOK, "Cases retrieved")
}

func (s *caseServiceImpl) ListByProvider(ctx context.Context, providerID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]claimsSchema.CaseRecordResponse] {
	offset := (page - 1) * pageSize
	records, err := s.caseRepo.ListByProvider(ctx, providerID, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]claimsSchema.CaseRecordResponse](http.StatusInternalServerError, "Failed to list cases", err)
	}
	return schema.NewServiceResponse(toCaseRecordResponses(records), http.StatusOK, "Cases retrieved")
}

func (s *caseServiceImpl) ListByStatus(ctx context.Context, status string, page, pageSize int) *schema.ServiceResponse[[]claimsSchema.CaseRecordResponse] {
	offset := (page - 1) * pageSize
	records, err := s.caseRepo.ListByStatus(ctx, status, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]claimsSchema.CaseRecordResponse](http.StatusInternalServerError, "Failed to list cases", err)
	}
	return schema.NewServiceResponse(toCaseRecordResponses(records), http.StatusOK, "Cases retrieved")
}

func (s *caseServiceImpl) AdmitCase(ctx context.Context, id uuid.UUID, req claimsSchema.AdmitCaseRequest, userID uuid.UUID) *schema.ServiceResponse[claimsSchema.CaseRecordResponse] {
	record, err := s.caseRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.CaseRecordResponse](http.StatusNotFound, "Case not found", err)
	}

	if record.Status != string(shared.CaseStatusScheduled) {
		return schema.NewServiceErrorResponse[claimsSchema.CaseRecordResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Cannot admit case in %s status; must be SCHEDULED", record.Status),
			nil,
		)
	}

	updated, err := s.caseRepo.Admit(ctx, id, req.AdmissionDate)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.CaseRecordResponse](http.StatusInternalServerError, "Failed to admit case", err)
	}

	s.logAudit(ctx, userID, string(shared.AuditEntityTypeCaseRecord), id, string(shared.AuditActionStateChange))
	return schema.NewServiceResponse(claimsSchema.ToCaseRecordResponse(updated), http.StatusOK, "Case admitted")
}

func (s *caseServiceImpl) UpdateCase(ctx context.Context, id uuid.UUID, req claimsSchema.UpdateCaseRequest, userID uuid.UUID) *schema.ServiceResponse[claimsSchema.CaseRecordResponse] {
	record, err := s.caseRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.CaseRecordResponse](http.StatusNotFound, "Case not found", err)
	}

	if record.Status != string(shared.CaseStatusAdmitted) && record.Status != string(shared.CaseStatusInTreatment) {
		return schema.NewServiceErrorResponse[claimsSchema.CaseRecordResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Cannot update case in %s status", record.Status),
			nil,
		)
	}

	updateRecord := &entity.CaseRecord{
		ID:             id,
		Diagnosis:      req.Diagnosis,
		TreatingDoctor: req.TreatingDoctor,
		RoomType:       req.RoomType,
		Notes:          req.Notes,
	}
	if req.EstimatedCost != nil {
		updateRecord.TotalEstimatedCost = *req.EstimatedCost
	}

	updated, err := s.caseRepo.Update(ctx, updateRecord)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.CaseRecordResponse](http.StatusInternalServerError, "Failed to update case", err)
	}

	s.logAudit(ctx, userID, string(shared.AuditEntityTypeCaseRecord), id, string(shared.AuditActionUpdate))
	return schema.NewServiceResponse(claimsSchema.ToCaseRecordResponse(updated), http.StatusOK, "Case updated")
}

func (s *caseServiceImpl) StartTreatment(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[claimsSchema.CaseRecordResponse] {
	record, err := s.caseRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.CaseRecordResponse](http.StatusNotFound, "Case not found", err)
	}

	if record.Status != string(shared.CaseStatusAdmitted) {
		return schema.NewServiceErrorResponse[claimsSchema.CaseRecordResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Cannot start treatment for case in %s status; must be ADMITTED", record.Status),
			nil,
		)
	}

	updated, err := s.caseRepo.UpdateStatus(ctx, id, string(shared.CaseStatusInTreatment))
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.CaseRecordResponse](http.StatusInternalServerError, "Failed to start treatment", err)
	}

	s.logAudit(ctx, userID, string(shared.AuditEntityTypeCaseRecord), id, string(shared.AuditActionStateChange))
	return schema.NewServiceResponse(claimsSchema.ToCaseRecordResponse(updated), http.StatusOK, "Case treatment started")
}

func (s *caseServiceImpl) DischargeCase(ctx context.Context, id uuid.UUID, req claimsSchema.DischargeCaseRequest, userID uuid.UUID) *schema.ServiceResponse[claimsSchema.CaseRecordResponse] {
	record, err := s.caseRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.CaseRecordResponse](http.StatusNotFound, "Case not found", err)
	}

	if record.Status != string(shared.CaseStatusAdmitted) && record.Status != string(shared.CaseStatusInTreatment) {
		return schema.NewServiceErrorResponse[claimsSchema.CaseRecordResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Cannot discharge case in %s status", record.Status),
			nil,
		)
	}

	updated, err := s.caseRepo.Discharge(ctx, id, req.ActualDischarge, req.ActualCost)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.CaseRecordResponse](http.StatusInternalServerError, "Failed to discharge case", err)
	}

	s.logAudit(ctx, userID, string(shared.AuditEntityTypeCaseRecord), id, string(shared.AuditActionStateChange))
	return schema.NewServiceResponse(claimsSchema.ToCaseRecordResponse(updated), http.StatusOK, "Case discharged")
}

func (s *caseServiceImpl) CloseCase(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[claimsSchema.CaseRecordResponse] {
	record, err := s.caseRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.CaseRecordResponse](http.StatusNotFound, "Case not found", err)
	}

	if record.Status != string(shared.CaseStatusDischarged) {
		return schema.NewServiceErrorResponse[claimsSchema.CaseRecordResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Cannot close case in %s status; must be DISCHARGED", record.Status),
			nil,
		)
	}

	updated, err := s.caseRepo.Close(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.CaseRecordResponse](http.StatusInternalServerError, "Failed to close case", err)
	}

	s.logAudit(ctx, userID, string(shared.AuditEntityTypeCaseRecord), id, string(shared.AuditActionStateChange))
	return schema.NewServiceResponse(claimsSchema.ToCaseRecordResponse(updated), http.StatusOK, "Case closed")
}

func (s *caseServiceImpl) CountByStatus(ctx context.Context, status string) *schema.ServiceResponse[int64] {
	count, err := s.caseRepo.CountByStatus(ctx, status)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to count cases", err)
	}
	return schema.NewServiceResponse(count, http.StatusOK, "Count retrieved")
}

func (s *caseServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}

func toCaseRecordResponses(records []*entity.CaseRecord) []claimsSchema.CaseRecordResponse {
	responses := make([]claimsSchema.CaseRecordResponse, len(records))
	for i, r := range records {
		responses[i] = claimsSchema.ToCaseRecordResponse(r)
	}
	return responses
}
