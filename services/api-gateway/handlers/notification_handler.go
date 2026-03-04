package handlers

import (
	"net/http"

	"github.com/bitbiz/hias-core/domains/notification/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	notifSvc service.NotificationService
}

func NewNotificationHandler(notifSvc service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notifSvc: notifSvc}
}

func (h *NotificationHandler) ListNotifications(ctx *gin.Context) {
	userID := getUserID(ctx)
	pagination := utils.GetPaginationParams(ctx)

	resp := h.notifSvc.ListByUser(ctx.Request.Context(), userID, pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *NotificationHandler) MarkRead(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	resp := h.notifSvc.MarkRead(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *NotificationHandler) GetUnreadCount(ctx *gin.Context) {
	userID := getUserID(ctx)

	resp := h.notifSvc.GetUnreadCount(ctx.Request.Context(), userID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
