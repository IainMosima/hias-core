package main

import (
	"context"
	"log"
	"time"

	"github.com/bitbiz/hias-core/configs"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/bitbiz/hias-core/infrastructures/documents"
	"github.com/bitbiz/hias-core/infrastructures/queue"
	reportingInfra "github.com/bitbiz/hias-core/infrastructures/reporting"
	"github.com/bitbiz/hias-core/infrastructures/repository"
	"github.com/bitbiz/hias-core/infrastructures/scheduler"
	schedulerTasks "github.com/bitbiz/hias-core/infrastructures/scheduler/tasks"
	"github.com/bitbiz/hias-core/services/analytics"
	"github.com/bitbiz/hias-core/services/api-gateway/handlers"
	"github.com/bitbiz/hias-core/services/api-gateway/routes"
	"github.com/bitbiz/hias-core/services/audit"
	awsSvc "github.com/bitbiz/hias-core/shared/aws"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bitbiz/hias-core/services/billing"
	"github.com/bitbiz/hias-core/services/claims"
	"github.com/bitbiz/hias-core/services/identity"
	"github.com/bitbiz/hias-core/services/notification"
	"github.com/bitbiz/hias-core/services/policy"
	"github.com/bitbiz/hias-core/services/preauth"
	"github.com/bitbiz/hias-core/services/product"
	"github.com/bitbiz/hias-core/services/provider"
	"github.com/bitbiz/hias-core/services/reinsurance"
	reportingSvc "github.com/bitbiz/hias-core/services/reporting"
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
	endorsementRepo := repository.NewEndorsementRepository(store)
	renewalRepo := repository.NewPolicyRenewalRepository(store)
	underwritingRepo := repository.NewUnderwritingRepository(store)
	policyDocumentRepo := repository.NewPolicyDocumentRepository(store)
	underwritingFlagRepo := repository.NewUnderwritingFlagRepository(store)
	underwritingRuleRepo := repository.NewUnderwritingRuleRepository(store)
	creditNoteRepo := repository.NewCreditNoteRepository(store)
	caseRecordRepo := repository.NewCaseRecordRepository(store)
	claimDocRepo := repository.NewClaimDocumentRepository(store)
	statementRepo := repository.NewProviderStatementRepository(store)

	// Reinsurance repositories
	treatyRepo := repository.NewTreatyRepository(store)
	treatyParticipantRepo := repository.NewTreatyParticipantRepository(store)
	treatyLayerRepo := repository.NewTreatyLayerRepository(store)
	profitCommissionRepo := repository.NewProfitCommissionRepository(store)
	reinsuranceCessionRepo := repository.NewCessionRepository(store)
	reinsuranceRecoveryRepo := repository.NewRecoveryRepository(store)
	recoveryWorkflowEventRepo := repository.NewRecoveryWorkflowEventRepository(store)
	bordereauRepo := repository.NewBordereauRepository(store)
	bordereauItemRepo := repository.NewBordereauItemRepository(store)
	reinsurerStatementRepo := repository.NewReinsurerStatementRepository(store)
	treatyAlertRepo := repository.NewTreatyAlertRepository(store)

	// Reporting repositories
	reportRepo := reportingInfra.NewReportRepository(store)
	reportDataRepo := reportingInfra.NewReportDataRepository(store)
	reportExporter := reportingInfra.NewReportExporter()

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
	premiumRuleSvc := product.NewPremiumRuleService(premiumRuleRepo, planRepo, auditSvc)
	providerNetworkSvc := product.NewProviderNetworkService(providerNetworkRepo, planRepo, auditSvc)

	// S3 service (optional — graceful degradation if unconfigured)
	var s3Svc awsSvc.S3Service
	if cfg.AWSS3Bucket != "" {
		awsCfg, awsErr := awsconfig.LoadDefaultConfig(context.TODO(), awsconfig.WithRegion(cfg.AWSS3Region))
		if awsErr != nil {
			log.Printf("Warning: Cannot load AWS config for S3: %v", awsErr)
		} else {
			s3Svc = awsSvc.NewS3Service(s3.NewFromConfig(awsCfg), cfg.AWSS3Bucket, cfg.AWSS3CDNDomain)
		}
	}

	// PDF generator and policy document service (created before policySvc to avoid circular deps)
	pdfGenerator := documents.NewPDFGenerator()
	policyDocSvc := policy.NewPolicyDocumentService(policyDocumentRepo, policyRepo, memberRepo, planRepo, benefitRepo, renewalRepo, preauthRepo, providerRepo, pdfGenerator, s3Svc, auditSvc, notifSvc)

	// Credit note service (created before policy services since they depend on it)
	creditNoteSvc := billing.NewCreditNoteService(creditNoteRepo, invoiceRepo, auditSvc)

	// Underwriting flag and rule services
	underwritingFlagSvc := policy.NewUnderwritingFlagService(underwritingFlagRepo, auditSvc)
	underwritingRuleSvc := policy.NewUnderwritingRuleService(underwritingRuleRepo, planRepo, auditSvc)

	// Policy services
	memberSvc := policy.NewMemberService(memberRepo, policyRepo, planRepo, premiumRuleRepo, premiumRuleSvc, underwritingFlagRepo, underwritingRuleRepo, creditNoteSvc, auditSvc)
	policySvc := policy.NewPolicyService(policyRepo, planRepo, memberRepo, premiumRuleSvc, policyDocSvc, creditNoteSvc, auditSvc)
	endorsementSvc := policy.NewEndorsementService(endorsementRepo, policyRepo, memberSvc, policySvc, auditSvc)
	renewalSvc := policy.NewRenewalService(renewalRepo, policyRepo, memberRepo, claimRepo, premiumRuleSvc, premiumRuleRepo, planRepo, underwritingFlagRepo, auditSvc)
	underwritingSvc := policy.NewUnderwritingService(underwritingRepo, policyRepo, memberRepo, underwritingRuleRepo, underwritingFlagRepo, auditSvc)

	// Provider services
	providerSvc := provider.NewProviderService(providerRepo, contractRepo, rateCardRepo, auditSvc)

	// Claims services
	fraudSvc := claims.NewFraudService(fraudFlagRepo, contractRepo, rateCardRepo, providerRepo)
	validatorSvc := claims.NewValidatorService(policyRepo, memberRepo, providerRepo)
	adjudicatorSvc := claims.NewAdjudicatorService(claimRepo, benefitRepo, exclusionRepo, policyRepo, memberRepo, providerRepo, providerNetworkRepo, fraudSvc, contractRepo, preauthRepo)
	claimSvc := claims.NewClaimService(claimRepo, lineItemRepo, adjudicatorSvc, validatorSvc, fraudSvc, adjRepo, fraudFlagRepo, claimDocRepo, preauthRepo, auditSvc)

	// Case management service
	caseSvc := claims.NewCaseService(caseRecordRepo, preauthRepo, auditSvc)

	// Pre-auth service
	preauthSvc := preauth.NewPreAuthService(preauthRepo)

	// Billing services
	billingSvc := billing.NewBillingService(invoiceRepo, policyRepo)
	paymentSvc := billing.NewPaymentService(paymentRepo, invoiceRepo)
	remittanceSvc := billing.NewRemittanceService(remittanceRepo, claimRepo, providerRepo)
	installmentSvc := billing.NewInstallmentService(installmentScheduleRepo, installmentRepo, policyRepo)
	statementSvc := billing.NewStatementService(statementRepo, claimRepo, auditSvc)

	// Analytics service
	analyticsSvc := analytics.NewAnalyticsService(analyticsRepo)

	// Sales services
	leadSvc := sales.NewLeadService(leadRepo, leadActivityRepo, auditSvc)
	quotationSvc := sales.NewQuotationService(quotationRepo, quotationVersionRepo, quotationDocumentRepo, approvalLimitRepo, leadRepo, auditSvc, premiumRuleSvc, notifSvc, policySvc, memberSvc, installmentSvc)
	approvalLimitSvc := sales.NewApprovalLimitService(approvalLimitRepo, auditSvc)

	// Reinsurance services
	treatySvc := reinsurance.NewTreatyService(treatyRepo, treatyParticipantRepo, treatyLayerRepo, profitCommissionRepo, auditSvc)
	cessionSvc := reinsurance.NewCessionService(reinsuranceCessionRepo, treatyRepo, treatyParticipantRepo, auditSvc)
	recoverySvc := reinsurance.NewRecoveryService(reinsuranceRecoveryRepo, recoveryWorkflowEventRepo, treatyRepo, treatyLayerRepo, treatyParticipantRepo, reinsuranceCessionRepo, auditSvc)
	bordereauSvc := reinsurance.NewBordereauService(bordereauRepo, bordereauItemRepo, reinsuranceCessionRepo, reinsuranceRecoveryRepo, auditSvc)
	reinsurerStatementSvc := reinsurance.NewReinsurerStatementService(reinsurerStatementRepo, reinsuranceCessionRepo, reinsuranceRecoveryRepo, treatyParticipantRepo, profitCommissionRepo, auditSvc)
	treatyAlertSvc := reinsurance.NewTreatyAlertService(treatyAlertRepo, treatyLayerRepo, reinsuranceRecoveryRepo, treatyRepo)

	// Reporting service
	reportSvc := reportingSvc.NewReportService(reportRepo, reportDataRepo, reportExporter, notifSvc, auditSvc, analyticsRepo)

	// Suppress unused variable warnings for services used internally
	_ = billingSvc

	// 7. Handlers
	h := routes.Handlers{
		Health:           handlers.NewHealthHandler(),
		Auth:             handlers.NewAuthHandler(authSvc),
		User:             handlers.NewUserHandler(userSvc),
		Plan:             handlers.NewPlanHandler(planSvc),
		Benefit:          handlers.NewBenefitHandler(benefitSvc),
		Exclusion:        handlers.NewExclusionHandler(exclusionSvc),
		PremiumRule:      handlers.NewPremiumRuleHandler(premiumRuleSvc),
		ProviderNetwork:  handlers.NewProviderNetworkHandler(providerNetworkSvc),
		Policy:           handlers.NewPolicyHandler(policySvc),
		Member:           handlers.NewMemberHandler(memberSvc),
		Provider:         handlers.NewProviderHandler(providerSvc),
		Contract:         handlers.NewContractHandler(contractRepo),
		RateCard:         handlers.NewRateCardHandler(rateCardRepo),
		Claim:            handlers.NewClaimHandler(claimSvc),
		PreAuth:          handlers.NewPreAuthHandler(preauthSvc),
		Invoice:          handlers.NewInvoiceHandler(invoiceRepo),
		Payment:          handlers.NewPaymentHandler(paymentSvc),
		Remittance:       handlers.NewRemittanceHandler(remittanceSvc),
		Installment:      handlers.NewInstallmentHandler(installmentSvc),
		Notification:     handlers.NewNotificationHandler(notifSvc),
		Audit:            handlers.NewAuditHandler(auditSvc),
		Analytics:        handlers.NewAnalyticsHandler(analyticsSvc),
		Lead:             handlers.NewLeadHandler(leadSvc),
		Quotation:        handlers.NewQuotationHandler(quotationSvc),
		ApprovalLimit:    handlers.NewApprovalLimitHandler(approvalLimitSvc),
		Endorsement:      handlers.NewEndorsementHandler(endorsementSvc),
		Renewal:          handlers.NewRenewalHandler(renewalSvc),
		Underwriting:     handlers.NewUnderwritingHandler(underwritingSvc),
		PolicyDocument:   handlers.NewPolicyDocumentHandler(policyDocSvc, claimSvc),
		UnderwritingFlag: handlers.NewUnderwritingFlagHandler(underwritingFlagSvc),
		UnderwritingRule: handlers.NewUnderwritingRuleHandler(underwritingRuleSvc),
		CreditNote:       handlers.NewCreditNoteHandler(creditNoteSvc),
		Case:             handlers.NewCaseHandler(caseSvc),
		Statement:        handlers.NewStatementHandler(statementSvc),

		// Reinsurance
		Treaty:               handlers.NewTreatyHandler(treatySvc),
		Cession:              handlers.NewCessionHandler(cessionSvc),
		Recovery:             handlers.NewRecoveryHandler(recoverySvc),
		Bordereau:            handlers.NewBordereauHandler(bordereauSvc),
		ReinsurerStatement:   handlers.NewReinsurerStatementHandler(reinsurerStatementSvc),
		TreatyAlert:          handlers.NewTreatyAlertHandler(treatyAlertSvc),
		ReinsuranceAnalytics: handlers.NewReinsuranceAnalyticsHandler(treatySvc, cessionSvc, recoverySvc, treatyAlertSvc),

		// Reporting
		Report: handlers.NewReportHandler(reportSvc),
	}

	// 8. Scheduler (optional — only if enabled)
	if cfg.SchedulerEnabled {
		schedulerMgr := scheduler.NewSchedulerManager()

		reportDistSchedule := cfg.ReportDistributionSchedule
		if reportDistSchedule == "" {
			reportDistSchedule = "*/5 * * * *"
		}
		reportCleanupSchedule := cfg.ReportCleanupSchedule
		if reportCleanupSchedule == "" {
			reportCleanupSchedule = "0 2 * * *"
		}

		reportDistTask := schedulerTasks.NewReportDistributionTask(reportDistSchedule, reportSvc)
		reportCleanupTask := schedulerTasks.NewReportCleanupTask(reportCleanupSchedule, reportRepo)

		if err := schedulerMgr.RegisterTask(reportDistTask); err != nil {
			log.Printf("Warning: Failed to register report distribution task: %v", err)
		}
		if err := schedulerMgr.RegisterTask(reportCleanupTask); err != nil {
			log.Printf("Warning: Failed to register report cleanup task: %v", err)
		}

		if err := schedulerMgr.Start(); err != nil {
			log.Printf("Warning: Failed to start scheduler: %v", err)
		} else {
			log.Printf("Scheduler started with %d tasks", len(schedulerMgr.GetRegisteredTasks()))
		}
	}

	// 9. Server
	server := NewServer(tokenMaker, cfg.AllowedOrigins)
	server.RegisterRoutes(h)

	// 10. Start
	if err := server.Start(cfg.HTTPServerAddress); err != nil {
		log.Fatalf("Cannot start server: %v", err)
	}
}
