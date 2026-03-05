-- name: GetClaimsVolume :one
SELECT
    COUNT(*) as total_claims,
    COUNT(CASE WHEN status = 'APPROVED' THEN 1 END) as approved_claims,
    COUNT(CASE WHEN status = 'REJECTED' THEN 1 END) as rejected_claims,
    COUNT(CASE WHEN status = 'MANUAL_REVIEW' THEN 1 END) as manual_review_claims,
    COUNT(CASE WHEN status = 'PAID' THEN 1 END) as paid_claims
FROM claims
WHERE created_at >= sqlc.arg('start_date') AND created_at <= sqlc.arg('end_date');

-- name: GetApprovalRate :one
SELECT
    CASE WHEN COUNT(*) > 0
        THEN (COUNT(CASE WHEN status IN ('APPROVED', 'PAID') THEN 1 END) * 100.0 / COUNT(*))
        ELSE 0
    END as approval_rate
FROM claims
WHERE status NOT IN ('RECEIVED', 'VALIDATED')
  AND created_at >= sqlc.arg('start_date') AND created_at <= sqlc.arg('end_date');

-- name: GetAverageTAT :one
SELECT
    COALESCE(AVG(EXTRACT(EPOCH FROM (updated_at - created_at)) / 3600), 0) as avg_tat_hours
FROM claims
WHERE status IN ('APPROVED', 'REJECTED', 'PAID')
  AND created_at >= sqlc.arg('start_date') AND created_at <= sqlc.arg('end_date');

-- name: GetLossRatio :one
SELECT
    CASE WHEN COALESCE(premium.total, 0) > 0
        THEN (COALESCE(SUM(c.approved_amount), 0) * 100.0 / premium.total)
        ELSE 0
    END as loss_ratio
FROM claims c
CROSS JOIN (
    SELECT COALESCE(SUM(pay.amount), 0) as total
    FROM payments pay
    WHERE pay.type = 'PREMIUM' AND pay.status = 'CONFIRMED'
      AND pay.created_at >= sqlc.arg('start_date') AND pay.created_at <= sqlc.arg('end_date')
) premium
WHERE c.status IN ('APPROVED', 'PAID')
  AND c.created_at >= sqlc.arg('start_date') AND c.created_at <= sqlc.arg('end_date')
GROUP BY premium.total;

-- name: GetFraudRate :one
SELECT
    CASE WHEN COUNT(DISTINCT c.id) > 0
        THEN (COUNT(DISTINCT ff.claim_id) * 100.0 / COUNT(DISTINCT c.id))
        ELSE 0
    END as fraud_rate
FROM claims c
LEFT JOIN fraud_flags ff ON ff.claim_id = c.id
WHERE c.created_at >= sqlc.arg('start_date') AND c.created_at <= sqlc.arg('end_date');

-- name: GetTopProviders :many
SELECT
    p.id,
    p.name,
    COUNT(c.id) as claim_count,
    COALESCE(SUM(c.total_amount), 0)::bigint as total_amount,
    COALESCE(SUM(c.approved_amount), 0)::bigint as total_approved
FROM providers p
LEFT JOIN claims c ON c.provider_id = p.id AND c.created_at >= sqlc.arg('start_date') AND c.created_at <= sqlc.arg('end_date')
GROUP BY p.id, p.name
ORDER BY claim_count DESC
LIMIT $1;

-- name: GetTotalPremiumCollected :one
SELECT COALESCE(SUM(amount), 0)::bigint as total_premium
FROM payments
WHERE type = 'PREMIUM' AND status = 'CONFIRMED'
  AND created_at >= sqlc.arg('start_date') AND created_at <= sqlc.arg('end_date');

-- name: GetTotalClaimsPaid :one
SELECT COALESCE(SUM(approved_amount), 0)::bigint as total_paid
FROM claims
WHERE status = 'PAID'
  AND created_at >= sqlc.arg('start_date') AND created_at <= sqlc.arg('end_date');