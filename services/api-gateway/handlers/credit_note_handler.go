package handlers

import (
	"net/http"

	billingService "github.com/bitbiz/hias-core/domains/billing/service"
	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreditNoteHandler struct {
	creditNoteSvc billingService.CreditNoteService
}

func NewCreditNoteHandler(creditNoteSvc billingService.CreditNoteService) *CreditNoteHandler {
	return &CreditNoteHandler{creditNoteSvc: creditNoteSvc}
}

// ListByPolicy godoc
// @Summary      List credit notes by policy
// @Description  List all credit notes associated with a policy
// @Tags         CreditNotes
// @Accept       json
// @Produce      json
// @Param        id path string true "Policy ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/credit-notes [get]
func (h *CreditNoteHandler) ListByPolicy(ctx *gin.Context) {
	policyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}
	resp := h.creditNoteSvc.ListByPolicy(ctx.Request.Context(), policyID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetCreditNote godoc
// @Summary      Get a credit note
// @Description  Get credit note details by ID
// @Tags         CreditNotes
// @Accept       json
// @Produce      json
// @Param        id path string true "Credit Note ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/credit-notes/{id} [get]
func (h *CreditNoteHandler) GetCreditNote(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid credit note ID")
		return
	}
	resp := h.creditNoteSvc.GetCreditNote(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ApproveCreditNote godoc
// @Summary      Approve a credit note
// @Description  Approve a pending credit note by ID
// @Tags         CreditNotes
// @Accept       json
// @Produce      json
// @Param        id path string true "Credit Note ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/credit-notes/{id}/approve [put]
func (h *CreditNoteHandler) ApproveCreditNote(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid credit note ID")
		return
	}
	resp := h.creditNoteSvc.ApproveCreditNote(ctx.Request.Context(), id, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ApplyCreditNote godoc
// @Summary      Apply a credit note
// @Description  Apply an approved credit note to an invoice
// @Tags         CreditNotes
// @Accept       json
// @Produce      json
// @Param        id path string true "Credit Note ID"
// @Param        request body policySchema.ApplyCreditNoteRequest true "Invoice to apply credit note to"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/credit-notes/{id}/apply [put]
func (h *CreditNoteHandler) ApplyCreditNote(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid credit note ID")
		return
	}
	var req policySchema.ApplyCreditNoteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	invoiceID, err := uuid.Parse(req.InvoiceID)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid invoice ID")
		return
	}
	resp := h.creditNoteSvc.ApplyCreditNote(ctx.Request.Context(), id, invoiceID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
