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

type reinsurerStatementServiceImpl struct {
	statementRepo   repository.ReinsurerStatementRepository
	cessionRepo     repository.CessionRepository
	recoveryRepo    repository.RecoveryRepository
	participantRepo repository.TreatyParticipantRepository
	profitCommRepo  repository.ProfitCommissionRepository
	auditSvc        auditService.AuditService
}

func NewReinsurerStatementService(
	statementRepo repository.ReinsurerStatementRepository,
	cessionRepo repository.CessionRepository,
	recoveryRepo repository.RecoveryRepository,
	participantRepo repository.TreatyParticipantRepository,
	profitCommRepo repository.ProfitCommissionRepository,
	auditSvc auditService.AuditService,
) service.ReinsurerStatementService {
	return &reinsurerStatementServiceImpl{
		statementRepo:   statementRepo,
		cessionRepo:     cessionRepo,
		recoveryRepo:    recoveryRepo,
		participantRepo: participantRepo,
		profitCommRepo:  profitCommRepo,
		auditSvc:        auditSvc,
	}
}

func (s *reinsurerStatementServiceImpl) GenerateStatement(ctx context.Context, req reinsuranceSchema.GenerateStatementRequest, createdBy uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.ReinsurerStatementResponse] {
	treatyID, _ := uuid.Parse(req.TreatyID)
	participantID, _ := uuid.Parse(req.ParticipantID)

	participant, err := s.participantRepo.GetByID(ctx, participantID)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.ReinsurerStatementResponse](http.StatusNotFound, "Participant not found", err)
	}

	totalCeded, _ := s.cessionRepo.GetTotalCededByTreatyAndPeriod(ctx, treatyID, req.PeriodStart, req.PeriodEnd)
	totalRecovered, _ := s.recoveryRepo.GetTotalRecoveredByTreaty(ctx, treatyID)

	// Calculate participant's share
	premiumCeded := totalCeded * int64(participant.SharePercentage) / 100
	claimsRecovered := totalRecovered * int64(participant.SharePercentage) / 100
	commissionDue := premiumCeded * int64(participant.CommissionRate) / 100
	netBalance := premiumCeded - claimsRecovered - commissionDue

	statement := &entity.ReinsurerStatement{
		StatementNumber:  utils.GenerateReinsurerStatementNumber(),
		TreatyID:         treatyID,
		ParticipantID:    participantID,
		PeriodStart:      req.PeriodStart,
		PeriodEnd:        req.PeriodEnd,
		PremiumCeded:     premiumCeded,
		ClaimsRecovered:  claimsRecovered,
		CommissionDue:    commissionDue,
		ProfitCommission: 0,
		NetBalance:       netBalance,
		Status:           string(shared.ReinsurerStatementStatusDraft),
		CreatedBy:        createdBy,
	}

	created, err := s.statementRepo.Create(ctx, statement)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.ReinsurerStatementResponse](http.StatusInternalServerError, "Failed to create statement", err)
	}

	logAudit(ctx, s.auditSvc, createdBy, string(shared.AuditEntityTypeReinsurerStatement), created.ID, string(shared.AuditActionCreate))

	resp := reinsuranceSchema.ToReinsurerStatementResponse(created)
	return schema.NewServiceResponse(resp, http.StatusCreated, "Reinsurer statement generated successfully")
}

func (s *reinsurerStatementServiceImpl) GetStatement(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.ReinsurerStatementResponse] {
	statement, err := s.statementRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.ReinsurerStatementResponse](http.StatusNotFound, "Statement not found", err)
	}

	resp := reinsuranceSchema.ToReinsurerStatementResponse(statement)
	return schema.NewServiceResponse(resp, http.StatusOK, "Statement retrieved successfully")
}

func (s *reinsurerStatementServiceImpl) ListByTreaty(ctx context.Context, treatyID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]reinsuranceSchema.ReinsurerStatementResponse] {
	offset := (page - 1) * pageSize
	statements, err := s.statementRepo.ListByTreaty(ctx, treatyID, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]reinsuranceSchema.ReinsurerStatementResponse](http.StatusInternalServerError, "Failed to list statements", err)
	}

	responses := make([]reinsuranceSchema.ReinsurerStatementResponse, len(statements))
	for i, st := range statements {
		responses[i] = reinsuranceSchema.ToReinsurerStatementResponse(st)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Statements retrieved successfully")
}

func (s *reinsurerStatementServiceImpl) IssueStatement(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.ReinsurerStatementResponse] {
	existing, err := s.statementRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.ReinsurerStatementResponse](http.StatusNotFound, "Statement not found", err)
	}
	if existing.Status != string(shared.ReinsurerStatementStatusDraft) {
		return schema.NewServiceErrorResponse[reinsuranceSchema.ReinsurerStatementResponse](http.StatusBadRequest, "Only DRAFT statements can be issued", fmt.Errorf("invalid status: %s", existing.Status))
	}

	updated, err := s.statementRepo.UpdateStatus(ctx, id, string(shared.ReinsurerStatementStatusIssued))
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.ReinsurerStatementResponse](http.StatusInternalServerError, "Failed to issue statement", err)
	}

	logAudit(ctx, s.auditSvc, userID, string(shared.AuditEntityTypeReinsurerStatement), id, string(shared.AuditActionStateChange))

	resp := reinsuranceSchema.ToReinsurerStatementResponse(updated)
	return schema.NewServiceResponse(resp, http.StatusOK, "Statement issued successfully")
}

func (s *reinsurerStatementServiceImpl) AcknowledgeStatement(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.ReinsurerStatementResponse] {
	existing, err := s.statementRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.ReinsurerStatementResponse](http.StatusNotFound, "Statement not found", err)
	}
	if existing.Status != string(shared.ReinsurerStatementStatusIssued) {
		return schema.NewServiceErrorResponse[reinsuranceSchema.ReinsurerStatementResponse](http.StatusBadRequest, "Only ISSUED statements can be acknowledged", fmt.Errorf("invalid status: %s", existing.Status))
	}

	updated, err := s.statementRepo.UpdateStatus(ctx, id, string(shared.ReinsurerStatementStatusAcknowledged))
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.ReinsurerStatementResponse](http.StatusInternalServerError, "Failed to acknowledge statement", err)
	}

	resp := reinsuranceSchema.ToReinsurerStatementResponse(updated)
	return schema.NewServiceResponse(resp, http.StatusOK, "Statement acknowledged successfully")
}

func (s *reinsurerStatementServiceImpl) SettleStatement(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.ReinsurerStatementResponse] {
	existing, err := s.statementRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.ReinsurerStatementResponse](http.StatusNotFound, "Statement not found", err)
	}
	if existing.Status != string(shared.ReinsurerStatementStatusAcknowledged) {
		return schema.NewServiceErrorResponse[reinsuranceSchema.ReinsurerStatementResponse](http.StatusBadRequest, "Only ACKNOWLEDGED statements can be settled", fmt.Errorf("invalid status: %s", existing.Status))
	}

	updated, err := s.statementRepo.UpdateStatus(ctx, id, string(shared.ReinsurerStatementStatusSettled))
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.ReinsurerStatementResponse](http.StatusInternalServerError, "Failed to settle statement", err)
	}

	logAudit(ctx, s.auditSvc, userID, string(shared.AuditEntityTypeReinsurerStatement), id, string(shared.AuditActionStateChange))

	resp := reinsuranceSchema.ToReinsurerStatementResponse(updated)
	return schema.NewServiceResponse(resp, http.StatusOK, "Statement settled successfully")
}

func (s *reinsurerStatementServiceImpl) CalculateProfitCommission(ctx context.Context, req reinsuranceSchema.CalculateProfitCommissionRequest) *schema.ServiceResponse[reinsuranceSchema.ProfitCommissionCalculationResponse] {
	treatyID, _ := uuid.Parse(req.TreatyID)

	premiumCeded, _ := s.cessionRepo.GetTotalCededByTreatyAndPeriod(ctx, treatyID, req.PeriodStart, req.PeriodEnd)
	claimsRecovered, _ := s.recoveryRepo.GetTotalRecoveredByTreaty(ctx, treatyID)

	var lossRatio float64
	if premiumCeded > 0 {
		lossRatio = float64(claimsRecovered) * 100 / float64(premiumCeded)
	}

	rules, _ := s.profitCommRepo.ListByTreaty(ctx, treatyID)

	var commissionRate float64
	var carryForward int64
	for _, rule := range rules {
		if lossRatio >= rule.LossRatioFrom && lossRatio <= rule.LossRatioTo {
			commissionRate = rule.CommissionRate
			break
		}
	}

	netProfit := premiumCeded - claimsRecovered
	// Apply carry forward if applicable
	for _, rule := range rules {
		if rule.CommissionType == string(shared.ProfitCommissionTypeCarryForward) && rule.CarryForwardBalance > 0 {
			netProfit -= rule.CarryForwardBalance
			break
		}
	}

	var commissionAmount int64
	if netProfit > 0 {
		commissionAmount = netProfit * int64(commissionRate) / 100
	} else {
		carryForward = -netProfit
	}

	resp := reinsuranceSchema.ProfitCommissionCalculationResponse{
		TreatyID:         treatyID,
		PremiumCeded:     premiumCeded,
		ClaimsRecovered:  claimsRecovered,
		LossRatio:        lossRatio,
		NetProfit:        netProfit,
		CommissionRate:   commissionRate,
		CommissionAmount: commissionAmount,
		CarryForward:     carryForward,
	}
	return schema.NewServiceResponse(resp, http.StatusOK, "Profit commission calculated successfully")
}
