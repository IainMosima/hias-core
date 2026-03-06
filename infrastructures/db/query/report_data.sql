-- Report Data Queries: cross-domain queries for pre-built reports

-- name: GetClaimsExperienceData :many
SELECT
    p.policy_number,
    p.policyholder_name,
    p.premium_amount AS total_premium,
    COALESCE(SUM(c.total_amount), 0)::bigint AS total_claims,
    COALESCE(SUM(CASE WHEN c.status IN ('APPROVED','PAID','PART_PAID','READY_FOR_PAYMENT') THEN c.approved_amount ELSE 0 END), 0)::bigint AS approved_claims,
    COALESCE(SUM(CASE WHEN c.status = 'REJECTED' THEN c.total_amount ELSE 0 END), 0)::bigint AS rejected_claims,
    CASE WHEN p.premium_amount > 0
        THEN ROUND(COALESCE(SUM(CASE WHEN c.status IN ('APPROVED','PAID','PART_PAID','READY_FOR_PAYMENT') THEN c.approved_amount ELSE 0 END), 0)::numeric / p.premium_amount * 100, 2)
        ELSE 0
    END AS loss_ratio,
    COUNT(c.id)::bigint AS claim_count,
    COALESCE(AVG(EXTRACT(EPOCH FROM (c.updated_at - c.created_at)) / 3600), 0)::numeric(10,2) AS avg_tat_hours
FROM policies p
LEFT JOIN claims c ON c.policy_id = p.id
    AND c.created_at BETWEEN sqlc.arg('start_date')::timestamptz AND sqlc.arg('end_date')::timestamptz
WHERE p.status IN ('ACTIVE', 'LAPSED', 'TERMINATED')
AND (sqlc.narg('policy_id')::uuid IS NULL OR p.id = sqlc.narg('policy_id'))
GROUP BY p.id, p.policy_number, p.policyholder_name, p.premium_amount
ORDER BY loss_ratio DESC;

-- name: GetClaimsRegisterData :many
SELECT
    c.claim_number,
    p.policy_number,
    m.name AS member_name,
    pr.name AS provider_name,
    c.claim_type,
    c.service_date,
    c.total_amount,
    c.approved_amount,
    c.status,
    c.created_at
FROM claims c
JOIN policies p ON p.id = c.policy_id
JOIN members m ON m.id = c.member_id
JOIN providers pr ON pr.id = c.provider_id
WHERE c.created_at BETWEEN sqlc.arg('start_date')::timestamptz AND sqlc.arg('end_date')::timestamptz
AND (sqlc.narg('claim_status')::varchar IS NULL OR c.status = sqlc.narg('claim_status'))
ORDER BY c.created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetPremiumDebtorsAgeingData :many
SELECT
    p.policy_number,
    p.policyholder_name,
    p.premium_amount AS total_premium,
    COALESCE(SUM(CASE WHEN pay.status = 'CONFIRMED' THEN pay.amount ELSE 0 END), 0)::bigint AS total_paid,
    (p.premium_amount - COALESCE(SUM(CASE WHEN pay.status = 'CONFIRMED' THEN pay.amount ELSE 0 END), 0))::bigint AS outstanding,
    COALESCE(SUM(CASE WHEN i.due_date >= NOW() THEN i.amount - COALESCE(
        (SELECT COALESCE(SUM(py.amount),0) FROM payments py WHERE py.invoice_id = i.id AND py.status = 'CONFIRMED'), 0
    ) ELSE 0 END), 0)::bigint AS current_bucket,
    COALESCE(SUM(CASE WHEN i.due_date < NOW() AND i.due_date >= NOW() - INTERVAL '30 days' THEN i.amount - COALESCE(
        (SELECT COALESCE(SUM(py.amount),0) FROM payments py WHERE py.invoice_id = i.id AND py.status = 'CONFIRMED'), 0
    ) ELSE 0 END), 0)::bigint AS days_30,
    COALESCE(SUM(CASE WHEN i.due_date < NOW() - INTERVAL '30 days' AND i.due_date >= NOW() - INTERVAL '60 days' THEN i.amount - COALESCE(
        (SELECT COALESCE(SUM(py.amount),0) FROM payments py WHERE py.invoice_id = i.id AND py.status = 'CONFIRMED'), 0
    ) ELSE 0 END), 0)::bigint AS days_60,
    COALESCE(SUM(CASE WHEN i.due_date < NOW() - INTERVAL '60 days' THEN i.amount - COALESCE(
        (SELECT COALESCE(SUM(py.amount),0) FROM payments py WHERE py.invoice_id = i.id AND py.status = 'CONFIRMED'), 0
    ) ELSE 0 END), 0)::bigint AS days_90_plus
FROM policies p
LEFT JOIN invoices i ON i.policy_id = p.id AND i.status IN ('PENDING', 'OVERDUE')
LEFT JOIN payments pay ON pay.invoice_id = i.id
WHERE p.status = 'ACTIVE'
GROUP BY p.id, p.policy_number, p.policyholder_name, p.premium_amount
HAVING (p.premium_amount - COALESCE(SUM(CASE WHEN pay.status = 'CONFIRMED' THEN pay.amount ELSE 0 END), 0)) > 0
ORDER BY outstanding DESC;

-- name: GetPremiumRegisterData :many
SELECT
    p.policy_number,
    p.policyholder_name,
    pl.name AS plan_name,
    p.premium_amount,
    pay.amount AS payment_amount,
    pay.paid_at AS payment_date,
    pay.status AS payment_status,
    pay.method AS payment_method
FROM payments pay
JOIN invoices i ON i.id = pay.invoice_id
JOIN policies p ON p.id = i.policy_id
JOIN plans pl ON pl.id = p.plan_id
WHERE pay.type = 'PREMIUM'
AND pay.created_at BETWEEN sqlc.arg('start_date')::timestamptz AND sqlc.arg('end_date')::timestamptz
ORDER BY pay.created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetMembershipData :many
SELECT
    m.member_number,
    m.name AS full_name,
    m.date_of_birth,
    m.gender,
    m.phone,
    m.email,
    m.relationship,
    p.policy_number,
    pl.name AS plan_name,
    m.status,
    m.created_at AS enrollment_date
FROM members m
JOIN policies p ON p.id = m.policy_id
JOIN plans pl ON pl.id = p.plan_id
WHERE (sqlc.narg('policy_id')::uuid IS NULL OR m.policy_id = sqlc.narg('policy_id'))
AND (sqlc.narg('member_status')::varchar IS NULL OR m.status = sqlc.narg('member_status'))
ORDER BY m.name
LIMIT $1 OFFSET $2;

-- name: GetProviderPerformanceData :many
SELECT
    pr.name AS provider_name,
    pr.tier,
    pr.status,
    COUNT(c.id)::bigint AS total_claims,
    COALESCE(SUM(c.total_amount), 0)::bigint AS total_claimed,
    COALESCE(SUM(c.approved_amount), 0)::bigint AS total_approved,
    CASE WHEN COUNT(c.id) > 0
        THEN ROUND(COUNT(CASE WHEN c.status = 'REJECTED' THEN 1 END)::numeric / COUNT(c.id) * 100, 2)
        ELSE 0
    END AS rejection_rate,
    CASE WHEN COUNT(c.id) > 0
        THEN (COALESCE(SUM(c.total_amount), 0) / COUNT(c.id))::bigint
        ELSE 0
    END AS avg_claim_amount,
    COALESCE((SELECT COUNT(*) FROM fraud_flags ff JOIN claims fc ON fc.id = ff.claim_id WHERE fc.provider_id = pr.id), 0)::bigint AS fraud_flag_count
FROM providers pr
LEFT JOIN claims c ON c.provider_id = pr.id
    AND c.created_at BETWEEN sqlc.arg('start_date')::timestamptz AND sqlc.arg('end_date')::timestamptz
GROUP BY pr.id, pr.name, pr.tier, pr.status
ORDER BY total_claims DESC;

-- name: GetLossRatioData :many
SELECT
    pl.name AS plan_name,
    COUNT(DISTINCT CASE WHEN p.status = 'ACTIVE' THEN p.id END)::bigint AS active_policies,
    COUNT(DISTINCT m.id)::bigint AS total_members,
    COALESCE(SUM(DISTINCT p.premium_amount), 0)::bigint AS earned_premium,
    COALESCE(SUM(CASE WHEN c.status IN ('APPROVED','PAID','PART_PAID','READY_FOR_PAYMENT') THEN c.approved_amount ELSE 0 END), 0)::bigint AS incurred_claims,
    CASE WHEN SUM(DISTINCT p.premium_amount) > 0
        THEN ROUND(COALESCE(SUM(CASE WHEN c.status IN ('APPROVED','PAID','PART_PAID','READY_FOR_PAYMENT') THEN c.approved_amount ELSE 0 END), 0)::numeric / NULLIF(SUM(DISTINCT p.premium_amount), 0) * 100, 2)
        ELSE 0
    END AS loss_ratio,
    0::numeric(10,2) AS expense_ratio,
    CASE WHEN SUM(DISTINCT p.premium_amount) > 0
        THEN ROUND(COALESCE(SUM(CASE WHEN c.status IN ('APPROVED','PAID','PART_PAID','READY_FOR_PAYMENT') THEN c.approved_amount ELSE 0 END), 0)::numeric / NULLIF(SUM(DISTINCT p.premium_amount), 0) * 100, 2)
        ELSE 0
    END AS combined_ratio
FROM plans pl
LEFT JOIN policies p ON p.plan_id = pl.id
LEFT JOIN members m ON m.policy_id = p.id AND m.status = 'ACTIVE'
LEFT JOIN claims c ON c.policy_id = p.id
    AND c.created_at BETWEEN sqlc.arg('start_date')::timestamptz AND sqlc.arg('end_date')::timestamptz
WHERE pl.status = 'ACTIVE'
GROUP BY pl.id, pl.name
ORDER BY loss_ratio DESC;

-- name: GetRenewalData :many
SELECT
    p.policy_number,
    p.policyholder_name,
    pl.name AS plan_name,
    p.end_date AS expiry_date,
    COALESCE(pr.status, 'NOT_INITIATED') AS renewal_status,
    p.premium_amount AS current_premium,
    COALESCE(pr.new_premium, p.premium_amount) AS proposed_premium,
    CASE WHEN p.premium_amount > 0
        THEN ROUND((COALESCE(pr.new_premium, p.premium_amount) - p.premium_amount)::numeric / p.premium_amount * 100, 2)
        ELSE 0
    END AS premium_change_pct,
    (SELECT COUNT(*) FROM members mm WHERE mm.policy_id = p.id AND mm.status = 'ACTIVE')::bigint AS member_count
FROM policies p
JOIN plans pl ON pl.id = p.plan_id
LEFT JOIN policy_renewals pr ON pr.policy_id = p.id
    AND pr.created_at = (SELECT MAX(pr2.created_at) FROM policy_renewals pr2 WHERE pr2.policy_id = p.id)
WHERE p.end_date BETWEEN sqlc.arg('start_date')::timestamptz AND sqlc.arg('end_date')::timestamptz
ORDER BY p.end_date;

-- name: DrillDownClaimsByPolicy :many
SELECT
    c.claim_number,
    m.name AS member_name,
    pr.name AS provider_name,
    c.claim_type,
    c.service_date,
    c.total_amount,
    c.approved_amount,
    c.co_pay_amount,
    c.status,
    c.created_at
FROM claims c
JOIN members m ON m.id = c.member_id
JOIN providers pr ON pr.id = c.provider_id
WHERE c.policy_id = $1
AND c.created_at BETWEEN sqlc.arg('start_date')::timestamptz AND sqlc.arg('end_date')::timestamptz
ORDER BY c.created_at DESC;

-- name: DrillDownPaymentsByPolicy :many
SELECT
    pay.reference_number,
    pay.amount,
    pay.method,
    pay.status,
    pay.paid_at,
    i.invoice_number,
    i.amount AS invoice_amount,
    i.due_date
FROM payments pay
JOIN invoices i ON i.id = pay.invoice_id
WHERE i.policy_id = $1
AND pay.created_at BETWEEN sqlc.arg('start_date')::timestamptz AND sqlc.arg('end_date')::timestamptz
ORDER BY pay.created_at DESC;

-- Dashboard KPI queries

-- name: GetOutstandingPremium :one
SELECT COALESCE(SUM(
    i.amount - COALESCE((SELECT COALESCE(SUM(py.amount),0) FROM payments py WHERE py.invoice_id = i.id AND py.status = 'CONFIRMED'), 0)
), 0)::bigint AS outstanding
FROM invoices i
WHERE i.status IN ('PENDING', 'OVERDUE');

-- name: GetSLABreachCount :one
SELECT COUNT(*)::bigint FROM claims
WHERE sla_breach_at IS NOT NULL AND sla_breach_at < NOW()
AND status NOT IN ('PAID', 'REJECTED');
