package routes

import (
	"github.com/bitbiz/hias-core/services/api-gateway/rest/handlers"
	"github.com/bitbiz/hias-core/services/api-gateway/rest/middleware"
	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/gin-gonic/gin"
)

func SetupPreAuthRoutes(router *gin.RouterGroup, preAuthHandler *handlers.PreAuthHandler, tokenMaker auth.TokenMaker) {
	preauths := router.Group("/preauths")
	preauths.Use(middleware.AuthMiddleware(tokenMaker))
	{
		// Provider can submit pre-auths; ClaimsOfficer and Admin can view/review
		preauths.POST("", middleware.RequireRole("Admin", "Provider"), preAuthHandler.SubmitPreAuth)
		preauths.GET("", middleware.RequireRole("Admin", "ClaimsOfficer", "Provider"), preAuthHandler.ListPreAuths)
		preauths.GET("/:id", middleware.RequireRole("Admin", "ClaimsOfficer", "Provider"), preAuthHandler.GetPreAuth)

		// Only ClaimsOfficer and Admin can approve/deny
		preauths.PUT("/:id/approve", middleware.RequireRole("Admin", "ClaimsOfficer"), preAuthHandler.ApprovePreAuth)
		preauths.PUT("/:id/deny", middleware.RequireRole("Admin", "ClaimsOfficer"), preAuthHandler.DenyPreAuth)
	}
}
