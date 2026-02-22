package routes

import (
	"github.com/bitbiz/hias-core/services/api-gateway/rest/handlers"
	"github.com/bitbiz/hias-core/services/api-gateway/rest/middleware"
	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/gin-gonic/gin"
)

func SetupProviderRoutes(router *gin.RouterGroup, providerHandler *handlers.ProviderHandler, tokenMaker auth.TokenMaker) {
	providers := router.Group("/providers")
	providers.Use(middleware.AuthMiddleware(tokenMaker))
	providers.Use(middleware.RequireRole("Admin", "Provider"))
	{
		providers.POST("", providerHandler.RegisterProvider)
		providers.GET("", providerHandler.ListProviders)
		providers.GET("/:id", providerHandler.GetProvider)
		providers.PUT("/:id/credential", providerHandler.CredentialProvider)
		providers.PUT("/:id/activate", providerHandler.ActivateProvider)
		providers.PUT("/:id/suspend", providerHandler.SuspendProvider)
	}
}
