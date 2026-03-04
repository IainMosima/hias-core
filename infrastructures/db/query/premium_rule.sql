-- name: CreatePremiumRule :one
INSERT INTO premium_rules (plan_id, calculation_type, relationship, rate_amount, discount_type, discount_value, min_members)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetPremiumRuleByID :one
SELECT * FROM premium_rules WHERE id = $1;

-- name: ListPremiumRulesByPlan :many
SELECT * FROM premium_rules WHERE plan_id = $1 ORDER BY created_at;

-- name: UpdatePremiumRule :one
UPDATE premium_rules SET
    calculation_type = COALESCE(sqlc.narg('calculation_type'), calculation_type),
    relationship = COALESCE(sqlc.narg('relationship'), relationship),
    rate_amount = COALESCE(sqlc.narg('rate_amount'), rate_amount),
    discount_type = COALESCE(sqlc.narg('discount_type'), discount_type),
    discount_value = COALESCE(sqlc.narg('discount_value'), discount_value),
    min_members = COALESCE(sqlc.narg('min_members'), min_members),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;

-- name: DeletePremiumRule :exec
DELETE FROM premium_rules WHERE id = $1;
