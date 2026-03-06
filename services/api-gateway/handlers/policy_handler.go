package handlers

import (
	"net/http"
	"time"

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

// CreatePolicy godoc
// @Summary      Create a new policy
// @Description  Create a new insurance policy
// @Tags         Policies
// @Accept       json
// @Produce      json
// @Param        request body policySchema.CreatePolicyRequest true "Policy creation details"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies [post]
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

// GetPolicy godoc
// @Summary      Get a policy by ID
// @Description  Retrieve a single policy by its unique identifier
// @Tags         Policies
// @Produce      json
// @Param        id path string true "Policy ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id} [get]
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

// ListPolicies godoc
// @Summary      List policies
// @Description  List all policies with pagination and optional date/search filters
// @Tags         Policies
// @Produce      json
// @Param        date_from query string false "Filter from date (RFC3339 or YYYY-MM-DD)"
// @Param        date_to query string false "Filter to date (RFC3339 or YYYY-MM-DD)"
// @Param        search query string false "Search by policy number, holder name, or email"
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies [get]
func (h *PolicyHandler) ListPolicies(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)

	dateFromStr := ctx.Query("date_from")
	dateToStr := ctx.Query("date_to")
	search := ctx.Query("search")

	if dateFromStr != "" || dateToStr != "" || search != "" {
		var dateFrom, dateTo *time.Time
		if dateFromStr != "" {
			if t, err := time.Parse(time.RFC3339, dateFromStr); err == nil {
				dateFrom = &t
			} else if t, err := time.Parse("2006-01-02", dateFromStr); err == nil {
				dateFrom = &t
			}
		}
		if dateToStr != "" {
			if t, err := time.Parse(time.RFC3339, dateToStr); err == nil {
				dateTo = &t
			} else if t, err := time.Parse("2006-01-02", dateToStr); err == nil {
				endOfDay := t.Add(24*time.Hour - time.Second)
				dateTo = &endOfDay
			}
		}

		resp := h.policySvc.ListPoliciesFiltered(ctx.Request.Context(), dateFrom, dateTo, search, pagination.Page, pagination.PageSize)
		if resp.Error != nil {
			utils.RespondError(ctx, resp.StatusCode, resp.Message)
			return
		}
		countResp := h.policySvc.CountPoliciesFiltered(ctx.Request.Context(), dateFrom, dateTo, search)
		utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
		return
	}

	resp := h.policySvc.ListPolicies(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	countResp := h.policySvc.GetTotalCount(ctx.Request.Context())
	utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
}

// ActivatePolicy godoc
// @Summary      Activate a policy
// @Description  Activate an insurance policy by its ID
// @Tags         Policies
// @Produce      json
// @Param        id path string true "Policy ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/activate [put]
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

// LapsePolicy godoc
// @Summary      Lapse a policy
// @Description  Mark a policy as lapsed due to non-payment or expiry
// @Tags         Policies
// @Produce      json
// @Param        id path string true "Policy ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/lapse [put]
func (h *PolicyHandler) LapsePolicy(ctx *gin.Context) {
	id, _ := uuid.Parse(ctx.Param("id"))
	resp := h.policySvc.LapsePolicy(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// TerminatePolicy godoc
// @Summary      Terminate a policy
// @Description  Permanently terminate an insurance policy
// @Tags         Policies
// @Produce      json
// @Param        id path string true "Policy ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/terminate [put]
func (h *PolicyHandler) TerminatePolicy(ctx *gin.Context) {
	id, _ := uuid.Parse(ctx.Param("id"))
	resp := h.policySvc.TerminatePolicy(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ReinstatePolicy godoc
// @Summary      Reinstate a policy
// @Description  Reinstate a previously lapsed or suspended policy
// @Tags         Policies
// @Produce      json
// @Param        id path string true "Policy ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/reinstate [put]
func (h *PolicyHandler) ReinstatePolicy(ctx *gin.Context) {
	id, _ := uuid.Parse(ctx.Param("id"))
	resp := h.policySvc.ReinstatePolicy(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// CalculateProrate godoc
// @Summary      Calculate prorated premium
// @Description  Calculate the prorated premium for a policy
// @Tags         Policies
// @Produce      json
// @Param        id path string true "Policy ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/prorate [get]
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

// SuspendPolicy godoc
// @Summary      Suspend a policy
// @Description  Temporarily suspend an insurance policy
// @Tags         Policies
// @Produce      json
// @Param        id path string true "Policy ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/suspend [put]
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

// UpdatePolicy godoc
// @Summary      Update a policy
// @Description  Update an existing insurance policy by its ID
// @Tags         Policies
// @Accept       json
// @Produce      json
// @Param        id path string true "Policy ID"
// @Param        request body policySchema.UpdatePolicyRequest true "Policy update details"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id} [put]
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

// ListPoliciesByStatus godoc
// @Summary      List policies by status
// @Description  List policies filtered by their status
// @Tags         Policies
// @Produce      json
// @Param        status query string true "Policy status to filter by"
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/by-status [get]
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

// ChangePlan godoc
// @Summary      Change policy plan
// @Description  Change the plan associated with an insurance policy
// @Tags         Policies
// @Accept       json
// @Produce      json
// @Param        id path string true "Policy ID"
// @Param        request body policySchema.ChangePlanRequest true "New plan details"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/change-plan [put]
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

// BulkActivate godoc
// @Summary      Bulk activate policies
// @Description  Activate multiple policies in a single request
// @Tags         Policies
// @Accept       json
// @Produce      json
// @Param        request body policySchema.BulkIDsRequest true "List of policy IDs to activate"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/bulk/activate [post]
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

// BulkLapse godoc
// @Summary      Bulk lapse policies
// @Description  Lapse multiple policies in a single request
// @Tags         Policies
// @Accept       json
// @Produce      json
// @Param        request body policySchema.BulkIDsRequest true "List of policy IDs to lapse"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/bulk/lapse [post]
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
