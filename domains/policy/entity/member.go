package entity

import (
	"github.com/google/uuid"
	"time"
)

type Member struct {
	ID                uuid.UUID  `json:"id"`
	PolicyID          uuid.UUID  `json:"policy_id"`
	NationalID        string     `json:"national_id"`
	Name              string     `json:"name"`
	DateOfBirth       time.Time  `json:"date_of_birth"`
	Gender            string     `json:"gender"`
	Relationship      string     `json:"relationship"` // principal, spouse, child, parent
	MemberNumber      string     `json:"member_number"`
	Phone             string     `json:"phone"`
	Email             string     `json:"email"`
	KRAPin            string     `json:"kra_pin"`
	County            string     `json:"county"`
	City              string     `json:"city"`
	Country           string     `json:"country"`
	Address           string     `json:"address"`
	Status            string     `json:"status"` // ACTIVE, SUSPENDED, REMOVED
	Verified          bool       `json:"verified"`
	VerifiedAt        *time.Time `json:"verified_at,omitempty"`
	CoverageStartDate *time.Time `json:"coverage_start_date,omitempty"`
	CoverageEndDate   *time.Time `json:"coverage_end_date,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
