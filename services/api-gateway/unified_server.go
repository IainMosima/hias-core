package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bitbiz/hias-core/configs"
	"github.com/bitbiz/hias-core/infrastructures/cache"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	infraNotif "github.com/bitbiz/hias-core/infrastructures/notifications"
	infraRepo "github.com/bitbiz/hias-core/infrastructures/repository"
	"github.com/bitbiz/hias-core/infrastructures/scheduler"
	"github.com/bitbiz/hias-core/infrastructures/scheduler/tasks"
	"github.com/bitbiz/hias-core/infrastructures/sse"
	"github.com/bitbiz/hias-core/infrastructures/websocket"
	"github.com/bitbiz/hias-core/services/analytics"
	"github.com/bitbiz/hias-core/services/api-gateway/rest"
	"github.com/bitbiz/hias-core/services/api-gateway/rest/handlers"
	"github.com/bitbiz/hias-core/services/audit"
	"github.com/bitbiz/hias-core/services/billing"
	"github.com/bitbiz/hias-core/services/claims"
	identityService "github.com/bitbiz/hias-core/services/identity"
	"github.com/bitbiz/hias-core/services/notification"
	"github.com/bitbiz/hias-core/services/policy"
	"github.com/bitbiz/hias-core/services/preauth"
	"github.com/bitbiz/hias-core/services/product"
	"github.com/bitbiz/hias-core/services/provider"
	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UnifiedServer struct {
	config     configs.Config
	store      db.Store
	cache      cache.CacheManager
	restServer *rest.RestServer
	scheduler  scheduler.SchedulerManager
	sseManager sse.SSEManager
	wsManager  websocket.WebSocketManager
}

func NewUnifiedServer(config configs.Config) (*UnifiedServer, error) {
	// ─── Database Connection ─────────────────────────────────────────
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to db: %w", err)
	}
	if err := connPool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("cannot ping db: %w", err)
	}
	log.Println("Connected to database")

	store := db.NewStore(connPool)

	// ─── Redis Cache ─────────────────────────────────────────────────
	var cacheManager cache.CacheManager
	if config.RedisURL != "" {
		redisClient, err := cache.NewRedisClient(cache.CacheConfig{RedisURL: config.RedisURL})
		if err != nil {
			log.Printf("Warning: Failed to connect to Redis: %v", err)
		} else {
			cacheManager = cache.NewRedisCacheManager(redisClient)
			log.Println("Connected to Redis")
		}
	}

	// ─── Token Maker ─────────────────────────────────────────────────
	tokenMaker, err := auth.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	// ─── Infrastructure ──────────────────────────────────────────────
	notifFactory := infraNotif.NewNotificationFactory(
		config.SMSAPIKey, config.SMSUsername, config.SMSSenderID,
		nil, // sesClient — would be initialized with AWS SDK in production
		config.SESFromEmail,
	)
	notifManager := infraNotif.NewNotificationManager(notifFactory)

	sseManager := sse.NewSSEManager()
	wsManager := websocket.NewWebSocketManager()

	// ─── Repositories ────────────────────────────────────────────────
	// Identity
	userRepo := infraRepo.NewUserRepository(store)
	roleRepo := infraRepo.NewRoleRepository(store)
	permissionRepo := infraRepo.NewPermissionRepository(store)

	// Product
	planRepo := infraRepo.NewPlanRepository(store)
	benefitRepo := infraRepo.NewBenefitRepository(store)
	exclusionRepo := infraRepo.NewExclusionRepository(store)

	// Policy
	policyRepo := infraRepo.NewPolicyRepository(store)
	memberRepo := infraRepo.NewMemberRepository(store)

	// Provider
	providerRepo := infraRepo.NewProviderRepository(store)
	contractRepo := infraRepo.NewContractRepository(store)
	rateCardRepo := infraRepo.NewRateCardRepository(store)

	// Claims
	claimRepo := infraRepo.NewClaimRepository(store)
	lineItemRepo := infraRepo.NewClaimLineItemRepository(store)
	adjudicationRepo := infraRepo.NewAdjudicationRepository(store)
	fraudFlagRepo := infraRepo.NewFraudFlagRepository(store)

	// Pre-Auth
	preauthRepo := infraRepo.NewPreAuthRepository(store)

	// Billing
	invoiceRepo := infraRepo.NewInvoiceRepository(store)
	paymentRepo := infraRepo.NewPaymentRepository(store)
	remittanceRepo := infraRepo.NewRemittanceRepository(store)

	// Notification, Analytics, Audit
	notificationRepo := infraRepo.NewNotificationRepository(store)
	analyticsRepo := infraRepo.NewAnalyticsRepository(store)
	auditRepo := infraRepo.NewAuditRepository(store)

	// ─── Services ────────────────────────────────────────────────────
	// Identity
	authService := identityService.NewAuthService(
		userRepo, roleRepo, permissionRepo,
		tokenMaker, nil, // cognitoService nil for local dev
		config.AccessTokenDuration,
	)
	userService := identityService.NewUserService(
		userRepo, roleRepo, permissionRepo,
		nil, // cognitoService
		tokenMaker,
	)

	// Product
	planService := product.NewPlanService(planRepo, benefitRepo, exclusionRepo)
	benefitService := product.NewBenefitService(benefitRepo, planRepo, exclusionRepo)

	// Policy
	policyService := policy.NewPolicyService(policyRepo, planRepo)
	memberService := policy.NewMemberService(memberRepo, policyRepo)

	// Provider
	providerService := provider.NewProviderService(providerRepo, contractRepo, rateCardRepo)

	// Claims (adjudication engine)
	fraudService := claims.NewFraudService(claimRepo, rateCardRepo, fraudFlagRepo)
	validatorService := claims.NewValidatorService(policyRepo, memberRepo, providerRepo)
	adjudicatorService := claims.NewAdjudicatorService(
		claimRepo, benefitRepo, exclusionRepo, policyRepo, memberRepo, fraudService,
	)
	claimService := claims.NewClaimService(
		claimRepo, lineItemRepo, adjudicationRepo,
		policyRepo, memberRepo, providerRepo,
		validatorService, adjudicatorService,
	)

	// Pre-Auth
	preauthService := preauth.NewPreAuthService(preauthRepo, policyRepo, memberRepo)

	// Billing
	billingService := billing.NewBillingService(invoiceRepo, policyRepo)
	paymentService := billing.NewPaymentService(paymentRepo, invoiceRepo)
	remittanceService := billing.NewRemittanceService(remittanceRepo, claimRepo, providerRepo, paymentRepo)

	// Notification, Analytics, Audit
	notificationService := notification.NewNotificationService(notificationRepo, notifManager)
	analyticsService := analytics.NewAnalyticsService(analyticsRepo)
	auditService := audit.NewAuditService(auditRepo)

	// ─── Handlers ────────────────────────────────────────────────────
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	planHandler := handlers.NewPlanHandler(planService, benefitService)
	policyHandler := handlers.NewPolicyHandler(policyService)
	memberHandler := handlers.NewMemberHandler(memberService)
	providerHandler := handlers.NewProviderHandler(providerService)
	claimHandler := handlers.NewClaimHandler(claimService)
	preAuthHandler := handlers.NewPreAuthHandler(preauthService)
	billingHandler := handlers.NewBillingHandler(billingService, paymentService)
	remittanceHandler := handlers.NewRemittanceHandler(remittanceService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)
	auditHandler := handlers.NewAuditHandler(auditService)

	// ─── REST Server ─────────────────────────────────────────────────
	restServer := rest.NewRestServer(rest.RestServerConfig{
		Config:              config,
		TokenMaker:          tokenMaker,
		AuthHandler:         authHandler,
		UserHandler:         userHandler,
		PlanHandler:         planHandler,
		PolicyHandler:       policyHandler,
		MemberHandler:       memberHandler,
		ProviderHandler:     providerHandler,
		ClaimHandler:        claimHandler,
		PreAuthHandler:      preAuthHandler,
		BillingHandler:      billingHandler,
		RemittanceHandler:   remittanceHandler,
		NotificationHandler: notificationHandler,
		AnalyticsHandler:    analyticsHandler,
		AuditHandler:        auditHandler,
		SSEManager:          sseManager,
		WSManager:           wsManager,
	})

	// ─── Scheduler ───────────────────────────────────────────────────
	schedConfig := scheduler.DefaultSchedulerConfig()
	schedulerMgr := scheduler.NewSchedulerManager()

	schedulerMgr.RegisterTask(tasks.NewBillingCycleTask(schedConfig.BillingCycleCron))
	schedulerMgr.RegisterTask(tasks.NewPaymentReminderTask(schedConfig.PaymentReminderCron))
	schedulerMgr.RegisterTask(tasks.NewPolicyLapseTask(schedConfig.PolicyLapseCron))
	schedulerMgr.RegisterTask(tasks.NewPreAuthExpiryTask(schedConfig.PreAuthExpiryCron))
	schedulerMgr.RegisterTask(tasks.NewRemittanceCycleTask(schedConfig.RemittanceCycleCron))
	schedulerMgr.RegisterTask(tasks.NewPaymentRetryTask(schedConfig.PaymentRetryCron))
	schedulerMgr.RegisterTask(tasks.NewReconciliationTask(schedConfig.ReconciliationCron))
	schedulerMgr.RegisterTask(tasks.NewNotificationRetryTask(schedConfig.NotificationRetryCron))

	return &UnifiedServer{
		config:     config,
		store:      store,
		cache:      cacheManager,
		restServer: restServer,
		scheduler:  schedulerMgr,
		sseManager: sseManager,
		wsManager:  wsManager,
	}, nil
}

func (s *UnifiedServer) Start() error {
	// Start real-time infrastructure
	s.sseManager.Start()
	s.wsManager.Start()

	// Start scheduler
	if err := s.scheduler.Start(); err != nil {
		log.Printf("Warning: Failed to start scheduler: %v", err)
	}

	// Start REST server (blocking)
	log.Printf("Starting REST server on %s", s.config.HTTPServerAddress)
	return s.restServer.Start(s.config.HTTPServerAddress)
}
