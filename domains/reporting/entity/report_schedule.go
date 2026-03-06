package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ReportSchedule struct {
	ID                 uuid.UUID       `json:"id"`
	ReportDefinitionID uuid.UUID       `json:"report_definition_id"`
	Name               string          `json:"name"`
	CronExpression     string          `json:"cron_expression"`
	Parameters         json.RawMessage `json:"parameters"`
	ExportFormat       string          `json:"export_format"`
	Recipients         []uuid.UUID     `json:"recipients"`
	IsActive           bool            `json:"is_active"`
	LastRunAt          *time.Time      `json:"last_run_at,omitempty"`
	NextRunAt          *time.Time      `json:"next_run_at,omitempty"`
	CreatedBy          uuid.UUID       `json:"created_by"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}
