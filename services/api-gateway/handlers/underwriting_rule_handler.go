package handlers

import (
	"net/http"

	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/bitbiz/hias-core/domains/policy/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UnderwritingRuleHandler struct {
	ruleSvc service.UnderwritingRuleService
}

func NewUnderwritingRuleHandler(ruleSvc service.UnderwritingRuleService) *UnderwritingRuleHandler {
	return &UnderwritingRuleHandler{ruleSvc: ruleSvc}
}

// ListRules godoc
// @Summary      List underwriting rules
// @Description  List all underwriting rules for a plan
// @Tags         UnderwritingRules
// @Accept       json
// @Produce      json
// @Param        id path string true "Plan ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/plans/{id}/underwriting-rules [get]
func (h *UnderwritingRuleHandler) ListRules(ctx *gin.Context) {
	planID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid plan ID")
		return
	}
	resp := h.ruleSvc.ListByPlan(ctx.Request.Context(), planID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// CreateRule godoc
// @Summary      Create an underwriting rule
// @Description  Create a new underwriting rule for a plan
// @Tags         UnderwritingRules
// @Accept       json
// @Produce      json
// @Param        id path string true "Plan ID"
// @Param        request body policySchema.CreateUnderwritingRuleRequest true "Rule details"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/plans/{id}/underwriting-rules [post]
func (h *UnderwritingRuleHandler) CreateRule(ctx *gin.Context) {
	var req policySchema.CreateUnderwritingRuleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	req.PlanID = ctx.Param("id")
	resp := h.ruleSvc.CreateRule(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// UpdateRule godoc
// @Summary      Update an underwriting rule
// @Description  Update an existing underwriting rule by ID
// @Tags         UnderwritingRules
// @Accept       json
// @Produce      json
// @Param        id path string true "Rule ID"
// @Param        request body policySchema.UpdateUnderwritingRuleRequest true "Updated rule details"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/underwriting-rules/{id} [put]
func (h *UnderwritingRuleHandler) UpdateRule(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid rule ID")
		return
	}
	var req policySchema.UpdateUnderwritingRuleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp := h.ruleSvc.UpdateRule(ctx.Request.Context(), id, req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// DeleteRule godoc
// @Summary      Delete an underwriting rule
// @Description  Delete an underwriting rule by ID
// @Tags         UnderwritingRules
// @Accept       json
// @Produce      json
// @Param        id path string true "Rule ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/underwriting-rules/{id} [delete]
func (h *UnderwritingRuleHandler) DeleteRule(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid rule ID")
		return
	}
	resp := h.ruleSvc.DeleteRule(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
