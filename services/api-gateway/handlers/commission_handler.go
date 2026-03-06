package handlers

import (
	"net/http"

	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	"github.com/bitbiz/hias-core/domains/billing/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommissionHandler struct {
	svc service.CommissionService
}

func NewCommissionHandler(svc service.CommissionService) *CommissionHandler {
	return &CommissionHandler{svc: svc}
}

// CreateRule godoc
// @Summary      Create a commission rule
// @Description  Create a new commission rule for a plan
// @Tags         Commissions
// @Accept       json
// @Produce      json
// @Param        request body billingSchema.CreateCommissionRuleRequest true "Create commission rule request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/commissions/rules [post]
func (h *CommissionHandler) CreateRule(ctx *gin.Context) {
	var req billingSchema.CreateCommissionRuleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp := h.svc.CreateRule(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListRulesByPlan godoc
// @Summary      List commission rules by plan
// @Description  Retrieve all commission rules for a specific plan
// @Tags         Commissions
// @Accept       json
// @Produce      json
// @Param        id path string true "Plan ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/commissions/rules/plan/{id} [get]
func (h *CommissionHandler) ListRulesByPlan(ctx *gin.Context) {
	planID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid plan ID")
		return
	}
	resp := h.svc.ListRulesByPlan(ctx.Request.Context(), planID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// CalculateCommission godoc
// @Summary      Calculate commission
// @Description  Calculate commission for a given premium and plan
// @Tags         Commissions
// @Accept       json
// @Produce      json
// @Param        request body billingSchema.CalculateCommissionRequest true "Calculate commission request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/commissions/calculate [post]
func (h *CommissionHandler) CalculateCommission(ctx *gin.Context) {
	var req billingSchema.CalculateCommissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp := h.svc.CalculateCommission(ctx.Request.Context(), req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListPayments godoc
// @Summary      List commission payments
// @Description  Retrieve a paginated list of commission payments
// @Tags         Commissions
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/commissions/payments [get]
func (h *CommissionHandler) ListPayments(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)
	resp := h.svc.ListPayments(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ProcessPayments godoc
// @Summary      Process commission payments
// @Description  Trigger processing of pending commission payments
// @Tags         Commissions
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/commissions/payments/process [post]
func (h *CommissionHandler) ProcessPayments(ctx *gin.Context) {
	resp := h.svc.ProcessPayments(ctx.Request.Context())
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
