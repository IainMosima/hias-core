-- name: CreateTreaty :one
INSERT INTO treaties (treaty_number, name, treaty_type, status, effective_date, expiry_date, retention_limit, currency, notes, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *;

-- name: GetTreatyByID :one
SELECT * FROM treaties WHERE id = $1;

-- name: GetTreatyByNumber :one
SELECT * FROM treaties WHERE treaty_number = $1;

-- name: ListTreaties :many
SELECT * FROM treaties ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: ListActiveTreaties :many
SELECT * FROM treaties WHERE status = 'ACTIVE' ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: ListTreatiesByStatus :many
SELECT * FROM treaties WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListTreatiesByType :many
SELECT * FROM treaties WHERE treaty_type = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: UpdateTreaty :one
UPDATE treaties SET
    name = $2, effective_date = $3, expiry_date = $4,
    retention_limit = $5, currency = $6, notes = $7, updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: UpdateTreatyStatus :one
UPDATE treaties SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: CountTreaties :one
SELECT COUNT(*) FROM treaties;

-- name: ListExpiringTreaties :many
SELECT * FROM treaties
WHERE status = 'ACTIVE'
  AND expiry_date <= NOW() + ($1::int || ' days')::interval
ORDER BY expiry_date ASC
LIMIT $2 OFFSET $3;

-- name: ListExpiredActiveTreaties :many
SELECT * FROM treaties
WHERE status = 'ACTIVE' AND expiry_date < NOW();
