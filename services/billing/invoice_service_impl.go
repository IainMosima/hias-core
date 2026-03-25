package billing

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	"github.com/bitbiz/hias-core/domains/billing/entity"
	billingRepo "github.com/bitbiz/hias-core/domains/billing/repository"
	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	"github.com/bitbiz/hias-core/domains/billing/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	policyRepo "github.com/bitbiz/hias-core/domains/policy/repository"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/google/uuid"
)

type invoiceServiceImpl struct {
	invoiceRepo billingRepo.InvoiceRepository
	policyRepo  policyRepo.PolicyRepository
	auditSvc    auditService.AuditService
}

func NewInvoiceService(
	invoiceRepo billingRepo.InvoiceRepository,
	policyRepo policyRepo.PolicyRepository,
	auditSvc auditService.AuditService,
) service.InvoiceService {
	return &invoiceServiceImpl{
		invoiceRepo: invoiceRepo,
		policyRepo:  policyRepo,
		auditSvc:    auditSvc,
	}
}

func (s *invoiceServiceImpl) GetInvoice(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[billingSchema.InvoiceResponse] {
	invoice, err := s.invoiceRepo.GetWithPolicy(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.InvoiceResponse](
			http.StatusNotFound, "Invoice not found", err,
		)
	}

	resp := billingSchema.ToInvoiceResponse(invoice)
	return schema.NewServiceResponse(resp, http.StatusOK, "Invoice retrieved")
}

func (s *invoiceServiceImpl) ListInvoices(ctx context.Context, dateFrom, dateTo *time.Time, page, pageSize int) *schema.ServiceResponse[billingSchema.InvoiceListResponse] {
	offset := (page - 1) * pageSize

	if dateFrom != nil || dateTo != nil {
		invoices, err := s.invoiceRepo.ListFilteredWithPolicy(ctx, dateFrom, dateTo, pageSize, offset)
		if err != nil {
			return schema.NewServiceErrorResponse[billingSchema.InvoiceListResponse](
				http.StatusInternalServerError, fmt.Sprintf("failed to list invoices: %v", err), err,
			)
		}

		responses := make([]billingSchema.InvoiceResponse, len(invoices))
		for i, inv := range invoices {
			responses[i] = billingSchema.ToInvoiceResponse(inv)
		}

		count, _ := s.invoiceRepo.CountFiltered(ctx, dateFrom, dateTo)
		return schema.NewServiceResponse(billingSchema.InvoiceListResponse{
			Invoices:   responses,
			TotalCount: count,
		}, http.StatusOK, "Invoices retrieved")
	}

	invoices, err := s.invoiceRepo.ListWithPolicy(ctx, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.InvoiceListResponse](
			http.StatusInternalServerError, fmt.Sprintf("failed to list invoices: %v", err), err,
		)
	}

	responses := make([]billingSchema.InvoiceResponse, len(invoices))
	for i, inv := range invoices {
		responses[i] = billingSchema.ToInvoiceResponse(inv)
	}

	count, _ := s.invoiceRepo.Count(ctx)
	return schema.NewServiceResponse(billingSchema.InvoiceListResponse{
		Invoices:   responses,
		TotalCount: count,
	}, http.StatusOK, "Invoices retrieved")
}

func (s *invoiceServiceImpl) CreateInvoice(ctx context.Context, req billingSchema.CreateInvoiceRequest) *schema.ServiceResponse[billingSchema.InvoiceResponse] {
	policyID, err := uuid.Parse(req.PolicyID)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.InvoiceResponse](
			http.StatusBadRequest, "Invalid policy ID", err,
		)
	}

	policy, err := s.policyRepo.GetByID(ctx, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.InvoiceResponse](
			http.StatusNotFound, "Policy not found", err,
		)
	}

	if shared.PolicyStatus(policy.Status) != shared.PolicyStatusActive {
		return schema.NewServiceErrorResponse[billingSchema.InvoiceResponse](
			http.StatusBadRequest, "Policy must be ACTIVE to create an invoice", fmt.Errorf("policy status is %s", policy.Status),
		)
	}

	currency := req.Currency
	if currency == "" {
		currency = "KES"
	}

	now := time.Now()
	billingStart := now
	billingEnd := now.AddDate(0, 1, 0)
	if req.BillingPeriodStart != nil {
		billingStart = *req.BillingPeriodStart
	}
	if req.BillingPeriodEnd != nil {
		billingEnd = *req.BillingPeriodEnd
	}

	invoice := &entity.Invoice{
		PolicyID:           policyID,
		InvoiceNumber:      utils.GenerateInvoiceNumber(),
		Amount:             req.Amount,
		Currency:           currency,
		DueDate:            req.DueDate,
		Status:             string(shared.InvoiceStatusPending),
		BillingPeriodStart: billingStart,
		BillingPeriodEnd:   billingEnd,
		Notes:              req.Notes,
	}

	created, err := s.invoiceRepo.Create(ctx, invoice)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.InvoiceResponse](
			http.StatusInternalServerError, fmt.Sprintf("failed to create invoice: %v", err), err,
		)
	}

	s.logAudit(ctx, uuid.Nil, string(shared.AuditEntityTypeInvoice), created.ID, "CREATE")

	resp := billingSchema.ToInvoiceResponse(created)
	return schema.NewServiceResponse(resp, http.StatusCreated, "Invoice created")
}

func (s *invoiceServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}
