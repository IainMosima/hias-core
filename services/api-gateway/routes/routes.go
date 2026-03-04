package routes

import (
	"github.com/bitbiz/hias-core/services/api-gateway/handlers"
	"github.com/bitbiz/hias-core/services/api-gateway/middleware"
	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Health          *handlers.HealthHandler
	Auth            *handlers.AuthHandler
	User            *handlers.UserHandler
	Plan            *handlers.PlanHandler
	Benefit         *handlers.BenefitHandler
	Exclusion       *handlers.ExclusionHandler
	PremiumRule     *handlers.PremiumRuleHandler
	ProviderNetwork *handlers.ProviderNetworkHandler
	Policy          *handlers.PolicyHandler
	Member          *handlers.MemberHandler
	Provider        *handlers.ProviderHandler
	Contract        *handlers.ContractHandler
	RateCard        *handlers.RateCardHandler
	Claim           *handlers.ClaimHandler
	PreAuth         *handlers.PreAuthHandler
	Invoice         *handlers.InvoiceHandler
	Payment         *handlers.PaymentHandler
	Remittance      *handlers.RemittanceHandler
	Installment     *handlers.InstallmentHandler
	Notification    *handlers.NotificationHandler
	Audit           *handlers.AuditHandler
	Analytics       *handlers.AnalyticsHandler
}

func RegisterRoutes(router *gin.Engine, h Handlers, tokenMaker auth.TokenMaker) {
	// Public routes
	router.GET("/health", h.Health.Health)
	router.GET("/ready", h.Health.Ready)

	// Public auth routes
	authPublic := router.Group("/api/v1/auth")
	{
		authPublic.POST("/login", h.Auth.Login)
		authPublic.POST("/register", h.Auth.Register)
		authPublic.POST("/refresh", h.Auth.RefreshToken)
	}

	// Public webhook
	router.POST("/api/v1/webhooks/mpesa", h.Payment.ProcessWebhook)

	// Authenticated routes
	authenticated := router.Group("/api/v1")
	authenticated.Use(middleware.AuthMiddleware(tokenMaker))
	{
		// Auth (authenticated)
		authenticated.POST("/auth/logout", h.Auth.Logout)

		// Users
		users := authenticated.Group("/users")
		{
			users.GET("", h.User.ListUsers)
			users.GET("/:id", h.User.GetUser)
			users.POST("", h.User.CreateUser)
			users.PUT("/:id", h.User.UpdateUser)
			users.PUT("/:id/role", h.User.AssignRole)
			users.PUT("/:id/status", h.User.UpdateStatus)
		}

		// Plans
		plans := authenticated.Group("/plans")
		{
			plans.GET("", h.Plan.ListPlans)
			plans.GET("/:id", h.Plan.GetPlan)
			plans.POST("", h.Plan.CreatePlan)
			plans.PUT("/:id", h.Plan.UpdatePlan)

			// Benefits (nested under plans)
			plans.GET("/:id/benefits", h.Benefit.ListBenefits)
			plans.POST("/:id/benefits", h.Benefit.CreateBenefit)

			// Exclusions (nested under plans)
			plans.GET("/:id/exclusions", h.Exclusion.ListExclusions)
			plans.POST("/:id/exclusions", h.Exclusion.CreateExclusion)

			// Premium Rules (nested under plans)
			plans.GET("/:id/premium-rules", h.PremiumRule.ListPremiumRules)
			plans.POST("/:id/premium-rules", h.PremiumRule.CreatePremiumRule)
			plans.POST("/:id/calculate-premium", h.PremiumRule.CalculatePremium)

			// Provider Networks (nested under plans)
			plans.GET("/:id/provider-networks", h.ProviderNetwork.ListProviderNetworks)
			plans.POST("/:id/provider-networks", h.ProviderNetwork.CreateProviderNetwork)
		}

		// Exclusions (standalone for update/delete)
		exclusions := authenticated.Group("/exclusions")
		{
			exclusions.PUT("/:id", h.Exclusion.UpdateExclusion)
			exclusions.DELETE("/:id", h.Exclusion.DeleteExclusion)
		}

		// Premium Rules (standalone for delete)
		premiumRules := authenticated.Group("/premium-rules")
		{
			premiumRules.DELETE("/:id", h.PremiumRule.DeletePremiumRule)
		}

		// Provider Networks (standalone for update/delete)
		providerNetworks := authenticated.Group("/provider-networks")
		{
			providerNetworks.PUT("/:id/status", h.ProviderNetwork.UpdateProviderNetworkStatus)
			providerNetworks.DELETE("/:id", h.ProviderNetwork.DeleteProviderNetwork)
		}

		// Policies
		policies := authenticated.Group("/policies")
		{
			policies.GET("", h.Policy.ListPolicies)
			policies.GET("/:id", h.Policy.GetPolicy)
			policies.POST("", h.Policy.CreatePolicy)
			policies.PUT("/:id/activate", h.Policy.ActivatePolicy)
			policies.PUT("/:id/lapse", h.Policy.LapsePolicy)
			policies.PUT("/:id/terminate", h.Policy.TerminatePolicy)
			policies.PUT("/:id/reinstate", h.Policy.ReinstatePolicy)
			policies.GET("/:id/prorate", h.Policy.CalculateProrate)

			// Members (nested under policies)
			policies.GET("/:id/members", h.Member.ListMembers)
			policies.POST("/:id/members", h.Member.EnrollMember)

			// Installments (nested under policies)
			policies.GET("/:id/installments", h.Installment.GetSchedulesByPolicy)
			policies.POST("/:id/installments", h.Installment.CreateSchedule)
		}

		// Members
		members := authenticated.Group("/members")
		{
			members.PUT("/:id/verify", h.Member.VerifyMember)
			members.GET("/:id/eligibility", h.Member.GetMemberEligibility)
		}

		// Providers
		providers := authenticated.Group("/providers")
		{
			providers.GET("", h.Provider.ListProviders)
			providers.GET("/:id", h.Provider.GetProvider)
			providers.POST("", h.Provider.RegisterProvider)
			providers.PUT("/:id", h.Provider.UpdateProvider)
			providers.PUT("/:id/credential", h.Provider.CredentialProvider)
			providers.PUT("/:id/activate", h.Provider.ActivateProvider)
			providers.PUT("/:id/suspend", h.Provider.SuspendProvider)
			providers.PUT("/:id/terminate", h.Provider.TerminateProvider)

			// Contracts (nested under providers)
			providers.GET("/:id/contracts", h.Contract.ListContracts)
			providers.POST("/:id/contracts", h.Contract.CreateContract)

			// Rate Cards (nested under providers)
			providers.GET("/:id/rate-cards", h.RateCard.ListRateCards)
			providers.POST("/:id/rate-cards", h.RateCard.CreateRateCard)
		}

		// Pre-Authorizations
		preauths := authenticated.Group("/preauths")
		{
			preauths.GET("", h.PreAuth.ListPreAuths)
			preauths.GET("/:id", h.PreAuth.GetPreAuth)
			preauths.POST("", h.PreAuth.SubmitPreAuth)
			preauths.PUT("/:id/review", h.PreAuth.ReviewPreAuth)
			preauths.PUT("/:id/approve", h.PreAuth.ApprovePreAuth)
			preauths.PUT("/:id/deny", h.PreAuth.DenyPreAuth)
		}

		// Claims
		claims := authenticated.Group("/claims")
		{
			claims.GET("", h.Claim.ListClaims)
			claims.GET("/:id", h.Claim.GetClaim)
			claims.POST("", h.Claim.SubmitClaim)
			claims.PUT("/:id/approve", h.Claim.ApproveClaim)
			claims.PUT("/:id/reject", h.Claim.RejectClaim)
		}

		// Installment Schedules (standalone)
		installments := authenticated.Group("/installments")
		{
			installments.GET("/schedule/:id", h.Installment.ListInstallmentsBySchedule)
			installments.PUT("/:id/pay", h.Installment.MarkInstallmentPaid)
		}

		// Invoices
		invoices := authenticated.Group("/invoices")
		{
			invoices.GET("", h.Invoice.ListInvoices)
			invoices.GET("/:id", h.Invoice.GetInvoice)
		}

		// Payments
		payments := authenticated.Group("/payments")
		{
			payments.GET("", h.Payment.ListPayments)
			payments.GET("/:id", h.Payment.GetPayment)
			payments.POST("", h.Payment.InitiatePayment)
			payments.PUT("/:id/retry", h.Payment.RetryPayment)
			payments.PUT("/:id/reconcile", h.Payment.ReconcilePayment)
		}

		// Remittances
		remittances := authenticated.Group("/remittances")
		{
			remittances.GET("", h.Remittance.ListRemittances)
			remittances.GET("/:id", h.Remittance.GetRemittance)
			remittances.POST("", h.Remittance.CreateRemittance)
		}

		// Notifications
		notifications := authenticated.Group("/notifications")
		{
			notifications.GET("", h.Notification.ListNotifications)
			notifications.PUT("/:id/read", h.Notification.MarkRead)
			notifications.GET("/unread-count", h.Notification.GetUnreadCount)
		}

		// Audit
		audit := authenticated.Group("/audit")
		{
			audit.GET("", h.Audit.ListEvents)
			audit.GET("/entity/:type/:id", h.Audit.ListByEntity)
			audit.GET("/user/:id", h.Audit.ListByUser)
		}

		// Analytics
		analytics := authenticated.Group("/analytics")
		{
			analytics.GET("/dashboard", h.Analytics.GetDashboard)
			analytics.GET("/kpis", h.Analytics.GetKPIs)
			analytics.GET("/export", h.Analytics.ExportCSV)
		}
	}
}
