package handlers

import (
	"github.com/bitbiz/hias-core/domains/analytics/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	analyticsService service.AnalyticsService
}

func NewAnalyticsHandler(analyticsService service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

func (h *AnalyticsHandler) GetDashboard(ctx *gin.Context) {
	period := ctx.DefaultQuery("period", "monthly")

	resp := h.analyticsService.GetDashboard(ctx.Request.Context(), period)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *AnalyticsHandler) GetKPIs(ctx *gin.Context) {
	period := ctx.DefaultQuery("period", "monthly")

	resp := h.analyticsService.GetKPIs(ctx.Request.Context(), period)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
