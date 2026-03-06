package handlers

import (
	"net/http"

	preauthSchema "github.com/bitbiz/hias-core/domains/preauth/schema"
	"github.com/bitbiz/hias-core/domains/preauth/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PreAuthHandler struct {
	preAuthSvc service.PreAuthService
}

func NewPreAuthHandler(preAuthSvc service.PreAuthService) *PreAuthHandler {
	return &PreAuthHandler{preAuthSvc: preAuthSvc}
}

// SubmitPreAuth godoc
// @Summary      Submit a pre-authorization request
// @Description  Submits a new pre-authorization request for review
// @Tags         PreAuthorizations
// @Accept       json
// @Produce      json
// @Param        request body preauthSchema.SubmitPreAuthRequest true "Submit pre-auth request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/preauths [post]
func (h *PreAuthHandler) SubmitPreAuth(ctx *gin.Context) {
	var req preauthSchema.SubmitPreAuthRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.preAuthSvc.SubmitPreAuth(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetPreAuth godoc
// @Summary      Get a pre-authorization by ID
// @Description  Retrieves the details of a specific pre-authorization request
// @Tags         PreAuthorizations
// @Accept       json
// @Produce      json
// @Param        id path string true "Pre-Auth ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/preauths/{id} [get]
func (h *PreAuthHandler) GetPreAuth(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid pre-auth ID")
		return
	}

	resp := h.preAuthSvc.GetPreAuth(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListPreAuths godoc
// @Summary      List all pre-authorizations
// @Description  Returns a paginated list of all pre-authorization requests
// @Tags         PreAuthorizations
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/preauths [get]
func (h *PreAuthHandler) ListPreAuths(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)
	resp := h.preAuthSvc.ListPreAuths(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	countResp := h.preAuthSvc.GetTotalCount(ctx.Request.Context())
	utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
}

// ReviewPreAuth godoc
// @Summary      Review a pre-authorization
// @Description  Submits a review for the specified pre-authorization request
// @Tags         PreAuthorizations
// @Accept       json
// @Produce      json
// @Param        id path string true "Pre-Auth ID"
// @Param        request body preauthSchema.ReviewPreAuthRequest true "Review pre-auth request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/preauths/{id}/review [put]
func (h *PreAuthHandler) ReviewPreAuth(ctx *gin.Context) {
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

	resp := h.preAuthSvc.ReviewPreAuth(ctx.Request.Context(), id, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ApprovePreAuth godoc
// @Summary      Approve a pre-authorization
// @Description  Approves the specified pre-authorization request
// @Tags         PreAuthorizations
// @Accept       json
// @Produce      json
// @Param        id path string true "Pre-Auth ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/preauths/{id}/approve [put]
func (h *PreAuthHandler) ApprovePreAuth(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid pre-auth ID")
		return
	}

	resp := h.preAuthSvc.ApprovePreAuth(ctx.Request.Context(), id, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// DenyPreAuth godoc
// @Summary      Deny a pre-authorization
// @Description  Denies the specified pre-authorization request with a reason
// @Tags         PreAuthorizations
// @Accept       json
// @Produce      json
// @Param        id path string true "Pre-Auth ID"
// @Param        request body object true "Deny request with reason field"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/preauths/{id}/deny [put]
func (h *PreAuthHandler) DenyPreAuth(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid pre-auth ID")
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.preAuthSvc.DenyPreAuth(ctx.Request.Context(), id, req.Reason, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
