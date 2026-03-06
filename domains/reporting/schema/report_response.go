package schema

import (
	"encoding/json"
	"time"

	"github.com/bitbiz/hias-core/domains/reporting/entity"
	"github.com/google/uuid"
)

type ReportDefinitionResponse struct {
	ID                uuid.UUID       `json:"id"`
	Code              string          `json:"code"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	Category          string          `json:"category"`
	ReportType        string          `json:"report_type"`
	DefaultParameters json.RawMessage `json:"default_parameters"`
	AllowedRoles      []string        `json:"allowed_roles"`
	Columns           json.RawMessage `json:"columns"`
	IsActive          bool            `json:"is_active"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

type ReportScheduleResponse struct {
	ID                 uuid.UUID       `json:"id"`
	ReportDefinitionID uuid.UUID       `json:"report_definition_id"`
	DefinitionName     string          `json:"definition_name"`
	Name               string          `json:"name"`
	CronExpression     string          `json:"cron_expression"`
	Parameters         json.RawMessage `json:"parameters"`
	ExportFormat       string          `json:"export_format"`
	Recipients         []uuid.UUID     `json:"recipients"`
	IsActive           bool            `json:"is_active"`
	LastRunAt          *time.Time      `json:"last_run_at,omitempty"`
	NextRunAt          *time.Time      `json:"next_run_at,omitempty"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

type GeneratedReportResponse struct {
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

type ReportPreviewResponse struct {
	Columns  json.RawMessage          `json:"columns"`
	Data     []map[string]interface{} `json:"data"`
	RowCount int                      `json:"row_count"`
	Summary  map[string]interface{}   `json:"summary,omitempty"`
}

type ManagementDashboardResponse struct {
	LossRatio          float64 `json:"loss_ratio"`
	ClaimsVolume       int64   `json:"claims_volume"`
	ApprovalRate       float64 `json:"approval_rate"`
	AvgTATHours        float64 `json:"avg_tat_hours"`
	TotalPremium       int64   `json:"total_premium"`
	TotalClaimsPaid    int64   `json:"total_claims_paid"`
	ActivePolicies     int64   `json:"active_policies"`
	TotalMembers       int64   `json:"total_members"`
	RenewalRate        float64 `json:"renewal_rate"`
	PremiumGrowth      float64 `json:"premium_growth_pct"`
	OutstandingPremium int64   `json:"outstanding_premium"`
	SLABreachCount     int64   `json:"sla_breach_count"`
}

func ToReportDefinitionResponse(d *entity.ReportDefinition) ReportDefinitionResponse {
	return ReportDefinitionResponse{
		ID:                d.ID,
		Code:              d.Code,
		Name:              d.Name,
		Description:       d.Description,
		Category:          d.Category,
		ReportType:        d.ReportType,
		DefaultParameters: d.DefaultParameters,
		AllowedRoles:      d.AllowedRoles,
		Columns:           d.Columns,
		IsActive:          d.IsActive,
		CreatedAt:         d.CreatedAt,
		UpdatedAt:         d.UpdatedAt,
	}
}

func ToReportScheduleResponse(s *entity.ReportSchedule, defName string) ReportScheduleResponse {
	return ReportScheduleResponse{
		ID:                 s.ID,
		ReportDefinitionID: s.ReportDefinitionID,
		DefinitionName:     defName,
		Name:               s.Name,
		CronExpression:     s.CronExpression,
		Parameters:         s.Parameters,
		ExportFormat:       s.ExportFormat,
		Recipients:         s.Recipients,
		IsActive:           s.IsActive,
		LastRunAt:          s.LastRunAt,
		NextRunAt:          s.NextRunAt,
		CreatedAt:          s.CreatedAt,
		UpdatedAt:          s.UpdatedAt,
	}
}

func ToGeneratedReportResponse(r *entity.GeneratedReport) GeneratedReportResponse {
	return GeneratedReportResponse{
		ID:                 r.ID,
		ReportDefinitionID: r.ReportDefinitionID,
		ScheduleID:         r.ScheduleID,
		ReportNumber:       r.ReportNumber,
		Name:               r.Name,
		Parameters:         r.Parameters,
		Format:             r.Format,
		Status:             r.Status,
		RowCount:           r.RowCount,
		FileSize:           r.FileSize,
		ErrorMessage:       r.ErrorMessage,
		GeneratedBy:        r.GeneratedBy,
		GeneratedAt:        r.GeneratedAt,
		ExpiresAt:          r.ExpiresAt,
		CreatedAt:          r.CreatedAt,
	}
}
