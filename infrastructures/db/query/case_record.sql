-- name: CreateCaseRecord :one
INSERT INTO case_records (case_number, preauth_id, policy_id, member_id, provider_id, status, admission_date, expected_discharge, diagnosis, treating_doctor, room_type, total_estimated_cost, notes, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING *;

-- name: GetCaseRecordByID :one
SELECT * FROM case_records WHERE id = $1;

-- name: GetCaseRecordByNumber :one
SELECT * FROM case_records WHERE case_number = $1;

-- name: GetCaseRecordByPreAuth :one
SELECT * FROM case_records WHERE preauth_id = $1;

-- name: ListCaseRecordsByPolicy :many
SELECT * FROM case_records WHERE policy_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListCaseRecordsByMember :many
SELECT * FROM case_records WHERE member_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListCaseRecordsByProvider :many
SELECT * FROM case_records WHERE provider_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListCaseRecordsByStatus :many
SELECT * FROM case_records WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: UpdateCaseRecord :one
UPDATE case_records SET
    diagnosis = COALESCE(sqlc.narg('diagnosis'), diagnosis),
    treating_doctor = COALESCE(sqlc.narg('treating_doctor'), treating_doctor),
    room_type = COALESCE(sqlc.narg('room_type'), room_type),
    total_estimated_cost = COALESCE(sqlc.narg('total_estimated_cost'), total_estimated_cost),
    notes = COALESCE(sqlc.narg('notes'), notes),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;

-- name: AdmitCaseRecord :one
UPDATE case_records SET
    status = 'ADMITTED',
    admission_date = $2,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: DischargeCaseRecord :one
UPDATE case_records SET
    status = 'DISCHARGED',
    actual_discharge = $2,
    total_actual_cost = $3,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: CloseCaseRecord :one
UPDATE case_records SET
    status = 'CLOSED',
    closed_at = NOW(),
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: UpdateCaseRecordStatus :one
UPDATE case_records SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: CountCaseRecordsByStatus :one
SELECT COUNT(*) FROM case_records WHERE status = $1;
