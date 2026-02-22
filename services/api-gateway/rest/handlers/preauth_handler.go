package handlers

import (
	"net/http"

	preauthSchema "github.com/bitbiz/hias-core/domains/preauth/schema"
	"github.com/bitbiz/hias-core/domains/preauth/service"
	"github.com/bitbiz/hias-core/services/api-gateway/rest/middleware"
	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PreAuthHandler struct {
	preauthService service.PreAuthService
}

func NewPreAuthHandler(preauthService service.PreAuthService) *PreAuthHandler {
	return &PreAuthHandler{preauthService: preauthService}
}

func (h *PreAuthHandler) SubmitPreAuth(ctx *gin.Context) {
	var req preauthSchema.SubmitPreAuthRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	payload := ctx.MustGet(middleware.AuthPayloadKey).(*auth.Payload)
	createdBy, _ := uuid.Parse(payload.UserID)

	resp := h.preauthService.SubmitPreAuth(ctx.Request.Context(), req, createdBy)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *PreAuthHandler) GetPreAuth(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid pre-auth ID")
		return
	}

	resp := h.preauthService.GetPreAuth(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *PreAuthHandler) ListPreAuths(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)

	resp := h.preauthService.ListPreAuths(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	countResp := h.preauthService.GetTotalCount(ctx.Request.Context())
	totalCount := int64(0)
	if countResp.Error == nil {
		totalCount = countResp.Data
	}

	utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, totalCount)
}

func (h *PreAuthHandler) ApprovePreAuth(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid pre-auth ID")
		return
	}

	payload := ctx.MustGet(middleware.AuthPayloadKey).(*auth.Payload)
	reviewedBy, _ := uuid.Parse(payload.UserID)

	resp := h.preauthService.ApprovePreAuth(ctx.Request.Context(), id, reviewedBy)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *PreAuthHandler) DenyPreAuth(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid pre-auth ID")
		return
	}

	var req preauthSchema.ReviewPreAuthRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	payload := ctx.MustGet(middleware.AuthPayloadKey).(*auth.Payload)
	reviewedBy, _ := uuid.Parse(payload.UserID)

	resp := h.preauthService.DenyPreAuth(ctx.Request.Context(), id, req.DenialReason, reviewedBy)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
