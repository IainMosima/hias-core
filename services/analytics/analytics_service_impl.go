package analytics

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bitbiz/hias-core/domains/analytics/repository"
	analyticsSvc "github.com/bitbiz/hias-core/domains/analytics/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	reportDataRepo "github.com/bitbiz/hias-core/domains/reporting/repository"
	reportingInfra "github.com/bitbiz/hias-core/infrastructures/reporting"
)

type analyticsServiceImpl struct {
	analyticsRepo  repository.AnalyticsRepository
	reportDataRepo reportDataRepo.ReportDataRepository
	exporter       reportingInfra.ReportExporter
}

func NewAnalyticsService(analyticsRepo repository.AnalyticsRepository, reportDataRepo reportDataRepo.ReportDataRepository, exporter reportingInfra.ReportExporter) analyticsSvc.AnalyticsService {
	return &analyticsServiceImpl{analyticsRepo: analyticsRepo, reportDataRepo: reportDataRepo, exporter: exporter}
}

func (s *analyticsServiceImpl) GetDashboard(ctx context.Context, period string) *schema.ServiceResponse[analyticsSvc.DashboardData] {
	start, end := parsePeriod(period)

	claimsVolume, err := s.analyticsRepo.GetClaimsVolume(ctx, start, end)
	if err != nil {
		log.Printf("Dashboard: failed to get claims volume: %v", err)
	}
	approvalRate, err := s.analyticsRepo.GetApprovalRate(ctx, start, end)
	if err != nil {
		log.Printf("Dashboard: failed to get approval rate: %v", err)
	}
	avgTAT, err := s.analyticsRepo.GetAverageTAT(ctx, start, end)
	if err != nil {
		log.Printf("Dashboard: failed to get average TAT: %v", err)
	}
	lossRatio, err := s.analyticsRepo.GetLossRatio(ctx, start, end)
	if err != nil {
		log.Printf("Dashboard: failed to get loss ratio: %v", err)
	}
	fraudRate, err := s.analyticsRepo.GetFraudRate(ctx, start, end)
	if err != nil {
		log.Printf("Dashboard: failed to get fraud rate: %v", err)
	}
	totalPremium, err := s.analyticsRepo.GetTotalPremiumCollected(ctx, start, end)
	if err != nil {
		log.Printf("Dashboard: failed to get total premium: %v", err)
	}
	totalClaimsPaid, err := s.analyticsRepo.GetTotalClaimsPaid(ctx, start, end)
	if err != nil {
		log.Printf("Dashboard: failed to get total claims paid: %v", err)
	}
	topProviders, err := s.analyticsRepo.GetTopProviders(ctx, start, end, 10)
	if err != nil {
		log.Printf("Dashboard: failed to get top providers: %v", err)
	}
	activePolicies, err := s.analyticsRepo.GetActivePolicyCount(ctx, start, end)
	if err != nil {
		log.Printf("Dashboard: failed to get active policies: %v", err)
	}
	lapsedPolicies, err := s.analyticsRepo.GetLapsedPolicyCount(ctx, start, end)
	if err != nil {
		log.Printf("Dashboard: failed to get lapsed policies: %v", err)
	}
	totalMembers, err := s.analyticsRepo.GetTotalMemberCount(ctx, start, end)
	if err != nil {
		log.Printf("Dashboard: failed to get total members: %v", err)
	}
	renewalRate, err := s.analyticsRepo.GetRenewalRate(ctx, start, end)
	if err != nil {
		log.Printf("Dashboard: failed to get renewal rate: %v", err)
	}

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

	approvalRate, err := s.analyticsRepo.GetApprovalRate(ctx, start, end)
	if err != nil {
		log.Printf("KPIs: failed to get approval rate: %v", err)
	}
	lossRatio, err := s.analyticsRepo.GetLossRatio(ctx, start, end)
	if err != nil {
		log.Printf("KPIs: failed to get loss ratio: %v", err)
	}
	avgTAT, err := s.analyticsRepo.GetAverageTAT(ctx, start, end)
	if err != nil {
		log.Printf("KPIs: failed to get average TAT: %v", err)
	}
	totalPremium, err := s.analyticsRepo.GetTotalPremiumCollected(ctx, start, end)
	if err != nil {
		log.Printf("KPIs: failed to get total premium: %v", err)
	}
	totalClaimsPaid, err := s.analyticsRepo.GetTotalClaimsPaid(ctx, start, end)
	if err != nil {
		log.Printf("KPIs: failed to get total claims paid: %v", err)
	}
	activePolicies, err := s.analyticsRepo.GetActivePolicyCount(ctx, start, end)
	if err != nil {
		log.Printf("KPIs: failed to get active policies: %v", err)
	}
	lapsedPolicies, err := s.analyticsRepo.GetLapsedPolicyCount(ctx, start, end)
	if err != nil {
		log.Printf("KPIs: failed to get lapsed policies: %v", err)
	}
	totalMembers, err := s.analyticsRepo.GetTotalMemberCount(ctx, start, end)
	if err != nil {
		log.Printf("KPIs: failed to get total members: %v", err)
	}
	renewalRate, err := s.analyticsRepo.GetRenewalRate(ctx, start, end)
	if err != nil {
		log.Printf("KPIs: failed to get renewal rate: %v", err)
	}

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
	start, end := parsePeriod(period)

	var data []map[string]interface{}
	var columns json.RawMessage
	var fetchErr error

	switch reportType {
	case "claims":
		data, fetchErr = s.reportDataRepo.GetClaimsRegisterData(ctx, start, end, "", 10000, 0)
		columns = json.RawMessage(`[{"name":"claim_number","label":"Claim Number","type":"string"},{"name":"member_name","label":"Member","type":"string"},{"name":"provider_name","label":"Provider","type":"string"},{"name":"claimed_amount","label":"Claimed Amount","type":"money"},{"name":"approved_amount","label":"Approved Amount","type":"money"},{"name":"status","label":"Status","type":"string"},{"name":"service_date","label":"Service Date","type":"date"}]`)
	case "premiums":
		data, fetchErr = s.reportDataRepo.GetPremiumRegisterData(ctx, start, end, 10000, 0)
		columns = json.RawMessage(`[{"name":"policy_number","label":"Policy Number","type":"string"},{"name":"policyholder","label":"Policyholder","type":"string"},{"name":"premium_amount","label":"Premium Amount","type":"money"},{"name":"status","label":"Status","type":"string"},{"name":"due_date","label":"Due Date","type":"date"}]`)
	case "members":
		data, fetchErr = s.reportDataRepo.GetMembershipData(ctx, nil, "", 10000, 0)
		columns = json.RawMessage(`[{"name":"member_number","label":"Member Number","type":"string"},{"name":"full_name","label":"Full Name","type":"string"},{"name":"relationship","label":"Relationship","type":"string"},{"name":"status","label":"Status","type":"string"},{"name":"policy_number","label":"Policy Number","type":"string"}]`)
	case "providers":
		data, fetchErr = s.reportDataRepo.GetProviderPerformanceData(ctx, start, end)
		columns = json.RawMessage(`[{"name":"provider_name","label":"Provider Name","type":"string"},{"name":"total_claims","label":"Total Claims","type":"number"},{"name":"total_amount","label":"Total Amount","type":"money"},{"name":"avg_tat","label":"Avg TAT (hrs)","type":"decimal"}]`)
	case "loss_ratio":
		data, fetchErr = s.reportDataRepo.GetLossRatioData(ctx, start, end)
		columns = json.RawMessage(`[{"name":"period","label":"Period","type":"string"},{"name":"premium_earned","label":"Premium Earned","type":"money"},{"name":"claims_incurred","label":"Claims Incurred","type":"money"},{"name":"loss_ratio","label":"Loss Ratio","type":"percentage"}]`)
	case "debtors":
		data, fetchErr = s.reportDataRepo.GetPremiumDebtorsAgeingData(ctx)
		columns = json.RawMessage(`[{"name":"policy_number","label":"Policy Number","type":"string"},{"name":"policyholder","label":"Policyholder","type":"string"},{"name":"outstanding","label":"Outstanding","type":"money"},{"name":"days_overdue","label":"Days Overdue","type":"number"},{"name":"ageing_band","label":"Ageing Band","type":"string"}]`)
	default:
		return schema.NewServiceErrorResponse[[]byte](http.StatusBadRequest, fmt.Sprintf("Unknown report type: %s", reportType), nil)
	}

	if fetchErr != nil {
		log.Printf("ExportCSV: failed to fetch %s data: %v", reportType, fetchErr)
		return schema.NewServiceErrorResponse[[]byte](http.StatusInternalServerError, "Failed to fetch report data", fetchErr)
	}

	csvBytes, err := s.exporter.ExportCSV(columns, data)
	if err != nil {
		return schema.NewServiceErrorResponse[[]byte](http.StatusInternalServerError, "Failed to generate CSV", err)
	}

	return schema.NewServiceResponse(csvBytes, http.StatusOK, "CSV export generated")
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
