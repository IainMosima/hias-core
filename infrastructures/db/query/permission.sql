-- name: CreatePermission :one
INSERT INTO permissions (resource, action, description) VALUES ($1, $2, $3) RETURNING *;

-- name: GetPermissionByID :one
SELECT * FROM permissions WHERE id = $1;

-- name: ListPermissions :many
SELECT * FROM permissions ORDER BY resource, action;

-- name: ListPermissionsByRole :many
SELECT p.* FROM permissions p
JOIN role_permissions rp ON rp.permission_id = p.id
WHERE rp.role_id = $1;

-- name: AssignPermissionToRole :one
INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2) RETURNING *;

-- name: RemovePermissionFromRole :exec
DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2;

-- name: ListRolePermissions :many
SELECT * FROM role_permissions WHERE role_id = $1;
