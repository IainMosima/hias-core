package routes

import (
	"github.com/bitbiz/hias-core/services/api-gateway/handlers"
	"github.com/bitbiz/hias-core/services/api-gateway/middleware"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Health           *handlers.HealthHandler
	Auth             *handlers.AuthHandler
	User             *handlers.UserHandler
	Plan             *handlers.PlanHandler
	Benefit          *handlers.BenefitHandler
	Exclusion        *handlers.ExclusionHandler
	PremiumRule      *handlers.PremiumRuleHandler
	ProviderNetwork  *handlers.ProviderNetworkHandler
	Policy           *handlers.PolicyHandler
	Member           *handlers.MemberHandler
	Provider         *handlers.ProviderHandler
	Contract         *handlers.ContractHandler
	RateCard         *handlers.RateCardHandler
	Claim            *handlers.ClaimHandler
	PreAuth          *handlers.PreAuthHandler
	Invoice          *handlers.InvoiceHandler
	Payment          *handlers.PaymentHandler
	Remittance       *handlers.RemittanceHandler
	Installment      *handlers.InstallmentHandler
	Notification     *handlers.NotificationHandler
	Audit            *handlers.AuditHandler
	Analytics        *handlers.AnalyticsHandler
	Lead             *handlers.LeadHandler
	Quotation        *handlers.QuotationHandler
	ApprovalLimit    *handlers.ApprovalLimitHandler
	Endorsement      *handlers.EndorsementHandler
	Renewal          *handlers.RenewalHandler
	Underwriting     *handlers.UnderwritingHandler
	PolicyDocument   *handlers.PolicyDocumentHandler
	UnderwritingFlag *handlers.UnderwritingFlagHandler
	UnderwritingRule *handlers.UnderwritingRuleHandler
	CreditNote       *handlers.CreditNoteHandler
	Case             *handlers.CaseHandler
	Statement        *handlers.StatementHandler

	// Reinsurance
	Treaty               *handlers.TreatyHandler
	Cession              *handlers.CessionHandler
	Recovery             *handlers.RecoveryHandler
	Bordereau            *handlers.BordereauHandler
	ReinsurerStatement   *handlers.ReinsurerStatementHandler
	TreatyAlert          *handlers.TreatyAlertHandler
	ReinsuranceAnalytics *handlers.ReinsuranceAnalyticsHandler

	// Reporting
	Report *handlers.ReportHandler
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
			plans.GET("/:id/rate-sheet", h.PremiumRule.GetRateSheet)

			// Provider Networks (nested under plans)
			plans.GET("/:id/provider-networks", h.ProviderNetwork.ListProviderNetworks)
			plans.POST("/:id/provider-networks", h.ProviderNetwork.CreateProviderNetwork)

			// Underwriting Rules (nested under plans)
			plans.GET("/:id/underwriting-rules", h.UnderwritingRule.ListRules)
			plans.POST("/:id/underwriting-rules", h.UnderwritingRule.CreateRule)
		}

		// Benefits (standalone for sub-benefits)
		benefits := authenticated.Group("/benefits")
		{
			benefits.GET("/:id/sub-benefits", h.Benefit.ListSubBenefits)
			benefits.POST("/:id/sub-benefits", h.Benefit.CreateSubBenefit)
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
			policies.GET("/by-status", h.Policy.ListPoliciesByStatus)
			policies.GET("/:id", h.Policy.GetPolicy)
			policies.POST("", h.Policy.CreatePolicy)
			policies.PUT("/:id", h.Policy.UpdatePolicy)
			policies.PUT("/:id/activate", h.Policy.ActivatePolicy)
			policies.PUT("/:id/lapse", h.Policy.LapsePolicy)
			policies.PUT("/:id/terminate", h.Policy.TerminatePolicy)
			policies.PUT("/:id/reinstate", h.Policy.ReinstatePolicy)
			policies.PUT("/:id/suspend", h.Policy.SuspendPolicy)
			policies.PUT("/:id/change-plan", h.Policy.ChangePlan)
			policies.GET("/:id/prorate", h.Policy.CalculateProrate)

			// Bulk policy operations
			policies.POST("/bulk/activate", h.Policy.BulkActivate)
			policies.POST("/bulk/lapse", h.Policy.BulkLapse)

			// Members (nested under policies)
			policies.GET("/:id/members", h.Member.ListMembers)
			policies.POST("/:id/members", h.Member.EnrollMember)
			policies.POST("/:id/members/bulk", h.Member.BulkEnrollMembers)
			policies.POST("/:id/members/import", h.Member.ImportMembersCSV)
			policies.POST("/:id/members/bulk-remove", h.Member.BulkRemoveMembers)

			// Endorsements (nested under policies)
			policies.GET("/:id/endorsements", h.Endorsement.ListEndorsements)
			policies.POST("/:id/endorsements", h.Endorsement.CreateEndorsement)

			// Renewals (nested under policies)
			policies.GET("/:id/renewals", h.Renewal.ListRenewals)
			policies.POST("/:id/renewals", h.Renewal.InitiateRenewal)

			// Underwriting (nested under policies)
			policies.GET("/:id/underwriting", h.Underwriting.ListAssessments)
			policies.POST("/:id/underwriting", h.Underwriting.SubmitAssessment)

			// Policy Documents (nested under policies)
			policies.GET("/:id/documents", h.PolicyDocument.ListDocuments)
			policies.POST("/:id/documents/welcome-letter", h.PolicyDocument.GenerateWelcomeLetter)
			policies.POST("/:id/documents/policy-schedule", h.PolicyDocument.GeneratePolicySchedule)
			policies.POST("/:id/documents/member-cards", h.PolicyDocument.BulkGenerateMemberCards)

			// Underwriting Flags (nested under policies)
			policies.GET("/:id/underwriting-flags", h.UnderwritingFlag.ListByPolicy)

			// Credit Notes (nested under policies)
			policies.GET("/:id/credit-notes", h.CreditNote.ListByPolicy)

			// Installments (nested under policies)
			policies.GET("/:id/installments", h.Installment.GetSchedulesByPolicy)
			policies.POST("/:id/installments", h.Installment.CreateSchedule)

			// Cases (nested under policies)
			policies.GET("/:id/cases", h.Case.ListByPolicy)
		}

		// Members
		members := authenticated.Group("/members")
		{
			members.GET("/:id", h.Member.GetMember)
			members.PUT("/:id", h.Member.UpdateMember)
			members.PUT("/:id/verify", h.Member.VerifyMember)
			members.PUT("/:id/suspend", h.Member.SuspendMember)
			members.PUT("/:id/reactivate", h.Member.ReactivateMember)
			members.DELETE("/:id", h.Member.RemoveMember)
			members.GET("/:id/eligibility", h.Member.GetMemberEligibility)
			members.POST("/:id/card", h.PolicyDocument.GenerateMemberCard)
			members.GET("/:id/underwriting-flags", h.UnderwritingFlag.ListByMember)
			members.GET("/:id/cases", h.Case.ListByMember)
		}

		// Endorsements (standalone)
		endorsements := authenticated.Group("/endorsements")
		{
			endorsements.GET("/:id", h.Endorsement.GetEndorsement)
			endorsements.PUT("/:id/approve", h.Endorsement.ApproveEndorsement)
			endorsements.PUT("/:id/reject", h.Endorsement.RejectEndorsement)
			endorsements.PUT("/:id/apply", h.Endorsement.ApplyEndorsement)
		}

		// Renewals (standalone)
		renewals := authenticated.Group("/renewals")
		{
			renewals.GET("/:id", h.Renewal.GetRenewal)
			renewals.PUT("/:id/approve", h.Renewal.ApproveRenewal)
			renewals.PUT("/:id/reject", h.Renewal.RejectRenewal)
			renewals.POST("/:id/complete", h.Renewal.CompleteRenewal)
			renewals.POST("/expire", middleware.RequireRole(string(shared.UserRoleAdmin)), h.Renewal.ExpireRenewals)
			renewals.POST("/bulk", h.Renewal.BulkInitiateRenewals)
		}

		// Underwriting (standalone)
		underwriting := authenticated.Group("/underwriting")
		{
			underwriting.GET("/:id", h.Underwriting.GetAssessment)
			underwriting.PUT("/:id/review", middleware.RequireRole(
				string(shared.UserRoleAdmin), string(shared.UserRoleUnderwriter)),
				h.Underwriting.ReviewAssessment)
		}

		// Underwriting Flags (standalone)
		uwFlags := authenticated.Group("/underwriting-flags")
		{
			uwFlags.GET("", h.UnderwritingFlag.ListOpen)
			uwFlags.GET("/count", h.UnderwritingFlag.CountOpen)
			uwFlags.GET("/:id", h.UnderwritingFlag.GetFlag)
			uwFlags.PUT("/:id/resolve", middleware.RequireRole(
				string(shared.UserRoleAdmin), string(shared.UserRoleUnderwriter)),
				h.UnderwritingFlag.ResolveFlag)
			uwFlags.PUT("/:id/override", middleware.RequireRole(
				string(shared.UserRoleAdmin), string(shared.UserRoleUnderwriter)),
				h.UnderwritingFlag.OverrideFlag)
		}

		// Underwriting Rules (standalone for update/delete)
		uwRules := authenticated.Group("/underwriting-rules")
		{
			uwRules.PUT("/:id", h.UnderwritingRule.UpdateRule)
			uwRules.DELETE("/:id", h.UnderwritingRule.DeleteRule)
		}

		// Credit Notes (standalone)
		creditNotes := authenticated.Group("/credit-notes")
		{
			creditNotes.GET("/:id", h.CreditNote.GetCreditNote)
			creditNotes.PUT("/:id/approve", middleware.RequireRole(
				string(shared.UserRoleAdmin)), h.CreditNote.ApproveCreditNote)
			creditNotes.PUT("/:id/apply", middleware.RequireRole(
				string(shared.UserRoleAdmin)), h.CreditNote.ApplyCreditNote)
		}

		// Policy Documents (standalone)
		policyDocs := authenticated.Group("/policy-documents")
		{
			policyDocs.GET("/:id", h.PolicyDocument.GetDocument)
			policyDocs.DELETE("/:id", h.PolicyDocument.DeleteDocument)
		}

		// Providers
		providers := authenticated.Group("/providers")
		{
			providers.GET("", h.Provider.ListProviders)
			providers.GET("/by-tier", h.Provider.ListByTier)
			providers.GET("/by-accreditation", h.Provider.ListByAccreditationStatus)
			providers.GET("/expiring-accreditations", h.Provider.ListExpiringAccreditations)
			providers.GET("/:id", h.Provider.GetProvider)
			providers.POST("", h.Provider.RegisterProvider)
			providers.PUT("/:id", h.Provider.UpdateProvider)
			providers.PUT("/:id/credential", h.Provider.CredentialProvider)
			providers.PUT("/:id/activate", h.Provider.ActivateProvider)
			providers.PUT("/:id/suspend", h.Provider.SuspendProvider)
			providers.PUT("/:id/terminate", h.Provider.TerminateProvider)
			providers.PUT("/:id/tier", h.Provider.UpdateTier)
			providers.PUT("/:id/accreditation", h.Provider.UpdateAccreditation)

			// Contracts (nested under providers)
			providers.GET("/:id/contracts", h.Contract.ListContracts)
			providers.POST("/:id/contracts", h.Contract.CreateContract)

			// Rate Cards (nested under providers)
			providers.GET("/:id/rate-cards", h.RateCard.ListRateCards)
			providers.POST("/:id/rate-cards", h.RateCard.CreateRateCard)
			providers.POST("/:id/rate-cards/bulk", h.RateCard.BulkCreateRateCards)

			// Provider Statements (nested under providers)
			providers.GET("/:id/statements", h.Statement.ListByProvider)
			providers.POST("/:id/statements", h.Statement.UploadStatement)

			// Cases (nested under providers)
			providers.GET("/:id/cases", h.Case.ListByProvider)
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
			preauths.POST("/:id/lou", h.PolicyDocument.GenerateLOU)
		}

		// Claims
		claims := authenticated.Group("/claims")
		{
			claims.GET("", h.Claim.ListClaims)
			claims.GET("/sla-breached", h.Claim.ListSLABreached)
			claims.GET("/:id", h.Claim.GetClaim)
			claims.POST("", h.Claim.SubmitClaim)
			claims.POST("/bulk", h.Claim.BulkSubmitClaims)
			claims.POST("/import-csv", middleware.RequireRole(
				string(shared.UserRoleAdmin),
				string(shared.UserRoleClaimsOfficer),
			), h.Claim.ImportClaimsCSV)
			claims.PUT("/:id/vet", middleware.RequireRole(
				string(shared.UserRoleAdmin),
				string(shared.UserRoleClaimsOfficer),
			), h.Claim.VetClaim)
			claims.PUT("/:id/approve", middleware.RequireRole(
				string(shared.UserRoleAdmin),
				string(shared.UserRoleManager),
			), h.Claim.ApproveClaim)
			claims.PUT("/:id/reject", middleware.RequireRole(
				string(shared.UserRoleAdmin),
				string(shared.UserRoleManager),
			), h.Claim.RejectClaim)
			claims.PUT("/:id/ready-for-payment", middleware.RequireRole(
				string(shared.UserRoleAdmin),
				string(shared.UserRoleFinance),
			), h.Claim.MarkReadyForPayment)
			claims.PUT("/:id/mark-paid", middleware.RequireRole(
				string(shared.UserRoleAdmin),
				string(shared.UserRoleFinance),
			), h.Claim.MarkPaid)
			claims.PUT("/:id/mark-part-paid", middleware.RequireRole(
				string(shared.UserRoleAdmin),
				string(shared.UserRoleFinance),
			), h.Claim.MarkPartPaid)

			// Claim Documents (nested under claims)
			claims.GET("/:id/documents", h.Claim.ListClaimDocuments)
			claims.POST("/:id/documents", h.Claim.UploadClaimDocument)

			// Decline Letter
			claims.POST("/:id/decline-letter", h.PolicyDocument.GenerateDeclineLetter)
		}

		// Claim Documents (standalone for delete)
		claimDocs := authenticated.Group("/claim-documents")
		{
			claimDocs.DELETE("/:id", h.Claim.DeleteClaimDocument)
		}

		// Cases (standalone)
		cases := authenticated.Group("/cases")
		{
			cases.GET("", h.Case.ListCases)
			cases.GET("/count", h.Case.CountByStatus)
			cases.GET("/:id", h.Case.GetCase)
			cases.POST("", h.Case.CreateCase)
			cases.PUT("/:id", h.Case.UpdateCase)
			cases.PUT("/:id/admit", h.Case.AdmitCase)
			cases.PUT("/:id/start-treatment", h.Case.StartTreatment)
			cases.PUT("/:id/discharge", h.Case.DischargeCase)
			cases.PUT("/:id/close", h.Case.CloseCase)
		}

		// Provider Statements (standalone)
		providerStatements := authenticated.Group("/provider-statements")
		{
			providerStatements.GET("/:id", h.Statement.GetStatement)
			providerStatements.GET("/:id/line-items", h.Statement.ListLineItems)
			providerStatements.POST("/:id/reconcile", h.Statement.ReconcileStatement)
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
			remittances.GET("/:id/export", h.Remittance.ExportPaymentFile)
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

		// Leads
		leads := authenticated.Group("/leads")
		{
			leads.GET("", h.Lead.ListLeads)
			leads.POST("", h.Lead.CreateLead)
			leads.GET("/due-follow-ups", h.Lead.GetDueFollowUps)
			leads.GET("/:id", h.Lead.GetLead)
			leads.PUT("/:id", h.Lead.UpdateLead)
			leads.PUT("/:id/status", h.Lead.UpdateLeadStatus)
			leads.GET("/:id/activities", h.Lead.ListActivities)
			leads.POST("/:id/activities", h.Lead.AddActivity)
			leads.GET("/:id/quotations", h.Quotation.ListQuotationsByLead)
		}

		// Quotations
		quotations := authenticated.Group("/quotations")
		{
			quotations.GET("", h.Quotation.ListQuotations)
			quotations.POST("", h.Quotation.CreateQuotation)
			quotations.POST("/expire", middleware.RequireRole(string(shared.UserRoleAdmin)), h.Quotation.ExpireQuotations)
			quotations.GET("/:id", h.Quotation.GetQuotation)
			quotations.PUT("/:id/issue", h.Quotation.IssueQuotation)
			quotations.PUT("/:id/accept", h.Quotation.AcceptQuotation)
			quotations.PUT("/:id/decline", h.Quotation.DeclineQuotation)
			quotations.PUT("/:id/send", h.Quotation.SendToClient)
			quotations.POST("/:id/convert", h.Quotation.ConvertToPolicy)

			// Versions (nested under quotations)
			quotations.GET("/:id/versions", h.Quotation.ListVersions)
			quotations.POST("/:id/versions", h.Quotation.CreateVersion)
			quotations.GET("/:id/versions/compare", h.Quotation.CompareVersions)
			quotations.GET("/:id/versions/:version", h.Quotation.GetVersion)
			quotations.PUT("/:id/versions/:version/submit-approval", h.Quotation.SubmitForApproval)
			quotations.PUT("/:id/versions/:version/approve", middleware.RequireRole(string(shared.UserRoleAdmin), string(shared.UserRoleUnderwriter), string(shared.UserRoleManager)), h.Quotation.ApproveVersion)
			quotations.PUT("/:id/versions/:version/reject", middleware.RequireRole(string(shared.UserRoleAdmin), string(shared.UserRoleUnderwriter), string(shared.UserRoleManager)), h.Quotation.RejectVersion)

			// Documents (nested under quotations)
			quotations.GET("/:id/documents", h.Quotation.ListDocuments)
			quotations.POST("/:id/documents", h.Quotation.UploadDocument)
		}

		// Quotation Documents (standalone for download/delete)
		quotationDocs := authenticated.Group("/quotation-documents")
		{
			quotationDocs.PUT("/:id", h.Quotation.UpdateDocument)
			quotationDocs.DELETE("/:id", h.Quotation.DeleteDocument)
		}

		// Approval Limits
		approvalLimits := authenticated.Group("/approval-limits")
		approvalLimits.Use(middleware.RequireRole(string(shared.UserRoleAdmin)))
		{
			approvalLimits.GET("", h.ApprovalLimit.ListLimits)
			approvalLimits.POST("", h.ApprovalLimit.CreateLimit)
			approvalLimits.PUT("/:id", h.ApprovalLimit.UpdateLimit)
		}

		// ===== Reinsurance =====

		// Treaties
		treaties := authenticated.Group("/treaties")
		{
			treaties.GET("", h.Treaty.ListTreaties)
			treaties.POST("", h.Treaty.CreateTreaty)
			treaties.POST("/expire", h.Treaty.ExpireOverdue)
			treaties.GET("/:id", h.Treaty.GetTreaty)
			treaties.PUT("/:id", h.Treaty.UpdateTreaty)
			treaties.PUT("/:id/activate", h.Treaty.ActivateTreaty)
			treaties.PUT("/:id/terminate", h.Treaty.TerminateTreaty)

			// Participants (nested under treaties)
			treaties.GET("/:id/participants", h.Treaty.ListParticipants)
			treaties.POST("/:id/participants", h.Treaty.AddParticipant)
			treaties.PUT("/:id/participants/:participantId", h.Treaty.UpdateParticipant)
			treaties.DELETE("/:id/participants/:participantId", h.Treaty.RemoveParticipant)

			// Layers (nested under treaties)
			treaties.GET("/:id/layers", h.Treaty.ListLayers)
			treaties.POST("/:id/layers", h.Treaty.AddLayer)
			treaties.PUT("/:id/layers/:layerId", h.Treaty.UpdateLayer)
			treaties.DELETE("/:id/layers/:layerId", h.Treaty.RemoveLayer)

			// Profit Commission Rules (nested under treaties)
			treaties.GET("/:id/profit-commission-rules", h.Treaty.ListProfitCommissionRules)
			treaties.POST("/:id/profit-commission-rules", h.Treaty.AddProfitCommissionRule)
			treaties.DELETE("/:id/profit-commission-rules/:ruleId", h.Treaty.RemoveProfitCommissionRule)

			// Treaty sub-resources
			treaties.GET("/:id/cessions", h.Cession.ListCessions)
			treaties.GET("/:id/recoveries", h.Recovery.ListRecoveries)
			treaties.GET("/:id/bordereaux", h.Bordereau.ListByTreaty)
			treaties.GET("/:id/statements", h.ReinsurerStatement.ListByTreaty)
			treaties.GET("/:id/alerts", h.TreatyAlert.ListByTreaty)
		}

		// Cessions
		cessions := authenticated.Group("/cessions")
		{
			cessions.POST("", h.Cession.CedePremium)
			cessions.POST("/auto-cede", h.Cession.AutoCede)
			cessions.GET("/:id", h.Cession.GetCession)
			cessions.PUT("/:id/book", h.Cession.BookCession)
			cessions.PUT("/:id/reverse", h.Cession.ReverseCession)
		}

		// Recoveries
		recoveries := authenticated.Group("/recoveries")
		{
			recoveries.POST("", h.Recovery.CreateRecovery)
			recoveries.GET("/outstanding", h.Recovery.ListOutstanding)
			recoveries.GET("/aged-analysis", h.Recovery.AgedAnalysis)
			recoveries.POST("/apply-for-claim/:claimId", h.Recovery.ApplyRecoveryForClaim)
			recoveries.GET("/:id", h.Recovery.GetRecovery)
			recoveries.PUT("/:id/acknowledge", h.Recovery.AcknowledgeRecovery)
			recoveries.PUT("/:id/request-info", h.Recovery.RequestInfo)
			recoveries.PUT("/:id/approve", h.Recovery.ApproveRecovery)
			recoveries.PUT("/:id/record-payment", h.Recovery.RecordPayment)
			recoveries.PUT("/:id/write-off", h.Recovery.WriteOffRecovery)
			recoveries.GET("/:id/workflow", h.Recovery.GetWorkflowEvents)
		}

		// Bordereaux
		bordereaux := authenticated.Group("/bordereaux")
		{
			bordereaux.POST("/premium", h.Bordereau.GeneratePremiumBordereau)
			bordereaux.POST("/claim", h.Bordereau.GenerateClaimBordereau)
			bordereaux.GET("/:id", h.Bordereau.GetBordereau)
			bordereaux.PUT("/:id/finalize", h.Bordereau.FinalizeBordereau)
			bordereaux.PUT("/:id/mark-sent", h.Bordereau.MarkSent)
			bordereaux.GET("/:id/items", h.Bordereau.ListItems)
		}

		// Reinsurer Statements
		reinsurerStatements := authenticated.Group("/reinsurer-statements")
		{
			reinsurerStatements.POST("", h.ReinsurerStatement.GenerateStatement)
			reinsurerStatements.POST("/profit-commission", h.ReinsurerStatement.CalculateProfitCommission)
			reinsurerStatements.GET("/:id", h.ReinsurerStatement.GetStatement)
			reinsurerStatements.PUT("/:id/issue", h.ReinsurerStatement.IssueStatement)
			reinsurerStatements.PUT("/:id/acknowledge", h.ReinsurerStatement.AcknowledgeStatement)
			reinsurerStatements.PUT("/:id/settle", h.ReinsurerStatement.SettleStatement)
		}

		// Treaty Alerts
		treatyAlerts := authenticated.Group("/treaty-alerts")
		{
			treatyAlerts.GET("", h.TreatyAlert.ListAlerts)
			treatyAlerts.GET("/unacknowledged", h.TreatyAlert.ListUnacknowledged)
			treatyAlerts.GET("/count", h.TreatyAlert.CountUnacknowledged)
			treatyAlerts.PUT("/:id/acknowledge", h.TreatyAlert.AcknowledgeAlert)
			treatyAlerts.POST("/check-limits/:treatyId", h.TreatyAlert.CheckTreatyLimits)
			treatyAlerts.POST("/check-catastrophe/:treatyId", h.TreatyAlert.CheckCatastropheThresholds)
			treatyAlerts.POST("/check-expiry", h.TreatyAlert.CheckExpiryWarnings)
		}

		// Reinsurance Analytics
		reinsuranceAnalytics := authenticated.Group("/analytics/reinsurance")
		{
			reinsuranceAnalytics.GET("", h.ReinsuranceAnalytics.GetReinsuranceDashboard)
		}

		// ===== Reporting =====
		reports := authenticated.Group("/reports")
		{
			reports.GET("/definitions", h.Report.ListDefinitions)
			reports.GET("/definitions/:id", h.Report.GetDefinition)
			reports.POST("/definitions/adhoc", middleware.RequireRole(
				string(shared.UserRoleAdmin), string(shared.UserRoleManager),
			), h.Report.CreateAdHocDefinition)
			reports.POST("/generate", h.Report.GenerateReport)
			reports.POST("/preview", h.Report.PreviewReport)
			reports.POST("/drilldown", h.Report.DrillDown)
			reports.GET("/generated", h.Report.ListGeneratedReports)
			reports.GET("/generated/:id", h.Report.GetGeneratedReport)
			reports.GET("/generated/:id/download", h.Report.DownloadReport)
			reports.POST("/schedules", middleware.RequireRole(
				string(shared.UserRoleAdmin), string(shared.UserRoleManager),
			), h.Report.CreateSchedule)
			reports.GET("/schedules", middleware.RequireRole(
				string(shared.UserRoleAdmin), string(shared.UserRoleManager),
			), h.Report.ListSchedules)
			reports.PUT("/schedules/:id", middleware.RequireRole(
				string(shared.UserRoleAdmin), string(shared.UserRoleManager),
			), h.Report.UpdateSchedule)
			reports.DELETE("/schedules/:id", middleware.RequireRole(
				string(shared.UserRoleAdmin), string(shared.UserRoleManager),
			), h.Report.DeleteSchedule)
			reports.GET("/dashboard", middleware.RequireRole(
				string(shared.UserRoleAdmin), string(shared.UserRoleManager), string(shared.UserRoleFinance),
			), h.Report.GetManagementDashboard)
		}
	}
}
