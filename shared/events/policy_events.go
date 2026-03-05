package events

import (
	"time"

	"github.com/google/uuid"
)

const (
	EventPolicyActivated  = "policy.activated"
	EventPolicyLapsed     = "policy.lapsed"
	EventPolicyTerminated = "policy.terminated"
	EventPolicyReinstated = "policy.reinstated"
	EventPolicySuspended  = "policy.suspended"
	EventPolicyRenewed    = "policy.renewed"
	EventPolicyUpgraded   = "policy.upgraded"
	EventPolicyDowngraded = "policy.downgraded"

	EventMemberEnrolled  = "member.enrolled"
	EventMemberRemoved   = "member.removed"
	EventMemberSuspended = "member.suspended"

	EventEndorsementCreated  = "endorsement.created"
	EventEndorsementApproved = "endorsement.approved"
	EventEndorsementApplied  = "endorsement.applied"

	EventRenewalInitiated = "renewal.initiated"
	EventRenewalCompleted = "renewal.completed"

	EventDocumentGenerated = "document.generated"
)

type PolicyActivatedEvent struct {
	PolicyID     uuid.UUID `json:"policy_id"`
	PolicyNumber string    `json:"policy_number"`
	PlanID       uuid.UUID `json:"plan_id"`
	Timestamp    time.Time `json:"timestamp"`
}

type PolicyLapsedEvent struct {
	PolicyID     uuid.UUID `json:"policy_id"`
	PolicyNumber string    `json:"policy_number"`
	Reason       string    `json:"reason"`
	Timestamp    time.Time `json:"timestamp"`
}

type PolicyTerminatedEvent struct {
	PolicyID     uuid.UUID `json:"policy_id"`
	PolicyNumber string    `json:"policy_number"`
	Reason       string    `json:"reason"`
	Timestamp    time.Time `json:"timestamp"`
}

type PolicyReinstatedEvent struct {
	PolicyID     uuid.UUID `json:"policy_id"`
	PolicyNumber string    `json:"policy_number"`
	Timestamp    time.Time `json:"timestamp"`
}

type PolicySuspendedEvent struct {
	PolicyID     uuid.UUID `json:"policy_id"`
	PolicyNumber string    `json:"policy_number"`
	Timestamp    time.Time `json:"timestamp"`
}

type PolicyRenewedEvent struct {
	PolicyID        uuid.UUID `json:"policy_id"`
	PolicyNumber    string    `json:"policy_number"`
	NewPolicyID     uuid.UUID `json:"new_policy_id"`
	NewPolicyNumber string    `json:"new_policy_number"`
	Timestamp       time.Time `json:"timestamp"`
}

type PolicyPlanChangedEvent struct {
	PolicyID     uuid.UUID `json:"policy_id"`
	PolicyNumber string    `json:"policy_number"`
	OldPlanID    uuid.UUID `json:"old_plan_id"`
	NewPlanID    uuid.UUID `json:"new_plan_id"`
	OldPremium   int64     `json:"old_premium"`
	NewPremium   int64     `json:"new_premium"`
	Timestamp    time.Time `json:"timestamp"`
}

type MemberEnrolledEvent struct {
	MemberID     uuid.UUID `json:"member_id"`
	MemberNumber string    `json:"member_number"`
	PolicyID     uuid.UUID `json:"policy_id"`
	Timestamp    time.Time `json:"timestamp"`
}

type MemberRemovedEvent struct {
	MemberID  uuid.UUID `json:"member_id"`
	PolicyID  uuid.UUID `json:"policy_id"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `json:"timestamp"`
}

type MemberSuspendedEvent struct {
	MemberID  uuid.UUID `json:"member_id"`
	PolicyID  uuid.UUID `json:"policy_id"`
	Timestamp time.Time `json:"timestamp"`
}

type EndorsementCreatedEvent struct {
	EndorsementID   uuid.UUID `json:"endorsement_id"`
	PolicyID        uuid.UUID `json:"policy_id"`
	EndorsementType string    `json:"endorsement_type"`
	Timestamp       time.Time `json:"timestamp"`
}

type EndorsementApprovedEvent struct {
	EndorsementID uuid.UUID `json:"endorsement_id"`
	PolicyID      uuid.UUID `json:"policy_id"`
	ApprovedBy    uuid.UUID `json:"approved_by"`
	Timestamp     time.Time `json:"timestamp"`
}

type EndorsementAppliedEvent struct {
	EndorsementID uuid.UUID `json:"endorsement_id"`
	PolicyID      uuid.UUID `json:"policy_id"`
	Timestamp     time.Time `json:"timestamp"`
}

type RenewalInitiatedEvent struct {
	RenewalID uuid.UUID `json:"renewal_id"`
	PolicyID  uuid.UUID `json:"policy_id"`
	Timestamp time.Time `json:"timestamp"`
}

type RenewalCompletedEvent struct {
	RenewalID   uuid.UUID `json:"renewal_id"`
	PolicyID    uuid.UUID `json:"policy_id"`
	NewPolicyID uuid.UUID `json:"new_policy_id"`
	Timestamp   time.Time `json:"timestamp"`
}

type DocumentGeneratedEvent struct {
	DocumentID   uuid.UUID `json:"document_id"`
	PolicyID     uuid.UUID `json:"policy_id"`
	DocumentType string    `json:"document_type"`
	Timestamp    time.Time `json:"timestamp"`
}
