-- name: CreatePolicy :one
INSERT INTO policies (plan_id, policyholder_name, policyholder_email, policyholder_phone, policy_number, status, start_date, end_date, premium_amount, currency, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING *;

-- name: GetPolicyByID :one
SELECT * FROM policies WHERE id = $1;

-- name: GetPolicyByNumber :one
SELECT * FROM policies WHERE policy_number = $1;

-- name: ListPolicies :many
SELECT * FROM policies ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: ListPoliciesByStatus :many
SELECT * FROM policies WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: CountPolicies :one
SELECT COUNT(*) FROM policies;

-- name: CountPoliciesByStatus :one
SELECT COUNT(*) FROM policies WHERE status = $1;

-- name: GetActivePoliciesForBilling :many
SELECT * FROM policies WHERE status = 'ACTIVE' AND end_date > NOW();

-- name: UpdatePolicyStatus :one
UPDATE policies SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: UpdatePolicy :one
UPDATE policies SET
    policyholder_name = COALESCE(sqlc.narg('policyholder_name'), policyholder_name),
    policyholder_email = COALESCE(sqlc.narg('policyholder_email'), policyholder_email),
    policyholder_phone = COALESCE(sqlc.narg('policyholder_phone'), policyholder_phone),
    start_date = COALESCE(sqlc.narg('start_date'), start_date),
    end_date = COALESCE(sqlc.narg('end_date'), end_date),
    premium_amount = COALESCE(sqlc.narg('premium_amount'), premium_amount),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;

-- name: GetLapsedPoliciesForTermination :many
SELECT * FROM policies WHERE status = 'LAPSED' AND updated_at < NOW() - INTERVAL '90 days';

-- name: GetOverduePoliciesForLapse :many
SELECT p.* FROM policies p
JOIN invoices i ON i.policy_id = p.id
WHERE p.status = 'ACTIVE' AND i.status = 'OVERDUE' AND i.due_date < NOW() - INTERVAL '30 days';
