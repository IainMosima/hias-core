-- name: CreateInstallmentSchedule :one
INSERT INTO installment_schedules (policy_id, frequency, total_installments, amount_per_installment, start_date, status, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetInstallmentScheduleByID :one
SELECT * FROM installment_schedules WHERE id = $1;

-- name: GetInstallmentScheduleByPolicy :many
SELECT * FROM installment_schedules WHERE policy_id = $1 ORDER BY created_at DESC;

-- name: UpdateInstallmentScheduleStatus :one
UPDATE installment_schedules SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING *;
