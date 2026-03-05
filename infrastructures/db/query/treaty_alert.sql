-- name: CreateTreatyAlert :one
INSERT INTO treaty_alerts (treaty_id, treaty_layer_id, alert_type, severity, message, threshold_value, current_value)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetTreatyAlertByID :one
SELECT * FROM treaty_alerts WHERE id = $1;

-- name: ListTreatyAlerts :many
SELECT * FROM treaty_alerts ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: ListTreatyAlertsByTreaty :many
SELECT * FROM treaty_alerts WHERE treaty_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListUnacknowledgedTreatyAlerts :many
SELECT * FROM treaty_alerts WHERE is_acknowledged = false ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: AcknowledgeTreatyAlert :one
UPDATE treaty_alerts SET is_acknowledged = true, acknowledged_by = $2, acknowledged_at = NOW() WHERE id = $1 RETURNING *;

-- name: CountUnacknowledgedTreatyAlerts :one
SELECT COUNT(*) FROM treaty_alerts WHERE is_acknowledged = false;
