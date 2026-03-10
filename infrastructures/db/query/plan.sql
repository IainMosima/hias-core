-- name: CreatePlan :one
INSERT INTO plans (name, type, segment, base_premium, premium_frequency, currency, status, description, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;

-- name: GetPlanByID :one
SELECT * FROM plans WHERE id = $1;

-- name: ListPlans :many
SELECT * FROM plans ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: ListPlansByStatus :many
SELECT * FROM plans WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListPlansBySegment :many
SELECT * FROM plans WHERE segment = $1 AND status != 'INACTIVE' ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: CountPlans :one
SELECT COUNT(*) FROM plans;

-- name: CountPlansByStatus :one
SELECT COUNT(*) FROM plans WHERE status = $1;

-- name: UpdatePlan :one
UPDATE plans SET
    name = COALESCE(sqlc.narg('name'), name),
    type = COALESCE(sqlc.narg('type'), type),
    segment = COALESCE(sqlc.narg('segment'), segment),
    base_premium = COALESCE(sqlc.narg('base_premium'), base_premium),
    premium_frequency = COALESCE(sqlc.narg('premium_frequency'), premium_frequency),
    description = COALESCE(sqlc.narg('description'), description),
    status = COALESCE(sqlc.narg('status'), status),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;

-- name: SoftDeletePlan :one
UPDATE plans SET status = 'INACTIVE', updated_at = NOW() WHERE id = $1 RETURNING *;
