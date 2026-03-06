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

// ListNotifications godoc
// @Summary      List notifications
// @Description  Retrieve a paginated list of notifications for the authenticated user
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/notifications [get]
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

// MarkRead godoc
// @Summary      Mark notification as read
// @Description  Mark a specific notification as read by its ID
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Param        id path string true "Notification ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/notifications/{id}/read [put]
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

// GetUnreadCount godoc
// @Summary      Get unread notification count
// @Description  Retrieve the count of unread notifications for the authenticated user
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/notifications/unread-count [get]
func (h *NotificationHandler) GetUnreadCount(ctx *gin.Context) {
	userID := getUserID(ctx)

	resp := h.notifSvc.GetUnreadCount(ctx.Request.Context(), userID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
