package routes

import (
	"github.com/bitbiz/hias-core/services/api-gateway/rest/handlers"
	"github.com/bitbiz/hias-core/services/api-gateway/rest/middleware"
	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/gin-gonic/gin"
)

func SetupAnalyticsRoutes(router *gin.RouterGroup, analyticsHandler *handlers.AnalyticsHandler, tokenMaker auth.TokenMaker) {
	analytics := router.Group("/analytics")
	analytics.Use(middleware.AuthMiddleware(tokenMaker))
	analytics.Use(middleware.RequireRole("Admin", "Finance"))
	{
		analytics.GET("/dashboard", analyticsHandler.GetDashboard)
		analytics.GET("/kpis", analyticsHandler.GetKPIs)
	}
}
