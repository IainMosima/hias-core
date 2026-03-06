package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ReportDataRepository interface {
	GetClaimsExperienceData(ctx context.Context, start, end time.Time, policyID *uuid.UUID) ([]map[string]interface{}, error)
	GetClaimsRegisterData(ctx context.Context, start, end time.Time, status string, limit, offset int) ([]map[string]interface{}, error)
	GetPremiumDebtorsAgeingData(ctx context.Context) ([]map[string]interface{}, error)
	GetPremiumRegisterData(ctx context.Context, start, end time.Time, limit, offset int) ([]map[string]interface{}, error)
	GetMembershipData(ctx context.Context, policyID *uuid.UUID, status string, limit, offset int) ([]map[string]interface{}, error)
	GetProviderPerformanceData(ctx context.Context, start, end time.Time) ([]map[string]interface{}, error)
	GetLossRatioData(ctx context.Context, start, end time.Time) ([]map[string]interface{}, error)
	GetRenewalData(ctx context.Context, start, end time.Time) ([]map[string]interface{}, error)
	DrillDownClaimsByPolicy(ctx context.Context, policyID uuid.UUID, start, end time.Time) ([]map[string]interface{}, error)
	DrillDownPaymentsByPolicy(ctx context.Context, policyID uuid.UUID, start, end time.Time) ([]map[string]interface{}, error)
	GetOutstandingPremium(ctx context.Context) (int64, error)
	GetSLABreachCount(ctx context.Context) (int64, error)
}
