package handlers

import (
	"net/http"

	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	"github.com/bitbiz/hias-core/domains/billing/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InstallmentHandler struct {
	installmentSvc service.InstallmentService
}

func NewInstallmentHandler(installmentSvc service.InstallmentService) *InstallmentHandler {
	return &InstallmentHandler{installmentSvc: installmentSvc}
}

// CreateSchedule godoc
// @Summary      Create an installment schedule
// @Description  Create a new installment payment schedule for a policy
// @Tags         Installments
// @Accept       json
// @Produce      json
// @Param        id path string true "Policy ID"
// @Param        request body billingSchema.CreateInstallmentScheduleRequest true "Schedule creation request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/installments [post]
func (h *InstallmentHandler) CreateSchedule(ctx *gin.Context) {
	var req billingSchema.CreateInstallmentScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Use policy ID from URL if nested under policies
	if policyID := ctx.Param("id"); policyID != "" {
		req.PolicyID = policyID
	}

	resp := h.installmentSvc.CreateSchedule(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetSchedulesByPolicy godoc
// @Summary      Get installment schedules by policy
// @Description  Retrieve all installment schedules associated with a specific policy
// @Tags         Installments
// @Accept       json
// @Produce      json
// @Param        id path string true "Policy ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/installments [get]
func (h *InstallmentHandler) GetSchedulesByPolicy(ctx *gin.Context) {
	policyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	resp := h.installmentSvc.GetSchedulesByPolicy(ctx.Request.Context(), policyID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListInstallmentsBySchedule godoc
// @Summary      List installments by schedule
// @Description  Retrieve all installments for a specific installment schedule
// @Tags         Installments
// @Accept       json
// @Produce      json
// @Param        id path string true "Schedule ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/installments/schedule/{id} [get]
func (h *InstallmentHandler) ListInstallmentsBySchedule(ctx *gin.Context) {
	scheduleID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid schedule ID")
		return
	}

	resp := h.installmentSvc.ListInstallmentsBySchedule(ctx.Request.Context(), scheduleID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// MarkInstallmentPaid godoc
// @Summary      Mark installment as paid
// @Description  Mark a specific installment as paid with the associated invoice
// @Tags         Installments
// @Accept       json
// @Produce      json
// @Param        id path string true "Installment ID"
// @Param        request body billingSchema.MarkInstallmentPaidRequest true "Payment details"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/installments/{id}/pay [put]
func (h *InstallmentHandler) MarkInstallmentPaid(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid installment ID")
		return
	}

	var req billingSchema.MarkInstallmentPaidRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	invoiceID, _ := uuid.Parse(req.InvoiceID)

	resp := h.installmentSvc.MarkInstallmentPaid(ctx.Request.Context(), id, invoiceID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
