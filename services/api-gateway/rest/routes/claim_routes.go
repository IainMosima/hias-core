package routes

import (
	"github.com/bitbiz/hias-core/services/api-gateway/rest/handlers"
	"github.com/bitbiz/hias-core/services/api-gateway/rest/middleware"
	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/gin-gonic/gin"
)

func SetupClaimRoutes(router *gin.RouterGroup, claimHandler *handlers.ClaimHandler, tokenMaker auth.TokenMaker) {
	claims := router.Group("/claims")
	claims.Use(middleware.AuthMiddleware(tokenMaker))
	{
		// Provider can submit claims; ClaimsOfficer and Admin can view/review
		claims.POST("", middleware.RequireRole("Admin", "Provider"), claimHandler.SubmitClaim)
		claims.GET("", middleware.RequireRole("Admin", "ClaimsOfficer", "Provider"), claimHandler.ListClaims)
		claims.GET("/:id", middleware.RequireRole("Admin", "ClaimsOfficer", "Provider"), claimHandler.GetClaim)

		// Only ClaimsOfficer and Admin can approve/reject
		claims.PUT("/:id/approve", middleware.RequireRole("Admin", "ClaimsOfficer"), claimHandler.ApproveClaim)
		claims.PUT("/:id/reject", middleware.RequireRole("Admin", "ClaimsOfficer"), claimHandler.RejectClaim)
	}
}
