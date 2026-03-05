package entity

import (
	"time"

	"github.com/google/uuid"
)

type CreditNote struct {
	ID                 uuid.UUID  `json:"id"`
	PolicyID           uuid.UUID  `json:"policy_id"`
	MemberID           uuid.UUID  `json:"member_id"`
	CreditNoteNumber   string     `json:"credit_note_number"`
	Amount             int64      `json:"amount"`
	Currency           string     `json:"currency"`
	Reason             string     `json:"reason"`
	Status             string     `json:"status"`
	AppliedToInvoiceID uuid.UUID  `json:"applied_to_invoice_id"`
	ApprovedBy         uuid.UUID  `json:"approved_by"`
	ApprovedAt         *time.Time `json:"approved_at,omitempty"`
	AppliedAt          *time.Time `json:"applied_at,omitempty"`
	CreatedBy          uuid.UUID  `json:"created_by"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}
