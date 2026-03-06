package reporting

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bitbiz/hias-core/domains/reporting/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/reporting/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type reportRepository struct {
	store db.Store
}

func NewReportRepository(store db.Store) domainRepo.ReportRepository {
	return &reportRepository{store: store}
}

// --- Definitions ---

func (r *reportRepository) CreateDefinition(ctx context.Context, def *entity.ReportDefinition) (*entity.ReportDefinition, error) {
	dbDef, err := r.store.CreateReportDefinition(ctx, db.CreateReportDefinitionParams{
		Code:              def.Code,
		Name:              def.Name,
		Description:       pgtype.Text{String: def.Description, Valid: def.Description != ""},
		Category:          def.Category,
		ReportType:        def.ReportType,
		QueryTemplate:     pgtype.Text{String: def.QueryTemplate, Valid: def.QueryTemplate != ""},
		DefaultParameters: def.DefaultParameters,
		AllowedRoles:      def.AllowedRoles,
		Columns:           def.Columns,
		IsActive:          pgtype.Bool{Bool: def.IsActive, Valid: true},
		CreatedBy:         uuidToPgtype(def.CreatedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create report definition: %w", err)
	}
	return sqlcDefinitionToDomain(dbDef), nil
}

func (r *reportRepository) GetDefinition(ctx context.Context, id uuid.UUID) (*entity.ReportDefinition, error) {
	dbDef, err := r.store.GetReportDefinition(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get report definition: %w", err)
	}
	return sqlcDefinitionToDomain(dbDef), nil
}

func (r *reportRepository) GetDefinitionByCode(ctx context.Context, code string) (*entity.ReportDefinition, error) {
	dbDef, err := r.store.GetReportDefinitionByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get report definition by code: %w", err)
	}
	return sqlcDefinitionToDomain(dbDef), nil
}

func (r *reportRepository) ListDefinitions(ctx context.Context, category, reportType string, limit, offset int) ([]*entity.ReportDefinition, error) {
	dbDefs, err := r.store.ListReportDefinitions(ctx, db.ListReportDefinitionsParams{
		Limit:      int32(limit),
		Offset:     int32(offset),
		Category:   pgtype.Text{String: category, Valid: category != ""},
		ReportType: pgtype.Text{String: reportType, Valid: reportType != ""},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list report definitions: %w", err)
	}
	defs := make([]*entity.ReportDefinition, len(dbDefs))
	for i, d := range dbDefs {
		defs[i] = sqlcDefinitionToDomain(d)
	}
	return defs, nil
}

func (r *reportRepository) UpdateDefinition(ctx context.Context, def *entity.ReportDefinition) (*entity.ReportDefinition, error) {
	dbDef, err := r.store.UpdateReportDefinition(ctx, db.UpdateReportDefinitionParams{
		ID:                def.ID,
		Name:              pgtype.Text{String: def.Name, Valid: def.Name != ""},
		Description:       pgtype.Text{String: def.Description, Valid: def.Description != ""},
		Category:          pgtype.Text{String: def.Category, Valid: def.Category != ""},
		DefaultParameters: def.DefaultParameters,
		AllowedRoles:      def.AllowedRoles,
		Columns:           def.Columns,
		IsActive:          pgtype.Bool{Bool: def.IsActive, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update report definition: %w", err)
	}
	return sqlcDefinitionToDomain(dbDef), nil
}

func (r *reportRepository) CountDefinitions(ctx context.Context) (int64, error) {
	count, err := r.store.CountReportDefinitions(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count report definitions: %w", err)
	}
	return count, nil
}

// --- Schedules ---

func (r *reportRepository) CreateSchedule(ctx context.Context, sched *entity.ReportSchedule) (*entity.ReportSchedule, error) {
	var nextRunAt pgtype.Timestamptz
	if sched.NextRunAt != nil {
		nextRunAt = pgtype.Timestamptz{Time: *sched.NextRunAt, Valid: true}
	}

	dbSched, err := r.store.CreateReportSchedule(ctx, db.CreateReportScheduleParams{
		ReportDefinitionID: sched.ReportDefinitionID,
		Name:               sched.Name,
		CronExpression:     sched.CronExpression,
		Parameters:         sched.Parameters,
		ExportFormat:       sched.ExportFormat,
		Recipients:         sched.Recipients,
		IsActive:           pgtype.Bool{Bool: sched.IsActive, Valid: true},
		NextRunAt:          nextRunAt,
		CreatedBy:          uuidToPgtype(sched.CreatedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create report schedule: %w", err)
	}
	return sqlcScheduleToDomain(dbSched), nil
}

func (r *reportRepository) GetSchedule(ctx context.Context, id uuid.UUID) (*entity.ReportSchedule, error) {
	dbSched, err := r.store.GetReportSchedule(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get report schedule: %w", err)
	}
	return sqlcScheduleToDomain(dbSched), nil
}

func (r *reportRepository) ListSchedulesByDefinition(ctx context.Context, defID uuid.UUID, limit, offset int) ([]*entity.ReportSchedule, error) {
	dbScheds, err := r.store.ListReportSchedulesByDefinition(ctx, db.ListReportSchedulesByDefinitionParams{
		ReportDefinitionID: defID,
		Limit:              int32(limit),
		Offset:             int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list report schedules: %w", err)
	}
	scheds := make([]*entity.ReportSchedule, len(dbScheds))
	for i, s := range dbScheds {
		scheds[i] = sqlcScheduleToDomain(s)
	}
	return scheds, nil
}

func (r *reportRepository) ListDueSchedules(ctx context.Context) ([]*entity.ReportSchedule, error) {
	dbScheds, err := r.store.ListDueSchedules(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list due schedules: %w", err)
	}
	scheds := make([]*entity.ReportSchedule, len(dbScheds))
	for i, s := range dbScheds {
		scheds[i] = sqlcScheduleToDomain(s)
	}
	return scheds, nil
}

func (r *reportRepository) UpdateSchedule(ctx context.Context, sched *entity.ReportSchedule) (*entity.ReportSchedule, error) {
	dbSched, err := r.store.UpdateReportSchedule(ctx, db.UpdateReportScheduleParams{
		ID:             sched.ID,
		Name:           pgtype.Text{String: sched.Name, Valid: sched.Name != ""},
		CronExpression: pgtype.Text{String: sched.CronExpression, Valid: sched.CronExpression != ""},
		Parameters:     sched.Parameters,
		ExportFormat:   pgtype.Text{String: sched.ExportFormat, Valid: sched.ExportFormat != ""},
		Recipients:     sched.Recipients,
		IsActive:       pgtype.Bool{Bool: sched.IsActive, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update report schedule: %w", err)
	}
	return sqlcScheduleToDomain(dbSched), nil
}

func (r *reportRepository) UpdateScheduleLastRun(ctx context.Context, id uuid.UUID, lastRun, nextRun time.Time) error {
	err := r.store.UpdateScheduleLastRun(ctx, db.UpdateScheduleLastRunParams{
		ID:        id,
		LastRunAt: pgtype.Timestamptz{Time: lastRun, Valid: true},
		NextRunAt: pgtype.Timestamptz{Time: nextRun, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to update schedule last run: %w", err)
	}
	return nil
}

func (r *reportRepository) DeleteSchedule(ctx context.Context, id uuid.UUID) error {
	err := r.store.DeleteReportSchedule(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete report schedule: %w", err)
	}
	return nil
}

// --- Generated Reports ---

func (r *reportRepository) CreateGenerated(ctx context.Context, report *entity.GeneratedReport) (*entity.GeneratedReport, error) {
	var schedID pgtype.UUID
	if report.ScheduleID != nil {
		schedID = pgtype.UUID{Bytes: *report.ScheduleID, Valid: true}
	}
	var expiresAt pgtype.Timestamptz
	if report.ExpiresAt != nil {
		expiresAt = pgtype.Timestamptz{Time: *report.ExpiresAt, Valid: true}
	}

	dbReport, err := r.store.CreateGeneratedReport(ctx, db.CreateGeneratedReportParams{
		ReportDefinitionID: report.ReportDefinitionID,
		ScheduleID:         schedID,
		ReportNumber:       report.ReportNumber,
		Name:               report.Name,
		Parameters:         report.Parameters,
		Format:             report.Format,
		Status:             report.Status,
		GeneratedBy:        uuidToPgtype(report.GeneratedBy),
		ExpiresAt:          expiresAt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create generated report: %w", err)
	}
	return sqlcGeneratedToDomain(dbReport), nil
}

func (r *reportRepository) GetGenerated(ctx context.Context, id uuid.UUID) (*entity.GeneratedReport, error) {
	dbRow, err := r.store.GetGeneratedReport(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get generated report: %w", err)
	}
	return sqlcGeneratedRowToDomain(dbRow), nil
}

func (r *reportRepository) ListGenerated(ctx context.Context, defID *uuid.UUID, status string, generatedBy *uuid.UUID, limit, offset int) ([]*entity.GeneratedReport, error) {
	var defIDPg pgtype.UUID
	if defID != nil {
		defIDPg = pgtype.UUID{Bytes: *defID, Valid: true}
	}
	var genByPg pgtype.UUID
	if generatedBy != nil {
		genByPg = pgtype.UUID{Bytes: *generatedBy, Valid: true}
	}

	dbRows, err := r.store.ListGeneratedReports(ctx, db.ListGeneratedReportsParams{
		Limit:              int32(limit),
		Offset:             int32(offset),
		ReportDefinitionID: defIDPg,
		Status:             pgtype.Text{String: status, Valid: status != ""},
		GeneratedBy:        genByPg,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list generated reports: %w", err)
	}
	reports := make([]*entity.GeneratedReport, len(dbRows))
	for i, row := range dbRows {
		reports[i] = sqlcListGeneratedRowToDomain(row)
	}
	return reports, nil
}

func (r *reportRepository) UpdateGeneratedStatus(ctx context.Context, id uuid.UUID, status string, rowCount int, fileSize int64, errorMsg string) (*entity.GeneratedReport, error) {
	dbRow, err := r.store.UpdateGeneratedReportStatus(ctx, db.UpdateGeneratedReportStatusParams{
		ID:           id,
		Status:       status,
		RowCount:     pgtype.Int4{Int32: int32(rowCount), Valid: true},
		FileSize:     pgtype.Int8{Int64: fileSize, Valid: true},
		ErrorMessage: pgtype.Text{String: errorMsg, Valid: errorMsg != ""},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update generated report status: %w", err)
	}
	return sqlcUpdateStatusRowToDomain(dbRow), nil
}

func (r *reportRepository) StoreReportFile(ctx context.Context, id uuid.UUID, data []byte, fileSize int64) error {
	err := r.store.StoreReportFile(ctx, db.StoreReportFileParams{
		ID:       id,
		FileData: data,
		FileSize: pgtype.Int8{Int64: fileSize, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to store report file: %w", err)
	}
	return nil
}

func (r *reportRepository) GetReportFile(ctx context.Context, id uuid.UUID) ([]byte, string, string, error) {
	row, err := r.store.GetReportFileData(ctx, id)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to get report file: %w", err)
	}
	return row.FileData, row.Format, row.ReportNumber, nil
}

func (r *reportRepository) CountGenerated(ctx context.Context, defID *uuid.UUID) (int64, error) {
	var defIDPg pgtype.UUID
	if defID != nil {
		defIDPg = pgtype.UUID{Bytes: *defID, Valid: true}
	}
	count, err := r.store.CountGeneratedReports(ctx, defIDPg)
	if err != nil {
		return 0, fmt.Errorf("failed to count generated reports: %w", err)
	}
	return count, nil
}

func (r *reportRepository) DeleteExpiredReports(ctx context.Context) error {
	err := r.store.DeleteExpiredReports(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete expired reports: %w", err)
	}
	return nil
}

// --- Conversion helpers ---

func uuidToPgtype(id uuid.UUID) pgtype.UUID {
	if id == uuid.Nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: id, Valid: true}
}

func pgtypeToUUID(id pgtype.UUID) uuid.UUID {
	if !id.Valid {
		return uuid.Nil
	}
	return uuid.UUID(id.Bytes)
}

func pgtypeToUUIDPtr(id pgtype.UUID) *uuid.UUID {
	if !id.Valid {
		return nil
	}
	u := uuid.UUID(id.Bytes)
	return &u
}

func pgtypeTimestamptzToTimePtr(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	result := t.Time
	return &result
}

func pgtypeTimestamptzToTime(t pgtype.Timestamptz) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

func pgtypeTextToString(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

func pgtypeBoolToBool(b pgtype.Bool) bool {
	if !b.Valid {
		return false
	}
	return b.Bool
}

func sqlcDefinitionToDomain(d db.ReportDefinition) *entity.ReportDefinition {
	return &entity.ReportDefinition{
		ID:                d.ID,
		Code:              d.Code,
		Name:              d.Name,
		Description:       pgtypeTextToString(d.Description),
		Category:          d.Category,
		ReportType:        d.ReportType,
		QueryTemplate:     pgtypeTextToString(d.QueryTemplate),
		DefaultParameters: json.RawMessage(d.DefaultParameters),
		AllowedRoles:      d.AllowedRoles,
		Columns:           d.Columns,
		IsActive:          pgtypeBoolToBool(d.IsActive),
		CreatedBy:         pgtypeToUUID(d.CreatedBy),
		CreatedAt:         pgtypeTimestamptzToTime(d.CreatedAt),
		UpdatedAt:         pgtypeTimestamptzToTime(d.UpdatedAt),
	}
}

func sqlcScheduleToDomain(s db.ReportSchedule) *entity.ReportSchedule {
	return &entity.ReportSchedule{
		ID:                 s.ID,
		ReportDefinitionID: s.ReportDefinitionID,
		Name:               s.Name,
		CronExpression:     s.CronExpression,
		Parameters:         json.RawMessage(s.Parameters),
		ExportFormat:       s.ExportFormat,
		Recipients:         s.Recipients,
		IsActive:           pgtypeBoolToBool(s.IsActive),
		LastRunAt:          pgtypeTimestamptzToTimePtr(s.LastRunAt),
		NextRunAt:          pgtypeTimestamptzToTimePtr(s.NextRunAt),
		CreatedBy:          pgtypeToUUID(s.CreatedBy),
		CreatedAt:          pgtypeTimestamptzToTime(s.CreatedAt),
		UpdatedAt:          pgtypeTimestamptzToTime(s.UpdatedAt),
	}
}

func sqlcGeneratedToDomain(g db.GeneratedReport) *entity.GeneratedReport {
	return &entity.GeneratedReport{
		ID:                 g.ID,
		ReportDefinitionID: g.ReportDefinitionID,
		ScheduleID:         pgtypeToUUIDPtr(g.ScheduleID),
		ReportNumber:       g.ReportNumber,
		Name:               g.Name,
		Parameters:         json.RawMessage(g.Parameters),
		Format:             g.Format,
		Status:             g.Status,
		RowCount:           int(g.RowCount.Int32),
		FileSize:           g.FileSize.Int64,
		ErrorMessage:       pgtypeTextToString(g.ErrorMessage),
		GeneratedBy:        pgtypeToUUID(g.GeneratedBy),
		GeneratedAt:        pgtypeTimestamptzToTime(g.GeneratedAt),
		ExpiresAt:          pgtypeTimestamptzToTimePtr(g.ExpiresAt),
		CreatedAt:          pgtypeTimestamptzToTime(g.CreatedAt),
	}
}

func sqlcGeneratedRowToDomain(g db.GetGeneratedReportRow) *entity.GeneratedReport {
	return &entity.GeneratedReport{
		ID:                 g.ID,
		ReportDefinitionID: g.ReportDefinitionID,
		ScheduleID:         pgtypeToUUIDPtr(g.ScheduleID),
		ReportNumber:       g.ReportNumber,
		Name:               g.Name,
		Parameters:         json.RawMessage(g.Parameters),
		Format:             g.Format,
		Status:             g.Status,
		RowCount:           int(g.RowCount.Int32),
		FileSize:           g.FileSize.Int64,
		ErrorMessage:       pgtypeTextToString(g.ErrorMessage),
		GeneratedBy:        pgtypeToUUID(g.GeneratedBy),
		GeneratedAt:        pgtypeTimestamptzToTime(g.GeneratedAt),
		ExpiresAt:          pgtypeTimestamptzToTimePtr(g.ExpiresAt),
		CreatedAt:          pgtypeTimestamptzToTime(g.CreatedAt),
	}
}

func sqlcListGeneratedRowToDomain(g db.ListGeneratedReportsRow) *entity.GeneratedReport {
	return &entity.GeneratedReport{
		ID:                 g.ID,
		ReportDefinitionID: g.ReportDefinitionID,
		ScheduleID:         pgtypeToUUIDPtr(g.ScheduleID),
		ReportNumber:       g.ReportNumber,
		Name:               g.Name,
		Parameters:         json.RawMessage(g.Parameters),
		Format:             g.Format,
		Status:             g.Status,
		RowCount:           int(g.RowCount.Int32),
		FileSize:           g.FileSize.Int64,
		ErrorMessage:       pgtypeTextToString(g.ErrorMessage),
		GeneratedBy:        pgtypeToUUID(g.GeneratedBy),
		GeneratedAt:        pgtypeTimestamptzToTime(g.GeneratedAt),
		ExpiresAt:          pgtypeTimestamptzToTimePtr(g.ExpiresAt),
		CreatedAt:          pgtypeTimestamptzToTime(g.CreatedAt),
	}
}

func sqlcUpdateStatusRowToDomain(g db.UpdateGeneratedReportStatusRow) *entity.GeneratedReport {
	return &entity.GeneratedReport{
		ID:                 g.ID,
		ReportDefinitionID: g.ReportDefinitionID,
		ScheduleID:         pgtypeToUUIDPtr(g.ScheduleID),
		ReportNumber:       g.ReportNumber,
		Name:               g.Name,
		Parameters:         json.RawMessage(g.Parameters),
		Format:             g.Format,
		Status:             g.Status,
		RowCount:           int(g.RowCount.Int32),
		FileSize:           g.FileSize.Int64,
		ErrorMessage:       pgtypeTextToString(g.ErrorMessage),
		GeneratedBy:        pgtypeToUUID(g.GeneratedBy),
		GeneratedAt:        pgtypeTimestamptzToTime(g.GeneratedAt),
		ExpiresAt:          pgtypeTimestamptzToTimePtr(g.ExpiresAt),
		CreatedAt:          pgtypeTimestamptzToTime(g.CreatedAt),
	}
}
