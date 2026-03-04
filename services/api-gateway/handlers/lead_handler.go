package handlers

import (
	"net/http"

	salesSchema "github.com/bitbiz/hias-core/domains/sales/schema"
	"github.com/bitbiz/hias-core/domains/sales/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LeadHandler struct {
	leadSvc service.LeadService
}

func NewLeadHandler(leadSvc service.LeadService) *LeadHandler {
	return &LeadHandler{leadSvc: leadSvc}
}

func (h *LeadHandler) CreateLead(ctx *gin.Context) {
	var req salesSchema.CreateLeadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.leadSvc.CreateLead(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *LeadHandler) GetLead(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid lead ID")
		return
	}

	resp := h.leadSvc.GetLead(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *LeadHandler) ListLeads(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)

	status := ctx.Query("status")
	assignedTo := ctx.Query("assigned_to")

	if status != "" {
		resp := h.leadSvc.ListLeadsByStatus(ctx.Request.Context(), status, pagination.Page, pagination.PageSize)
		if resp.Error != nil {
			utils.RespondError(ctx, resp.StatusCode, resp.Message)
			return
		}
		countResp := h.leadSvc.GetTotalCount(ctx.Request.Context())
		utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
		return
	}

	if assignedTo == "me" {
		resp := h.leadSvc.ListMyLeads(ctx.Request.Context(), getUserID(ctx), pagination.Page, pagination.PageSize)
		if resp.Error != nil {
			utils.RespondError(ctx, resp.StatusCode, resp.Message)
			return
		}
		countResp := h.leadSvc.GetTotalCount(ctx.Request.Context())
		utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
		return
	}

	resp := h.leadSvc.ListLeads(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	countResp := h.leadSvc.GetTotalCount(ctx.Request.Context())
	utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
}

func (h *LeadHandler) UpdateLead(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid lead ID")
		return
	}

	var req salesSchema.UpdateLeadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.leadSvc.UpdateLead(ctx.Request.Context(), id, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *LeadHandler) UpdateLeadStatus(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid lead ID")
		return
	}

	var req salesSchema.UpdateLeadStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.leadSvc.UpdateLeadStatus(ctx.Request.Context(), id, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *LeadHandler) AddActivity(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid lead ID")
		return
	}

	var req salesSchema.CreateLeadActivityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.leadSvc.AddActivity(ctx.Request.Context(), id, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *LeadHandler) ListActivities(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid lead ID")
		return
	}

	resp := h.leadSvc.ListActivities(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *LeadHandler) GetDueFollowUps(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)

	resp := h.leadSvc.GetDueFollowUps(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
