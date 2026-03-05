package entity

import (
	"time"

	"github.com/google/uuid"
)

type CaseRecord struct {
	ID                 uuid.UUID  `json:"id"`
	CaseNumber         string     `json:"case_number"`
	PreAuthID          uuid.UUID  `json:"preauth_id"`
	PolicyID           uuid.UUID  `json:"policy_id"`
	MemberID           uuid.UUID  `json:"member_id"`
	ProviderID         uuid.UUID  `json:"provider_id"`
	Status             string     `json:"status"`
	AdmissionDate      *time.Time `json:"admission_date,omitempty"`
	ExpectedDischarge  *time.Time `json:"expected_discharge,omitempty"`
	ActualDischarge    *time.Time `json:"actual_discharge,omitempty"`
	Diagnosis          string     `json:"diagnosis,omitempty"`
	TreatingDoctor     string     `json:"treating_doctor,omitempty"`
	RoomType           string     `json:"room_type,omitempty"`
	TotalEstimatedCost int64      `json:"total_estimated_cost"`
	TotalActualCost    int64      `json:"total_actual_cost"`
	Notes              string     `json:"notes,omitempty"`
	ClosedAt           *time.Time `json:"closed_at,omitempty"`
	CreatedBy          uuid.UUID  `json:"created_by"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}
