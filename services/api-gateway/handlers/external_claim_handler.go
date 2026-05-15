package handlers

import (
	"net/http"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	claimsSchema "github.com/bitbiz/hias-core/domains/claims/schema"
	"github.com/bitbiz/hias-core/domains/claims/service"
	"github.com/bitbiz/hias-core/services/api-gateway/middleware"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ExternalClaimHandler struct {
	intakeSvc service.ClaimIntakeService
}

func NewExternalClaimHandler(intakeSvc service.ClaimIntakeService) *ExternalClaimHandler {
	return &ExternalClaimHandler{intakeSvc: intakeSvc}
}

func getAPIPartner(ctx *gin.Context) *entity.APIPartner {
	val, exists := ctx.Get(middleware.APIPartnerKey)
	if !exists {
		return nil
	}
	partner, ok := val.(*entity.APIPartner)
	if !ok {
		return nil
	}
	return partner
}

func (h *ExternalClaimHandler) SubmitExternalClaim(ctx *gin.Context) {
	partner := getAPIPartner(ctx)
	if partner == nil {
		utils.RespondError(ctx, http.StatusUnauthorized, "API partner not found in context")
		return
	}

	var req claimsSchema.ExternalClaimRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.intakeSvc.SubmitExternal(ctx.Request.Context(), req, partner)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *ExternalClaimHandler) GetExternalClaimStatus(ctx *gin.Context) {
	partner := getAPIPartner(ctx)
	if partner == nil {
		utils.RespondError(ctx, http.StatusUnauthorized, "API partner not found in context")
		return
	}

	claimID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid claim ID")
		return
	}

	resp := h.intakeSvc.GetExternalStatus(ctx.Request.Context(), claimID, partner.ID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
