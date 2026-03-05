package billing

import (
	"context"
	"log"
	"math"
	"net/http"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	"github.com/bitbiz/hias-core/domains/billing/entity"
	billingRepo "github.com/bitbiz/hias-core/domains/billing/repository"
	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	"github.com/bitbiz/hias-core/domains/billing/service"
	claimRepo "github.com/bitbiz/hias-core/domains/claims/repository"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/google/uuid"
)

type statementServiceImpl struct {
	statementRepo billingRepo.ProviderStatementRepository
	claimRepo     claimRepo.ClaimRepository
	auditSvc      auditService.AuditService
}

func NewStatementService(
	statementRepo billingRepo.ProviderStatementRepository,
	claimRepo claimRepo.ClaimRepository,
	auditSvc auditService.AuditService,
) service.StatementService {
	return &statementServiceImpl{
		statementRepo: statementRepo,
		claimRepo:     claimRepo,
		auditSvc:      auditSvc,
	}
}

func (s *statementServiceImpl) UploadStatement(ctx context.Context, req billingSchema.UploadStatementRequest, createdBy uuid.UUID) *schema.ServiceResponse[billingSchema.ProviderStatementResponse] {
	providerID, _ := uuid.Parse(req.ProviderID)

	// Calculate total claimed
	var totalClaimed int64
	for _, li := range req.LineItems {
		totalClaimed += li.ClaimedAmount
	}

	stmt := &entity.ProviderStatement{
		ProviderID:      providerID,
		StatementNumber: utils.GenerateStatementNumber(),
		PeriodStart:     req.PeriodStart,
		PeriodEnd:       req.PeriodEnd,
		TotalClaimed:    totalClaimed,
		Status:          string(shared.StatementStatusUploaded),
		FileName:        req.FileName,
		S3Key:           req.S3Key,
		CreatedBy:       createdBy,
	}

	created, err := s.statementRepo.Create(ctx, stmt)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.ProviderStatementResponse](http.StatusInternalServerError, "Failed to create statement", err)
	}

	// Create line items
	for _, li := range req.LineItems {
		item := &entity.StatementLineItem{
			StatementID:   created.ID,
			ClaimNumber:   li.ClaimNumber,
			ServiceDate:   li.ServiceDate,
			MemberName:    li.MemberName,
			ProcedureCode: li.ProcedureCode,
			ClaimedAmount: li.ClaimedAmount,
			MatchStatus:   string(shared.MatchStatusUnmatched),
		}
		_, lineErr := s.statementRepo.CreateLineItem(ctx, item)
		if lineErr != nil {
			log.Printf("Failed to create statement line item: %v", lineErr)
		}
	}

	s.logAudit(ctx, createdBy, string(shared.AuditEntityTypeProviderStatement), created.ID, string(shared.AuditActionCreate))
	return schema.NewServiceResponse(billingSchema.ToProviderStatementResponse(created), http.StatusCreated, "Statement uploaded")
}

func (s *statementServiceImpl) GetStatement(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[billingSchema.ProviderStatementResponse] {
	stmt, err := s.statementRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.ProviderStatementResponse](http.StatusNotFound, "Statement not found", err)
	}
	return schema.NewServiceResponse(billingSchema.ToProviderStatementResponse(stmt), http.StatusOK, "Statement retrieved")
}

func (s *statementServiceImpl) ListByProvider(ctx context.Context, providerID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]billingSchema.ProviderStatementResponse] {
	offset := (page - 1) * pageSize
	stmts, err := s.statementRepo.ListByProvider(ctx, providerID, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]billingSchema.ProviderStatementResponse](http.StatusInternalServerError, "Failed to list statements", err)
	}

	responses := make([]billingSchema.ProviderStatementResponse, len(stmts))
	for i, stmt := range stmts {
		responses[i] = billingSchema.ToProviderStatementResponse(stmt)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Statements retrieved")
}

func (s *statementServiceImpl) ListLineItems(ctx context.Context, statementID uuid.UUID) *schema.ServiceResponse[[]billingSchema.StatementLineItemResponse] {
	items, err := s.statementRepo.ListLineItems(ctx, statementID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]billingSchema.StatementLineItemResponse](http.StatusInternalServerError, "Failed to list line items", err)
	}

	responses := make([]billingSchema.StatementLineItemResponse, len(items))
	for i, item := range items {
		responses[i] = billingSchema.ToStatementLineItemResponse(item)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Line items retrieved")
}

func (s *statementServiceImpl) ReconcileStatement(ctx context.Context, id uuid.UUID, reconciledBy uuid.UUID) *schema.ServiceResponse[billingSchema.ProviderStatementResponse] {
	// Fetch line items
	items, err := s.statementRepo.ListLineItems(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.ProviderStatementResponse](http.StatusInternalServerError, "Failed to fetch line items", err)
	}

	var totalMatched, totalDiscrepancy int64
	var matchedCount, unmatchedCount int
	tolerance := int64(100) // 1 KES tolerance

	for _, item := range items {
		if item.ClaimNumber == "" {
			unmatchedCount++
			continue
		}

		// Try to match by claim number
		claim, claimErr := s.claimRepo.GetByNumber(ctx, item.ClaimNumber)
		if claimErr != nil {
			unmatchedCount++
			s.statementRepo.MatchLineItem(ctx, item.ID, uuid.Nil, string(shared.MatchStatusUnmatched), 0, "Claim not found")
			continue
		}

		discrepancy := item.ClaimedAmount - claim.ApprovedAmount
		matchStatus := string(shared.MatchStatusMatched)
		if math.Abs(float64(discrepancy)) > float64(tolerance) {
			totalDiscrepancy += discrepancy
		} else {
			discrepancy = 0
		}

		totalMatched += item.ClaimedAmount
		matchedCount++
		s.statementRepo.MatchLineItem(ctx, item.ID, claim.ID, matchStatus, discrepancy, "")

		// Update claim payment status based on match
		if discrepancy <= 0 {
			s.claimRepo.UpdateStatus(ctx, claim.ID, string(shared.ClaimStatusPaid))
		} else {
			s.claimRepo.UpdateStatus(ctx, claim.ID, string(shared.ClaimStatusPartPaid))
		}
	}

	// Update statement totals
	reconciled, err := s.statementRepo.Reconcile(ctx, id, totalMatched, totalDiscrepancy, matchedCount, unmatchedCount, reconciledBy)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.ProviderStatementResponse](http.StatusInternalServerError, "Failed to reconcile statement", err)
	}

	s.logAudit(ctx, reconciledBy, string(shared.AuditEntityTypeProviderStatement), id, string(shared.AuditActionStateChange))
	return schema.NewServiceResponse(billingSchema.ToProviderStatementResponse(reconciled), http.StatusOK, "Statement reconciled")
}

func (s *statementServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}
