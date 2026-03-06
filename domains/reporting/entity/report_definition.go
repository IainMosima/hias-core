package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ReportDefinition struct {
	ID                uuid.UUID       `json:"id"`
	Code              string          `json:"code"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	Category          string          `json:"category"`
	ReportType        string          `json:"report_type"`
	QueryTemplate     string          `json:"query_template,omitempty"`
	DefaultParameters json.RawMessage `json:"default_parameters"`
	AllowedRoles      []string        `json:"allowed_roles"`
	Columns           json.RawMessage `json:"columns"`
	IsActive          bool            `json:"is_active"`
	CreatedBy         uuid.UUID       `json:"created_by"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}
