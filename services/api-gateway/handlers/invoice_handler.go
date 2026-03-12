package handlers

import (
	"net/http"
	"time"

	"github.com/bitbiz/hias-core/domains/billing/entity"
	"github.com/bitbiz/hias-core/domains/billing/repository"
	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	policyRepo "github.com/bitbiz/hias-core/domains/policy/repository"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InvoiceHandler struct {
	invoiceRepo repository.InvoiceRepository
	policyRepo  policyRepo.PolicyRepository
}

func NewInvoiceHandler(invoiceRepo repository.InvoiceRepository, policyRepo policyRepo.PolicyRepository) *InvoiceHandler {
	return &InvoiceHandler{invoiceRepo: invoiceRepo, policyRepo: policyRepo}
}

// GetInvoice godoc
// @Summary      Get an invoice by ID
// @Description  Retrieve a single invoice by its unique identifier
// @Tags         Invoices
// @Accept       json
// @Produce      json
// @Param        id path string true "Invoice ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/invoices/{id} [get]
func (h *InvoiceHandler) GetInvoice(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid invoice ID")
		return
	}

	invoice, err := h.invoiceRepo.GetWithPolicy(ctx.Request.Context(), id)
	if err != nil {
		utils.RespondError(ctx, http.StatusNotFound, "Invoice not found")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, "Invoice retrieved", billingSchema.ToInvoiceResponse(invoice))
}

// ListInvoices godoc
// @Summary      List all invoices
// @Description  Retrieve a paginated list of invoices with optional date range filtering
// @Tags         Invoices
// @Accept       json
// @Produce      json
// @Param        date_from query string false "Filter from date (RFC3339 or YYYY-MM-DD)"
// @Param        date_to query string false "Filter to date (RFC3339 or YYYY-MM-DD)"
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/invoices [get]
func (h *InvoiceHandler) ListInvoices(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)
	offset := (pagination.Page - 1) * pagination.PageSize

	dateFromStr := ctx.Query("date_from")
	dateToStr := ctx.Query("date_to")

	if dateFromStr != "" || dateToStr != "" {
		var dateFrom, dateTo *time.Time
		if dateFromStr != "" {
			if t, err := time.Parse(time.RFC3339, dateFromStr); err == nil {
				dateFrom = &t
			} else if t, err := time.Parse("2006-01-02", dateFromStr); err == nil {
				dateFrom = &t
			}
		}
		if dateToStr != "" {
			if t, err := time.Parse(time.RFC3339, dateToStr); err == nil {
				dateTo = &t
			} else if t, err := time.Parse("2006-01-02", dateToStr); err == nil {
				endOfDay := t.Add(24*time.Hour - time.Second)
				dateTo = &endOfDay
			}
		}

		invoices, err := h.invoiceRepo.ListFilteredWithPolicy(ctx.Request.Context(), dateFrom, dateTo, pagination.PageSize, offset)
		if err != nil {
			utils.RespondError(ctx, http.StatusInternalServerError, "Failed to list invoices")
			return
		}

		responses := make([]billingSchema.InvoiceResponse, len(invoices))
		for i, inv := range invoices {
			responses[i] = billingSchema.ToInvoiceResponse(inv)
		}

		count, _ := h.invoiceRepo.CountFiltered(ctx.Request.Context(), dateFrom, dateTo)
		utils.RespondPaginated(ctx, "Invoices retrieved", responses, pagination.Page, pagination.PageSize, count)
		return
	}

	invoices, err := h.invoiceRepo.ListWithPolicy(ctx.Request.Context(), pagination.PageSize, offset)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "Failed to list invoices")
		return
	}

	responses := make([]billingSchema.InvoiceResponse, len(invoices))
	for i, inv := range invoices {
		responses[i] = billingSchema.ToInvoiceResponse(inv)
	}

	count, _ := h.invoiceRepo.Count(ctx.Request.Context())
	utils.RespondPaginated(ctx, "Invoices retrieved", responses, pagination.Page, pagination.PageSize, count)
}

// CreateInvoice godoc
// @Summary      Create a new invoice
// @Description  Create an invoice manually for a given policy
// @Tags         Invoices
// @Accept       json
// @Produce      json
// @Param        request body billingSchema.CreateInvoiceRequest true "Create Invoice Request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/invoices [post]
func (h *InvoiceHandler) CreateInvoice(ctx *gin.Context) {
	var req billingSchema.CreateInvoiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	policyID, err := uuid.Parse(req.PolicyID)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	policy, err := h.policyRepo.GetByID(ctx.Request.Context(), policyID)
	if err != nil {
		utils.RespondError(ctx, http.StatusNotFound, "Policy not found")
		return
	}

	if shared.PolicyStatus(policy.Status) != shared.PolicyStatusActive {
		utils.RespondError(ctx, http.StatusBadRequest, "Policy must be ACTIVE to create an invoice")
		return
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

	created, err := h.invoiceRepo.Create(ctx.Request.Context(), invoice)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "Failed to create invoice")
		return
	}

	utils.RespondSuccess(ctx, http.StatusCreated, "Invoice created", billingSchema.ToInvoiceResponse(created))
}

// DownloadInvoice godoc
// @Summary      Download invoice PDF
// @Description  Download an invoice as a PDF document (not yet implemented)
// @Tags         Invoices
// @Produce      json
// @Param        id path string true "Invoice ID"
// @Success      200 {file} file
// @Failure      501 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/invoices/{id}/download [get]
func (h *InvoiceHandler) DownloadInvoice(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{"message": "Invoice PDF generation not yet available"})
}
