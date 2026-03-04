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

func (h *AuditHandler) ListEvents(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)
	resp := h.auditSvc.ListEvents(ctx.Request.Context(), pagination.Page, pagination.PageSize)
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
	resp := h.auditSvc.ListByEntity(ctx.Request.Context(), entityType, entityID, pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

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
