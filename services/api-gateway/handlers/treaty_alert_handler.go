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

// ListAlerts godoc
// @Summary      List treaty alerts
// @Description  Retrieve a paginated list of all treaty alerts
// @Tags         TreatyAlerts
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaty-alerts [get]
func (h *TreatyAlertHandler) ListAlerts(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)

	resp := h.alertSvc.ListAlerts(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListByTreaty godoc
// @Summary      List alerts by treaty
// @Description  Retrieve a paginated list of alerts for a specific treaty
// @Tags         TreatyAlerts
// @Accept       json
// @Produce      json
// @Param        treaty query string true "Treaty ID"
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaties/{id}/alerts [get]
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

// ListUnacknowledged godoc
// @Summary      List unacknowledged alerts
// @Description  Retrieve a paginated list of unacknowledged treaty alerts
// @Tags         TreatyAlerts
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaty-alerts/unacknowledged [get]
func (h *TreatyAlertHandler) ListUnacknowledged(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)

	resp := h.alertSvc.ListUnacknowledged(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// CountUnacknowledged godoc
// @Summary      Count unacknowledged alerts
// @Description  Get the count of unacknowledged treaty alerts
// @Tags         TreatyAlerts
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaty-alerts/count [get]
func (h *TreatyAlertHandler) CountUnacknowledged(ctx *gin.Context) {
	resp := h.alertSvc.CountUnacknowledged(ctx.Request.Context())
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// AcknowledgeAlert godoc
// @Summary      Acknowledge a treaty alert
// @Description  Acknowledge a treaty alert by ID
// @Tags         TreatyAlerts
// @Accept       json
// @Produce      json
// @Param        id path string true "Alert ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaty-alerts/{id}/acknowledge [put]
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

// CheckExpiryWarnings godoc
// @Summary      Check expiry warnings
// @Description  Check for treaties nearing expiry and generate alerts
// @Tags         TreatyAlerts
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaty-alerts/check-expiry [post]
func (h *TreatyAlertHandler) CheckExpiryWarnings(ctx *gin.Context) {
	resp := h.alertSvc.CheckExpiryWarnings(ctx.Request.Context())
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// CheckTreatyLimits godoc
// @Summary      Check treaty limits
// @Description  Check if a treaty is approaching or exceeding its limits
// @Tags         TreatyAlerts
// @Accept       json
// @Produce      json
// @Param        treatyId path string true "Treaty ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaty-alerts/check-limits/{treatyId} [post]
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

// CheckCatastropheThresholds godoc
// @Summary      Check catastrophe thresholds
// @Description  Check if a treaty has reached catastrophe thresholds
// @Tags         TreatyAlerts
// @Accept       json
// @Produce      json
// @Param        treatyId path string true "Treaty ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaty-alerts/check-catastrophe/{treatyId} [post]
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
