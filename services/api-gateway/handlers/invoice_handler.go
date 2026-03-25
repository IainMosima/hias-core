package handlers

import (
	"net/http"
	"time"

	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	"github.com/bitbiz/hias-core/domains/billing/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InvoiceHandler struct {
	invoiceSvc service.InvoiceService
}

func NewInvoiceHandler(invoiceSvc service.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{invoiceSvc: invoiceSvc}
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

	resp := h.invoiceSvc.GetInvoice(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
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

	var dateFrom, dateTo *time.Time
	if dateFromStr := ctx.Query("date_from"); dateFromStr != "" {
		if t, err := time.Parse(time.RFC3339, dateFromStr); err == nil {
			dateFrom = &t
		} else if t, err := time.Parse("2006-01-02", dateFromStr); err == nil {
			dateFrom = &t
		}
	}
	if dateToStr := ctx.Query("date_to"); dateToStr != "" {
		if t, err := time.Parse(time.RFC3339, dateToStr); err == nil {
			dateTo = &t
		} else if t, err := time.Parse("2006-01-02", dateToStr); err == nil {
			endOfDay := t.Add(24*time.Hour - time.Second)
			dateTo = &endOfDay
		}
	}

	resp := h.invoiceSvc.ListInvoices(ctx.Request.Context(), dateFrom, dateTo, pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondPaginated(ctx, resp.Message, resp.Data.Invoices, pagination.Page, pagination.PageSize, resp.Data.TotalCount)
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

	resp := h.invoiceSvc.CreateInvoice(ctx.Request.Context(), req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
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
