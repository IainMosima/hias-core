package entity

import (
	"github.com/google/uuid"
	"time"
)

type Provider struct {
	ID                  uuid.UUID  `json:"id"`
	Name                string     `json:"name"`
	Type                string     `json:"type"` // hospital, clinic, pharmacy, lab
	LicenseNumber       string     `json:"license_number"`
	Status              string     `json:"status"` // PENDING, CREDENTIALING, ACTIVE, SUSPENDED, TERMINATED
	County              string     `json:"county"`
	Address             string     `json:"address"`
	Phone               string     `json:"phone"`
	Email               string     `json:"email"`
	ContactPerson       string     `json:"contact_person"`
	Tier                string     `json:"tier"`
	AccreditationStatus string     `json:"accreditation_status"`
	AccreditationExpiry *time.Time `json:"accreditation_expiry"`
	AccreditationBody   string     `json:"accreditation_body"`
	UserID              uuid.UUID  `json:"user_id"`
	CreatedBy           uuid.UUID  `json:"created_by"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}
