package entity

import (
	"time"
	"github.com/google/uuid"
)

type Plan struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Type         string    `json:"type"` // individual, group
	BasePremium  int64     `json:"base_premium"`
	Currency     string    `json:"currency"`
	Status       string    `json:"status"`
	Description  string    `json:"description"`
	CreatedBy    uuid.UUID `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
