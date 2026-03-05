-- name: CreateTreatyLayer :one
INSERT INTO treaty_layers (treaty_id, layer_number, attachment_point, layer_limit, deductible_amount, premium_rate, aggregate_limit)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetTreatyLayerByID :one
SELECT * FROM treaty_layers WHERE id = $1;

-- name: ListTreatyLayersByTreaty :many
SELECT * FROM treaty_layers WHERE treaty_id = $1 ORDER BY layer_number ASC;

-- name: UpdateTreatyLayer :one
UPDATE treaty_layers SET
    attachment_point = $2, layer_limit = $3, deductible_amount = $4,
    premium_rate = $5, aggregate_limit = $6, updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: UpdateTreatyLayerAggregateUsed :one
UPDATE treaty_layers SET aggregate_used = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: DeleteTreatyLayer :exec
DELETE FROM treaty_layers WHERE id = $1;
