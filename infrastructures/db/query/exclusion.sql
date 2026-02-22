-- name: CreateExclusion :one
INSERT INTO exclusions (plan_id, description, type, icd_codes)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetExclusionByID :one
SELECT * FROM exclusions WHERE id = $1;

-- name: ListExclusionsByPlan :many
SELECT * FROM exclusions WHERE plan_id = $1 ORDER BY type;

-- name: ListExclusionsByType :many
SELECT * FROM exclusions WHERE plan_id = $1 AND type = $2;

-- name: UpdateExclusion :one
UPDATE exclusions SET
    description = COALESCE(sqlc.narg('description'), description),
    type = COALESCE(sqlc.narg('type'), type),
    icd_codes = COALESCE(sqlc.narg('icd_codes'), icd_codes),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;

-- name: DeleteExclusion :exec
DELETE FROM exclusions WHERE id = $1;
