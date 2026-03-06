package handlers

import (
	"net/http"

	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	"github.com/bitbiz/hias-core/domains/billing/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RefundHandler struct {
	svc service.RefundService
}

func NewRefundHandler(svc service.RefundService) *RefundHandler {
	return &RefundHandler{svc: svc}
}

// RequestRefund godoc
// @Summary      Request a refund
// @Description  Create a new refund request for a policy
// @Tags         Refunds
// @Accept       json
// @Produce      json
// @Param        request body billingSchema.CreateRefundRequest true "Create refund request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/refunds [post]
func (h *RefundHandler) RequestRefund(ctx *gin.Context) {
	var req billingSchema.CreateRefundRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp := h.svc.RequestRefund(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ApproveRefund godoc
// @Summary      Approve a refund
// @Description  Approve a pending refund request by ID
// @Tags         Refunds
// @Accept       json
// @Produce      json
// @Param        id path string true "Refund ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/refunds/{id}/approve [put]
func (h *RefundHandler) ApproveRefund(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid refund ID")
		return
	}
	resp := h.svc.ApproveRefund(ctx.Request.Context(), id, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ProcessRefund godoc
// @Summary      Process a refund
// @Description  Process an approved refund by ID
// @Tags         Refunds
// @Accept       json
// @Produce      json
// @Param        id path string true "Refund ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/refunds/{id}/process [put]
func (h *RefundHandler) ProcessRefund(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid refund ID")
		return
	}
	resp := h.svc.ProcessRefund(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListRefunds godoc
// @Summary      List refunds for a policy
// @Description  Retrieve a paginated list of refunds for a specific policy
// @Tags         Refunds
// @Accept       json
// @Produce      json
// @Param        id path string true "Policy ID"
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/refunds [get]
func (h *RefundHandler) ListRefunds(ctx *gin.Context) {
	policyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}
	pagination := utils.GetPaginationParams(ctx)
	resp := h.svc.ListRefunds(ctx.Request.Context(), policyID, pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
