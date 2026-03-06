package entity

import (
	"time"

	"github.com/google/uuid"
)

type ClaimStatusHistory struct {
	ID              uuid.UUID
	ClaimID         uuid.UUID
	FromStatus      string
	ToStatus        string
	Action          string
	Notes           string
	PerformedBy     uuid.UUID
	PerformedByName string
	CreatedAt       time.Time
}
