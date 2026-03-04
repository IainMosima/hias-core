package entity

import (
	"github.com/google/uuid"
	"time"
)

type RateCard struct {
	ID            uuid.UUID `json:"id"`
	ProviderID    uuid.UUID `json:"provider_id"`
	ProcedureCode string    `json:"procedure_code"`
	ProcedureName string    `json:"procedure_name"`
	RateAmount    int64     `json:"rate_amount"`
	EffectiveDate time.Time `json:"effective_date"`
	AgeFrom       int       `json:"age_from"`
	AgeTo         int       `json:"age_to"`
	Gender        string    `json:"gender,omitempty"`
	Relationship  string    `json:"relationship,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
