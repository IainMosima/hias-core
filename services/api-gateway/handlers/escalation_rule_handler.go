package handlers

import (
	"net/http"

	claimsSchema "github.com/bitbiz/hias-core/domains/claims/schema"
	"github.com/bitbiz/hias-core/domains/claims/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EscalationRuleHandler struct {
	svc service.EscalationRuleService
}

func NewEscalationRuleHandler(svc service.EscalationRuleService) *EscalationRuleHandler {
	return &EscalationRuleHandler{svc: svc}
}

// CreateRule godoc
// @Summary      Create an escalation rule
// @Description  Create a new escalation rule for claims processing
// @Tags         EscalationRules
// @Accept       json
// @Produce      json
// @Param        request body claimsSchema.CreateEscalationRuleRequest true "Create escalation rule request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/escalation-rules [post]
func (h *EscalationRuleHandler) CreateRule(ctx *gin.Context) {
	var req claimsSchema.CreateEscalationRuleRequest
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
// @Summary      Get an escalation rule
// @Description  Retrieve a single escalation rule by ID
// @Tags         EscalationRules
// @Accept       json
// @Produce      json
// @Param        id path string true "Escalation Rule ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/escalation-rules/{id} [get]
func (h *EscalationRuleHandler) GetRule(ctx *gin.Context) {
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
// @Summary      List escalation rules
// @Description  Retrieve a paginated list of escalation rules
// @Tags         EscalationRules
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/escalation-rules [get]
func (h *EscalationRuleHandler) ListRules(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)
	resp := h.svc.ListRules(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// UpdateRule godoc
// @Summary      Update an escalation rule
// @Description  Update an existing escalation rule by ID
// @Tags         EscalationRules
// @Accept       json
// @Produce      json
// @Param        id path string true "Escalation Rule ID"
// @Param        request body claimsSchema.UpdateEscalationRuleRequest true "Update escalation rule request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/escalation-rules/{id} [put]
func (h *EscalationRuleHandler) UpdateRule(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid rule ID")
		return
	}
	var req claimsSchema.UpdateEscalationRuleRequest
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
// @Summary      Delete an escalation rule
// @Description  Delete an escalation rule by ID
// @Tags         EscalationRules
// @Accept       json
// @Produce      json
// @Param        id path string true "Escalation Rule ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/escalation-rules/{id} [delete]
func (h *EscalationRuleHandler) DeleteRule(ctx *gin.Context) {
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
