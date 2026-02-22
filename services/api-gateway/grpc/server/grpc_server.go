package server

import (
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// GRPCServer wraps the gRPC server with health checks and graceful shutdown.
type GRPCServer struct {
	server       *grpc.Server
	healthServer *health.Server
	address      string
	environment  string
}

// NewGRPCServer creates a new gRPC server.
func NewGRPCServer(address, environment string) *GRPCServer {
	opts := []grpc.ServerOption{
		grpc.ConnectionTimeout(30 * time.Second),
		grpc.MaxRecvMsgSize(4 * 1024 * 1024), // 4MB
		grpc.MaxSendMsgSize(4 * 1024 * 1024), // 4MB
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     5 * time.Minute,
			MaxConnectionAge:      30 * time.Minute,
			MaxConnectionAgeGrace: 5 * time.Second,
			Time:                  10 * time.Second,
			Timeout:               3 * time.Second,
		}),
	}

	server := grpc.NewServer(opts...)
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(server, healthServer)

	if environment != "production" {
		reflection.Register(server)
	}

	return &GRPCServer{
		server:       server,
		healthServer: healthServer,
		address:      address,
		environment:  environment,
	}
}

// GetServer returns the underlying gRPC server for service registration.
func (s *GRPCServer) GetServer() *grpc.Server {
	return s.server
}

// Start begins listening for gRPC connections.
func (s *GRPCServer) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.address, err)
	}

	s.healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	log.Printf("gRPC server listening on %s", s.address)
	return s.server.Serve(listener)
}

// Stop gracefully shuts down the gRPC server.
func (s *GRPCServer) Stop() {
	s.healthServer.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
	s.server.GracefulStop()
	log.Println("gRPC server stopped")
}
