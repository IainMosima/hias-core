package entity

import (
	"time"
	"github.com/google/uuid"
)

type Invoice struct {
	ID                 uuid.UUID `json:"id"`
	PolicyID           uuid.UUID `json:"policy_id"`
	InvoiceNumber      string    `json:"invoice_number"`
	Amount             int64     `json:"amount"`
	Currency           string    `json:"currency"`
	DueDate            time.Time `json:"due_date"`
	Status             string    `json:"status"` // PENDING, PAID, OVERDUE, CANCELLED
	BillingPeriodStart time.Time `json:"billing_period_start"`
	BillingPeriodEnd   time.Time `json:"billing_period_end"`
	Notes              string    `json:"notes"`
	CreatedBy          uuid.UUID `json:"created_by"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
