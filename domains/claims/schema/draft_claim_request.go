package schema

import "time"

type DraftClaimRequest struct {
	PolicyID       string            `json:"policy_id"`
	MemberID       string            `json:"member_id"`
	ProviderID     string            `json:"provider_id"`
	PreAuthID      string            `json:"preauth_id,omitempty"`
	DiagnosisCodes []string          `json:"diagnosis_codes"`
	ServiceDate    *time.Time        `json:"service_date"`
	Notes          string            `json:"notes"`
	ClaimType      string            `json:"claim_type"`
	LineItems      []LineItemRequest `json:"line_items"`
}
