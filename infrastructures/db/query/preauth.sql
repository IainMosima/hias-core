-- name: CreatePreAuth :one
INSERT INTO preauthorizations (policy_id, member_id, provider_id, procedure_codes, diagnosis_codes, estimated_cost, status, notes, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;

-- name: GetPreAuthByID :one
SELECT * FROM preauthorizations WHERE id = $1;

-- name: GetPreAuthByAuthCode :one
SELECT * FROM preauthorizations WHERE auth_code = $1;

-- name: ListPreAuths :many
SELECT * FROM preauthorizations ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: ListPreAuthsByStatus :many
SELECT * FROM preauthorizations WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListPreAuthsByProvider :many
SELECT * FROM preauthorizations WHERE provider_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListPreAuthsByMember :many
SELECT * FROM preauthorizations WHERE member_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: CountPreAuths :one
SELECT COUNT(*) FROM preauthorizations;

-- name: UpdatePreAuthStatus :one
UPDATE preauthorizations SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: ApprovePreAuth :one
UPDATE preauthorizations SET
    status = 'APPROVED',
    auth_code = $2,
    approved_amount = $3,
    validity_start = $4,
    validity_end = $5,
    reviewed_by = $6,
    reviewed_at = NOW(),
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: DenyPreAuth :one
UPDATE preauthorizations SET
    status = 'DENIED',
    denial_reason = $2,
    reviewed_by = $3,
    reviewed_at = NOW(),
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: GetExpiringPreAuths :many
SELECT * FROM preauthorizations WHERE status = 'APPROVED' AND validity_end < NOW();
