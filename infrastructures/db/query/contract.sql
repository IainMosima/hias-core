-- name: CreateContract :one
INSERT INTO contracts (provider_id, start_date, end_date, terms, status, created_by)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetContractByID :one
SELECT * FROM contracts WHERE id = $1;

-- name: ListContractsByProvider :many
SELECT * FROM contracts WHERE provider_id = $1 ORDER BY start_date DESC;

-- name: GetActiveContractByProvider :one
SELECT * FROM contracts WHERE provider_id = $1 AND status = 'ACTIVE' ORDER BY start_date DESC LIMIT 1;

-- name: UpdateContractStatus :one
UPDATE contracts SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING *;
