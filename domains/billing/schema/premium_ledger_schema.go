package schema

import (
	"github.com/bitbiz/hias-core/domains/billing/entity"
	"github.com/google/uuid"
	"time"
)

type CreatePremiumLedgerRequest struct {
	PolicyID        string    `json:"policy_id" binding:"required,uuid"`
	EntryType       string    `json:"entry_type" binding:"required,oneof=DEBIT CREDIT"`
	Amount          int64     `json:"amount" binding:"required,min=1"`
	Description     string    `json:"description"`
	ReferenceNumber string    `json:"reference_number" binding:"required"`
	EffectiveDate   time.Time `json:"effective_date" binding:"required"`
}

type PremiumLedgerResponse struct {
	ID              uuid.UUID `json:"id"`
	PolicyID        uuid.UUID `json:"policy_id"`
	EntryType       string    `json:"entry_type"`
	Amount          int64     `json:"amount"`
	Description     string    `json:"description"`
	ReferenceNumber string    `json:"reference_number"`
	EffectiveDate   time.Time `json:"effective_date"`
	BalanceAfter    int64     `json:"balance_after"`
	CreatedAt       time.Time `json:"created_at"`
}

func ToPremiumLedgerResponse(e *entity.PremiumLedgerEntry) PremiumLedgerResponse {
	return PremiumLedgerResponse{
		ID: e.ID, PolicyID: e.PolicyID, EntryType: e.EntryType,
		Amount: e.Amount, Description: e.Description,
		ReferenceNumber: e.ReferenceNumber, EffectiveDate: e.EffectiveDate,
		BalanceAfter: e.BalanceAfter, CreatedAt: e.CreatedAt,
	}
}

type PremiumBalanceResponse struct {
	PolicyID uuid.UUID `json:"policy_id"`
	Balance  int64     `json:"balance"`
}
