package routes

import (
	"github.com/bitbiz/hias-core/services/api-gateway/rest/handlers"
	"github.com/bitbiz/hias-core/services/api-gateway/rest/middleware"
	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/gin-gonic/gin"
)

func SetupNotificationRoutes(router *gin.RouterGroup, notificationHandler *handlers.NotificationHandler, tokenMaker auth.TokenMaker) {
	notifications := router.Group("/notifications")
	notifications.Use(middleware.AuthMiddleware(tokenMaker))
	{
		notifications.GET("", notificationHandler.ListNotifications)
		notifications.PUT("/:id/read", notificationHandler.MarkRead)
	}
}
