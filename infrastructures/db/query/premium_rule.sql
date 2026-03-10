-- name: CreatePremiumRule :one
INSERT INTO premium_rules (plan_id, calculation_type, relationship, rate_amount, discount_type, discount_value, min_members, min_age, max_age, rule_type, effective_from, effective_to, sort_order)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING *;

-- name: GetPremiumRuleByID :one
SELECT * FROM premium_rules WHERE id = $1;

-- name: ListPremiumRulesByPlan :many
SELECT * FROM premium_rules WHERE plan_id = $1 ORDER BY sort_order, created_at;

-- name: ListEffectivePremiumRulesByPlan :many
SELECT * FROM premium_rules
WHERE plan_id = $1 AND effective_from <= $2 AND (effective_to IS NULL OR effective_to >= $3)
ORDER BY sort_order, created_at;

-- name: UpdatePremiumRule :one
UPDATE premium_rules SET
    calculation_type = COALESCE(sqlc.narg('calculation_type'), calculation_type),
    relationship = COALESCE(sqlc.narg('relationship'), relationship),
    rate_amount = COALESCE(sqlc.narg('rate_amount'), rate_amount),
    discount_type = COALESCE(sqlc.narg('discount_type'), discount_type),
    discount_value = COALESCE(sqlc.narg('discount_value'), discount_value),
    min_members = COALESCE(sqlc.narg('min_members'), min_members),
    min_age = COALESCE(sqlc.narg('min_age'), min_age),
    max_age = COALESCE(sqlc.narg('max_age'), max_age),
    rule_type = COALESCE(sqlc.narg('rule_type'), rule_type),
    effective_from = COALESCE(sqlc.narg('effective_from'), effective_from),
    effective_to = COALESCE(sqlc.narg('effective_to'), effective_to),
    sort_order = COALESCE(sqlc.narg('sort_order'), sort_order),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;

-- name: DeletePremiumRule :exec
DELETE FROM premium_rules WHERE id = $1;
