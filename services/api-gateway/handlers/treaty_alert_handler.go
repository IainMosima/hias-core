package handlers

import (
	"net/http"

	"github.com/bitbiz/hias-core/domains/reinsurance/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TreatyAlertHandler struct {
	alertSvc service.TreatyAlertService
}

func NewTreatyAlertHandler(alertSvc service.TreatyAlertService) *TreatyAlertHandler {
	return &TreatyAlertHandler{alertSvc: alertSvc}
}

func (h *TreatyAlertHandler) ListAlerts(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)

	resp := h.alertSvc.ListAlerts(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *TreatyAlertHandler) ListByTreaty(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)

	treatyIDStr := ctx.Query("treaty")
	if treatyIDStr == "" {
		utils.RespondError(ctx, http.StatusBadRequest, "treaty query parameter is required")
		return
	}

	treatyID, err := uuid.Parse(treatyIDStr)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid treaty ID")
		return
	}

	resp := h.alertSvc.ListByTreaty(ctx.Request.Context(), treatyID, pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *TreatyAlertHandler) ListUnacknowledged(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)

	resp := h.alertSvc.ListUnacknowledged(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *TreatyAlertHandler) CountUnacknowledged(ctx *gin.Context) {
	resp := h.alertSvc.CountUnacknowledged(ctx.Request.Context())
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *TreatyAlertHandler) AcknowledgeAlert(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid alert ID")
		return
	}

	resp := h.alertSvc.AcknowledgeAlert(ctx.Request.Context(), id, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *TreatyAlertHandler) CheckExpiryWarnings(ctx *gin.Context) {
	resp := h.alertSvc.CheckExpiryWarnings(ctx.Request.Context())
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *TreatyAlertHandler) CheckTreatyLimits(ctx *gin.Context) {
	treatyID, err := uuid.Parse(ctx.Param("treatyId"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid treaty ID")
		return
	}

	resp := h.alertSvc.CheckTreatyLimits(ctx.Request.Context(), treatyID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *TreatyAlertHandler) CheckCatastropheThresholds(ctx *gin.Context) {
	treatyID, err := uuid.Parse(ctx.Param("treatyId"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid treaty ID")
		return
	}

	resp := h.alertSvc.CheckCatastropheThresholds(ctx.Request.Context(), treatyID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
