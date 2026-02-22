package handlers

import (
	"net/http"

	"github.com/bitbiz/hias-core/domains/policy/service"
	"github.com/bitbiz/hias-core/services/api-gateway/rest/middleware"
	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
)

type PolicyHandler struct {
	policyService service.PolicyService
}

func NewPolicyHandler(policyService service.PolicyService) *PolicyHandler {
	return &PolicyHandler{policyService: policyService}
}

func (h *PolicyHandler) CreatePolicy(ctx *gin.Context) {
	var req policySchema.CreatePolicyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	payload := ctx.MustGet(middleware.AuthPayloadKey).(*auth.Payload)
	createdBy, _ := uuid.Parse(payload.UserID)

	resp := h.policyService.CreatePolicy(ctx.Request.Context(), req, createdBy)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *PolicyHandler) GetPolicy(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	resp := h.policyService.GetPolicy(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *PolicyHandler) ListPolicies(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)

	resp := h.policyService.ListPolicies(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	countResp := h.policyService.GetTotalCount(ctx.Request.Context())
	totalCount := int64(0)
	if countResp.Error == nil {
		totalCount = countResp.Data
	}

	utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, totalCount)
}

func (h *PolicyHandler) ActivatePolicy(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	resp := h.policyService.ActivatePolicy(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *PolicyHandler) LapsePolicy(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	resp := h.policyService.LapsePolicy(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *PolicyHandler) TerminatePolicy(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	resp := h.policyService.TerminatePolicy(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
