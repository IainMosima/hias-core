package schema

import "encoding/json"

type GenerateReportRequest struct {
	ReportCode string          `json:"report_code" binding:"required"`
	Parameters json.RawMessage `json:"parameters"`
	Format     string          `json:"format" binding:"required,oneof=CSV XLSX PDF"`
}

type PreviewReportRequest struct {
	ReportCode string          `json:"report_code" binding:"required"`
	Parameters json.RawMessage `json:"parameters"`
}

type CreateAdHocReportRequest struct {
	Name         string          `json:"name" binding:"required"`
	Description  string          `json:"description"`
	Category     string          `json:"category" binding:"required"`
	Columns      json.RawMessage `json:"columns" binding:"required"`
	Filters      json.RawMessage `json:"filters" binding:"required"`
	AllowedRoles []string        `json:"allowed_roles" binding:"required"`
}

type CreateScheduleRequest struct {
	ReportDefinitionID string          `json:"report_definition_id" binding:"required,uuid"`
	Name               string          `json:"name" binding:"required"`
	CronExpression     string          `json:"cron_expression" binding:"required"`
	Parameters         json.RawMessage `json:"parameters"`
	ExportFormat       string          `json:"export_format" binding:"required,oneof=CSV XLSX PDF"`
	Recipients         []string        `json:"recipients" binding:"required,min=1"`
}

type UpdateScheduleRequest struct {
	Name           string          `json:"name"`
	CronExpression string          `json:"cron_expression"`
	Parameters     json.RawMessage `json:"parameters"`
	ExportFormat   string          `json:"export_format" binding:"omitempty,oneof=CSV XLSX PDF"`
	Recipients     []string        `json:"recipients"`
	IsActive       *bool           `json:"is_active"`
}

type DrillDownRequest struct {
	ReportCode string `json:"report_code" binding:"required"`
	EntityID   string `json:"entity_id" binding:"required,uuid"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Format     string `json:"format" binding:"required,oneof=CSV XLSX PDF"`
}
