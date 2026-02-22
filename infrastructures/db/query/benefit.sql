-- name: CreateBenefit :one
INSERT INTO benefits (plan_id, name, category, annual_limit, co_pay_type, co_pay_value, waiting_period_days)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetBenefitByID :one
SELECT * FROM benefits WHERE id = $1;

-- name: ListBenefitsByPlan :many
SELECT * FROM benefits WHERE plan_id = $1 ORDER BY category, name;

-- name: ListBenefitsByCategory :many
SELECT * FROM benefits WHERE plan_id = $1 AND category = $2;

-- name: UpdateBenefit :one
UPDATE benefits SET
    name = COALESCE(sqlc.narg('name'), name),
    category = COALESCE(sqlc.narg('category'), category),
    annual_limit = COALESCE(sqlc.narg('annual_limit'), annual_limit),
    co_pay_type = COALESCE(sqlc.narg('co_pay_type'), co_pay_type),
    co_pay_value = COALESCE(sqlc.narg('co_pay_value'), co_pay_value),
    waiting_period_days = COALESCE(sqlc.narg('waiting_period_days'), waiting_period_days),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;

-- name: DeleteBenefit :exec
DELETE FROM benefits WHERE id = $1;
