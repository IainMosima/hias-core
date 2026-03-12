-- name: CreateAPIPartner :one
INSERT INTO api_partners (name, partner_type, api_key, api_secret_hash, provider_id, is_active, rate_limit_per_minute, allowed_claim_types, webhook_url, contact_email, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING *;

-- name: GetAPIPartnerByID :one
SELECT * FROM api_partners WHERE id = $1;

-- name: GetAPIPartnerByAPIKey :one
SELECT * FROM api_partners WHERE api_key = $1;

-- name: ListAPIPartners :many
SELECT * FROM api_partners ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: UpdateAPIPartner :one
UPDATE api_partners SET
    name = $2,
    partner_type = $3,
    provider_id = $4,
    rate_limit_per_minute = $5,
    allowed_claim_types = $6,
    webhook_url = $7,
    contact_email = $8,
    metadata = $9,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: DeactivateAPIPartner :exec
UPDATE api_partners SET is_active = false, updated_at = NOW() WHERE id = $1;

-- name: UpdateAPIPartnerAPIKey :one
UPDATE api_partners SET api_key = $2, api_secret_hash = $3, updated_at = NOW() WHERE id = $1 RETURNING *;
