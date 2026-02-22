package entity

import (
	"encoding/json"
	"time"
	"github.com/google/uuid"
)

type AdjudicationDecision struct {
	ID                   uuid.UUID       `json:"id"`
	ClaimID              uuid.UUID       `json:"claim_id"`
	Decision             string          `json:"decision"` // APPROVE, REJECT, MANUAL_REVIEW
	PayableAmount        int64           `json:"payable_amount"`
	MemberResponsibility int64           `json:"member_responsibility"`
	Reasons              json.RawMessage `json:"reasons"`
	RuleResults          json.RawMessage `json:"rule_results"`
	AdjudicatedBy        uuid.UUID       `json:"adjudicated_by"`
	AdjudicatedAt        time.Time       `json:"adjudicated_at"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}

type AdjudicationResult struct {
	Decision             string       `json:"decision"`
	PayableAmount        int64        `json:"payable_amount"`
	MemberResponsibility int64        `json:"member_responsibility"`
	Reasons              []RuleResult `json:"reasons"`
}

type RuleResult struct {
	Category string `json:"category"` // eligibility, coverage, limits, fraud
	Rule     string `json:"rule"`
	Result   string `json:"result"` // PASS, FAIL, FLAG
	Details  string `json:"details"`
}
