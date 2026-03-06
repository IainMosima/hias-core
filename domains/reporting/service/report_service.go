package service

import (
	"context"
	"encoding/json"

	"github.com/bitbiz/hias-core/domains/identity/schema"
	reportSchema "github.com/bitbiz/hias-core/domains/reporting/schema"
	"github.com/google/uuid"
)

type ReportService interface {
	// Definitions
	ListDefinitions(ctx context.Context, category, role string, page, pageSize int) *schema.ServiceResponse[[]reportSchema.ReportDefinitionResponse]
	GetDefinition(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[reportSchema.ReportDefinitionResponse]
	CreateAdHocDefinition(ctx context.Context, req reportSchema.CreateAdHocReportRequest, createdBy uuid.UUID) *schema.ServiceResponse[reportSchema.ReportDefinitionResponse]

	// Generate + Preview
	GenerateReport(ctx context.Context, req reportSchema.GenerateReportRequest, generatedBy uuid.UUID, role string) *schema.ServiceResponse[reportSchema.GeneratedReportResponse]
	PreviewReport(ctx context.Context, reportCode string, params json.RawMessage, role string) *schema.ServiceResponse[reportSchema.ReportPreviewResponse]
	DrillDown(ctx context.Context, req reportSchema.DrillDownRequest, generatedBy uuid.UUID, role string) *schema.ServiceResponse[reportSchema.GeneratedReportResponse]

	// Generated Reports
	ListGeneratedReports(ctx context.Context, defID *uuid.UUID, page, pageSize int, userID uuid.UUID) *schema.ServiceResponse[[]reportSchema.GeneratedReportResponse]
	GetGeneratedReport(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[reportSchema.GeneratedReportResponse]
	DownloadReport(ctx context.Context, id uuid.UUID) ([]byte, string, string, error) // data, format, reportNumber, error

	// Schedules
	CreateSchedule(ctx context.Context, req reportSchema.CreateScheduleRequest, createdBy uuid.UUID) *schema.ServiceResponse[reportSchema.ReportScheduleResponse]
	ListSchedules(ctx context.Context, defID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]reportSchema.ReportScheduleResponse]
	UpdateSchedule(ctx context.Context, id uuid.UUID, req reportSchema.UpdateScheduleRequest) *schema.ServiceResponse[reportSchema.ReportScheduleResponse]
	DeleteSchedule(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[string]

	// Management Dashboard
	GetManagementDashboard(ctx context.Context, period string) *schema.ServiceResponse[reportSchema.ManagementDashboardResponse]

	// Scheduled execution
	ExecuteDueSchedules(ctx context.Context) error
}
