-- name: CreateEndorsement :one
INSERT INTO endorsements (policy_id, endorsement_type, status, effective_date, changes, reason, premium_adjustment, requested_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: GetEndorsementByID :one
SELECT * FROM endorsements WHERE id = $1;

-- name: ListEndorsementsByPolicy :many
SELECT * FROM endorsements WHERE policy_id = $1 ORDER BY created_at DESC;

-- name: UpdateEndorsementStatus :one
UPDATE endorsements SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: UpdateEndorsement :one
UPDATE endorsements SET
    status = COALESCE(sqlc.narg('status'), status),
    approved_by = COALESCE(sqlc.narg('approved_by'), approved_by),
    approved_at = COALESCE(sqlc.narg('approved_at'), approved_at),
    applied_at = COALESCE(sqlc.narg('applied_at'), applied_at),
    reason = COALESCE(sqlc.narg('reason'), reason),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;
