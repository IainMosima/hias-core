-- name: CreateBordereau :one
INSERT INTO bordereaux (bordereau_number, treaty_id, bordereau_type, period_start, period_end, total_gross, total_ceded, total_commission, item_count, status, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING *;

-- name: GetBordereauByID :one
SELECT * FROM bordereaux WHERE id = $1;

-- name: GetBordereauByNumber :one
SELECT * FROM bordereaux WHERE bordereau_number = $1;

-- name: ListBordereauxByTreaty :many
SELECT * FROM bordereaux WHERE treaty_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: UpdateBordereau :one
UPDATE bordereaux SET
    total_gross = $2, total_ceded = $3, total_commission = $4,
    item_count = $5, updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: UpdateBordereauStatus :one
UPDATE bordereaux SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: CountBordereaux :one
SELECT COUNT(*) FROM bordereaux;
