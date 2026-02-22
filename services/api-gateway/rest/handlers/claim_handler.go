package handlers

import (
	"net/http"

	claimsSchema "github.com/bitbiz/hias-core/domains/claims/schema"
	"github.com/bitbiz/hias-core/domains/claims/service"
	"github.com/bitbiz/hias-core/services/api-gateway/rest/middleware"
	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ClaimHandler struct {
	claimService service.ClaimService
}

func NewClaimHandler(claimService service.ClaimService) *ClaimHandler {
	return &ClaimHandler{claimService: claimService}
}

func (h *ClaimHandler) SubmitClaim(ctx *gin.Context) {
	var req claimsSchema.SubmitClaimRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	payload := ctx.MustGet(middleware.AuthPayloadKey).(*auth.Payload)
	createdBy, _ := uuid.Parse(payload.UserID)

	resp := h.claimService.SubmitClaim(ctx.Request.Context(), req, createdBy)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *ClaimHandler) GetClaim(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid claim ID")
		return
	}

	resp := h.claimService.GetClaim(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *ClaimHandler) ListClaims(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)

	resp := h.claimService.ListClaims(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	countResp := h.claimService.GetTotalCount(ctx.Request.Context())
	totalCount := int64(0)
	if countResp.Error == nil {
		totalCount = countResp.Data
	}

	utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, totalCount)
}

func (h *ClaimHandler) ApproveClaim(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid claim ID")
		return
	}

	payload := ctx.MustGet(middleware.AuthPayloadKey).(*auth.Payload)
	approvedBy, _ := uuid.Parse(payload.UserID)

	resp := h.claimService.ApproveClaim(ctx.Request.Context(), id, approvedBy)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *ClaimHandler) RejectClaim(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid claim ID")
		return
	}

	var req claimsSchema.ReviewClaimRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	payload := ctx.MustGet(middleware.AuthPayloadKey).(*auth.Payload)
	rejectedBy, _ := uuid.Parse(payload.UserID)

	resp := h.claimService.RejectClaim(ctx.Request.Context(), id, req.Reason, rejectedBy)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
