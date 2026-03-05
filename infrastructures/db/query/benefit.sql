-- name: CreateBenefit :one
INSERT INTO benefits (plan_id, name, category, annual_limit, co_pay_type, co_pay_value, waiting_period_days, sub_limit_type, sub_limit_value, min_age, max_age, waiting_period_type, deductible_amount)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING *;

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
    sub_limit_type = COALESCE(sqlc.narg('sub_limit_type'), sub_limit_type),
    sub_limit_value = COALESCE(sqlc.narg('sub_limit_value'), sub_limit_value),
    min_age = COALESCE(sqlc.narg('min_age'), min_age),
    max_age = COALESCE(sqlc.narg('max_age'), max_age),
    waiting_period_type = COALESCE(sqlc.narg('waiting_period_type'), waiting_period_type),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;

-- name: DeleteBenefit :exec
DELETE FROM benefits WHERE id = $1;

-- name: CreateBenefitWithParent :one
INSERT INTO benefits (plan_id, parent_benefit_id, name, category, annual_limit, co_pay_type, co_pay_value, waiting_period_days, sub_limit_type, sub_limit_value, min_age, max_age, waiting_period_type, deductible_amount)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING *;

-- name: ListSubBenefits :many
SELECT * FROM benefits WHERE parent_benefit_id = $1 ORDER BY name;
