-- name: CreatePolicyRenewal :one
INSERT INTO policy_renewals (policy_id, status, renewal_date, new_premium, premium_change_reason, new_plan_id, expires_at, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: GetPolicyRenewalByID :one
SELECT * FROM policy_renewals WHERE id = $1;

-- name: GetPolicyRenewalByPolicyID :one
SELECT * FROM policy_renewals WHERE policy_id = $1 ORDER BY created_at DESC LIMIT 1;

-- name: ListPendingRenewals :many
SELECT * FROM policy_renewals WHERE status = 'PENDING' ORDER BY renewal_date ASC;

-- name: ListExpiredRenewals :many
SELECT * FROM policy_renewals WHERE status = 'PENDING' AND expires_at IS NOT NULL AND expires_at < NOW();

-- name: UpdatePolicyRenewalStatus :one
UPDATE policy_renewals SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: UpdatePolicyRenewal :one
UPDATE policy_renewals SET
    status = COALESCE(sqlc.narg('status'), status),
    renewed_policy_id = COALESCE(sqlc.narg('renewed_policy_id'), renewed_policy_id),
    new_premium = COALESCE(sqlc.narg('new_premium'), new_premium),
    premium_change_reason = COALESCE(sqlc.narg('premium_change_reason'), premium_change_reason),
    approved_by = COALESCE(sqlc.narg('approved_by'), approved_by),
    approved_at = COALESCE(sqlc.narg('approved_at'), approved_at),
    completed_at = COALESCE(sqlc.narg('completed_at'), completed_at),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;
