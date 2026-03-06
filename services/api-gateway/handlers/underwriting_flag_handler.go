package handlers

import (
	"net/http"
	"strconv"

	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/bitbiz/hias-core/domains/policy/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UnderwritingFlagHandler struct {
	flagSvc service.UnderwritingFlagService
}

func NewUnderwritingFlagHandler(flagSvc service.UnderwritingFlagService) *UnderwritingFlagHandler {
	return &UnderwritingFlagHandler{flagSvc: flagSvc}
}

// ListByPolicy godoc
// @Summary      List underwriting flags by policy
// @Description  List all underwriting flags associated with a policy
// @Tags         UnderwritingFlags
// @Accept       json
// @Produce      json
// @Param        id path string true "Policy ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/underwriting-flags [get]
func (h *UnderwritingFlagHandler) ListByPolicy(ctx *gin.Context) {
	policyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}
	resp := h.flagSvc.ListByPolicy(ctx.Request.Context(), policyID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListByMember godoc
// @Summary      List underwriting flags by member
// @Description  List all underwriting flags associated with a member
// @Tags         UnderwritingFlags
// @Accept       json
// @Produce      json
// @Param        id path string true "Member ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/members/{id}/underwriting-flags [get]
func (h *UnderwritingFlagHandler) ListByMember(ctx *gin.Context) {
	memberID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid member ID")
		return
	}
	resp := h.flagSvc.ListByMember(ctx.Request.Context(), memberID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetFlag godoc
// @Summary      Get an underwriting flag
// @Description  Get underwriting flag details by ID
// @Tags         UnderwritingFlags
// @Accept       json
// @Produce      json
// @Param        id path string true "Flag ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/underwriting-flags/{id} [get]
func (h *UnderwritingFlagHandler) GetFlag(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid flag ID")
		return
	}
	resp := h.flagSvc.GetFlag(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListOpen godoc
// @Summary      List open underwriting flags
// @Description  List all open/unresolved underwriting flags with pagination
// @Tags         UnderwritingFlags
// @Accept       json
// @Produce      json
// @Param        limit query int false "Limit" default(50)
// @Param        offset query int false "Offset" default(0)
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/underwriting-flags [get]
func (h *UnderwritingFlagHandler) ListOpen(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	resp := h.flagSvc.ListOpen(ctx.Request.Context(), int32(limit), int32(offset))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// CountOpen godoc
// @Summary      Count open underwriting flags
// @Description  Get the count of all open/unresolved underwriting flags
// @Tags         UnderwritingFlags
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/underwriting-flags/count [get]
func (h *UnderwritingFlagHandler) CountOpen(ctx *gin.Context) {
	resp := h.flagSvc.CountOpen(ctx.Request.Context())
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ResolveFlag godoc
// @Summary      Resolve an underwriting flag
// @Description  Resolve an open underwriting flag with resolution details
// @Tags         UnderwritingFlags
// @Accept       json
// @Produce      json
// @Param        id path string true "Flag ID"
// @Param        request body policySchema.ResolveFlagRequest true "Resolution details"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/underwriting-flags/{id}/resolve [put]
func (h *UnderwritingFlagHandler) ResolveFlag(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid flag ID")
		return
	}
	var req policySchema.ResolveFlagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp := h.flagSvc.ResolveFlag(ctx.Request.Context(), id, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// OverrideFlag godoc
// @Summary      Override an underwriting flag
// @Description  Override an underwriting flag with justification
// @Tags         UnderwritingFlags
// @Accept       json
// @Produce      json
// @Param        id path string true "Flag ID"
// @Param        request body policySchema.OverrideFlagRequest true "Override details"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/underwriting-flags/{id}/override [put]
func (h *UnderwritingFlagHandler) OverrideFlag(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid flag ID")
		return
	}
	var req policySchema.OverrideFlagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp := h.flagSvc.OverrideFlag(ctx.Request.Context(), id, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
