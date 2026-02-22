package entity

import (
	"time"
	"github.com/google/uuid"
)

type ClaimLineItem struct {
	ID             uuid.UUID `json:"id"`
	ClaimID        uuid.UUID `json:"claim_id"`
	ProcedureCode  string    `json:"procedure_code"`
	ProcedureName  string    `json:"procedure_name"`
	DiagnosisCode  string    `json:"diagnosis_code"`
	Quantity       int       `json:"quantity"`
	UnitPrice      int64     `json:"unit_price"`
	TotalPrice     int64     `json:"total_price"`
	ApprovedAmount int64     `json:"approved_amount"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
