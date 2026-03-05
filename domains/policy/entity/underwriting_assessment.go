package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type UnderwritingAssessment struct {
	ID                  uuid.UUID       `json:"id"`
	PolicyID            uuid.UUID       `json:"policy_id"`
	MemberID            uuid.UUID       `json:"member_id"`
	Status              string          `json:"status"`
	Questionnaire       json.RawMessage `json:"questionnaire"`
	MedicalDeclarations json.RawMessage `json:"medical_declarations"`
	RiskScore           int             `json:"risk_score"`
	RiskFlags           json.RawMessage `json:"risk_flags"`
	DecisionReason      string          `json:"decision_reason"`
	AssessedBy          uuid.UUID       `json:"assessed_by"`
	AssessedAt          *time.Time      `json:"assessed_at,omitempty"`
	CreatedBy           uuid.UUID       `json:"created_by"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
}
