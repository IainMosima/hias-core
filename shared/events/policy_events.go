package events

import (
	"time"

	"github.com/google/uuid"
)

const (
	EventPolicyActivated  = "policy.activated"
	EventPolicyLapsed     = "policy.lapsed"
	EventPolicyTerminated = "policy.terminated"
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
