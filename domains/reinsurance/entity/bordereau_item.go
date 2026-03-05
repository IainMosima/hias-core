package entity

import (
	"time"

	"github.com/google/uuid"
)

type BordereauItem struct {
	ID               uuid.UUID `json:"id"`
	BordereauID      uuid.UUID `json:"bordereau_id"`
	CessionID        uuid.UUID `json:"cession_id,omitempty"`
	RecoveryID       uuid.UUID `json:"recovery_id,omitempty"`
	PolicyNumber     string    `json:"policy_number,omitempty"`
	ClaimNumber      string    `json:"claim_number,omitempty"`
	GrossAmount      int64     `json:"gross_amount"`
	CededAmount      int64     `json:"ceded_amount"`
	CommissionAmount int64     `json:"commission_amount"`
	CreatedAt        time.Time `json:"created_at"`
}
