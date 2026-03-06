package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	"github.com/bitbiz/hias-core/services/api-gateway/middleware"
	"github.com/bitbiz/hias-core/services/api-gateway/routes"
	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router     *gin.Engine
	httpServer *http.Server
	tokenMaker auth.TokenMaker
}

func NewServer(tokenMaker auth.TokenMaker, allowedOrigins string, auditSvc auditService.AuditService) *Server {
	router := gin.Default()

	// Global middleware
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.CORSMiddleware(allowedOrigins))
	router.Use(middleware.RateLimiterMiddleware(120, 60))
	router.Use(middleware.MetricsMiddleware())
	router.Use(middleware.AuditMiddleware(auditSvc))

	return &Server{
		router:     router,
		tokenMaker: tokenMaker,
	}
}

func (s *Server) RegisterRoutes(h routes.Handlers) {
	routes.RegisterRoutes(s.router, h, s.tokenMaker)
}

func (s *Server) Start(address string) error {
	s.httpServer = &http.Server{
		Addr:         address,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		log.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}

		log.Println("Server exited gracefully")
	}()

	log.Printf("Server starting on %s", address)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Router() *gin.Engine {
	return s.router
}
