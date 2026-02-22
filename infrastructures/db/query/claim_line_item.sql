-- name: CreateClaimLineItem :one
INSERT INTO claim_line_items (claim_id, procedure_code, procedure_name, diagnosis_code, quantity, unit_price, total_price)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetClaimLineItemByID :one
SELECT * FROM claim_line_items WHERE id = $1;

-- name: ListClaimLineItems :many
SELECT * FROM claim_line_items WHERE claim_id = $1 ORDER BY created_at;

-- name: UpdateClaimLineItemApprovedAmount :one
UPDATE claim_line_items SET approved_amount = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: DeleteClaimLineItem :exec
DELETE FROM claim_line_items WHERE id = $1;
