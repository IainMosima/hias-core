package events

import (
	"time"

	"github.com/google/uuid"
)

const (
	EventPreAuthSubmitted = "preauth.submitted"
	EventPreAuthApproved  = "preauth.approved"
	EventPreAuthDenied    = "preauth.denied"
)

type PreAuthSubmittedEvent struct {
	PreAuthID   uuid.UUID `json:"preauth_id"`
	PolicyID    uuid.UUID `json:"policy_id"`
	MemberID    uuid.UUID `json:"member_id"`
	ProviderID  uuid.UUID `json:"provider_id"`
	Timestamp   time.Time `json:"timestamp"`
}

type PreAuthApprovedEvent struct {
	PreAuthID      uuid.UUID `json:"preauth_id"`
	AuthCode       string    `json:"auth_code"`
	ApprovedAmount int64     `json:"approved_amount"`
	ValidityEnd    time.Time `json:"validity_end"`
	Timestamp      time.Time `json:"timestamp"`
}

type PreAuthDeniedEvent struct {
	PreAuthID uuid.UUID `json:"preauth_id"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `json:"timestamp"`
}
