package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bitbiz/hias-core/domains/product/schema"
	"github.com/bitbiz/hias-core/domains/product/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PremiumRuleHandler struct {
	premiumRuleSvc service.PremiumRuleService
}

func NewPremiumRuleHandler(premiumRuleSvc service.PremiumRuleService) *PremiumRuleHandler {
	return &PremiumRuleHandler{premiumRuleSvc: premiumRuleSvc}
}

// CreatePremiumRule godoc
// @Summary      Create a premium rule for a plan
// @Description  Create a new premium rule associated with the specified plan
// @Tags         PremiumRules
// @Accept       json
// @Produce      json
// @Param        id path string true "Plan ID"
// @Param        request body schema.CreatePremiumRuleRequest true "Premium rule creation payload"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/plans/{id}/premium-rules [post]
func (h *PremiumRuleHandler) CreatePremiumRule(ctx *gin.Context) {
	planID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid plan ID")
		return
	}

	var req schema.CreatePremiumRuleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.premiumRuleSvc.CreatePremiumRule(ctx.Request.Context(), planID, req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListPremiumRules godoc
// @Summary      List premium rules for a plan
// @Description  Retrieve all premium rules associated with the specified plan
// @Tags         PremiumRules
// @Accept       json
// @Produce      json
// @Param        id path string true "Plan ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/plans/{id}/premium-rules [get]
func (h *PremiumRuleHandler) ListPremiumRules(ctx *gin.Context) {
	planID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid plan ID")
		return
	}

	resp := h.premiumRuleSvc.ListPremiumRulesByPlan(ctx.Request.Context(), planID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// DeletePremiumRule godoc
// @Summary      Delete a premium rule
// @Description  Delete a premium rule by its unique identifier
// @Tags         PremiumRules
// @Accept       json
// @Produce      json
// @Param        id path string true "Premium Rule ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/premium-rules/{id} [delete]
func (h *PremiumRuleHandler) DeletePremiumRule(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid premium rule ID")
		return
	}

	resp := h.premiumRuleSvc.DeletePremiumRule(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetPremiumRule godoc
// @Summary      Get a premium rule by ID
// @Description  Retrieve a single premium rule by its unique identifier
// @Tags         PremiumRules
// @Accept       json
// @Produce      json
// @Param        id path string true "Premium Rule ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/premium-rules/{id} [get]
func (h *PremiumRuleHandler) GetPremiumRule(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid premium rule ID")
		return
	}

	resp := h.premiumRuleSvc.GetPremiumRule(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// UpdatePremiumRule godoc
// @Summary      Update a premium rule
// @Description  Update an existing premium rule by its unique identifier
// @Tags         PremiumRules
// @Accept       json
// @Produce      json
// @Param        id path string true "Premium Rule ID"
// @Param        request body schema.UpdatePremiumRuleRequest true "Premium rule update payload"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/premium-rules/{id} [put]
func (h *PremiumRuleHandler) UpdatePremiumRule(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid premium rule ID")
		return
	}

	var req schema.UpdatePremiumRuleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.premiumRuleSvc.UpdatePremiumRule(ctx.Request.Context(), id, req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetRateSheet godoc
// @Summary      Get rate sheet for a plan
// @Description  Retrieve the rate sheet associated with the specified plan
// @Tags         PremiumRules
// @Accept       json
// @Produce      json
// @Param        id path string true "Plan ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/plans/{id}/rate-sheet [get]
func (h *PremiumRuleHandler) GetRateSheet(ctx *gin.Context) {
	planID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid plan ID")
		return
	}

	resp := h.premiumRuleSvc.GetRateSheet(ctx.Request.Context(), planID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

type CalculatePremiumRequest struct {
	MemberCount   int      `json:"member_count" binding:"required,min=1"`
	Relationships []string `json:"relationships" binding:"required"`
}

// CalculatePremium godoc
// @Summary      Calculate premium for a plan
// @Description  Calculate the premium for a plan based on member count and relationships
// @Tags         PremiumRules
// @Accept       json
// @Produce      json
// @Param        id path string true "Plan ID"
// @Param        request body CalculatePremiumRequest true "Premium calculation payload"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/plans/{id}/calculate-premium [post]
func (h *PremiumRuleHandler) CalculatePremium(ctx *gin.Context) {
	planID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid plan ID")
		return
	}

	var req CalculatePremiumRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.premiumRuleSvc.CalculatePremium(ctx.Request.Context(), planID, req.MemberCount, req.Relationships)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

type PremiumBreakdownRequest struct {
	ProposedMembers json.RawMessage `json:"proposed_members" binding:"required"`
}

func (h *PremiumRuleHandler) CalculatePremiumBreakdown(ctx *gin.Context) {
	planID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid plan ID")
		return
	}

	var req PremiumBreakdownRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.premiumRuleSvc.CalculatePremiumBreakdown(ctx.Request.Context(), planID, req.ProposedMembers)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
