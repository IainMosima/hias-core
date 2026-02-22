package billing

import (
	"context"
	"fmt"
	"net/http"
	"time"

	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	"github.com/bitbiz/hias-core/domains/billing/entity"
	"github.com/bitbiz/hias-core/domains/billing/repository"
	"github.com/bitbiz/hias-core/domains/billing/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	policyDomainRepo "github.com/bitbiz/hias-core/domains/policy/repository"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/google/uuid"
)

type billingServiceImpl struct {
	invoiceRepo repository.InvoiceRepository
	policyRepo  policyDomainRepo.PolicyRepository
}

func NewBillingService(
	invoiceRepo repository.InvoiceRepository,
	policyRepo policyDomainRepo.PolicyRepository,
) service.BillingService {
	return &billingServiceImpl{
		invoiceRepo: invoiceRepo,
		policyRepo:  policyRepo,
	}
}

// GenerateInvoice creates a new invoice for a specific policy based on the
// plan's base premium amount. The billing period covers the current month.
func (s *billingServiceImpl) GenerateInvoice(ctx context.Context, policyID string) *schema.ServiceResponse[billingSchema.InvoiceResponse] {
	polID, err := uuid.Parse(policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.InvoiceResponse](http.StatusBadRequest, "Invalid policy ID", err)
	}

	// Fetch the policy to get premium amount
	policy, err := s.policyRepo.GetByID(ctx, polID)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.InvoiceResponse](http.StatusNotFound, "Policy not found", err)
	}

	if policy.Status != string(shared.PolicyStatusActive) {
		return schema.NewServiceErrorResponse[billingSchema.InvoiceResponse](
			http.StatusConflict,
			fmt.Sprintf("Cannot generate invoice for policy in status %s", policy.Status),
			fmt.Errorf("policy %s is not active", policyID),
		)
	}

	// Calculate billing period (current month)
	now := time.Now()
	billingStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	billingEnd := billingStart.AddDate(0, 1, -1) // Last day of the month
	dueDate := billingStart.AddDate(0, 0, 15)     // Due 15th of the month

	invoiceNumber := generateInvoiceNumber()

	invoice := &entity.Invoice{
		PolicyID:           polID,
		InvoiceNumber:      invoiceNumber,
		Amount:             policy.PremiumAmount,
		Currency:           policy.Currency,
		DueDate:            dueDate,
		Status:             string(shared.InvoiceStatusPending),
		BillingPeriodStart: billingStart,
		BillingPeriodEnd:   billingEnd,
		Notes:              fmt.Sprintf("Premium invoice for policy %s - %s %d", policy.PolicyNumber, now.Month().String(), now.Year()),
	}

	invoice, err = s.invoiceRepo.Create(ctx, invoice)
	if err != nil {
		utils.LogError("Failed to create invoice for policy %s: %v", policyID, err)
		return schema.NewServiceErrorResponse[billingSchema.InvoiceResponse](http.StatusInternalServerError, "Failed to create invoice", err)
	}

	utils.LogInfo("Invoice %s created for policy %s, amount: %d %s", invoiceNumber, policy.PolicyNumber, policy.PremiumAmount, policy.Currency)

	return schema.NewServiceResponse(billingSchema.ToInvoiceResponse(invoice), http.StatusCreated, "Invoice generated")
}

// RunBillingCycle iterates through all active policies and generates
// invoices for each. This is intended to be run by a scheduled job
// (e.g., monthly billing cycle).
func (s *billingServiceImpl) RunBillingCycle(ctx context.Context) *schema.ServiceResponse[int] {
	policies, err := s.policyRepo.GetActivePoliciesForBilling(ctx)
	if err != nil {
		utils.LogError("Failed to get active policies for billing: %v", err)
		return schema.NewServiceErrorResponse[int](http.StatusInternalServerError, "Failed to get active policies", err)
	}

	generatedCount := 0
	for _, policy := range policies {
		resp := s.GenerateInvoice(ctx, policy.ID.String())
		if resp.Error != nil {
			utils.LogError("Failed to generate invoice for policy %s: %v", policy.PolicyNumber, resp.Error)
			continue
		}
		generatedCount++
	}

	utils.LogInfo("Billing cycle completed: generated %d invoices for %d active policies", generatedCount, len(policies))

	return schema.NewServiceResponse(generatedCount, http.StatusOK, fmt.Sprintf("Billing cycle completed: %d invoices generated", generatedCount))
}

// SendReminder marks an invoice reminder as sent. In a full implementation,
// this would trigger a notification (email/SMS) to the policyholder.
func (s *billingServiceImpl) SendReminder(ctx context.Context, invoiceID string) *schema.ServiceResponse[string] {
	invID, err := uuid.Parse(invoiceID)
	if err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusBadRequest, "Invalid invoice ID", err)
	}

	invoice, err := s.invoiceRepo.GetByID(ctx, invID)
	if err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusNotFound, "Invoice not found", err)
	}

	if invoice.Status == string(shared.InvoiceStatusPaid) {
		return schema.NewServiceResponse("Invoice already paid", http.StatusOK, "No reminder needed")
	}

	// In a full implementation, this would:
	// 1. Look up the policyholder's contact info
	// 2. Send an SMS/email reminder via the notification service
	// 3. Log the reminder in the audit trail
	utils.LogInfo("Payment reminder sent for invoice %s (amount: %d %s, due: %s)",
		invoice.InvoiceNumber, invoice.Amount, invoice.Currency, invoice.DueDate.Format("2006-01-02"))

	return schema.NewServiceResponse(
		fmt.Sprintf("Reminder sent for invoice %s", invoice.InvoiceNumber),
		http.StatusOK,
		"Payment reminder sent",
	)
}

// HandleOverdue finds all overdue invoices and updates their status.
// For policies with overdue invoices, the policy may be flagged for lapse.
// This is intended to be run by a scheduled job.
func (s *billingServiceImpl) HandleOverdue(ctx context.Context) *schema.ServiceResponse[int] {
	overdueInvoices, err := s.invoiceRepo.ListOverdue(ctx)
	if err != nil {
		utils.LogError("Failed to get overdue invoices: %v", err)
		return schema.NewServiceErrorResponse[int](http.StatusInternalServerError, "Failed to get overdue invoices", err)
	}

	overdueCount := 0
	for _, invoice := range overdueInvoices {
		if invoice.Status == string(shared.InvoiceStatusOverdue) {
			// Already marked overdue; skip
			continue
		}

		_, err := s.invoiceRepo.UpdateStatus(ctx, invoice.ID, string(shared.InvoiceStatusOverdue))
		if err != nil {
			utils.LogError("Failed to mark invoice %s as overdue: %v", invoice.InvoiceNumber, err)
			continue
		}
		overdueCount++

		utils.LogWarn("Invoice %s marked OVERDUE (policy: %s, amount: %d %s, due: %s)",
			invoice.InvoiceNumber, invoice.PolicyID, invoice.Amount, invoice.Currency,
			invoice.DueDate.Format("2006-01-02"))
	}

	utils.LogInfo("Overdue handling completed: %d invoices marked overdue", overdueCount)

	return schema.NewServiceResponse(overdueCount, http.StatusOK, fmt.Sprintf("%d invoices marked as overdue", overdueCount))
}

// generateInvoiceNumber creates a unique invoice number.
// Format: INV-YYYY-<8 hex chars>
func generateInvoiceNumber() string {
	id := uuid.New()
	return fmt.Sprintf("INV-%d-%s", time.Now().Year(), id.String()[:8])
}
