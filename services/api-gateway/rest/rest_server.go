package rest

import (
	"github.com/bitbiz/hias-core/configs"
	"github.com/bitbiz/hias-core/infrastructures/sse"
	"github.com/bitbiz/hias-core/infrastructures/websocket"
	"github.com/bitbiz/hias-core/services/api-gateway/rest/handlers"
	"github.com/bitbiz/hias-core/services/api-gateway/rest/middleware"
	"github.com/bitbiz/hias-core/services/api-gateway/rest/routes"
	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
)

// RestServerConfig holds all dependencies needed by the REST server.
type RestServerConfig struct {
	Config              configs.Config
	TokenMaker          auth.TokenMaker
	AuthHandler         *handlers.AuthHandler
	UserHandler         *handlers.UserHandler
	PlanHandler         *handlers.PlanHandler
	PolicyHandler       *handlers.PolicyHandler
	MemberHandler       *handlers.MemberHandler
	ProviderHandler     *handlers.ProviderHandler
	ClaimHandler        *handlers.ClaimHandler
	PreAuthHandler      *handlers.PreAuthHandler
	BillingHandler      *handlers.BillingHandler
	RemittanceHandler   *handlers.RemittanceHandler
	NotificationHandler *handlers.NotificationHandler
	AnalyticsHandler    *handlers.AnalyticsHandler
	AuditHandler        *handlers.AuditHandler
	SSEManager          sse.SSEManager
	WSManager           websocket.WebSocketManager
}

type RestServer struct {
	cfg    RestServerConfig
	router *gin.Engine
}

func NewRestServer(cfg RestServerConfig) *RestServer {
	server := &RestServer{cfg: cfg}
	server.setupRouter()
	return server
}

func (s *RestServer) setupRouter() {
	if s.cfg.Config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Global middleware
	router.Use(middleware.CORSMiddleware(s.cfg.Config.AllowedOrigins))
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.RateLimitMiddleware(100))
	router.Use(middleware.AuditMiddleware())

	// Health check
	router.GET("/health", func(ctx *gin.Context) {
		utils.RespondSuccess(ctx, 200, "OK", gin.H{"status": "healthy"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	tm := s.cfg.TokenMaker

	// Public routes
	routes.SetupAuthRoutes(v1, s.cfg.AuthHandler)

	// Protected routes — Identity
	routes.SetupUserRoutes(v1, s.cfg.UserHandler, tm)

	// Protected routes — Product
	routes.SetupPlanRoutes(v1, s.cfg.PlanHandler, tm)

	// Protected routes — Policy
	routes.SetupPolicyRoutes(v1, s.cfg.PolicyHandler, tm)
	routes.SetupMemberRoutes(v1, s.cfg.MemberHandler, tm)

	// Protected routes — Provider
	routes.SetupProviderRoutes(v1, s.cfg.ProviderHandler, tm)

	// Protected routes — Claims
	routes.SetupClaimRoutes(v1, s.cfg.ClaimHandler, tm)

	// Protected routes — Pre-Auth
	routes.SetupPreAuthRoutes(v1, s.cfg.PreAuthHandler, tm)

	// Protected routes — Billing (includes public webhook)
	routes.SetupBillingRoutes(v1, s.cfg.BillingHandler, tm)
	routes.SetupRemittanceRoutes(v1, s.cfg.RemittanceHandler, tm)

	// Protected routes — Notification, Analytics, Audit
	routes.SetupNotificationRoutes(v1, s.cfg.NotificationHandler, tm)
	routes.SetupAnalyticsRoutes(v1, s.cfg.AnalyticsHandler, tm)
	routes.SetupAuditRoutes(v1, s.cfg.AuditHandler, tm)

	s.router = router
}

func (s *RestServer) GetRouter() *gin.Engine {
	return s.router
}

func (s *RestServer) Start(address string) error {
	return s.router.Run(address)
}
