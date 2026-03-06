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

// CreateLead godoc
// @Summary      Create a new lead
// @Description  Create a new sales lead with the provided details
// @Tags         Leads
// @Accept       json
// @Produce      json
// @Param        request body salesSchema.CreateLeadRequest true "Lead creation request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/leads [post]
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

// GetLead godoc
// @Summary      Get a lead by ID
// @Description  Retrieve a specific lead by its unique identifier
// @Tags         Leads
// @Accept       json
// @Produce      json
// @Param        id path string true "Lead ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/leads/{id} [get]
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

// ListLeads godoc
// @Summary      List leads
// @Description  Retrieve a paginated list of leads, optionally filtered by status or assignment
// @Tags         Leads
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Param        status query string false "Filter by lead status"
// @Param        assigned_to query string false "Filter by assignment (use 'me' for current user)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/leads [get]
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

// UpdateLead godoc
// @Summary      Update a lead
// @Description  Update an existing lead's details by its ID
// @Tags         Leads
// @Accept       json
// @Produce      json
// @Param        id path string true "Lead ID"
// @Param        request body salesSchema.UpdateLeadRequest true "Lead update request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/leads/{id} [put]
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

// UpdateLeadStatus godoc
// @Summary      Update lead status
// @Description  Update the status of an existing lead by its ID
// @Tags         Leads
// @Accept       json
// @Produce      json
// @Param        id path string true "Lead ID"
// @Param        request body salesSchema.UpdateLeadStatusRequest true "Status update request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/leads/{id}/status [put]
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

// AddActivity godoc
// @Summary      Add activity to a lead
// @Description  Add a new activity record to an existing lead
// @Tags         Leads
// @Accept       json
// @Produce      json
// @Param        id path string true "Lead ID"
// @Param        request body salesSchema.CreateLeadActivityRequest true "Activity creation request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/leads/{id}/activities [post]
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

// ListActivities godoc
// @Summary      List lead activities
// @Description  Retrieve all activities associated with a specific lead
// @Tags         Leads
// @Accept       json
// @Produce      json
// @Param        id path string true "Lead ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/leads/{id}/activities [get]
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

// GetDueFollowUps godoc
// @Summary      Get due follow-ups
// @Description  Retrieve a paginated list of leads with follow-ups that are due
// @Tags         Leads
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/leads/due-follow-ups [get]
func (h *LeadHandler) GetDueFollowUps(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)

	resp := h.leadSvc.GetDueFollowUps(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
