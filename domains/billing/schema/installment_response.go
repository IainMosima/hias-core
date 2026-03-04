package schema

import (
	"time"

	"github.com/bitbiz/hias-core/domains/billing/entity"
	"github.com/google/uuid"
)

type InstallmentScheduleResponse struct {
	ID                   uuid.UUID             `json:"id"`
	PolicyID             uuid.UUID             `json:"policy_id"`
	Frequency            string                `json:"frequency"`
	TotalInstallments    int                   `json:"total_installments"`
	AmountPerInstallment int64                 `json:"amount_per_installment"`
	StartDate            time.Time             `json:"start_date"`
	Status               string                `json:"status"`
	CreatedAt            time.Time             `json:"created_at"`
	Installments         []InstallmentResponse `json:"installments,omitempty"`
}

type InstallmentResponse struct {
	ID                uuid.UUID  `json:"id"`
	ScheduleID        uuid.UUID  `json:"schedule_id"`
	InstallmentNumber int        `json:"installment_number"`
	DueDate           time.Time  `json:"due_date"`
	Amount            int64      `json:"amount"`
	Status            string     `json:"status"`
	PaidAt            *time.Time `json:"paid_at,omitempty"`
	InvoiceID         uuid.UUID  `json:"invoice_id,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
}

func ToInstallmentScheduleResponse(s *entity.InstallmentSchedule) InstallmentScheduleResponse {
	return InstallmentScheduleResponse{
		ID: s.ID, PolicyID: s.PolicyID, Frequency: s.Frequency,
		TotalInstallments: s.TotalInstallments, AmountPerInstallment: s.AmountPerInstallment,
		StartDate: s.StartDate, Status: s.Status, CreatedAt: s.CreatedAt,
	}
}

func ToInstallmentResponse(i *entity.Installment) InstallmentResponse {
	return InstallmentResponse{
		ID: i.ID, ScheduleID: i.ScheduleID, InstallmentNumber: i.InstallmentNumber,
		DueDate: i.DueDate, Amount: i.Amount, Status: i.Status,
		PaidAt: i.PaidAt, InvoiceID: i.InvoiceID, CreatedAt: i.CreatedAt,
	}
}
