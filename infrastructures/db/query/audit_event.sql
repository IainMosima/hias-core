-- name: CreateAuditEvent :one
INSERT INTO audit_events (user_id, entity_type, entity_id, action, old_value, new_value, ip_address, user_agent)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

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
