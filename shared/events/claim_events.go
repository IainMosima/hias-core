package events

import (
	"time"

	"github.com/google/uuid"
)

const (
	EventClaimSubmitted = "claim.submitted"
	EventClaimApproved  = "claim.approved"
	EventClaimRejected  = "claim.rejected"
	EventClaimPaid      = "claim.paid"
)

type ClaimSubmittedEvent struct {
	ClaimID     uuid.UUID `json:"claim_id"`
	ClaimNumber string    `json:"claim_number"`
	PolicyID    uuid.UUID `json:"policy_id"`
	MemberID    uuid.UUID `json:"member_id"`
	ProviderID  uuid.UUID `json:"provider_id"`
	TotalAmount int64     `json:"total_amount"`
	Timestamp   time.Time `json:"timestamp"`
}

type ClaimApprovedEvent struct {
	ClaimID        uuid.UUID `json:"claim_id"`
	ClaimNumber    string    `json:"claim_number"`
	ApprovedAmount int64     `json:"approved_amount"`
	ProviderID     uuid.UUID `json:"provider_id"`
	Timestamp      time.Time `json:"timestamp"`
}

type ClaimRejectedEvent struct {
	ClaimID     uuid.UUID `json:"claim_id"`
	ClaimNumber string    `json:"claim_number"`
	Reason      string    `json:"reason"`
	Timestamp   time.Time `json:"timestamp"`
}

type ClaimPaidEvent struct {
	ClaimID     uuid.UUID `json:"claim_id"`
	ClaimNumber string    `json:"claim_number"`
	Amount      int64     `json:"amount"`
	PaymentID   uuid.UUID `json:"payment_id"`
	Timestamp   time.Time `json:"timestamp"`
}
