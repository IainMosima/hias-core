package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type GeneratedReport struct {
	ID                 uuid.UUID       `json:"id"`
	ReportDefinitionID uuid.UUID       `json:"report_definition_id"`
	ScheduleID         *uuid.UUID      `json:"schedule_id,omitempty"`
	ReportNumber       string          `json:"report_number"`
	Name               string          `json:"name"`
	Parameters         json.RawMessage `json:"parameters"`
	Format             string          `json:"format"`
	Status             string          `json:"status"`
	RowCount           int             `json:"row_count"`
	FileSize           int64           `json:"file_size"`
	ErrorMessage       string          `json:"error_message,omitempty"`
	GeneratedBy        uuid.UUID       `json:"generated_by"`
	GeneratedAt        time.Time       `json:"generated_at"`
	ExpiresAt          *time.Time      `json:"expires_at,omitempty"`
	CreatedAt          time.Time       `json:"created_at"`
}
