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
