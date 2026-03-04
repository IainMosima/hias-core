package main

import (
	"context"
	"log"
	"time"

	"github.com/bitbiz/hias-core/configs"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/bitbiz/hias-core/infrastructures/queue"
	"github.com/bitbiz/hias-core/infrastructures/repository"
	"github.com/bitbiz/hias-core/services/analytics"
	"github.com/bitbiz/hias-core/services/api-gateway/handlers"
	"github.com/bitbiz/hias-core/services/api-gateway/routes"
	"github.com/bitbiz/hias-core/services/audit"
	"github.com/bitbiz/hias-core/services/billing"
	"github.com/bitbiz/hias-core/services/claims"
	"github.com/bitbiz/hias-core/services/identity"
	"github.com/bitbiz/hias-core/services/notification"
	"github.com/bitbiz/hias-core/services/policy"
	"github.com/bitbiz/hias-core/services/preauth"
	"github.com/bitbiz/hias-core/services/product"
	"github.com/bitbiz/hias-core/services/provider"
	"github.com/bitbiz/hias-core/services/sales"
	"github.com/bitbiz/hias-core/shared/auth"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// 1. Load configuration
	cfg, _, err := configs.LoadConfig("./configs")
	if err != nil {
		log.Printf("Warning: Could not load local config: %v", err)
	}

	// Try SSM parameters overlay
	if err := configs.LoadSSMParameters(&cfg); err != nil {
		log.Printf("Warning: Could not load SSM parameters: %v", err)
	}

	// 2. Database connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connPool, err := pgxpool.New(ctx, cfg.DBSource)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}
	defer connPool.Close()

	if err := connPool.Ping(ctx); err != nil {
		log.Fatalf("Cannot ping database: %v", err)
	}
	log.Println("Connected to database")

	store := db.NewStore(connPool)

	// 3. Auth infrastructure
	tokenMaker, err := auth.NewPasetoMaker(cfg.TokenSymmetricKey)
	if err != nil {
		log.Fatalf("Cannot create token maker: %v", err)
	}

	// 4. Queue manager (optional — gracefully degrade if unavailable)
	var queueManager queue.QueueManager
	queueConfig := queue.QueueConfig{
		AWSRegion:                  cfg.AWSRegion,
		AWSEndpointURL:             cfg.AWSEndpointURL,
		DocumentProcessingQueueURL: cfg.AWSSQSDocumentProcessingQueueURL,
		ExtractionResultsQueueURL:  cfg.AWSSQSExtractionResultsQueueURL,
		ClaimProcessingQueueURL:    cfg.AWSSQSClaimProcessingQueueURL,
		PaymentEventsQueueURL:      cfg.AWSSQSPaymentEventsQueueURL,
		NotificationEventsQueueURL: cfg.AWSSQSNotificationEventsQueueURL,
	}

	queueFactory := queue.NewQueueFactory(queueConfig)
	sqsPublisher, err := queueFactory.CreatePublisher()
	if err != nil {
		log.Printf("Warning: Cannot create SQS publisher: %v — notifications will be degraded", err)
	} else {
		watermillPub := queue.NewWatermillPublisher(sqsPublisher, queueConfig)
		sqsSubscriber, subErr := queueFactory.CreateSubscriber()
		if subErr != nil {
			log.Printf("Warning: Cannot create SQS subscriber: %v", subErr)
		} else {
			queueManager, err = queue.NewWatermillQueueManager(watermillPub, sqsSubscriber, queueConfig)
			if err != nil {
				log.Printf("Warning: Cannot create queue manager: %v", err)
			}
		}
	}

	// 5. Repositories
	userRepo := repository.NewUserRepository(store)
	roleRepo := repository.NewRoleRepository(store)
	permissionRepo := repository.NewPermissionRepository(store)
	planRepo := repository.NewPlanRepository(store)
	benefitRepo := repository.NewBenefitRepository(store)
	exclusionRepo := repository.NewExclusionRepository(store)
	policyRepo := repository.NewPolicyRepository(store)
	memberRepo := repository.NewMemberRepository(store)
	providerRepo := repository.NewProviderRepository(store)
	contractRepo := repository.NewContractRepository(store)
	rateCardRepo := repository.NewRateCardRepository(store)
	claimRepo := repository.NewClaimRepository(store)
	lineItemRepo := repository.NewClaimLineItemRepository(store)
	adjRepo := repository.NewAdjudicationRepository(store)
	fraudFlagRepo := repository.NewFraudFlagRepository(store)
	preauthRepo := repository.NewPreAuthRepository(store)
	invoiceRepo := repository.NewInvoiceRepository(store)
	paymentRepo := repository.NewPaymentRepository(store)
	remittanceRepo := repository.NewRemittanceRepository(store)
	notifRepo := repository.NewNotificationRepository(store)
	auditRepo := repository.NewAuditRepository(store)
	analyticsRepo := repository.NewAnalyticsRepository(store)
	premiumRuleRepo := repository.NewPremiumRuleRepository(store)
	providerNetworkRepo := repository.NewProviderNetworkRepository(store)
	installmentScheduleRepo := repository.NewInstallmentScheduleRepository(store)
	installmentRepo := repository.NewInstallmentRepository(store)
	leadRepo := repository.NewLeadRepository(store)
	leadActivityRepo := repository.NewLeadActivityRepository(store)
	quotationRepo := repository.NewQuotationRepository(store)
	quotationVersionRepo := repository.NewQuotationVersionRepository(store)
	quotationDocumentRepo := repository.NewQuotationDocumentRepository(store)
	approvalLimitRepo := repository.NewApprovalLimitRepository(store)

	// 6. Services (bottom-up dependency order)
	auditSvc := audit.NewAuditService(auditRepo)
	notifSvc := notification.NewNotificationService(notifRepo, queueManager)

	// Identity services
	authSvc := identity.NewAuthService(
		userRepo, roleRepo, permissionRepo, tokenMaker,
		identity.AuthServiceConfig{
			AccessTokenDuration:  cfg.AccessTokenDuration,
			RefreshTokenDuration: cfg.RefreshTokenDuration,
		},
	)
	userSvc := identity.NewUserService(userRepo, roleRepo, auditSvc)

	// Product services
	planSvc := product.NewPlanService(planRepo, auditSvc)
	benefitSvc := product.NewBenefitService(benefitRepo, auditSvc)
	exclusionSvc := product.NewExclusionService(exclusionRepo, planRepo, auditSvc)

	// Policy services
	policySvc := policy.NewPolicyService(policyRepo, planRepo, auditSvc)
	memberSvc := policy.NewMemberService(memberRepo, policyRepo, auditSvc)

	// Provider services
	providerSvc := provider.NewProviderService(providerRepo, contractRepo, rateCardRepo, auditSvc)

	// Claims services
	fraudSvc := claims.NewFraudService(fraudFlagRepo)
	validatorSvc := claims.NewValidatorService(policyRepo, memberRepo, providerRepo)
	adjudicatorSvc := claims.NewAdjudicatorService(claimRepo, benefitRepo, exclusionRepo, policyRepo, memberRepo, providerRepo, providerNetworkRepo, fraudSvc)
	claimSvc := claims.NewClaimService(claimRepo, lineItemRepo, adjudicatorSvc, validatorSvc, fraudSvc, adjRepo, fraudFlagRepo, auditSvc)

	// Pre-auth service
	preauthSvc := preauth.NewPreAuthService(preauthRepo)

	// Billing services
	billingSvc := billing.NewBillingService(invoiceRepo, policyRepo)
	paymentSvc := billing.NewPaymentService(paymentRepo, invoiceRepo)
	remittanceSvc := billing.NewRemittanceService(remittanceRepo, claimRepo, providerRepo)
	installmentSvc := billing.NewInstallmentService(installmentScheduleRepo, installmentRepo, policyRepo)

	// Product services (new)
	premiumRuleSvc := product.NewPremiumRuleService(premiumRuleRepo, planRepo, auditSvc)
	providerNetworkSvc := product.NewProviderNetworkService(providerNetworkRepo, planRepo, auditSvc)

	// Analytics service
	analyticsSvc := analytics.NewAnalyticsService(analyticsRepo)

	// Sales services
	leadSvc := sales.NewLeadService(leadRepo, leadActivityRepo, auditSvc)
	quotationSvc := sales.NewQuotationService(quotationRepo, quotationVersionRepo, quotationDocumentRepo, approvalLimitRepo, leadRepo, auditSvc, premiumRuleSvc, notifSvc, policySvc, memberSvc, installmentSvc)
	approvalLimitSvc := sales.NewApprovalLimitService(approvalLimitRepo, auditSvc)

	// Suppress unused variable warnings for services used internally
	_ = billingSvc

	// 7. Handlers
	h := routes.Handlers{
		Health:          handlers.NewHealthHandler(),
		Auth:            handlers.NewAuthHandler(authSvc),
		User:            handlers.NewUserHandler(userSvc),
		Plan:            handlers.NewPlanHandler(planSvc),
		Benefit:         handlers.NewBenefitHandler(benefitSvc),
		Exclusion:       handlers.NewExclusionHandler(exclusionSvc),
		PremiumRule:     handlers.NewPremiumRuleHandler(premiumRuleSvc),
		ProviderNetwork: handlers.NewProviderNetworkHandler(providerNetworkSvc),
		Policy:          handlers.NewPolicyHandler(policySvc),
		Member:          handlers.NewMemberHandler(memberSvc),
		Provider:        handlers.NewProviderHandler(providerSvc),
		Contract:        handlers.NewContractHandler(contractRepo),
		RateCard:        handlers.NewRateCardHandler(rateCardRepo),
		Claim:           handlers.NewClaimHandler(claimSvc),
		PreAuth:         handlers.NewPreAuthHandler(preauthSvc),
		Invoice:         handlers.NewInvoiceHandler(invoiceRepo),
		Payment:         handlers.NewPaymentHandler(paymentSvc),
		Remittance:      handlers.NewRemittanceHandler(remittanceSvc),
		Installment:     handlers.NewInstallmentHandler(installmentSvc),
		Notification:    handlers.NewNotificationHandler(notifSvc),
		Audit:           handlers.NewAuditHandler(auditSvc),
		Analytics:       handlers.NewAnalyticsHandler(analyticsSvc),
		Lead:            handlers.NewLeadHandler(leadSvc),
		Quotation:       handlers.NewQuotationHandler(quotationSvc),
		ApprovalLimit:   handlers.NewApprovalLimitHandler(approvalLimitSvc),
	}

	// 8. Server
	server := NewServer(tokenMaker, cfg.AllowedOrigins)
	server.RegisterRoutes(h)

	// 9. Start
	if err := server.Start(cfg.HTTPServerAddress); err != nil {
		log.Fatalf("Cannot start server: %v", err)
	}
}
