-- name: CreateMember :one
INSERT INTO members (policy_id, national_id, name, date_of_birth, gender, relationship, member_number, phone, email, kra_pin, county, city, country, address)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING *;

-- name: GetMemberByID :one
SELECT * FROM members WHERE id = $1;

-- name: GetMemberByNumber :one
SELECT * FROM members WHERE member_number = $1;

-- name: GetMemberByNationalID :one
SELECT * FROM members WHERE national_id = $1;

-- name: ListMembersByPolicy :many
SELECT * FROM members WHERE policy_id = $1 ORDER BY relationship, name;

-- name: CountMembersByPolicy :one
SELECT COUNT(*) FROM members WHERE policy_id = $1;

-- name: VerifyMember :one
UPDATE members SET verified = TRUE, verified_at = NOW(), updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: UpdateMember :one
UPDATE members SET
    name = COALESCE(sqlc.narg('name'), name),
    phone = COALESCE(sqlc.narg('phone'), phone),
    email = COALESCE(sqlc.narg('email'), email),
    kra_pin = COALESCE(sqlc.narg('kra_pin'), kra_pin),
    county = COALESCE(sqlc.narg('county'), county),
    city = COALESCE(sqlc.narg('city'), city),
    country = COALESCE(sqlc.narg('country'), country),
    address = COALESCE(sqlc.narg('address'), address),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;

-- name: DeleteMember :exec
DELETE FROM members WHERE id = $1;

-- name: UpdateMemberStatus :one
UPDATE members SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: ListActiveMembersByPolicy :many
SELECT * FROM members WHERE policy_id = $1 AND status = 'ACTIVE' ORDER BY relationship, name;

-- name: CountActiveMembersByPolicy :one
SELECT COUNT(*) FROM members WHERE policy_id = $1 AND status = 'ACTIVE';

-- name: ListMembersFiltered :many
SELECT * FROM members
WHERE ($1::text = '' OR (name ILIKE '%' || $1 || '%' OR national_id ILIKE '%' || $1 || '%' OR email ILIKE '%' || $1 || '%' OR phone ILIKE '%' || $1 || '%'))
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountMembersFiltered :one
SELECT COUNT(*) FROM members
WHERE ($1::text = '' OR (name ILIKE '%' || $1 || '%' OR national_id ILIKE '%' || $1 || '%' OR email ILIKE '%' || $1 || '%' OR phone ILIKE '%' || $1 || '%'));

-- name: ActivatePendingMembersByPolicy :exec
UPDATE members SET status = 'ACTIVE', updated_at = NOW()
WHERE policy_id = $1 AND status = 'PENDING';

-- name: ListPendingMembersByPolicy :many
SELECT * FROM members WHERE policy_id = $1 AND status = 'PENDING';

-- name: UpdateMemberCoverageDates :one
UPDATE members SET
    coverage_start_date = COALESCE(sqlc.narg('coverage_start_date'), coverage_start_date),
    coverage_end_date = COALESCE(sqlc.narg('coverage_end_date'), coverage_end_date),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;
