package schema

import "time"

type CreateInvoiceRequest struct {
	PolicyID           string     `json:"policy_id" binding:"required,uuid"`
	Amount             int64      `json:"amount" binding:"required,min=1"`
	DueDate            time.Time  `json:"due_date" binding:"required"`
	Currency           string     `json:"currency"`
	BillingPeriodStart *time.Time `json:"billing_period_start"`
	BillingPeriodEnd   *time.Time `json:"billing_period_end"`
	Notes              string     `json:"notes"`
}

type InitiatePaymentRequest struct {
	InvoiceID       string `json:"invoice_id,omitempty"`
	ClaimID         string `json:"claim_id,omitempty"`
	Amount          int64  `json:"amount" binding:"required,min=1"`
	Method          string `json:"method" binding:"required,oneof=MPESA BANK_TRANSFER"`
	Phone           string `json:"phone"`
	ReferenceNumber string `json:"reference_number"`
}

type MpesaWebhookRequest struct {
	Body interface{} `json:"Body"`
}
