-- name: CreateMember :one
INSERT INTO members (policy_id, national_id, name, date_of_birth, gender, relationship, member_number, phone, email, kra_pin, county, address)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING *;

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
    address = COALESCE(sqlc.narg('address'), address),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;

-- name: DeleteMember :exec
DELETE FROM members WHERE id = $1;
