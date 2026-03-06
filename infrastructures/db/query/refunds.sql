-- name: CreateRefund :one
INSERT INTO refunds (policy_id, credit_note_id, amount, currency, status, reason, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetRefundByID :one
SELECT * FROM refunds WHERE id = $1;

-- name: ListRefundsByPolicy :many
SELECT * FROM refunds WHERE policy_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ApproveRefund :one
UPDATE refunds SET status = 'APPROVED', approved_by = $2, approved_at = NOW(), updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: ProcessRefund :one
UPDATE refunds SET status = 'PROCESSED', processed_at = NOW(), updated_at = NOW()
WHERE id = $1 RETURNING *;
