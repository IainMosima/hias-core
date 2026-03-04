package handlers

import (
	"net/http"

	"github.com/bitbiz/hias-core/domains/analytics/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	analyticsSvc service.AnalyticsService
}

func NewAnalyticsHandler(analyticsSvc service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsSvc: analyticsSvc}
}

func (h *AnalyticsHandler) GetDashboard(ctx *gin.Context) {
	period := ctx.DefaultQuery("period", "month")

	resp := h.analyticsSvc.GetDashboard(ctx.Request.Context(), period)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *AnalyticsHandler) GetKPIs(ctx *gin.Context) {
	period := ctx.DefaultQuery("period", "month")

	resp := h.analyticsSvc.GetKPIs(ctx.Request.Context(), period)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *AnalyticsHandler) ExportCSV(ctx *gin.Context) {
	reportType := ctx.DefaultQuery("type", "claims")
	period := ctx.DefaultQuery("period", "month")

	resp := h.analyticsSvc.ExportCSV(ctx.Request.Context(), reportType, period)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	ctx.Header("Content-Type", "text/csv")
	ctx.Header("Content-Disposition", "attachment; filename="+reportType+"_report.csv")
	ctx.Data(http.StatusOK, "text/csv", resp.Data)
}
