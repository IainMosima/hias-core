package entity

import (
	"time"

	"github.com/google/uuid"
)

type InstallmentSchedule struct {
	ID                   uuid.UUID `json:"id"`
	PolicyID             uuid.UUID `json:"policy_id"`
	Frequency            string    `json:"frequency"` // monthly, quarterly, semi_annual, annual
	TotalInstallments    int       `json:"total_installments"`
	AmountPerInstallment int64     `json:"amount_per_installment"`
	StartDate            time.Time `json:"start_date"`
	Status               string    `json:"status"`
	CreatedBy            uuid.UUID `json:"created_by,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}
