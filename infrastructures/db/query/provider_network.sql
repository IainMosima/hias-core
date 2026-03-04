-- name: CreateProviderNetwork :one
INSERT INTO provider_networks (plan_id, provider_id, benefit_category, status)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetProviderNetworkByID :one
SELECT * FROM provider_networks WHERE id = $1;

-- name: ListProviderNetworksByPlan :many
SELECT * FROM provider_networks WHERE plan_id = $1 ORDER BY created_at;

-- name: ListProviderNetworksByProvider :many
SELECT * FROM provider_networks WHERE provider_id = $1 ORDER BY created_at;

-- name: CheckProviderNetworkEligibility :one
SELECT COUNT(*) FROM provider_networks
WHERE plan_id = $1 AND provider_id = $2 AND status = 'ACTIVE'
AND (benefit_category IS NULL OR benefit_category = $3);

-- name: UpdateProviderNetworkStatus :one
UPDATE provider_networks SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: DeleteProviderNetwork :exec
DELETE FROM provider_networks WHERE id = $1;
