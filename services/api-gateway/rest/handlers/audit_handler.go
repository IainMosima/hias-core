package handlers

import (
	"net/http"

	"github.com/bitbiz/hias-core/domains/audit/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuditHandler struct {
	auditService service.AuditService
}

func NewAuditHandler(auditService service.AuditService) *AuditHandler {
	return &AuditHandler{auditService: auditService}
}

func (h *AuditHandler) ListAuditEvents(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)

	resp := h.auditService.ListEvents(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *AuditHandler) ListByEntity(ctx *gin.Context) {
	entityType := ctx.Param("type")
	entityID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid entity ID")
		return
	}

	pagination := utils.GetPaginationParams(ctx)

	resp := h.auditService.ListByEntity(ctx.Request.Context(), entityType, entityID, pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
