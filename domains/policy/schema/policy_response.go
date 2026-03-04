package schema

import (
	"github.com/bitbiz/hias-core/domains/policy/entity"
	"github.com/google/uuid"
	"time"
)

type PolicyResponse struct {
	ID                uuid.UUID `json:"id"`
	PlanID            uuid.UUID `json:"plan_id"`
	PolicyholderName  string    `json:"policyholder_name"`
	PolicyholderEmail string    `json:"policyholder_email"`
	PolicyholderPhone string    `json:"policyholder_phone"`
	PolicyNumber      string    `json:"policy_number"`
	Status            string    `json:"status"`
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`
	PremiumAmount     int64     `json:"premium_amount"`
	Currency          string    `json:"currency"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type MemberResponse struct {
	ID           uuid.UUID  `json:"id"`
	PolicyID     uuid.UUID  `json:"policy_id"`
	NationalID   string     `json:"national_id"`
	Name         string     `json:"name"`
	DateOfBirth  time.Time  `json:"date_of_birth"`
	Gender       string     `json:"gender"`
	Relationship string     `json:"relationship"`
	MemberNumber string     `json:"member_number"`
	Phone        string     `json:"phone"`
	Email        string     `json:"email"`
	KRAPin       string     `json:"kra_pin"`
	County       string     `json:"county"`
	Address      string     `json:"address"`
	Verified     bool       `json:"verified"`
	VerifiedAt   *time.Time `json:"verified_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

func ToPolicyResponse(p *entity.Policy) PolicyResponse {
	return PolicyResponse{
		ID: p.ID, PlanID: p.PlanID, PolicyholderName: p.PolicyholderName,
		PolicyholderEmail: p.PolicyholderEmail, PolicyholderPhone: p.PolicyholderPhone,
		PolicyNumber: p.PolicyNumber, Status: p.Status, StartDate: p.StartDate,
		EndDate: p.EndDate, PremiumAmount: p.PremiumAmount, Currency: p.Currency,
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}

func ToMemberResponse(m *entity.Member) MemberResponse {
	return MemberResponse{
		ID: m.ID, PolicyID: m.PolicyID, NationalID: m.NationalID,
		Name: m.Name, DateOfBirth: m.DateOfBirth, Gender: m.Gender,
		Relationship: m.Relationship, MemberNumber: m.MemberNumber,
		Phone: m.Phone, Email: m.Email, KRAPin: m.KRAPin,
		County: m.County, Address: m.Address, Verified: m.Verified,
		VerifiedAt: m.VerifiedAt, CreatedAt: m.CreatedAt,
	}
}
