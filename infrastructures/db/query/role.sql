-- name: CreateRole :one
INSERT INTO roles (name, description) VALUES ($1, $2) RETURNING *;

-- name: GetRoleByID :one
SELECT * FROM roles WHERE id = $1;

-- name: GetRoleByName :one
SELECT * FROM roles WHERE name = $1;

-- name: ListRoles :many
SELECT * FROM roles ORDER BY name;

-- name: UpdateRole :one
UPDATE roles SET name = $2, description = $3, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: DeleteRole :exec
DELETE FROM roles WHERE id = $1;
