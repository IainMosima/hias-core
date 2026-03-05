package entity

import (
	"time"

	"github.com/google/uuid"
)

type ReinsurerStatement struct {
	ID               uuid.UUID `json:"id"`
	StatementNumber  string    `json:"statement_number"`
	TreatyID         uuid.UUID `json:"treaty_id"`
	ParticipantID    uuid.UUID `json:"participant_id"`
	PeriodStart      time.Time `json:"period_start"`
	PeriodEnd        time.Time `json:"period_end"`
	PremiumCeded     int64     `json:"premium_ceded"`
	ClaimsRecovered  int64     `json:"claims_recovered"`
	CommissionDue    int64     `json:"commission_due"`
	ProfitCommission int64     `json:"profit_commission"`
	NetBalance       int64     `json:"net_balance"`
	Status           string    `json:"status"`
	CreatedBy        uuid.UUID `json:"created_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
