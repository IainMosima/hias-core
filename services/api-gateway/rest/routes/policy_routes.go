package routes

import (
	"github.com/bitbiz/hias-core/services/api-gateway/rest/handlers"
	"github.com/bitbiz/hias-core/services/api-gateway/rest/middleware"
	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/gin-gonic/gin"
)

func SetupPolicyRoutes(router *gin.RouterGroup, policyHandler *handlers.PolicyHandler, tokenMaker auth.TokenMaker) {
	policies := router.Group("/policies")
	policies.Use(middleware.AuthMiddleware(tokenMaker))
	policies.Use(middleware.RequireRole("Admin", "Underwriter"))
	{
		policies.POST("", policyHandler.CreatePolicy)
		policies.GET("", policyHandler.ListPolicies)
		policies.GET("/:id", policyHandler.GetPolicy)
		policies.PUT("/:id/activate", policyHandler.ActivatePolicy)
		policies.PUT("/:id/lapse", policyHandler.LapsePolicy)
		policies.PUT("/:id/terminate", policyHandler.TerminatePolicy)
	}
}
