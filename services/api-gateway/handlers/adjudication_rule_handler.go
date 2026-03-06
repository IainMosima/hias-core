package handlers

import (
	"net/http"

	claimsSchema "github.com/bitbiz/hias-core/domains/claims/schema"
	"github.com/bitbiz/hias-core/domains/claims/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdjudicationRuleHandler struct {
	svc service.AdjudicationRuleService
}

func NewAdjudicationRuleHandler(svc service.AdjudicationRuleService) *AdjudicationRuleHandler {
	return &AdjudicationRuleHandler{svc: svc}
}

// CreateRule godoc
// @Summary      Create an adjudication rule
// @Description  Create a new adjudication rule for claims processing
// @Tags         AdjudicationRules
// @Accept       json
// @Produce      json
// @Param        request body claimsSchema.CreateAdjudicationRuleRequest true "Create adjudication rule request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/adjudication-rules [post]
func (h *AdjudicationRuleHandler) CreateRule(ctx *gin.Context) {
	var req claimsSchema.CreateAdjudicationRuleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp := h.svc.CreateRule(ctx.Request.Context(), req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetRule godoc
// @Summary      Get an adjudication rule
// @Description  Retrieve a single adjudication rule by ID
// @Tags         AdjudicationRules
// @Accept       json
// @Produce      json
// @Param        id path string true "Adjudication Rule ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/adjudication-rules/{id} [get]
func (h *AdjudicationRuleHandler) GetRule(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid rule ID")
		return
	}
	resp := h.svc.GetRule(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListRules godoc
// @Summary      List adjudication rules
// @Description  Retrieve a paginated list of adjudication rules
// @Tags         AdjudicationRules
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/adjudication-rules [get]
func (h *AdjudicationRuleHandler) ListRules(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)
	resp := h.svc.ListRules(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// UpdateRule godoc
// @Summary      Update an adjudication rule
// @Description  Update an existing adjudication rule by ID
// @Tags         AdjudicationRules
// @Accept       json
// @Produce      json
// @Param        id path string true "Adjudication Rule ID"
// @Param        request body claimsSchema.UpdateAdjudicationRuleRequest true "Update adjudication rule request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/adjudication-rules/{id} [put]
func (h *AdjudicationRuleHandler) UpdateRule(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid rule ID")
		return
	}
	var req claimsSchema.UpdateAdjudicationRuleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp := h.svc.UpdateRule(ctx.Request.Context(), id, req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// DeleteRule godoc
// @Summary      Delete an adjudication rule
// @Description  Delete an adjudication rule by ID
// @Tags         AdjudicationRules
// @Accept       json
// @Produce      json
// @Param        id path string true "Adjudication Rule ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/adjudication-rules/{id} [delete]
func (h *AdjudicationRuleHandler) DeleteRule(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid rule ID")
		return
	}
	resp := h.svc.DeleteRule(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
