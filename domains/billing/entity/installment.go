package entity

import (
	"time"

	"github.com/google/uuid"
)

type Installment struct {
	ID                uuid.UUID  `json:"id"`
	ScheduleID        uuid.UUID  `json:"schedule_id"`
	InstallmentNumber int        `json:"installment_number"`
	DueDate           time.Time  `json:"due_date"`
	Amount            int64      `json:"amount"`
	Status            string     `json:"status"` // PENDING, PAID, OVERDUE
	PaidAt            *time.Time `json:"paid_at,omitempty"`
	InvoiceID         uuid.UUID  `json:"invoice_id,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
