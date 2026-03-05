-- name: CreateProvider :one
INSERT INTO providers (name, type, license_number, status, county, address, phone, email, contact_person, user_id, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING *;

-- name: GetProviderByID :one
SELECT * FROM providers WHERE id = $1;

-- name: GetProviderByLicense :one
SELECT * FROM providers WHERE license_number = $1;

-- name: GetProviderByUserID :one
SELECT * FROM providers WHERE user_id = $1;

-- name: ListProviders :many
SELECT * FROM providers ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: ListProvidersByStatus :many
SELECT * FROM providers WHERE status = $1 ORDER BY name LIMIT $2 OFFSET $3;

-- name: CountProviders :one
SELECT COUNT(*) FROM providers;

-- name: UpdateProviderStatus :one
UPDATE providers SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: UpdateProviderTier :one
UPDATE providers SET tier = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: ListProvidersByTier :many
SELECT * FROM providers WHERE tier = $1 ORDER BY name LIMIT $2 OFFSET $3;

-- name: UpdateProvider :one
UPDATE providers SET
    name = COALESCE(sqlc.narg('name'), name),
    county = COALESCE(sqlc.narg('county'), county),
    address = COALESCE(sqlc.narg('address'), address),
    phone = COALESCE(sqlc.narg('phone'), phone),
    email = COALESCE(sqlc.narg('email'), email),
    contact_person = COALESCE(sqlc.narg('contact_person'), contact_person),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;

-- name: UpdateAccreditation :one
UPDATE providers SET
    accreditation_status = $2,
    accreditation_expiry = $3,
    accreditation_body = $4,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: ListProvidersByAccreditationStatus :many
SELECT * FROM providers WHERE accreditation_status = $1 ORDER BY name LIMIT $2 OFFSET $3;

-- name: ListExpiringAccreditations :many
SELECT * FROM providers WHERE accreditation_status = 'ACCREDITED'
AND accreditation_expiry <= NOW() + make_interval(days => $1)
ORDER BY accreditation_expiry LIMIT $2 OFFSET $3;
