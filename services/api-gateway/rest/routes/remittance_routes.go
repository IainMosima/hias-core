package routes

import (
	"github.com/bitbiz/hias-core/services/api-gateway/rest/handlers"
	"github.com/bitbiz/hias-core/services/api-gateway/rest/middleware"
	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/gin-gonic/gin"
)

func SetupRemittanceRoutes(router *gin.RouterGroup, remittanceHandler *handlers.RemittanceHandler, tokenMaker auth.TokenMaker) {
	remittances := router.Group("/remittances")
	remittances.Use(middleware.AuthMiddleware(tokenMaker))
	remittances.Use(middleware.RequireRole("Admin", "Finance"))
	{
		remittances.POST("", remittanceHandler.CreateRemittance)
		remittances.GET("", remittanceHandler.ListRemittances)
		remittances.GET("/:id", remittanceHandler.GetRemittance)
	}
}
