package handlers

import (
	"net/http"

	"github.com/bitbiz/hias-core/domains/billing/repository"
	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InvoiceHandler struct {
	invoiceRepo repository.InvoiceRepository
}

func NewInvoiceHandler(invoiceRepo repository.InvoiceRepository) *InvoiceHandler {
	return &InvoiceHandler{invoiceRepo: invoiceRepo}
}

func (h *InvoiceHandler) GetInvoice(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid invoice ID")
		return
	}

	invoice, err := h.invoiceRepo.GetByID(ctx.Request.Context(), id)
	if err != nil {
		utils.RespondError(ctx, http.StatusNotFound, "Invoice not found")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, "Invoice retrieved", billingSchema.ToInvoiceResponse(invoice))
}

func (h *InvoiceHandler) ListInvoices(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)
	offset := (pagination.Page - 1) * pagination.PageSize

	invoices, err := h.invoiceRepo.List(ctx.Request.Context(), pagination.PageSize, offset)
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
