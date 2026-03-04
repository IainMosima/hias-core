package schema

import "time"

type CreateInstallmentScheduleRequest struct {
	PolicyID  string    `json:"policy_id" binding:"required,uuid"`
	Frequency string    `json:"frequency" binding:"required"`
	StartDate time.Time `json:"start_date"`
}

type MarkInstallmentPaidRequest struct {
	InvoiceID string `json:"invoice_id"`
}
