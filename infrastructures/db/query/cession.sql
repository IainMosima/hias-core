-- name: CreateCession :one
INSERT INTO cessions (cession_number, treaty_id, policy_id, treaty_layer_id, cession_type, gross_amount, ceded_amount, retained_amount, commission_amount, share_percentage, status, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING *;

-- name: GetCessionByID :one
SELECT * FROM cessions WHERE id = $1;

-- name: GetCessionByNumber :one
SELECT * FROM cessions WHERE cession_number = $1;

-- name: ListCessionsByTreaty :many
SELECT * FROM cessions WHERE treaty_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListCessionsByPolicy :many
SELECT * FROM cessions WHERE policy_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListCessionsByTreatyAndPeriod :many
SELECT * FROM cessions
WHERE treaty_id = $1 AND created_at >= $2 AND created_at <= $3
ORDER BY created_at DESC LIMIT $4 OFFSET $5;

-- name: ListBookedCessionsByTreatyAndPeriod :many
SELECT * FROM cessions
WHERE treaty_id = $1 AND status = 'BOOKED' AND created_at >= $2 AND created_at <= $3
ORDER BY created_at DESC;

-- name: UpdateCessionStatus :one
UPDATE cessions SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: GetTotalCededByTreaty :one
SELECT COALESCE(SUM(ceded_amount), 0)::bigint as total_ceded
FROM cessions WHERE treaty_id = $1 AND status = 'BOOKED';

-- name: GetTotalCededByTreatyAndPeriod :one
SELECT COALESCE(SUM(ceded_amount), 0)::bigint as total_ceded
FROM cessions WHERE treaty_id = $1 AND status = 'BOOKED' AND created_at >= $2 AND created_at <= $3;

-- name: CountCessions :one
SELECT COUNT(*) FROM cessions;

-- name: GetTotalCededAmountAll :one
SELECT COALESCE(SUM(ceded_amount), 0)::bigint as total_ceded
FROM cessions WHERE status = 'BOOKED';

-- name: GetTotalGrossAmountAll :one
SELECT COALESCE(SUM(gross_amount), 0)::bigint as total_gross
FROM cessions WHERE status = 'BOOKED';
