package routes

import (
	"github.com/bitbiz/hias-core/services/api-gateway/rest/handlers"
	"github.com/bitbiz/hias-core/services/api-gateway/rest/middleware"
	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/gin-gonic/gin"
)

func SetupAuditRoutes(router *gin.RouterGroup, auditHandler *handlers.AuditHandler, tokenMaker auth.TokenMaker) {
	audit := router.Group("/audit")
	audit.Use(middleware.AuthMiddleware(tokenMaker))
	audit.Use(middleware.RequireRole("Admin"))
	{
		audit.GET("", auditHandler.ListAuditEvents)
		audit.GET("/entity/:type/:id", auditHandler.ListByEntity)
	}
}
