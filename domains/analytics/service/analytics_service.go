package service

import (
	"context"
	"github.com/bitbiz/hias-core/domains/identity/schema"
)

type DashboardData struct {
	ClaimsVolume    interface{} `json:"claims_volume"`
	ApprovalRate    float64     `json:"approval_rate"`
	AverageTAT      float64     `json:"average_tat_hours"`
	LossRatio       float64     `json:"loss_ratio"`
	FraudRate       float64     `json:"fraud_rate"`
	TotalPremium    int64       `json:"total_premium_collected"`
	TotalClaimsPaid int64       `json:"total_claims_paid"`
	TopProviders    interface{} `json:"top_providers"`
	ActivePolicies  int64       `json:"active_policies"`
	LapsedPolicies  int64       `json:"lapsed_policies"`
	TotalMembers    int64       `json:"total_members"`
	RenewalRate     float64     `json:"renewal_rate"`
	TotalDocuments  int64       `json:"total_documents"`
	DocumentStats   interface{} `json:"document_stats"`
}

type AnalyticsService interface {
	GetDashboard(ctx context.Context, period string) *schema.ServiceResponse[DashboardData]
	GetKPIs(ctx context.Context, period string) *schema.ServiceResponse[interface{}]
	ExportCSV(ctx context.Context, reportType, period string) *schema.ServiceResponse[[]byte]
}
