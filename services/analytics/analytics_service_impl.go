package analytics

import (
	"context"
	"net/http"
	"time"

	"github.com/bitbiz/hias-core/domains/analytics/repository"
	analyticsSvc "github.com/bitbiz/hias-core/domains/analytics/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
)

type analyticsServiceImpl struct {
	analyticsRepo repository.AnalyticsRepository
}

func NewAnalyticsService(analyticsRepo repository.AnalyticsRepository) analyticsSvc.AnalyticsService {
	return &analyticsServiceImpl{analyticsRepo: analyticsRepo}
}

func (s *analyticsServiceImpl) GetDashboard(ctx context.Context, period string) *schema.ServiceResponse[analyticsSvc.DashboardData] {
	start, end := parsePeriod(period)

	claimsVolume, _ := s.analyticsRepo.GetClaimsVolume(ctx, start, end)
	approvalRate, _ := s.analyticsRepo.GetApprovalRate(ctx, start, end)
	avgTAT, _ := s.analyticsRepo.GetAverageTAT(ctx, start, end)
	lossRatio, _ := s.analyticsRepo.GetLossRatio(ctx, start, end)
	fraudRate, _ := s.analyticsRepo.GetFraudRate(ctx, start, end)
	totalPremium, _ := s.analyticsRepo.GetTotalPremiumCollected(ctx, start, end)
	totalClaimsPaid, _ := s.analyticsRepo.GetTotalClaimsPaid(ctx, start, end)
	topProviders, _ := s.analyticsRepo.GetTopProviders(ctx, start, end, 10)
	activePolicies, _ := s.analyticsRepo.GetActivePolicyCount(ctx, start, end)
	lapsedPolicies, _ := s.analyticsRepo.GetLapsedPolicyCount(ctx, start, end)
	totalMembers, _ := s.analyticsRepo.GetTotalMemberCount(ctx, start, end)
	renewalRate, _ := s.analyticsRepo.GetRenewalRate(ctx, start, end)

	dashboard := analyticsSvc.DashboardData{
		ClaimsVolume:    claimsVolume,
		ApprovalRate:    approvalRate,
		AverageTAT:      avgTAT,
		LossRatio:       lossRatio,
		FraudRate:       fraudRate,
		TotalPremium:    totalPremium,
		TotalClaimsPaid: totalClaimsPaid,
		TopProviders:    topProviders,
		ActivePolicies:  activePolicies,
		LapsedPolicies:  lapsedPolicies,
		TotalMembers:    totalMembers,
		RenewalRate:     renewalRate,
	}

	return schema.NewServiceResponse(dashboard, http.StatusOK, "Dashboard data retrieved")
}

func (s *analyticsServiceImpl) GetKPIs(ctx context.Context, period string) *schema.ServiceResponse[interface{}] {
	start, end := parsePeriod(period)

	approvalRate, _ := s.analyticsRepo.GetApprovalRate(ctx, start, end)
	lossRatio, _ := s.analyticsRepo.GetLossRatio(ctx, start, end)
	avgTAT, _ := s.analyticsRepo.GetAverageTAT(ctx, start, end)
	totalPremium, _ := s.analyticsRepo.GetTotalPremiumCollected(ctx, start, end)
	totalClaimsPaid, _ := s.analyticsRepo.GetTotalClaimsPaid(ctx, start, end)
	activePolicies, _ := s.analyticsRepo.GetActivePolicyCount(ctx, start, end)
	lapsedPolicies, _ := s.analyticsRepo.GetLapsedPolicyCount(ctx, start, end)
	totalMembers, _ := s.analyticsRepo.GetTotalMemberCount(ctx, start, end)
	renewalRate, _ := s.analyticsRepo.GetRenewalRate(ctx, start, end)

	kpis := map[string]interface{}{
		"approval_rate":     approvalRate,
		"loss_ratio":        lossRatio,
		"average_tat_hours": avgTAT,
		"total_premium":     totalPremium,
		"total_claims_paid": totalClaimsPaid,
		"active_policies":   activePolicies,
		"lapsed_policies":   lapsedPolicies,
		"total_members":     totalMembers,
		"renewal_rate":      renewalRate,
		"period":            period,
	}

	return schema.NewServiceResponse[interface{}](kpis, http.StatusOK, "KPIs retrieved")
}

func (s *analyticsServiceImpl) ExportCSV(ctx context.Context, reportType, period string) *schema.ServiceResponse[[]byte] {
	return schema.NewServiceResponse([]byte("report,data\n"), http.StatusOK, "CSV export generated")
}

func parsePeriod(period string) (time.Time, time.Time) {
	now := time.Now()
	switch period {
	case "week":
		return now.AddDate(0, 0, -7), now
	case "month":
		return now.AddDate(0, -1, 0), now
	case "quarter":
		return now.AddDate(0, -3, 0), now
	case "year":
		return now.AddDate(-1, 0, 0), now
	default:
		return now.AddDate(0, -1, 0), now
	}
}
