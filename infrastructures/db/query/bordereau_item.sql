-- name: CreateBordereauItem :one
INSERT INTO bordereau_items (bordereau_id, cession_id, recovery_id, policy_number, claim_number, gross_amount, ceded_amount, commission_amount)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: ListBordereauItemsByBordereau :many
SELECT * FROM bordereau_items WHERE bordereau_id = $1 ORDER BY created_at ASC;

-- name: DeleteBordereauItemsByBordereau :exec
DELETE FROM bordereau_items WHERE bordereau_id = $1;
