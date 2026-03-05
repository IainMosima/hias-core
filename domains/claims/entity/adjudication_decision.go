package entity

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
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
	DeductibleApplied    int64        `json:"deductible_applied"`
	CoPayApplied         int64        `json:"co_pay_applied"`
	SubLimitApplied      int64        `json:"sub_limit_applied,omitempty"`
	BenefitCategory      string       `json:"benefit_category,omitempty"`
	Reasons              []RuleResult `json:"reasons"`
}

type RuleResult struct {
	Category string `json:"category"` // eligibility, coverage, limits, fraud
	Rule     string `json:"rule"`
	Result   string `json:"result"` // PASS, FAIL, FLAG
	Details  string `json:"details"`
}
