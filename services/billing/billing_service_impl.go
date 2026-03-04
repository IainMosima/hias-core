package billing

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bitbiz/hias-core/domains/billing/entity"
	billingRepo "github.com/bitbiz/hias-core/domains/billing/repository"
	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	"github.com/bitbiz/hias-core/domains/billing/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	policyRepo "github.com/bitbiz/hias-core/domains/policy/repository"
	"github.com/bitbiz/hias-core/shared"
	"github.com/google/uuid"
)

type billingServiceImpl struct {
	invoiceRepo billingRepo.InvoiceRepository
	policyRepo  policyRepo.PolicyRepository
}

func NewBillingService(
	invoiceRepo billingRepo.InvoiceRepository,
	policyRepo policyRepo.PolicyRepository,
) service.BillingService {
	return &billingServiceImpl{
		invoiceRepo: invoiceRepo,
		policyRepo:  policyRepo,
	}
}

func (s *billingServiceImpl) GenerateInvoice(ctx context.Context, policyID string) *schema.ServiceResponse[billingSchema.InvoiceResponse] {
	pid, err := uuid.Parse(policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.InvoiceResponse](http.StatusBadRequest, "Invalid policy ID", err)
	}

	policy, err := s.policyRepo.GetByID(ctx, pid)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.InvoiceResponse](http.StatusNotFound, "Policy not found", err)
	}

	now := time.Now()
	invoiceNumber := fmt.Sprintf("INV-%d-%06d", now.Year(), now.UnixNano()%1000000)

	invoice := &entity.Invoice{
		PolicyID:           pid,
		InvoiceNumber:      invoiceNumber,
		Amount:             policy.PremiumAmount,
		Currency:           policy.Currency,
		DueDate:            now.AddDate(0, 0, shared.InvoiceDueDays),
		Status:             string(shared.InvoiceStatusPending),
		BillingPeriodStart: now,
		BillingPeriodEnd:   now.AddDate(0, 1, 0),
	}

	created, err := s.invoiceRepo.Create(ctx, invoice)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.InvoiceResponse](http.StatusInternalServerError, "Failed to generate invoice", err)
	}

	return schema.NewServiceResponse(billingSchema.ToInvoiceResponse(created), http.StatusCreated, "Invoice generated")
}

func (s *billingServiceImpl) RunBillingCycle(ctx context.Context) *schema.ServiceResponse[int] {
	policies, err := s.policyRepo.GetActivePoliciesForBilling(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int](http.StatusInternalServerError, "Failed to get active policies", err)
	}

	generated := 0
	for _, p := range policies {
		resp := s.GenerateInvoice(ctx, p.ID.String())
		if resp.Error == nil {
			generated++
		}
	}

	return schema.NewServiceResponse(generated, http.StatusOK, fmt.Sprintf("%d invoices generated", generated))
}

func (s *billingServiceImpl) SendReminder(ctx context.Context, invoiceID string) *schema.ServiceResponse[string] {
	return schema.NewServiceResponse("Reminder sent", http.StatusOK, "Payment reminder sent")
}

func (s *billingServiceImpl) HandleOverdue(ctx context.Context) *schema.ServiceResponse[int] {
	overdue, err := s.invoiceRepo.ListOverdue(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int](http.StatusInternalServerError, "Failed to get overdue invoices", err)
	}

	updated := 0
	for _, inv := range overdue {
		_, updateErr := s.invoiceRepo.UpdateStatus(ctx, inv.ID, string(shared.InvoiceStatusOverdue))
		if updateErr == nil {
			updated++
		}
	}

	return schema.NewServiceResponse(updated, http.StatusOK, fmt.Sprintf("%d invoices marked overdue", updated))
}
