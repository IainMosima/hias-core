package schema

import (
	"github.com/bitbiz/hias-core/domains/billing/entity"
	"github.com/google/uuid"
	"time"
)

type CreateRefundRequest struct {
	PolicyID     string `json:"policy_id" binding:"required,uuid"`
	CreditNoteID string `json:"credit_note_id,omitempty"`
	Amount       int64  `json:"amount" binding:"required,min=1"`
	Reason       string `json:"reason" binding:"required"`
}

type RefundResponse struct {
	ID           uuid.UUID  `json:"id"`
	PolicyID     uuid.UUID  `json:"policy_id"`
	CreditNoteID uuid.UUID  `json:"credit_note_id,omitempty"`
	Amount       int64      `json:"amount"`
	Currency     string     `json:"currency"`
	Status       string     `json:"status"`
	Reason       string     `json:"reason"`
	ApprovedBy   uuid.UUID  `json:"approved_by,omitempty"`
	ApprovedAt   *time.Time `json:"approved_at,omitempty"`
	ProcessedAt  *time.Time `json:"processed_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

func ToRefundResponse(r *entity.Refund) RefundResponse {
	return RefundResponse{
		ID: r.ID, PolicyID: r.PolicyID, CreditNoteID: r.CreditNoteID,
		Amount: r.Amount, Currency: r.Currency, Status: r.Status,
		Reason: r.Reason, ApprovedBy: r.ApprovedBy,
		ApprovedAt: r.ApprovedAt, ProcessedAt: r.ProcessedAt,
		CreatedAt: r.CreatedAt,
	}
}
