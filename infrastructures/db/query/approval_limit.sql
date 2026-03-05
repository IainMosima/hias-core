-- name: GetApprovalLimitByRole :one
SELECT * FROM approval_limits WHERE role_name = $1 AND is_active = true;

-- name: ListApprovalLimits :many
SELECT * FROM approval_limits WHERE is_active = true ORDER BY role_name;

-- name: CreateApprovalLimit :one
INSERT INTO approval_limits (role_name, max_discount_percentage, max_discount_amount, max_loading_percentage, max_loading_amount, escalation_role)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: UpdateApprovalLimit :one
UPDATE approval_limits SET
    max_discount_percentage = COALESCE(sqlc.narg('max_discount_percentage'), max_discount_percentage),
    max_discount_amount = COALESCE(sqlc.narg('max_discount_amount'), max_discount_amount),
    max_loading_percentage = COALESCE(sqlc.narg('max_loading_percentage'), max_loading_percentage),
    max_loading_amount = COALESCE(sqlc.narg('max_loading_amount'), max_loading_amount),
    escalation_role = COALESCE(NULLIF(sqlc.arg('escalation_role')::text, ''), escalation_role),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;

-- name: DeleteApprovalLimit :exec
UPDATE approval_limits SET is_active = false, updated_at = NOW() WHERE id = $1;
