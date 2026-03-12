package schema

import (
	"github.com/bitbiz/hias-core/domains/billing/entity"
	"github.com/google/uuid"
	"time"
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
	PolicyNumber       string    `json:"policy_number,omitempty"`
	PolicyholderName   string    `json:"policyholder_name,omitempty"`
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
	WHTRate              float64   `json:"wht_rate"`
	WHTAmount            int64     `json:"wht_amount"`
	NetAmount            int64     `json:"net_amount"`
	CreatedAt            time.Time `json:"created_at"`
}

func ToInvoiceResponse(i *entity.Invoice) InvoiceResponse {
	return InvoiceResponse{
		ID: i.ID, PolicyID: i.PolicyID, InvoiceNumber: i.InvoiceNumber,
		Amount: i.Amount, Currency: i.Currency, DueDate: i.DueDate,
		Status: i.Status, BillingPeriodStart: i.BillingPeriodStart,
		BillingPeriodEnd: i.BillingPeriodEnd,
		PolicyNumber:     i.PolicyNumber, PolicyholderName: i.PolicyholderName,
		CreatedAt: i.CreatedAt,
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
		PeriodStart: r.PeriodStart, PeriodEnd: r.PeriodEnd,
		WHTRate: r.WHTRate, WHTAmount: r.WHTAmount, NetAmount: r.NetAmount,
		CreatedAt: r.CreatedAt,
	}
}

type CreditNoteResponse struct {
	ID                 uuid.UUID  `json:"id"`
	PolicyID           uuid.UUID  `json:"policy_id"`
	MemberID           uuid.UUID  `json:"member_id,omitempty"`
	CreditNoteNumber   string     `json:"credit_note_number"`
	Amount             int64      `json:"amount"`
	Currency           string     `json:"currency"`
	Reason             string     `json:"reason"`
	Status             string     `json:"status"`
	AppliedToInvoiceID uuid.UUID  `json:"applied_to_invoice_id,omitempty"`
	ApprovedBy         uuid.UUID  `json:"approved_by,omitempty"`
	ApprovedAt         *time.Time `json:"approved_at,omitempty"`
	AppliedAt          *time.Time `json:"applied_at,omitempty"`
	CreatedBy          uuid.UUID  `json:"created_by"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

func ToCreditNoteResponse(cn *entity.CreditNote) CreditNoteResponse {
	return CreditNoteResponse{
		ID: cn.ID, PolicyID: cn.PolicyID, MemberID: cn.MemberID,
		CreditNoteNumber: cn.CreditNoteNumber, Amount: cn.Amount,
		Currency: cn.Currency, Reason: cn.Reason, Status: cn.Status,
		AppliedToInvoiceID: cn.AppliedToInvoiceID, ApprovedBy: cn.ApprovedBy,
		ApprovedAt: cn.ApprovedAt, AppliedAt: cn.AppliedAt,
		CreatedBy: cn.CreatedBy, CreatedAt: cn.CreatedAt, UpdatedAt: cn.UpdatedAt,
	}
}

type ProviderStatementResponse struct {
	ID               uuid.UUID  `json:"id"`
	ProviderID       uuid.UUID  `json:"provider_id"`
	StatementNumber  string     `json:"statement_number"`
	PeriodStart      time.Time  `json:"period_start"`
	PeriodEnd        time.Time  `json:"period_end"`
	TotalClaimed     int64      `json:"total_claimed"`
	TotalMatched     int64      `json:"total_matched"`
	TotalDiscrepancy int64      `json:"total_discrepancy"`
	MatchedCount     int        `json:"matched_count"`
	UnmatchedCount   int        `json:"unmatched_count"`
	Status           string     `json:"status"`
	FileName         string     `json:"file_name,omitempty"`
	ReconciledAt     *time.Time `json:"reconciled_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}

type StatementLineItemResponse struct {
	ID                uuid.UUID  `json:"id"`
	StatementID       uuid.UUID  `json:"statement_id"`
	ClaimNumber       string     `json:"claim_number,omitempty"`
	ServiceDate       *time.Time `json:"service_date,omitempty"`
	MemberName        string     `json:"member_name,omitempty"`
	ProcedureCode     string     `json:"procedure_code,omitempty"`
	ClaimedAmount     int64      `json:"claimed_amount"`
	MatchedClaimID    uuid.UUID  `json:"matched_claim_id,omitempty"`
	MatchStatus       string     `json:"match_status"`
	DiscrepancyAmount int64      `json:"discrepancy_amount"`
	Notes             string     `json:"notes,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
}

type UploadStatementRequest struct {
	ProviderID  string                   `json:"provider_id" binding:"required,uuid"`
	PeriodStart time.Time                `json:"period_start" binding:"required"`
	PeriodEnd   time.Time                `json:"period_end" binding:"required"`
	FileName    string                   `json:"file_name"`
	S3Key       string                   `json:"s3_key"`
	LineItems   []StatementLineItemInput `json:"line_items" binding:"required,min=1"`
}

type StatementLineItemInput struct {
	ClaimNumber   string     `json:"claim_number"`
	ServiceDate   *time.Time `json:"service_date,omitempty"`
	MemberName    string     `json:"member_name"`
	ProcedureCode string     `json:"procedure_code"`
	ClaimedAmount int64      `json:"claimed_amount" binding:"required,min=0"`
}

type PaymentExportResponse struct {
	RemittanceID   uuid.UUID            `json:"remittance_id"`
	ProviderID     uuid.UUID            `json:"provider_id"`
	ProviderName   string               `json:"provider_name"`
	TotalAmount    int64                `json:"total_amount"`
	Currency       string               `json:"currency"`
	PeriodStart    time.Time            `json:"period_start"`
	PeriodEnd      time.Time            `json:"period_end"`
	Claims         []PaymentExportClaim `json:"claims"`
	PaymentFileCSV string               `json:"payment_file_csv,omitempty"`
}

type PaymentExportClaim struct {
	ClaimNumber string    `json:"claim_number"`
	Amount      int64     `json:"amount"`
	ServiceDate time.Time `json:"service_date"`
}

func ToProviderStatementResponse(s *entity.ProviderStatement) ProviderStatementResponse {
	return ProviderStatementResponse{
		ID: s.ID, ProviderID: s.ProviderID, StatementNumber: s.StatementNumber,
		PeriodStart: s.PeriodStart, PeriodEnd: s.PeriodEnd,
		TotalClaimed: s.TotalClaimed, TotalMatched: s.TotalMatched,
		TotalDiscrepancy: s.TotalDiscrepancy, MatchedCount: s.MatchedCount,
		UnmatchedCount: s.UnmatchedCount, Status: s.Status,
		FileName: s.FileName, ReconciledAt: s.ReconciledAt,
		CreatedAt: s.CreatedAt,
	}
}

func ToStatementLineItemResponse(item *entity.StatementLineItem) StatementLineItemResponse {
	return StatementLineItemResponse{
		ID: item.ID, StatementID: item.StatementID, ClaimNumber: item.ClaimNumber,
		ServiceDate: item.ServiceDate, MemberName: item.MemberName,
		ProcedureCode: item.ProcedureCode, ClaimedAmount: item.ClaimedAmount,
		MatchedClaimID: item.MatchedClaimID, MatchStatus: item.MatchStatus,
		DiscrepancyAmount: item.DiscrepancyAmount, Notes: item.Notes,
		CreatedAt: item.CreatedAt,
	}
}
