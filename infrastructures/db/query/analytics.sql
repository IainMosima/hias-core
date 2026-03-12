-- name: GetClaimsVolume :one
SELECT
    COUNT(*) as total_claims,
    COUNT(CASE WHEN status = 'APPROVED' THEN 1 END) as approved_claims,
    COUNT(CASE WHEN status = 'REJECTED' THEN 1 END) as rejected_claims,
    COUNT(CASE WHEN status = 'MANUAL_REVIEW' THEN 1 END) as manual_review_claims,
    COUNT(CASE WHEN status IN ('PAID', 'PART_PAID') THEN 1 END) as paid_claims
FROM claims
WHERE created_at >= sqlc.arg('start_date') AND created_at <= sqlc.arg('end_date');

-- name: GetApprovalRate :one
SELECT (
    CASE WHEN COUNT(*) > 0
        THEN (COUNT(CASE WHEN status IN ('APPROVED', 'PAID') THEN 1 END) * 100 / COUNT(*))
        ELSE 0
    END
)::bigint as approval_rate
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
SELECT COALESCE(
    (SELECT CASE WHEN p.premium_total > 0
        THEN (c.claims_total * 100 / p.premium_total)::bigint
        ELSE 0 END
    FROM (
        SELECT COALESCE(SUM(COALESCE(NULLIF(cl.vetted_amount, 0), cl.approved_amount)), 0) as claims_total
        FROM claims cl WHERE cl.status IN ('APPROVED', 'PAID', 'PART_PAID', 'VETTED')
        AND cl.created_at >= sqlc.arg('start_date') AND cl.created_at <= sqlc.arg('end_date')
    ) c,
    (
        SELECT COALESCE(SUM(pay.amount), 0) as premium_total
        FROM payments pay WHERE pay.type = 'PREMIUM' AND pay.status = 'CONFIRMED'
        AND pay.created_at >= sqlc.arg('start_date') AND pay.created_at <= sqlc.arg('end_date')
    ) p),
0)::bigint as loss_ratio;

-- name: GetFraudRate :one
SELECT (
    CASE WHEN COUNT(DISTINCT c.id) > 0
        THEN (COUNT(DISTINCT ff.claim_id) * 100 / COUNT(DISTINCT c.id))
        ELSE 0
    END
)::bigint as fraud_rate
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
SELECT COALESCE(SUM(COALESCE(NULLIF(vetted_amount, 0), approved_amount)), 0)::bigint as total_paid
FROM claims
WHERE status IN ('PAID', 'PART_PAID')
  AND created_at >= sqlc.arg('start_date') AND created_at <= sqlc.arg('end_date');

-- name: GetDocumentStats :one
SELECT
    COUNT(*) as total_documents,
    COUNT(CASE WHEN status = 'ACTIVE' THEN 1 END) as active_documents,
    COUNT(CASE WHEN status = 'PENDING_UPLOAD' THEN 1 END) as pending_documents,
    COUNT(CASE WHEN status = 'UPLOAD_FAILED' THEN 1 END) as failed_documents
FROM documents
WHERE deleted_at IS NULL
  AND created_at >= sqlc.arg('start_date') AND created_at <= sqlc.arg('end_date');

-- name: GetTotalDocumentCount :one
SELECT COUNT(*)::bigint as total_documents
FROM documents
WHERE status = 'ACTIVE' AND deleted_at IS NULL
  AND created_at >= sqlc.arg('start_date') AND created_at <= sqlc.arg('end_date');