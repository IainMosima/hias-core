package entity

import (
	"time"

	"github.com/google/uuid"
)

type ReinsuranceRecovery struct {
	ID                uuid.UUID `json:"id"`
	RecoveryNumber    string    `json:"recovery_number"`
	ClaimID           uuid.UUID `json:"claim_id"`
	TreatyID          uuid.UUID `json:"treaty_id"`
	TreatyLayerID     uuid.UUID `json:"treaty_layer_id,omitempty"`
	CessionID         uuid.UUID `json:"cession_id,omitempty"`
	GrossClaimAmount  int64     `json:"gross_claim_amount"`
	RecoverableAmount int64     `json:"recoverable_amount"`
	RecoveredAmount   int64     `json:"recovered_amount"`
	OutstandingAmount int64     `json:"outstanding_amount"`
	Status            string    `json:"status"`
	WorkflowStatus    string    `json:"workflow_status"`
	Notes             string    `json:"notes,omitempty"`
	CreatedBy         uuid.UUID `json:"created_by"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
