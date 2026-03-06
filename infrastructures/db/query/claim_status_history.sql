-- name: CreateClaimStatusHistory :one
INSERT INTO claim_status_history (claim_id, from_status, to_status, action, notes, performed_by)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: ListClaimTimeline :many
SELECT
    csh.id,
    csh.claim_id,
    csh.from_status,
    csh.to_status,
    csh.action,
    csh.notes,
    csh.performed_by,
    COALESCE(u.name, '') AS performed_by_name,
    csh.created_at
FROM claim_status_history csh
LEFT JOIN users u ON u.id = csh.performed_by
WHERE csh.claim_id = $1
ORDER BY csh.created_at ASC;
