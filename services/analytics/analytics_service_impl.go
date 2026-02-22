package analytics

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bitbiz/hias-core/domains/analytics/repository"
	analyticsSchema "github.com/bitbiz/hias-core/domains/analytics/schema"
	"github.com/bitbiz/hias-core/domains/analytics/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
)

type analyticsServiceImpl struct {
	analyticsRepo repository.AnalyticsRepository
}

func NewAnalyticsService(
	analyticsRepo repository.AnalyticsRepository,
) service.AnalyticsService {
	return &analyticsServiceImpl{
		analyticsRepo: analyticsRepo,
	}
}

func (s *analyticsServiceImpl) GetDashboard(ctx context.Context, period string) *schema.ServiceResponse[service.DashboardData] {
	start, end := parsePeriod(period)

	claimsVolume, err := s.analyticsRepo.GetClaimsVolume(ctx, start, end)
	if err != nil {
		return schema.NewServiceErrorResponse[service.DashboardData](http.StatusInternalServerError, "Failed to get claims volume", err)
	}

	approvalRate, err := s.analyticsRepo.GetApprovalRate(ctx, start, end)
	if err != nil {
		return schema.NewServiceErrorResponse[service.DashboardData](http.StatusInternalServerError, "Failed to get approval rate", err)
	}

	avgTAT, err := s.analyticsRepo.GetAverageTAT(ctx, start, end)
	if err != nil {
		return schema.NewServiceErrorResponse[service.DashboardData](http.StatusInternalServerError, "Failed to get average TAT", err)
	}

	lossRatio, err := s.analyticsRepo.GetLossRatio(ctx, start, end)
	if err != nil {
		return schema.NewServiceErrorResponse[service.DashboardData](http.StatusInternalServerError, "Failed to get loss ratio", err)
	}

	fraudRate, err := s.analyticsRepo.GetFraudRate(ctx, start, end)
	if err != nil {
		return schema.NewServiceErrorResponse[service.DashboardData](http.StatusInternalServerError, "Failed to get fraud rate", err)
	}

	topProviders, err := s.analyticsRepo.GetTopProviders(ctx, start, end, 10)
	if err != nil {
		return schema.NewServiceErrorResponse[service.DashboardData](http.StatusInternalServerError, "Failed to get top providers", err)
	}

	totalPremium, err := s.analyticsRepo.GetTotalPremiumCollected(ctx, start, end)
	if err != nil {
		return schema.NewServiceErrorResponse[service.DashboardData](http.StatusInternalServerError, "Failed to get total premium collected", err)
	}

	totalClaimsPaid, err := s.analyticsRepo.GetTotalClaimsPaid(ctx, start, end)
	if err != nil {
		return schema.NewServiceErrorResponse[service.DashboardData](http.StatusInternalServerError, "Failed to get total claims paid", err)
	}

	dashboard := service.DashboardData{
		ClaimsVolume:    analyticsSchema.ToClaimsVolumeResponse(claimsVolume),
		ApprovalRate:    approvalRate,
		AverageTAT:     avgTAT,
		LossRatio:       lossRatio,
		FraudRate:       fraudRate,
		TotalPremium:    totalPremium,
		TotalClaimsPaid: totalClaimsPaid,
		TopProviders:    analyticsSchema.ToTopProviderResponseList(topProviders),
	}

	return schema.NewServiceResponse(dashboard, http.StatusOK, "Dashboard data retrieved")
}

func (s *analyticsServiceImpl) GetKPIs(ctx context.Context, period string) *schema.ServiceResponse[interface{}] {
	start, end := parsePeriod(period)

	approvalRate, err := s.analyticsRepo.GetApprovalRate(ctx, start, end)
	if err != nil {
		return schema.NewServiceErrorResponse[interface{}](http.StatusInternalServerError, "Failed to get approval rate", err)
	}

	avgTAT, err := s.analyticsRepo.GetAverageTAT(ctx, start, end)
	if err != nil {
		return schema.NewServiceErrorResponse[interface{}](http.StatusInternalServerError, "Failed to get average TAT", err)
	}

	lossRatio, err := s.analyticsRepo.GetLossRatio(ctx, start, end)
	if err != nil {
		return schema.NewServiceErrorResponse[interface{}](http.StatusInternalServerError, "Failed to get loss ratio", err)
	}

	fraudRate, err := s.analyticsRepo.GetFraudRate(ctx, start, end)
	if err != nil {
		return schema.NewServiceErrorResponse[interface{}](http.StatusInternalServerError, "Failed to get fraud rate", err)
	}

	totalPremium, err := s.analyticsRepo.GetTotalPremiumCollected(ctx, start, end)
	if err != nil {
		return schema.NewServiceErrorResponse[interface{}](http.StatusInternalServerError, "Failed to get total premium collected", err)
	}

	totalClaimsPaid, err := s.analyticsRepo.GetTotalClaimsPaid(ctx, start, end)
	if err != nil {
		return schema.NewServiceErrorResponse[interface{}](http.StatusInternalServerError, "Failed to get total claims paid", err)
	}

	kpis := analyticsSchema.KPIResponse{
		ApprovalRate:    approvalRate,
		AverageTAT:     avgTAT,
		LossRatio:       lossRatio,
		FraudRate:       fraudRate,
		TotalPremium:    totalPremium,
		TotalClaimsPaid: totalClaimsPaid,
	}

	return schema.NewServiceResponse[interface{}](kpis, http.StatusOK, "KPIs retrieved")
}

func (s *analyticsServiceImpl) ExportCSV(ctx context.Context, reportType, period string) *schema.ServiceResponse[[]byte] {
	start, end := parsePeriod(period)

	var buf strings.Builder
	writer := csv.NewWriter(&buf)

	switch reportType {
	case "claims":
		claimsVolume, err := s.analyticsRepo.GetClaimsVolume(ctx, start, end)
		if err != nil {
			return schema.NewServiceErrorResponse[[]byte](http.StatusInternalServerError, "Failed to export claims data", err)
		}

		_ = writer.Write([]string{"Metric", "Value"})
		_ = writer.Write([]string{"Total Claims", fmt.Sprintf("%d", claimsVolume.TotalClaims)})
		_ = writer.Write([]string{"Approved Claims", fmt.Sprintf("%d", claimsVolume.ApprovedClaims)})
		_ = writer.Write([]string{"Rejected Claims", fmt.Sprintf("%d", claimsVolume.RejectedClaims)})
		_ = writer.Write([]string{"Manual Review Claims", fmt.Sprintf("%d", claimsVolume.ManualReviewClaims)})
		_ = writer.Write([]string{"Paid Claims", fmt.Sprintf("%d", claimsVolume.PaidClaims)})

	case "providers":
		topProviders, err := s.analyticsRepo.GetTopProviders(ctx, start, end, 50)
		if err != nil {
			return schema.NewServiceErrorResponse[[]byte](http.StatusInternalServerError, "Failed to export providers data", err)
		}

		_ = writer.Write([]string{"Provider ID", "Provider Name", "Claim Count", "Total Amount (KES)", "Total Approved (KES)"})
		for _, p := range topProviders {
			_ = writer.Write([]string{
				p.ID.String(),
				p.Name,
				fmt.Sprintf("%d", p.ClaimCount),
				fmt.Sprintf("%.2f", float64(p.TotalAmount)/100),
				fmt.Sprintf("%.2f", float64(p.TotalApproved)/100),
			})
		}

	case "kpis":
		approvalRate, _ := s.analyticsRepo.GetApprovalRate(ctx, start, end)
		avgTAT, _ := s.analyticsRepo.GetAverageTAT(ctx, start, end)
		lossRatio, _ := s.analyticsRepo.GetLossRatio(ctx, start, end)
		fraudRate, _ := s.analyticsRepo.GetFraudRate(ctx, start, end)
		totalPremium, _ := s.analyticsRepo.GetTotalPremiumCollected(ctx, start, end)
		totalClaimsPaid, _ := s.analyticsRepo.GetTotalClaimsPaid(ctx, start, end)

		_ = writer.Write([]string{"KPI", "Value"})
		_ = writer.Write([]string{"Approval Rate (%)", fmt.Sprintf("%.2f", approvalRate)})
		_ = writer.Write([]string{"Average TAT (hours)", fmt.Sprintf("%.2f", avgTAT)})
		_ = writer.Write([]string{"Loss Ratio", fmt.Sprintf("%.4f", lossRatio)})
		_ = writer.Write([]string{"Fraud Rate (%)", fmt.Sprintf("%.2f", fraudRate)})
		_ = writer.Write([]string{"Total Premium Collected (KES)", fmt.Sprintf("%.2f", float64(totalPremium)/100)})
		_ = writer.Write([]string{"Total Claims Paid (KES)", fmt.Sprintf("%.2f", float64(totalClaimsPaid)/100)})

	default:
		return schema.NewServiceErrorResponse[[]byte](http.StatusBadRequest, "Invalid report type. Use: claims, providers, kpis", nil)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return schema.NewServiceErrorResponse[[]byte](http.StatusInternalServerError, "Failed to write CSV", err)
	}

	return schema.NewServiceResponse([]byte(buf.String()), http.StatusOK, "CSV exported successfully")
}

// parsePeriod converts a period string (e.g., "30d", "7d", "1y", "mtd", "ytd") into a start and end time.
func parsePeriod(period string) (time.Time, time.Time) {
	now := time.Now().UTC()
	end := now

	switch period {
	case "7d":
		return now.AddDate(0, 0, -7), end
	case "30d":
		return now.AddDate(0, 0, -30), end
	case "90d":
		return now.AddDate(0, 0, -90), end
	case "1y":
		return now.AddDate(-1, 0, 0), end
	case "mtd":
		return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC), end
	case "ytd":
		return time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC), end
	default:
		// Default to last 30 days
		return now.AddDate(0, 0, -30), end
	}
}
