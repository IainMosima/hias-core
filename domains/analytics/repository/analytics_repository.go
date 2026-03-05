package repository

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type AnalyticsRepository interface {
	GetClaimsVolume(ctx context.Context, start, end time.Time) (*ClaimsVolume, error)
	GetApprovalRate(ctx context.Context, start, end time.Time) (float64, error)
	GetAverageTAT(ctx context.Context, start, end time.Time) (float64, error)
	GetLossRatio(ctx context.Context, start, end time.Time) (float64, error)
	GetFraudRate(ctx context.Context, start, end time.Time) (float64, error)
	GetTopProviders(ctx context.Context, start, end time.Time, limit int) ([]*TopProvider, error)
	GetTotalPremiumCollected(ctx context.Context, start, end time.Time) (int64, error)
	GetTotalClaimsPaid(ctx context.Context, start, end time.Time) (int64, error)
	GetActivePolicyCount(ctx context.Context, start, end time.Time) (int64, error)
	GetLapsedPolicyCount(ctx context.Context, start, end time.Time) (int64, error)
	GetTotalMemberCount(ctx context.Context, start, end time.Time) (int64, error)
	GetRenewalRate(ctx context.Context, start, end time.Time) (float64, error)
}

type ClaimsVolume struct {
	TotalClaims        int64 `json:"total_claims"`
	ApprovedClaims     int64 `json:"approved_claims"`
	RejectedClaims     int64 `json:"rejected_claims"`
	ManualReviewClaims int64 `json:"manual_review_claims"`
	PaidClaims         int64 `json:"paid_claims"`
}

type TopProvider struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	ClaimCount    int64     `json:"claim_count"`
	TotalAmount   int64     `json:"total_amount"`
	TotalApproved int64     `json:"total_approved"`
}
