package server

import (
	"context"
	"log"

	claimsService "github.com/bitbiz/hias-core/domains/claims/service"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ClaimsGRPCServer implements the gRPC ClaimService.
// In production, this would implement the generated protobuf interface.
// For now it demonstrates the pattern.
type ClaimsGRPCServer struct {
	claimService claimsService.ClaimService
}

// NewClaimsGRPCServer creates a new gRPC claims server.
func NewClaimsGRPCServer(claimService claimsService.ClaimService) *ClaimsGRPCServer {
	return &ClaimsGRPCServer{
		claimService: claimService,
	}
}

// GetClaimStatus returns the current status of a claim via gRPC.
func (s *ClaimsGRPCServer) GetClaimStatus(ctx context.Context, claimID string) (interface{}, error) {
	id, err := uuid.Parse(claimID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid claim ID: %v", err)
	}

	resp := s.claimService.GetClaim(ctx, id)
	if resp.Error != nil {
		log.Printf("gRPC GetClaimStatus error: %v", resp.Error)
		return nil, status.Errorf(codes.NotFound, "%s", resp.Message)
	}

	return resp.Data, nil
}

// ListClaims lists claims for a provider via gRPC.
func (s *ClaimsGRPCServer) ListClaims(ctx context.Context, providerID string, page, pageSize int) (interface{}, error) {
	if providerID != "" {
		pid, err := uuid.Parse(providerID)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid provider ID: %v", err)
		}
		resp := s.claimService.ListClaimsByProvider(ctx, pid, page, pageSize)
		if resp.Error != nil {
			return nil, status.Errorf(codes.Internal, "%s", resp.Message)
		}
		return resp.Data, nil
	}

	resp := s.claimService.ListClaims(ctx, page, pageSize)
	if resp.Error != nil {
		return nil, status.Errorf(codes.Internal, "%s", resp.Message)
	}
	return resp.Data, nil
}
