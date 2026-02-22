package routes

import (
	"github.com/bitbiz/hias-core/services/api-gateway/rest/handlers"
	"github.com/bitbiz/hias-core/services/api-gateway/rest/middleware"
	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/gin-gonic/gin"
)

func SetupBillingRoutes(router *gin.RouterGroup, billingHandler *handlers.BillingHandler, tokenMaker auth.TokenMaker) {
	// Public webhook endpoint — no auth middleware
	router.POST("/payments/webhook", billingHandler.ProcessPaymentWebhook)

	// Invoice routes — Finance and Admin
	invoices := router.Group("/invoices")
	invoices.Use(middleware.AuthMiddleware(tokenMaker))
	invoices.Use(middleware.RequireRole("Admin", "Finance"))
	{
		invoices.POST("", billingHandler.CreateInvoice)
		invoices.GET("", billingHandler.ListInvoices)
		invoices.GET("/:id", billingHandler.GetInvoice)
	}

	// Payment routes — Finance and Admin
	payments := router.Group("/payments")
	payments.Use(middleware.AuthMiddleware(tokenMaker))
	payments.Use(middleware.RequireRole("Admin", "Finance"))
	{
		payments.POST("", billingHandler.InitiatePayment)
		payments.GET("", billingHandler.ListPayments)
		payments.GET("/:id", billingHandler.GetPayment)
	}
}
