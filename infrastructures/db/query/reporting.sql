-- Report Definitions CRUD

-- name: CreateReportDefinition :one
INSERT INTO report_definitions (code, name, description, category, report_type, query_template, default_parameters, allowed_roles, columns, is_active, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING *;

-- name: GetReportDefinition :one
SELECT * FROM report_definitions WHERE id = $1;

-- name: GetReportDefinitionByCode :one
SELECT * FROM report_definitions WHERE code = $1;

-- name: ListReportDefinitions :many
SELECT * FROM report_definitions
WHERE (sqlc.narg('category')::varchar IS NULL OR category = sqlc.narg('category'))
AND (sqlc.narg('report_type')::varchar IS NULL OR report_type = sqlc.narg('report_type'))
AND is_active = true
ORDER BY category, name
LIMIT $1 OFFSET $2;

-- name: CountReportDefinitions :one
SELECT COUNT(*) FROM report_definitions WHERE is_active = true;

-- name: UpdateReportDefinition :one
UPDATE report_definitions SET
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    category = COALESCE(sqlc.narg('category'), category),
    default_parameters = COALESCE(sqlc.narg('default_parameters'), default_parameters),
    allowed_roles = COALESCE(sqlc.narg('allowed_roles'), allowed_roles),
    columns = COALESCE(sqlc.narg('columns'), columns),
    is_active = COALESCE(sqlc.narg('is_active'), is_active),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;

-- Report Schedules CRUD

-- name: CreateReportSchedule :one
INSERT INTO report_schedules (report_definition_id, name, cron_expression, parameters, export_format, recipients, is_active, next_run_at, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;

-- name: GetReportSchedule :one
SELECT * FROM report_schedules WHERE id = $1;

-- name: ListReportSchedulesByDefinition :many
SELECT * FROM report_schedules
WHERE report_definition_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListDueSchedules :many
SELECT * FROM report_schedules
WHERE is_active = true AND next_run_at <= NOW()
ORDER BY next_run_at;

-- name: UpdateReportSchedule :one
UPDATE report_schedules SET
    name = COALESCE(sqlc.narg('name'), name),
    cron_expression = COALESCE(sqlc.narg('cron_expression'), cron_expression),
    parameters = COALESCE(sqlc.narg('parameters'), parameters),
    export_format = COALESCE(sqlc.narg('export_format'), export_format),
    recipients = COALESCE(sqlc.narg('recipients'), recipients),
    is_active = COALESCE(sqlc.narg('is_active'), is_active),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;

-- name: UpdateScheduleLastRun :exec
UPDATE report_schedules SET
    last_run_at = $2,
    next_run_at = $3,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteReportSchedule :exec
DELETE FROM report_schedules WHERE id = $1;

-- Generated Reports CRUD

-- name: CreateGeneratedReport :one
INSERT INTO generated_reports (report_definition_id, schedule_id, report_number, name, parameters, format, status, generated_by, expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;

-- name: GetGeneratedReport :one
SELECT id, report_definition_id, schedule_id, report_number, name, parameters, format, status, row_count, file_size, error_message, generated_by, generated_at, expires_at, created_at
FROM generated_reports WHERE id = $1;

-- name: GetGeneratedReportWithData :one
SELECT * FROM generated_reports WHERE id = $1;

-- name: ListGeneratedReports :many
SELECT id, report_definition_id, schedule_id, report_number, name, parameters, format, status, row_count, file_size, error_message, generated_by, generated_at, expires_at, created_at
FROM generated_reports
WHERE (sqlc.narg('report_definition_id')::uuid IS NULL OR report_definition_id = sqlc.narg('report_definition_id'))
AND (sqlc.narg('status')::varchar IS NULL OR status = sqlc.narg('status'))
AND (sqlc.narg('generated_by')::uuid IS NULL OR generated_by = sqlc.narg('generated_by'))
ORDER BY generated_at DESC
LIMIT $1 OFFSET $2;

-- name: CountGeneratedReports :one
SELECT COUNT(*) FROM generated_reports
WHERE (sqlc.narg('report_definition_id')::uuid IS NULL OR report_definition_id = sqlc.narg('report_definition_id'));

-- name: UpdateGeneratedReportStatus :one
UPDATE generated_reports SET
    status = $2,
    row_count = $3,
    file_size = $4,
    error_message = $5
WHERE id = $1 RETURNING id, report_definition_id, schedule_id, report_number, name, parameters, format, status, row_count, file_size, error_message, generated_by, generated_at, expires_at, created_at;

-- name: StoreReportFile :exec
UPDATE generated_reports SET
    file_data = $2,
    file_size = $3,
    status = 'COMPLETED'
WHERE id = $1;

-- name: GetReportFileData :one
SELECT file_data, format, report_number FROM generated_reports WHERE id = $1;

-- name: DeleteExpiredReports :exec
DELETE FROM generated_reports WHERE expires_at IS NOT NULL AND expires_at < NOW();
