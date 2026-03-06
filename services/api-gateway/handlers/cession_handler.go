package handlers

import (
	"net/http"

	reinsuranceSchema "github.com/bitbiz/hias-core/domains/reinsurance/schema"
	"github.com/bitbiz/hias-core/domains/reinsurance/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CessionHandler struct {
	cessionSvc service.CessionService
}

func NewCessionHandler(cessionSvc service.CessionService) *CessionHandler {
	return &CessionHandler{cessionSvc: cessionSvc}
}

// CedePremium godoc
// @Summary      Cede a premium
// @Description  Create a new cession by ceding a premium to a reinsurer
// @Tags         Cessions
// @Accept       json
// @Produce      json
// @Param        request body reinsuranceSchema.CedePremiumRequest true "Cede premium request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/cessions [post]
func (h *CessionHandler) CedePremium(ctx *gin.Context) {
	var req reinsuranceSchema.CedePremiumRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.cessionSvc.CedePremium(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetCession godoc
// @Summary      Get a cession
// @Description  Retrieve a single cession by ID
// @Tags         Cessions
// @Accept       json
// @Produce      json
// @Param        id path string true "Cession ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/cessions/{id} [get]
func (h *CessionHandler) GetCession(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid cession ID")
		return
	}

	resp := h.cessionSvc.GetCession(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListCessions godoc
// @Summary      List cessions
// @Description  Retrieve a paginated list of cessions filtered by treaty or policy
// @Tags         Cessions
// @Accept       json
// @Produce      json
// @Param        treaty query string false "Treaty ID"
// @Param        policy query string false "Policy ID"
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaties/{id}/cessions [get]
func (h *CessionHandler) ListCessions(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)

	treatyIDStr := ctx.Query("treaty")
	policyIDStr := ctx.Query("policy")

	if treatyIDStr != "" {
		treatyID, err := uuid.Parse(treatyIDStr)
		if err != nil {
			utils.RespondError(ctx, http.StatusBadRequest, "Invalid treaty ID")
			return
		}
		resp := h.cessionSvc.ListByTreaty(ctx.Request.Context(), treatyID, pagination.Page, pagination.PageSize)
		if resp.Error != nil {
			utils.RespondError(ctx, resp.StatusCode, resp.Message)
			return
		}
		countResp := h.cessionSvc.GetCessionCount(ctx.Request.Context())
		utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
		return
	}

	if policyIDStr != "" {
		policyID, err := uuid.Parse(policyIDStr)
		if err != nil {
			utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
			return
		}
		resp := h.cessionSvc.ListByPolicy(ctx.Request.Context(), policyID, pagination.Page, pagination.PageSize)
		if resp.Error != nil {
			utils.RespondError(ctx, resp.StatusCode, resp.Message)
			return
		}
		countResp := h.cessionSvc.GetCessionCount(ctx.Request.Context())
		utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
		return
	}

	utils.RespondError(ctx, http.StatusBadRequest, "Query parameter 'treaty' or 'policy' is required")
}

// BookCession godoc
// @Summary      Book a cession
// @Description  Book a pending cession by ID
// @Tags         Cessions
// @Accept       json
// @Produce      json
// @Param        id path string true "Cession ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/cessions/{id}/book [put]
func (h *CessionHandler) BookCession(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid cession ID")
		return
	}

	resp := h.cessionSvc.BookCession(ctx.Request.Context(), id, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ReverseCession godoc
// @Summary      Reverse a cession
// @Description  Reverse a booked cession by ID
// @Tags         Cessions
// @Accept       json
// @Produce      json
// @Param        id path string true "Cession ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/cessions/{id}/reverse [put]
func (h *CessionHandler) ReverseCession(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid cession ID")
		return
	}

	resp := h.cessionSvc.ReverseCession(ctx.Request.Context(), id, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// AutoCede godoc
// @Summary      Auto-cede a policy premium
// @Description  Automatically cede a policy premium across applicable treaties
// @Tags         Cessions
// @Accept       json
// @Produce      json
// @Param        request body reinsuranceSchema.AutoCedePolicyPremiumRequest true "Auto-cede request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/cessions/auto-cede [post]
func (h *CessionHandler) AutoCede(ctx *gin.Context) {
	var req reinsuranceSchema.AutoCedePolicyPremiumRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.cessionSvc.AutoCedePolicyPremium(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
