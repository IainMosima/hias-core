package routes

import (
	"github.com/bitbiz/hias-core/services/api-gateway/rest/handlers"
	"github.com/bitbiz/hias-core/services/api-gateway/rest/middleware"
	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/gin-gonic/gin"
)

func SetupMemberRoutes(router *gin.RouterGroup, memberHandler *handlers.MemberHandler, tokenMaker auth.TokenMaker) {
	// Member routes nested under policies
	policyMembers := router.Group("/policies")
	policyMembers.Use(middleware.AuthMiddleware(tokenMaker))
	policyMembers.Use(middleware.RequireRole("Admin", "Underwriter"))
	{
		policyMembers.POST("/:policyId/members", memberHandler.EnrollMember)
		policyMembers.GET("/:policyId/members", memberHandler.ListMembers)
	}

	// Standalone member routes
	members := router.Group("/members")
	members.Use(middleware.AuthMiddleware(tokenMaker))
	members.Use(middleware.RequireRole("Admin", "Underwriter"))
	{
		members.GET("/:id", memberHandler.GetMember)
		members.PUT("/:id/verify", memberHandler.VerifyMember)
	}
}
