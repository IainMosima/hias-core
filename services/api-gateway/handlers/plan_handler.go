package handlers

import (
	"net/http"

	"github.com/bitbiz/hias-core/domains/product/schema"
	"github.com/bitbiz/hias-core/domains/product/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PlanHandler struct {
	planSvc service.PlanService
}

func NewPlanHandler(planSvc service.PlanService) *PlanHandler {
	return &PlanHandler{planSvc: planSvc}
}

// CreatePlan godoc
// @Summary      Create a new plan
// @Description  Create a new insurance plan with the provided details
// @Tags         Plans
// @Accept       json
// @Produce      json
// @Param        request body schema.CreatePlanRequest true "Plan creation payload"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/plans [post]
func (h *PlanHandler) CreatePlan(ctx *gin.Context) {
	var req schema.CreatePlanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.planSvc.CreatePlan(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetPlan godoc
// @Summary      Get a plan by ID
// @Description  Retrieve a single insurance plan by its unique identifier
// @Tags         Plans
// @Accept       json
// @Produce      json
// @Param        id path string true "Plan ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/plans/{id} [get]
func (h *PlanHandler) GetPlan(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid plan ID")
		return
	}

	resp := h.planSvc.GetPlan(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListPlans godoc
// @Summary      List all plans
// @Description  Retrieve a paginated list of insurance plans
// @Tags         Plans
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/plans [get]
func (h *PlanHandler) ListPlans(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)
	resp := h.planSvc.ListPlans(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	countResp := h.planSvc.GetTotalCount(ctx.Request.Context())
	utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
}

// UpdatePlan godoc
// @Summary      Update a plan
// @Description  Update an existing insurance plan by its unique identifier
// @Tags         Plans
// @Accept       json
// @Produce      json
// @Param        id path string true "Plan ID"
// @Param        request body schema.UpdatePlanRequest true "Plan update payload"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/plans/{id} [put]
func (h *PlanHandler) UpdatePlan(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid plan ID")
		return
	}

	var req schema.UpdatePlanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.planSvc.UpdatePlan(ctx.Request.Context(), id, req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
