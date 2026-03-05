-- name: CreateRecoveryWorkflowEvent :one
INSERT INTO recovery_workflow_events (recovery_id, from_status, to_status, event_type, notes, performed_by)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: ListRecoveryWorkflowEventsByRecovery :many
SELECT * FROM recovery_workflow_events WHERE recovery_id = $1 ORDER BY created_at ASC;
