-- name: CreateUnderwritingFlag :one
INSERT INTO underwriting_flags (assessment_id, policy_id, member_id, flag_type, severity, details, status)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetUnderwritingFlagByID :one
SELECT * FROM underwriting_flags WHERE id = $1;

-- name: ListUnderwritingFlagsByPolicy :many
SELECT * FROM underwriting_flags WHERE policy_id = $1 ORDER BY created_at DESC;

-- name: ListUnderwritingFlagsByMember :many
SELECT * FROM underwriting_flags WHERE member_id = $1 ORDER BY created_at DESC;

-- name: ListUnderwritingFlagsByAssessment :many
SELECT * FROM underwriting_flags WHERE assessment_id = $1 ORDER BY created_at DESC;

-- name: ListOpenUnderwritingFlags :many
SELECT * FROM underwriting_flags WHERE status = 'OPEN' ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: CountOpenUnderwritingFlags :one
SELECT COUNT(*) FROM underwriting_flags WHERE status = 'OPEN';

-- name: ResolveUnderwritingFlag :one
UPDATE underwriting_flags SET
    status = 'RESOLVED',
    resolved_by = sqlc.arg('resolved_by'),
    resolved_at = NOW(),
    resolution = sqlc.arg('resolution'),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;

-- name: OverrideUnderwritingFlag :one
UPDATE underwriting_flags SET
    status = 'OVERRIDDEN',
    resolved_by = sqlc.arg('resolved_by'),
    resolved_at = NOW(),
    resolution = sqlc.arg('resolution'),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;
