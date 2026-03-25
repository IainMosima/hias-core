-- name: CreateAuditEvent :one
INSERT INTO audit_events (user_id, entity_type, entity_id, action, old_value, new_value, ip_address, user_agent)
VALUES (
  sqlc.arg('user_id'),
  sqlc.arg('entity_type'),
  sqlc.arg('entity_id'),
  sqlc.arg('action'),
  sqlc.arg('old_value')::jsonb,
  sqlc.arg('new_value')::jsonb,
  sqlc.arg('ip_address'),
  sqlc.arg('user_agent')
) RETURNING *;

-- name: GetAuditEventByID :one
SELECT * FROM audit_events WHERE id = $1;

-- name: ListAuditEvents :many
SELECT * FROM audit_events ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: ListAuditEventsByEntity :many
SELECT * FROM audit_events WHERE entity_type = $1 AND entity_id = $2 ORDER BY created_at DESC LIMIT $3 OFFSET $4;

-- name: ListAuditEventsByUser :many
SELECT * FROM audit_events WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: CountAuditEvents :one
SELECT COUNT(*) FROM audit_events;

-- name: CountAuditEventsByEntity :one
SELECT COUNT(*) FROM audit_events WHERE entity_type = $1 AND entity_id = $2;

-- name: CountAuditEventsByUser :one
SELECT COUNT(*) FROM audit_events WHERE user_id = $1;
