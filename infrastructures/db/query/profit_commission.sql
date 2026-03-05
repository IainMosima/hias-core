-- name: CreateProfitCommission :one
INSERT INTO profit_commissions (treaty_id, commission_type, loss_ratio_from, loss_ratio_to, commission_rate, carry_forward_years, carry_forward_balance, period_start, period_end, calculated_amount)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *;

-- name: GetProfitCommissionByID :one
SELECT * FROM profit_commissions WHERE id = $1;

-- name: ListProfitCommissionsByTreaty :many
SELECT * FROM profit_commissions WHERE treaty_id = $1 ORDER BY loss_ratio_from ASC;

-- name: UpdateProfitCommission :one
UPDATE profit_commissions SET
    commission_type = $2, loss_ratio_from = $3, loss_ratio_to = $4,
    commission_rate = $5, carry_forward_years = $6, carry_forward_balance = $7,
    period_start = $8, period_end = $9, calculated_amount = $10, updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: DeleteProfitCommission :exec
DELETE FROM profit_commissions WHERE id = $1;
