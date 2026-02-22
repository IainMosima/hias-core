package repository

import (
	"context"
	"fmt"
	"time"

	domainRepo "github.com/bitbiz/hias-core/domains/analytics/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
)

type analyticsRepository struct {
	store db.Store
}

func NewAnalyticsRepository(store db.Store) domainRepo.AnalyticsRepository {
	return &analyticsRepository{store: store}
}

func (r *analyticsRepository) GetClaimsVolume(ctx context.Context, start, end time.Time) (*domainRepo.ClaimsVolume, error) {
	result, err := r.store.GetClaimsVolume(ctx, db.GetClaimsVolumeParams{
		CreatedAt:   start,
		CreatedAt_2: end,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get claims volume: %w", err)
	}
	return &domainRepo.ClaimsVolume{
		TotalClaims:        result.TotalClaims,
		ApprovedClaims:     result.ApprovedClaims,
		RejectedClaims:     result.RejectedClaims,
		ManualReviewClaims: result.ManualReviewClaims,
		PaidClaims:         result.PaidClaims,
	}, nil
}

func (r *analyticsRepository) GetApprovalRate(ctx context.Context, start, end time.Time) (float64, error) {
	rate, err := r.store.GetApprovalRate(ctx, db.GetApprovalRateParams{
		CreatedAt:   start,
		CreatedAt_2: end,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get approval rate: %w", err)
	}
	return float64(rate), nil
}

func (r *analyticsRepository) GetAverageTAT(ctx context.Context, start, end time.Time) (float64, error) {
	tat, err := r.store.GetAverageTAT(ctx, db.GetAverageTATParams{
		CreatedAt:   start,
		CreatedAt_2: end,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get average turnaround time: %w", err)
	}
	// GetAverageTAT returns interface{}, attempt type assertion
	switch v := tat.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case string:
		return 0, nil
	default:
		return 0, nil
	}
}

func (r *analyticsRepository) GetLossRatio(ctx context.Context, start, end time.Time) (float64, error) {
	ratio, err := r.store.GetLossRatio(ctx, db.GetLossRatioParams{
		StartDate: start,
		EndDate:   end,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get loss ratio: %w", err)
	}
	return float64(ratio), nil
}

func (r *analyticsRepository) GetFraudRate(ctx context.Context, start, end time.Time) (float64, error) {
	rate, err := r.store.GetFraudRate(ctx, db.GetFraudRateParams{
		CreatedAt:   start,
		CreatedAt_2: end,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get fraud rate: %w", err)
	}
	return float64(rate), nil
}

func (r *analyticsRepository) GetTopProviders(ctx context.Context, start, end time.Time, limit int) ([]*domainRepo.TopProvider, error) {
	dbProviders, err := r.store.GetTopProviders(ctx, db.GetTopProvidersParams{
		CreatedAt:   start,
		CreatedAt_2: end,
		Limit:       int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get top providers: %w", err)
	}
	providers := make([]*domainRepo.TopProvider, len(dbProviders))
	for i, p := range dbProviders {
		providers[i] = &domainRepo.TopProvider{
			ID:            p.ID,
			Name:          p.Name,
			ClaimCount:    p.ClaimCount,
			TotalAmount:   p.TotalAmount,
			TotalApproved: p.TotalApproved,
		}
	}
	return providers, nil
}

func (r *analyticsRepository) GetTotalPremiumCollected(ctx context.Context, start, end time.Time) (int64, error) {
	amount, err := r.store.GetTotalPremiumCollected(ctx, db.GetTotalPremiumCollectedParams{
		CreatedAt:   start,
		CreatedAt_2: end,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get total premium collected: %w", err)
	}
	return amount, nil
}

func (r *analyticsRepository) GetTotalClaimsPaid(ctx context.Context, start, end time.Time) (int64, error) {
	amount, err := r.store.GetTotalClaimsPaid(ctx, db.GetTotalClaimsPaidParams{
		CreatedAt:   start,
		CreatedAt_2: end,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get total claims paid: %w", err)
	}
	return amount, nil
}
