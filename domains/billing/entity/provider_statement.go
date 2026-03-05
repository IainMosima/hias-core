package entity

import (
	"time"

	"github.com/google/uuid"
)

type ProviderStatement struct {
	ID               uuid.UUID  `json:"id"`
	ProviderID       uuid.UUID  `json:"provider_id"`
	StatementNumber  string     `json:"statement_number"`
	PeriodStart      time.Time  `json:"period_start"`
	PeriodEnd        time.Time  `json:"period_end"`
	TotalClaimed     int64      `json:"total_claimed"`
	TotalMatched     int64      `json:"total_matched"`
	TotalDiscrepancy int64      `json:"total_discrepancy"`
	MatchedCount     int        `json:"matched_count"`
	UnmatchedCount   int        `json:"unmatched_count"`
	Status           string     `json:"status"`
	FileName         string     `json:"file_name,omitempty"`
	S3Key            string     `json:"s3_key,omitempty"`
	ReconciledBy     uuid.UUID  `json:"reconciled_by,omitempty"`
	ReconciledAt     *time.Time `json:"reconciled_at,omitempty"`
	CreatedBy        uuid.UUID  `json:"created_by"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type StatementLineItem struct {
	ID                uuid.UUID  `json:"id"`
	StatementID       uuid.UUID  `json:"statement_id"`
	ClaimNumber       string     `json:"claim_number,omitempty"`
	ServiceDate       *time.Time `json:"service_date,omitempty"`
	MemberName        string     `json:"member_name,omitempty"`
	ProcedureCode     string     `json:"procedure_code,omitempty"`
	ClaimedAmount     int64      `json:"claimed_amount"`
	MatchedClaimID    uuid.UUID  `json:"matched_claim_id,omitempty"`
	MatchStatus       string     `json:"match_status"`
	DiscrepancyAmount int64      `json:"discrepancy_amount"`
	Notes             string     `json:"notes,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
}
