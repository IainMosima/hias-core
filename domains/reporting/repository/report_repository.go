package repository

import (
	"context"
	"time"

	"github.com/bitbiz/hias-core/domains/reporting/entity"
	"github.com/google/uuid"
)

type ReportRepository interface {
	// Definitions
	CreateDefinition(ctx context.Context, def *entity.ReportDefinition) (*entity.ReportDefinition, error)
	GetDefinition(ctx context.Context, id uuid.UUID) (*entity.ReportDefinition, error)
	GetDefinitionByCode(ctx context.Context, code string) (*entity.ReportDefinition, error)
	ListDefinitions(ctx context.Context, category, reportType string, limit, offset int) ([]*entity.ReportDefinition, error)
	UpdateDefinition(ctx context.Context, def *entity.ReportDefinition) (*entity.ReportDefinition, error)
	CountDefinitions(ctx context.Context) (int64, error)

	// Schedules
	CreateSchedule(ctx context.Context, sched *entity.ReportSchedule) (*entity.ReportSchedule, error)
	GetSchedule(ctx context.Context, id uuid.UUID) (*entity.ReportSchedule, error)
	ListSchedulesByDefinition(ctx context.Context, defID uuid.UUID, limit, offset int) ([]*entity.ReportSchedule, error)
	ListDueSchedules(ctx context.Context) ([]*entity.ReportSchedule, error)
	UpdateSchedule(ctx context.Context, sched *entity.ReportSchedule) (*entity.ReportSchedule, error)
	UpdateScheduleLastRun(ctx context.Context, id uuid.UUID, lastRun, nextRun time.Time) error
	DeleteSchedule(ctx context.Context, id uuid.UUID) error

	// Generated Reports
	CreateGenerated(ctx context.Context, report *entity.GeneratedReport) (*entity.GeneratedReport, error)
	GetGenerated(ctx context.Context, id uuid.UUID) (*entity.GeneratedReport, error)
	ListGenerated(ctx context.Context, defID *uuid.UUID, status string, generatedBy *uuid.UUID, limit, offset int) ([]*entity.GeneratedReport, error)
	UpdateGeneratedStatus(ctx context.Context, id uuid.UUID, status string, rowCount int, fileSize int64, errorMsg string) (*entity.GeneratedReport, error)
	StoreReportFile(ctx context.Context, id uuid.UUID, data []byte, fileSize int64) error
	GetReportFile(ctx context.Context, id uuid.UUID) ([]byte, string, string, error) // data, format, reportNumber
	CountGenerated(ctx context.Context, defID *uuid.UUID) (int64, error)
	DeleteExpiredReports(ctx context.Context) error
}
