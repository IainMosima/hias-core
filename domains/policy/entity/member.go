package entity

import (
	"time"
	"github.com/google/uuid"
)

type Member struct {
	ID           uuid.UUID  `json:"id"`
	PolicyID     uuid.UUID  `json:"policy_id"`
	NationalID   string     `json:"national_id"`
	Name         string     `json:"name"`
	DateOfBirth  time.Time  `json:"date_of_birth"`
	Gender       string     `json:"gender"`
	Relationship string     `json:"relationship"` // principal, spouse, child, parent
	MemberNumber string     `json:"member_number"`
	Phone        string     `json:"phone"`
	Email        string     `json:"email"`
	Verified     bool       `json:"verified"`
	VerifiedAt   *time.Time `json:"verified_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
