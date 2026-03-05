package handlers

import (
	"net/http"

	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/bitbiz/hias-core/domains/policy/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PolicyHandler struct {
	policySvc service.PolicyService
}

func NewPolicyHandler(policySvc service.PolicyService) *PolicyHandler {
	return &PolicyHandler{policySvc: policySvc}
}

func (h *PolicyHandler) CreatePolicy(ctx *gin.Context) {
	var req policySchema.CreatePolicyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.policySvc.CreatePolicy(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *PolicyHandler) GetPolicy(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	resp := h.policySvc.GetPolicy(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *PolicyHandler) ListPolicies(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)
	resp := h.policySvc.ListPolicies(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	countResp := h.policySvc.GetTotalCount(ctx.Request.Context())
	utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
}

func (h *PolicyHandler) ActivatePolicy(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	resp := h.policySvc.ActivatePolicy(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *PolicyHandler) LapsePolicy(ctx *gin.Context) {
	id, _ := uuid.Parse(ctx.Param("id"))
	resp := h.policySvc.LapsePolicy(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *PolicyHandler) TerminatePolicy(ctx *gin.Context) {
	id, _ := uuid.Parse(ctx.Param("id"))
	resp := h.policySvc.TerminatePolicy(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *PolicyHandler) ReinstatePolicy(ctx *gin.Context) {
	id, _ := uuid.Parse(ctx.Param("id"))
	resp := h.policySvc.ReinstatePolicy(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *PolicyHandler) CalculateProrate(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	resp := h.policySvc.CalculateProratedPremium(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *PolicyHandler) SuspendPolicy(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	resp := h.policySvc.SuspendPolicy(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *PolicyHandler) UpdatePolicy(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	var req policySchema.UpdatePolicyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.policySvc.UpdatePolicy(ctx.Request.Context(), id, req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *PolicyHandler) ListPoliciesByStatus(ctx *gin.Context) {
	status := ctx.Query("status")
	if status == "" {
		utils.RespondError(ctx, http.StatusBadRequest, "status query parameter is required")
		return
	}

	pagination := utils.GetPaginationParams(ctx)
	resp := h.policySvc.ListPoliciesByStatus(ctx.Request.Context(), status, pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *PolicyHandler) ChangePlan(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	var req policySchema.ChangePlanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.policySvc.ChangePlan(ctx.Request.Context(), id, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *PolicyHandler) BulkActivate(ctx *gin.Context) {
	var req policySchema.BulkIDsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var ids []uuid.UUID
	for _, idStr := range req.IDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID: "+idStr)
			return
		}
		ids = append(ids, id)
	}

	resp := h.policySvc.BulkActivate(ctx.Request.Context(), ids)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *PolicyHandler) BulkLapse(ctx *gin.Context) {
	var req policySchema.BulkIDsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var ids []uuid.UUID
	for _, idStr := range req.IDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID: "+idStr)
			return
		}
		ids = append(ids, id)
	}

	resp := h.policySvc.BulkLapse(ctx.Request.Context(), ids)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
