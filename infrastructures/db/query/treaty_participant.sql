-- name: CreateTreatyParticipant :one
INSERT INTO treaty_participants (treaty_id, reinsurer_name, share_percentage, commission_rate, is_lead)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetTreatyParticipantByID :one
SELECT * FROM treaty_participants WHERE id = $1;

-- name: ListTreatyParticipantsByTreaty :many
SELECT * FROM treaty_participants WHERE treaty_id = $1 ORDER BY is_lead DESC, reinsurer_name ASC;

-- name: UpdateTreatyParticipant :one
UPDATE treaty_participants SET
    reinsurer_name = $2, share_percentage = $3, commission_rate = $4, is_lead = $5, updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: DeleteTreatyParticipant :exec
DELETE FROM treaty_participants WHERE id = $1;

-- name: GetTotalShareByTreaty :one
SELECT COALESCE(SUM(share_percentage), 0)::numeric(5,2) as total_share
FROM treaty_participants WHERE treaty_id = $1;
