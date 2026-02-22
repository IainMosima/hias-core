-- name: CreateRateCard :one
INSERT INTO rate_cards (provider_id, procedure_code, procedure_name, rate_amount, effective_date)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetRateCardByID :one
SELECT * FROM rate_cards WHERE id = $1;

-- name: GetRateByProviderAndProcedure :one
SELECT * FROM rate_cards WHERE provider_id = $1 AND procedure_code = $2 ORDER BY effective_date DESC LIMIT 1;

-- name: ListRateCardsByProvider :many
SELECT * FROM rate_cards WHERE provider_id = $1 ORDER BY procedure_code;

-- name: UpdateRateCard :one
UPDATE rate_cards SET rate_amount = $2, effective_date = $3, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: DeleteRateCard :exec
DELETE FROM rate_cards WHERE id = $1;
