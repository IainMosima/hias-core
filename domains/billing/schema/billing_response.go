package schema

import (
	"time"
	"github.com/bitbiz/hias-core/domains/billing/entity"
	"github.com/google/uuid"
)

type InvoiceResponse struct {
	ID                 uuid.UUID `json:"id"`
	PolicyID           uuid.UUID `json:"policy_id"`
	InvoiceNumber      string    `json:"invoice_number"`
	Amount             int64     `json:"amount"`
	Currency           string    `json:"currency"`
	DueDate            time.Time `json:"due_date"`
	Status             string    `json:"status"`
	BillingPeriodStart time.Time `json:"billing_period_start"`
	BillingPeriodEnd   time.Time `json:"billing_period_end"`
	CreatedAt          time.Time `json:"created_at"`
}

type PaymentResponse struct {
	ID              uuid.UUID  `json:"id"`
	Type            string     `json:"type"`
	Amount          int64      `json:"amount"`
	Currency        string     `json:"currency"`
	Method          string     `json:"method"`
	ReferenceNumber string     `json:"reference_number"`
	Status          string     `json:"status"`
	RetryCount      int        `json:"retry_count"`
	PaidAt          *time.Time `json:"paid_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type RemittanceResponse struct {
	ID                   uuid.UUID `json:"id"`
	ProviderID           uuid.UUID `json:"provider_id"`
	TotalAmount          int64     `json:"total_amount"`
	Currency             string    `json:"currency"`
	Status               string    `json:"status"`
	RemittanceAdviceSent bool      `json:"remittance_advice_sent"`
	PeriodStart          time.Time `json:"period_start"`
	PeriodEnd            time.Time `json:"period_end"`
	CreatedAt            time.Time `json:"created_at"`
}

func ToInvoiceResponse(i *entity.Invoice) InvoiceResponse {
	return InvoiceResponse{
		ID: i.ID, PolicyID: i.PolicyID, InvoiceNumber: i.InvoiceNumber,
		Amount: i.Amount, Currency: i.Currency, DueDate: i.DueDate,
		Status: i.Status, BillingPeriodStart: i.BillingPeriodStart,
		BillingPeriodEnd: i.BillingPeriodEnd, CreatedAt: i.CreatedAt,
	}
}

func ToPaymentResponse(p *entity.Payment) PaymentResponse {
	return PaymentResponse{
		ID: p.ID, Type: p.Type, Amount: p.Amount, Currency: p.Currency,
		Method: p.Method, ReferenceNumber: p.ReferenceNumber, Status: p.Status,
		RetryCount: p.RetryCount, PaidAt: p.PaidAt, CreatedAt: p.CreatedAt,
	}
}

func ToRemittanceResponse(r *entity.Remittance) RemittanceResponse {
	return RemittanceResponse{
		ID: r.ID, ProviderID: r.ProviderID, TotalAmount: r.TotalAmount,
		Currency: r.Currency, Status: r.Status, RemittanceAdviceSent: r.RemittanceAdviceSent,
		PeriodStart: r.PeriodStart, PeriodEnd: r.PeriodEnd, CreatedAt: r.CreatedAt,
	}
}
