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

// GetDashboard godoc
// @Summary      Get analytics dashboard
// @Description  Retrieve the analytics dashboard data for a given period
// @Tags         Analytics
// @Accept       json
// @Produce      json
// @Param        period query string false "Period (e.g. month, quarter, year)" default(month)
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/analytics/dashboard [get]
func (h *AnalyticsHandler) GetDashboard(ctx *gin.Context) {
	period := ctx.DefaultQuery("period", "month")

	resp := h.analyticsSvc.GetDashboard(ctx.Request.Context(), period)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetKPIs godoc
// @Summary      Get KPIs
// @Description  Retrieve key performance indicators for a given period
// @Tags         Analytics
// @Accept       json
// @Produce      json
// @Param        period query string false "Period (e.g. month, quarter, year)" default(month)
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/analytics/kpis [get]
func (h *AnalyticsHandler) GetKPIs(ctx *gin.Context) {
	period := ctx.DefaultQuery("period", "month")

	resp := h.analyticsSvc.GetKPIs(ctx.Request.Context(), period)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ExportCSV godoc
// @Summary      Export analytics as CSV
// @Description  Export analytics report data as a CSV file download
// @Tags         Analytics
// @Accept       json
// @Produce      text/csv
// @Param        type query string false "Report type (e.g. claims)" default(claims)
// @Param        period query string false "Period (e.g. month, quarter, year)" default(month)
// @Success      200 {file} file
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/analytics/export [get]
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
