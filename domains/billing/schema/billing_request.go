package schema

type InitiatePaymentRequest struct {
	InvoiceID string `json:"invoice_id,omitempty"`
	ClaimID   string `json:"claim_id,omitempty"`
	Amount    int64  `json:"amount" binding:"required,min=1"`
	Method    string `json:"method" binding:"required,oneof=MPESA BANK_TRANSFER"`
	Phone     string `json:"phone"`
}

type MpesaWebhookRequest struct {
	Body interface{} `json:"Body"`
}
