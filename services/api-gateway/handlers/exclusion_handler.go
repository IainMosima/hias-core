package handlers

import (
	"net/http"

	"github.com/bitbiz/hias-core/domains/product/schema"
	"github.com/bitbiz/hias-core/domains/product/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ExclusionHandler struct {
	exclusionSvc service.ExclusionService
}

func NewExclusionHandler(exclusionSvc service.ExclusionService) *ExclusionHandler {
	return &ExclusionHandler{exclusionSvc: exclusionSvc}
}

// CreateExclusion godoc
// @Summary      Create an exclusion for a plan
// @Description  Create a new exclusion associated with the specified plan
// @Tags         Exclusions
// @Accept       json
// @Produce      json
// @Param        id path string true "Plan ID"
// @Param        request body schema.CreateExclusionRequest true "Exclusion creation payload"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/plans/{id}/exclusions [post]
func (h *ExclusionHandler) CreateExclusion(ctx *gin.Context) {
	planID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid plan ID")
		return
	}

	var req schema.CreateExclusionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.exclusionSvc.CreateExclusion(ctx.Request.Context(), planID, req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListExclusions godoc
// @Summary      List exclusions for a plan
// @Description  Retrieve all exclusions associated with the specified plan
// @Tags         Exclusions
// @Accept       json
// @Produce      json
// @Param        id path string true "Plan ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/plans/{id}/exclusions [get]
func (h *ExclusionHandler) ListExclusions(ctx *gin.Context) {
	planID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid plan ID")
		return
	}

	resp := h.exclusionSvc.ListExclusionsByPlan(ctx.Request.Context(), planID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// UpdateExclusion godoc
// @Summary      Update an exclusion
// @Description  Update an existing exclusion by its unique identifier
// @Tags         Exclusions
// @Accept       json
// @Produce      json
// @Param        id path string true "Exclusion ID"
// @Param        request body schema.UpdateExclusionRequest true "Exclusion update payload"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/exclusions/{id} [put]
func (h *ExclusionHandler) UpdateExclusion(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid exclusion ID")
		return
	}

	var req schema.UpdateExclusionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.exclusionSvc.UpdateExclusion(ctx.Request.Context(), id, req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// DeleteExclusion godoc
// @Summary      Delete an exclusion
// @Description  Delete an exclusion by its unique identifier
// @Tags         Exclusions
// @Accept       json
// @Produce      json
// @Param        id path string true "Exclusion ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/exclusions/{id} [delete]
func (h *ExclusionHandler) DeleteExclusion(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid exclusion ID")
		return
	}

	resp := h.exclusionSvc.DeleteExclusion(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
