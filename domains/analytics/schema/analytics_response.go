package schema

import (
	"github.com/bitbiz/hias-core/domains/analytics/repository"
	"github.com/google/uuid"
)

type DashboardResponse struct {
	ClaimsVolume    *ClaimsVolumeResponse  `json:"claims_volume"`
	ApprovalRate    float64                `json:"approval_rate"`
	AverageTAT      float64                `json:"average_tat_hours"`
	LossRatio       float64                `json:"loss_ratio"`
	FraudRate       float64                `json:"fraud_rate"`
	TotalPremium    int64                  `json:"total_premium_collected"`
	TotalClaimsPaid int64                  `json:"total_claims_paid"`
	TopProviders    []TopProviderResponse  `json:"top_providers"`
	DocumentStats   *DocumentStatsResponse `json:"document_stats"`
}

type ClaimsVolumeResponse struct {
	TotalClaims        int64 `json:"total_claims"`
	ApprovedClaims     int64 `json:"approved_claims"`
	RejectedClaims     int64 `json:"rejected_claims"`
	ManualReviewClaims int64 `json:"manual_review_claims"`
	PaidClaims         int64 `json:"paid_claims"`
}

type TopProviderResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	ClaimCount    int64     `json:"claim_count"`
	TotalAmount   int64     `json:"total_amount"`
	TotalApproved int64     `json:"total_approved"`
}

type DocumentStatsResponse struct {
	TotalDocuments   int64 `json:"total_documents"`
	ActiveDocuments  int64 `json:"active_documents"`
	PendingDocuments int64 `json:"pending_documents"`
	FailedDocuments  int64 `json:"failed_documents"`
}

type KPIResponse struct {
	ApprovalRate    float64 `json:"approval_rate"`
	AverageTAT      float64 `json:"average_tat_hours"`
	LossRatio       float64 `json:"loss_ratio"`
	FraudRate       float64 `json:"fraud_rate"`
	TotalPremium    int64   `json:"total_premium_collected"`
	TotalClaimsPaid int64   `json:"total_claims_paid"`
}

func ToClaimsVolumeResponse(cv *repository.ClaimsVolume) *ClaimsVolumeResponse {
	if cv == nil {
		return nil
	}
	return &ClaimsVolumeResponse{
		TotalClaims:        cv.TotalClaims,
		ApprovedClaims:     cv.ApprovedClaims,
		RejectedClaims:     cv.RejectedClaims,
		ManualReviewClaims: cv.ManualReviewClaims,
		PaidClaims:         cv.PaidClaims,
	}
}

func ToDocumentStatsResponse(ds *repository.DocumentStats) *DocumentStatsResponse {
	if ds == nil {
		return nil
	}
	return &DocumentStatsResponse{
		TotalDocuments:   ds.TotalDocuments,
		ActiveDocuments:  ds.ActiveDocuments,
		PendingDocuments: ds.PendingDocuments,
		FailedDocuments:  ds.FailedDocuments,
	}
}

func ToTopProviderResponseList(providers []*repository.TopProvider) []TopProviderResponse {
	responses := make([]TopProviderResponse, len(providers))
	for i, p := range providers {
		responses[i] = TopProviderResponse{
			ID:            p.ID,
			Name:          p.Name,
			ClaimCount:    p.ClaimCount,
			TotalAmount:   p.TotalAmount,
			TotalApproved: p.TotalApproved,
		}
	}
	return responses
}
