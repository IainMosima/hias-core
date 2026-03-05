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
