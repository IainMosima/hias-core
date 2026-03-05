package billing

import (
	"context"
	"log"
	"net/http"
	"strings"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	billingEntity "github.com/bitbiz/hias-core/domains/billing/entity"
	billingRepo "github.com/bitbiz/hias-core/domains/billing/repository"
	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	billingService "github.com/bitbiz/hias-core/domains/billing/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/google/uuid"
)

type creditNoteServiceImpl struct {
	creditNoteRepo billingRepo.CreditNoteRepository
	invoiceRepo    billingRepo.InvoiceRepository
	auditSvc       auditService.AuditService
}

func NewCreditNoteService(
	creditNoteRepo billingRepo.CreditNoteRepository,
	invoiceRepo billingRepo.InvoiceRepository,
	auditSvc auditService.AuditService,
) billingService.CreditNoteService {
	return &creditNoteServiceImpl{
		creditNoteRepo: creditNoteRepo,
		invoiceRepo:    invoiceRepo,
		auditSvc:       auditSvc,
	}
}

func (s *creditNoteServiceImpl) CreateCreditNote(ctx context.Context, policyID, memberID uuid.UUID, amount int64, reason string, createdBy uuid.UUID) *schema.ServiceResponse[billingSchema.CreditNoteResponse] {
	cn := &billingEntity.CreditNote{
		PolicyID:         policyID,
		MemberID:         memberID,
		CreditNoteNumber: utils.GenerateCreditNoteNumber(),
		Amount:           amount,
		Currency:         string(shared.CurrencyKES),
		Reason:           reason,
		Status:           string(shared.CreditNoteStatusDraft),
		CreatedBy:        createdBy,
	}

	created, err := s.creditNoteRepo.Create(ctx, cn)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.CreditNoteResponse](http.StatusInternalServerError, "Failed to create credit note", err)
	}

	// Auto-approve system-initiated pro-rata refund credit notes
	if strings.Contains(reason, "Pro-rata refund") {
		approved, approveErr := s.creditNoteRepo.Approve(ctx, created.ID, createdBy)
		if approveErr != nil {
			log.Printf("Warning: failed to auto-approve pro-rata credit note: %v", approveErr)
		} else {
			created = approved
		}
	}

	s.logAudit(ctx, createdBy, string(shared.AuditEntityTypeCreditNote), created.ID, string(shared.AuditActionCreate))

	return schema.NewServiceResponse(billingSchema.ToCreditNoteResponse(created), http.StatusCreated, "Credit note created")
}

func (s *creditNoteServiceImpl) GetCreditNote(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[billingSchema.CreditNoteResponse] {
	cn, err := s.creditNoteRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.CreditNoteResponse](http.StatusNotFound, "Credit note not found", err)
	}
	return schema.NewServiceResponse(billingSchema.ToCreditNoteResponse(cn), http.StatusOK, "Credit note retrieved")
}

func (s *creditNoteServiceImpl) ListByPolicy(ctx context.Context, policyID uuid.UUID) *schema.ServiceResponse[[]billingSchema.CreditNoteResponse] {
	cns, err := s.creditNoteRepo.ListByPolicy(ctx, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]billingSchema.CreditNoteResponse](http.StatusInternalServerError, "Failed to list credit notes", err)
	}
	responses := make([]billingSchema.CreditNoteResponse, len(cns))
	for i, cn := range cns {
		responses[i] = billingSchema.ToCreditNoteResponse(cn)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Credit notes retrieved")
}

func (s *creditNoteServiceImpl) ApproveCreditNote(ctx context.Context, id uuid.UUID, approvedBy uuid.UUID) *schema.ServiceResponse[billingSchema.CreditNoteResponse] {
	cn, err := s.creditNoteRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.CreditNoteResponse](http.StatusNotFound, "Credit note not found", err)
	}
	if cn.Status != string(shared.CreditNoteStatusDraft) {
		return schema.NewServiceErrorResponse[billingSchema.CreditNoteResponse](http.StatusBadRequest, "Only DRAFT credit notes can be approved", nil)
	}

	approved, err := s.creditNoteRepo.Approve(ctx, id, approvedBy)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.CreditNoteResponse](http.StatusInternalServerError, "Failed to approve credit note", err)
	}

	s.logAudit(ctx, approvedBy, string(shared.AuditEntityTypeCreditNote), id, string(shared.AuditActionStateChange))

	return schema.NewServiceResponse(billingSchema.ToCreditNoteResponse(approved), http.StatusOK, "Credit note approved")
}

func (s *creditNoteServiceImpl) ApplyCreditNote(ctx context.Context, id uuid.UUID, invoiceID uuid.UUID) *schema.ServiceResponse[billingSchema.CreditNoteResponse] {
	cn, err := s.creditNoteRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.CreditNoteResponse](http.StatusNotFound, "Credit note not found", err)
	}
	if cn.Status != string(shared.CreditNoteStatusApproved) {
		return schema.NewServiceErrorResponse[billingSchema.CreditNoteResponse](http.StatusBadRequest, "Only APPROVED credit notes can be applied", nil)
	}

	// Validate invoice belongs to same policy
	if s.invoiceRepo != nil {
		invoice, invErr := s.invoiceRepo.GetByID(ctx, invoiceID)
		if invErr != nil {
			return schema.NewServiceErrorResponse[billingSchema.CreditNoteResponse](http.StatusNotFound, "Invoice not found", invErr)
		}
		if invoice.PolicyID != cn.PolicyID {
			return schema.NewServiceErrorResponse[billingSchema.CreditNoteResponse](http.StatusBadRequest, "Invoice does not belong to the same policy", nil)
		}
	}

	applied, err := s.creditNoteRepo.Apply(ctx, id, invoiceID)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.CreditNoteResponse](http.StatusInternalServerError, "Failed to apply credit note", err)
	}

	s.logAudit(ctx, uuid.Nil, string(shared.AuditEntityTypeCreditNote), id, string(shared.AuditActionStateChange))

	return schema.NewServiceResponse(billingSchema.ToCreditNoteResponse(applied), http.StatusOK, "Credit note applied")
}

func (s *creditNoteServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}
