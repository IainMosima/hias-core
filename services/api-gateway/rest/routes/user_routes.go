package routes

import (
	"github.com/bitbiz/hias-core/services/api-gateway/rest/handlers"
	"github.com/bitbiz/hias-core/services/api-gateway/rest/middleware"
	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.RouterGroup, userHandler *handlers.UserHandler, tokenMaker auth.TokenMaker) {
	users := router.Group("/users")
	users.Use(middleware.AuthMiddleware(tokenMaker))
	users.Use(middleware.RequireRole("Admin"))
	{
		users.POST("", userHandler.CreateUser)
		users.GET("", userHandler.ListUsers)
		users.GET("/:id", userHandler.GetUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.PUT("/:id/role", userHandler.AssignRole)
		users.PUT("/:id/status", userHandler.UpdateStatus)
	}
}
