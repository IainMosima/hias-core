package events

import (
	"time"

	"github.com/google/uuid"
)

const (
	EventPaymentInitiated = "payment.initiated"
	EventPaymentConfirmed = "payment.confirmed"
	EventPaymentFailed    = "payment.failed"
)

type PaymentInitiatedEvent struct {
	PaymentID   uuid.UUID `json:"payment_id"`
	InvoiceID   uuid.UUID `json:"invoice_id,omitempty"`
	ClaimID     uuid.UUID `json:"claim_id,omitempty"`
	Amount      int64     `json:"amount"`
	Method      string    `json:"method"`
	Reference   string    `json:"reference"`
	Timestamp   time.Time `json:"timestamp"`
}

type PaymentConfirmedEvent struct {
	PaymentID   uuid.UUID `json:"payment_id"`
	Reference   string    `json:"reference"`
	Amount      int64     `json:"amount"`
	Timestamp   time.Time `json:"timestamp"`
}

type PaymentFailedEvent struct {
	PaymentID  uuid.UUID `json:"payment_id"`
	Reference  string    `json:"reference"`
	Reason     string    `json:"reason"`
	RetryCount int       `json:"retry_count"`
	Timestamp  time.Time `json:"timestamp"`
}
