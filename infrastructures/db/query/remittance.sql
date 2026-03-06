-- name: CreateRemittance :one
INSERT INTO remittances (provider_id, claim_ids, total_amount, currency, status, period_start, period_end, wht_rate, wht_amount, net_amount, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING *;

-- name: GetRemittanceByID :one
SELECT * FROM remittances WHERE id = $1;

-- name: ListRemittances :many
SELECT * FROM remittances ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: ListRemittancesByProvider :many
SELECT * FROM remittances WHERE provider_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListRemittancesByStatus :many
SELECT * FROM remittances WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: CountRemittances :one
SELECT COUNT(*) FROM remittances;

-- name: UpdateRemittanceStatus :one
UPDATE remittances SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: MarkRemittanceAdviceSent :one
UPDATE remittances SET remittance_advice_sent = TRUE, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: SetRemittancePayment :one
UPDATE remittances SET payment_id = $2, updated_at = NOW() WHERE id = $1 RETURNING *;
