package handlers

import (
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

type CalculatePremiumRequest struct {
	MemberCount   int      `json:"member_count" binding:"required,min=1"`
	Relationships []string `json:"relationships" binding:"required"`
}

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
