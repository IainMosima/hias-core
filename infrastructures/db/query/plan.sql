-- name: CreatePlan :one
INSERT INTO plans (name, type, segment, base_premium, currency, status, description, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: GetPlanByID :one
SELECT * FROM plans WHERE id = $1;

-- name: ListPlans :many
SELECT * FROM plans ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: ListPlansByStatus :many
SELECT * FROM plans WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListPlansBySegment :many
SELECT * FROM plans WHERE segment = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: CountPlans :one
SELECT COUNT(*) FROM plans;

-- name: UpdatePlan :one
UPDATE plans SET
    name = COALESCE(sqlc.narg('name'), name),
    type = COALESCE(sqlc.narg('type'), type),
    segment = COALESCE(sqlc.narg('segment'), segment),
    base_premium = COALESCE(sqlc.narg('base_premium'), base_premium),
    description = COALESCE(sqlc.narg('description'), description),
    status = COALESCE(sqlc.narg('status'), status),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;
