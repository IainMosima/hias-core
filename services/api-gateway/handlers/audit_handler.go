package handlers

import (
	"net/http"

	"github.com/bitbiz/hias-core/domains/audit/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuditHandler struct {
	auditSvc service.AuditService
}

func NewAuditHandler(auditSvc service.AuditService) *AuditHandler {
	return &AuditHandler{auditSvc: auditSvc}
}

// ListEvents godoc
// @Summary      List audit events
// @Description  Retrieve a paginated list of all audit events
// @Tags         Audit
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/audit [get]
func (h *AuditHandler) ListEvents(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)
	resp := h.auditSvc.ListEvents(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListByEntity godoc
// @Summary      List audit events by entity
// @Description  Retrieve a paginated list of audit events for a specific entity type and ID
// @Tags         Audit
// @Accept       json
// @Produce      json
// @Param        type path string true "Entity type"
// @Param        id path string true "Entity ID"
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/audit/entity/{type}/{id} [get]
func (h *AuditHandler) ListByEntity(ctx *gin.Context) {
	entityType := ctx.Param("type")
	entityID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid entity ID")
		return
	}

	pagination := utils.GetPaginationParams(ctx)
	resp := h.auditSvc.ListByEntity(ctx.Request.Context(), entityType, entityID, pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListByUser godoc
// @Summary      List audit events by user
// @Description  Retrieve a paginated list of audit events for a specific user
// @Tags         Audit
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/audit/user/{id} [get]
func (h *AuditHandler) ListByUser(ctx *gin.Context) {
	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid user ID")
		return
	}

	pagination := utils.GetPaginationParams(ctx)
	resp := h.auditSvc.ListByUser(ctx.Request.Context(), userID, pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
