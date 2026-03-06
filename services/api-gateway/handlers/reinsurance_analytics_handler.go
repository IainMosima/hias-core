package handlers

import (
	"net/http"

	reinsuranceSchema "github.com/bitbiz/hias-core/domains/reinsurance/schema"
	"github.com/bitbiz/hias-core/domains/reinsurance/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
)

type ReinsuranceAnalyticsHandler struct {
	treatySvc   service.TreatyService
	cessionSvc  service.CessionService
	recoverySvc service.RecoveryService
	alertSvc    service.TreatyAlertService
}

func NewReinsuranceAnalyticsHandler(
	treatySvc service.TreatyService,
	cessionSvc service.CessionService,
	recoverySvc service.RecoveryService,
	alertSvc service.TreatyAlertService,
) *ReinsuranceAnalyticsHandler {
	return &ReinsuranceAnalyticsHandler{
		treatySvc:   treatySvc,
		cessionSvc:  cessionSvc,
		recoverySvc: recoverySvc,
		alertSvc:    alertSvc,
	}
}

// GetReinsuranceDashboard godoc
// @Summary      Get reinsurance dashboard
// @Description  Retrieve aggregated reinsurance analytics including treaty counts, ceded premiums, recoveries, and alerts
// @Tags         ReinsuranceAnalytics
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/analytics/reinsurance [get]
func (h *ReinsuranceAnalyticsHandler) GetReinsuranceDashboard(ctx *gin.Context) {
	c := ctx.Request.Context()

	treatyCountResp := h.treatySvc.GetTreatyCount(c)
	if treatyCountResp.Error != nil {
		utils.RespondError(ctx, treatyCountResp.StatusCode, treatyCountResp.Message)
		return
	}

	totalCededResp := h.cessionSvc.GetTotalCededAmount(c)
	if totalCededResp.Error != nil {
		utils.RespondError(ctx, totalCededResp.StatusCode, totalCededResp.Message)
		return
	}

	totalGrossResp := h.cessionSvc.GetTotalGrossAmount(c)
	if totalGrossResp.Error != nil {
		utils.RespondError(ctx, totalGrossResp.StatusCode, totalGrossResp.Message)
		return
	}

	totalRecoverableResp := h.recoverySvc.GetTotalRecoverableAmount(c)
	if totalRecoverableResp.Error != nil {
		utils.RespondError(ctx, totalRecoverableResp.StatusCode, totalRecoverableResp.Message)
		return
	}

	totalRecoveredResp := h.recoverySvc.GetTotalRecoveredAmount(c)
	if totalRecoveredResp.Error != nil {
		utils.RespondError(ctx, totalRecoveredResp.StatusCode, totalRecoveredResp.Message)
		return
	}

	alertCountResp := h.alertSvc.CountUnacknowledged(c)
	if alertCountResp.Error != nil {
		utils.RespondError(ctx, alertCountResp.StatusCode, alertCountResp.Message)
		return
	}

	totalCeded := totalCededResp.Data
	totalGross := totalGrossResp.Data
	totalRecoverable := totalRecoverableResp.Data
	totalRecovered := totalRecoveredResp.Data
	totalOutstanding := totalRecoverable - totalRecovered

	var cessionRatio float64
	if totalGross > 0 {
		cessionRatio = float64(totalCeded) / float64(totalGross)
	}

	var recoverySuccessRate float64
	if totalRecoverable > 0 {
		recoverySuccessRate = float64(totalRecovered) / float64(totalRecoverable)
	}

	dashboard := reinsuranceSchema.ReinsuranceDashboardResponse{
		ActiveTreatyCount:    treatyCountResp.Data,
		TotalCededPremiums:   totalCeded,
		TotalRecoverable:     totalRecoverable,
		TotalRecovered:       totalRecovered,
		TotalOutstanding:     totalOutstanding,
		CessionRatio:         cessionRatio,
		RecoverySuccessRate:  recoverySuccessRate,
		UnacknowledgedAlerts: alertCountResp.Data,
	}

	utils.RespondSuccess(ctx, http.StatusOK, "Reinsurance dashboard retrieved successfully", dashboard)
}
