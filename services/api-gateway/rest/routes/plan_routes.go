package routes

import (
	"github.com/bitbiz/hias-core/services/api-gateway/rest/handlers"
	"github.com/bitbiz/hias-core/services/api-gateway/rest/middleware"
	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/gin-gonic/gin"
)

func SetupPlanRoutes(router *gin.RouterGroup, planHandler *handlers.PlanHandler, tokenMaker auth.TokenMaker) {
	plans := router.Group("/plans")
	plans.Use(middleware.AuthMiddleware(tokenMaker))
	plans.Use(middleware.RequireRole("Admin"))
	{
		plans.POST("", planHandler.CreatePlan)
		plans.GET("", planHandler.ListPlans)
		plans.GET("/:id", planHandler.GetPlan)
		plans.PUT("/:id", planHandler.UpdatePlan)

		// Benefits
		plans.POST("/:id/benefits", planHandler.AddBenefit)
		plans.GET("/:id/benefits", planHandler.ListBenefits)
	}
}
