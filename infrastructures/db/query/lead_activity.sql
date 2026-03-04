-- name: CreateLeadActivity :one
INSERT INTO lead_activities (lead_id, activity_type, description, scheduled_at, completed_at, created_by)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: ListLeadActivitiesByLead :many
SELECT * FROM lead_activities WHERE lead_id = $1 ORDER BY created_at DESC;
