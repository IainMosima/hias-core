-- name: CreateCommissionRule :one
INSERT INTO commission_rules (plan_id, intermediary_id, rate_pct, flat_amount, effective_from, effective_to, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetCommissionRuleByID :one
SELECT * FROM commission_rules WHERE id = $1;

-- name: ListCommissionRulesByPlan :many
SELECT * FROM commission_rules WHERE plan_id = $1 ORDER BY created_at DESC;

-- name: CreateCommissionPayment :one
INSERT INTO commission_payments (policy_id, intermediary_id, commission_rule_id, amount, currency, status, period_start, period_end, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: ListCommissionPayments :many
SELECT * FROM commission_payments ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: ListCommissionPaymentsByIntermediary :many
SELECT * FROM commission_payments WHERE intermediary_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: UpdateCommissionPaymentStatus :one
UPDATE commission_payments SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING *;
