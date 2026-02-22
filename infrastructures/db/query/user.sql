-- name: CreateUser :one
INSERT INTO users (cognito_sub, email, name, phone, national_id, role_id, status, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByCognitoSub :one
SELECT * FROM users WHERE cognito_sub = $1;

-- name: GetUserByNationalID :one
SELECT * FROM users WHERE national_id = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: UpdateUser :one
UPDATE users SET
    name = COALESCE(sqlc.narg('name'), name),
    phone = COALESCE(sqlc.narg('phone'), phone),
    national_id = COALESCE(sqlc.narg('national_id'), national_id),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;

-- name: UpdateUserStatus :one
UPDATE users SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: UpdateUserRole :one
UPDATE users SET role_id = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: UpdateUserCognitoSub :one
UPDATE users SET cognito_sub = $2, updated_at = NOW() WHERE id = $1 RETURNING *;
