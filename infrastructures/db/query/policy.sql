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

-- name: ListPoliciesExpiringSoon :many
SELECT * FROM policies WHERE status = 'ACTIVE' AND end_date BETWEEN NOW() AND NOW() + make_interval(days => sqlc.arg('days')::int) ORDER BY end_date ASC;

-- name: UpdatePolicyPlanAndPremium :one
UPDATE policies SET plan_id = $2, premium_amount = $3, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: GetActivePolicyCount :one
SELECT COUNT(*) FROM policies WHERE status = 'ACTIVE' AND created_at BETWEEN sqlc.arg('start_date') AND sqlc.arg('end_date');

-- name: GetLapsedPolicyCount :one
SELECT COUNT(*) FROM policies WHERE status = 'LAPSED' AND created_at BETWEEN sqlc.arg('start_date') AND sqlc.arg('end_date');

-- name: GetTotalActiveMemberCount :one
SELECT COUNT(*) FROM members WHERE status = 'ACTIVE' AND created_at BETWEEN sqlc.arg('start_date') AND sqlc.arg('end_date');

-- name: GetRenewalRate :one
SELECT CASE WHEN COUNT(*) = 0 THEN 0 ELSE
    ROUND(COUNT(CASE WHEN pr.status = 'COMPLETED' THEN 1 END)::numeric / COUNT(*)::numeric * 100, 2)
END FROM policy_renewals pr WHERE pr.created_at BETWEEN sqlc.arg('start_date') AND sqlc.arg('end_date');

-- name: ListPoliciesFiltered :many
SELECT * FROM policies
WHERE (sqlc.narg('date_from')::timestamptz IS NULL OR created_at >= sqlc.narg('date_from'))
  AND (sqlc.narg('date_to')::timestamptz IS NULL OR created_at <= sqlc.narg('date_to'))
  AND (sqlc.arg('search')::text = '' OR (policy_number ILIKE '%' || sqlc.arg('search') || '%' OR policyholder_name ILIKE '%' || sqlc.arg('search') || '%' OR policyholder_email ILIKE '%' || sqlc.arg('search') || '%'))
ORDER BY created_at DESC
LIMIT sqlc.arg('query_limit') OFFSET sqlc.arg('query_offset');

-- name: CountPoliciesFiltered :one
SELECT COUNT(*) FROM policies
WHERE (sqlc.narg('date_from')::timestamptz IS NULL OR created_at >= sqlc.narg('date_from'))
  AND (sqlc.narg('date_to')::timestamptz IS NULL OR created_at <= sqlc.narg('date_to'))
  AND (sqlc.arg('search')::text = '' OR (policy_number ILIKE '%' || sqlc.arg('search') || '%' OR policyholder_name ILIKE '%' || sqlc.arg('search') || '%' OR policyholder_email ILIKE '%' || sqlc.arg('search') || '%'));

-- name: ActivatePolicyWithTimestamp :one
UPDATE policies SET status = 'ACTIVE', activated_at = NOW(), updated_at = NOW()
WHERE id = $1 RETURNING *;
