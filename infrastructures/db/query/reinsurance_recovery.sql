-- name: CreateReinsuranceRecovery :one
INSERT INTO reinsurance_recoveries (recovery_number, claim_id, treaty_id, treaty_layer_id, cession_id, gross_claim_amount, recoverable_amount, recovered_amount, outstanding_amount, status, workflow_status, notes, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING *;

-- name: GetReinsuranceRecoveryByID :one
SELECT * FROM reinsurance_recoveries WHERE id = $1;

-- name: GetReinsuranceRecoveryByNumber :one
SELECT * FROM reinsurance_recoveries WHERE recovery_number = $1;

-- name: ListReinsuranceRecoveriesByClaim :many
SELECT * FROM reinsurance_recoveries WHERE claim_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListReinsuranceRecoveriesByTreaty :many
SELECT * FROM reinsurance_recoveries WHERE treaty_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListOutstandingRecoveries :many
SELECT * FROM reinsurance_recoveries WHERE outstanding_amount > 0 AND status NOT IN ('PAID', 'WRITTEN_OFF')
ORDER BY created_at ASC LIMIT $1 OFFSET $2;

-- name: UpdateReinsuranceRecoveryStatus :one
UPDATE reinsurance_recoveries SET status = $2, workflow_status = $3, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: UpdateReinsuranceRecoveryAmounts :one
UPDATE reinsurance_recoveries SET recovered_amount = $2, outstanding_amount = $3, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: GetTotalRecoverableByTreaty :one
SELECT COALESCE(SUM(recoverable_amount), 0)::bigint as total_recoverable
FROM reinsurance_recoveries WHERE treaty_id = $1;

-- name: GetTotalRecoveredByTreaty :one
SELECT COALESCE(SUM(recovered_amount), 0)::bigint as total_recovered
FROM reinsurance_recoveries WHERE treaty_id = $1;

-- name: GetTotalRecoveredByTreatyAndPeriod :one
SELECT COALESCE(SUM(recovered_amount), 0)::bigint as total_recovered
FROM reinsurance_recoveries WHERE treaty_id = $1 AND created_at >= $2 AND created_at <= $3;

-- name: CountReinsuranceRecoveries :one
SELECT COUNT(*) FROM reinsurance_recoveries;

-- name: CountOutstandingRecoveries :one
SELECT COUNT(*) FROM reinsurance_recoveries WHERE outstanding_amount > 0 AND status NOT IN ('PAID', 'WRITTEN_OFF');

-- name: GetAgedRecoveryAnalysis :many
SELECT
    CASE
        WHEN NOW() - created_at <= INTERVAL '30 days' THEN '0-30 days'
        WHEN NOW() - created_at <= INTERVAL '60 days' THEN '31-60 days'
        WHEN NOW() - created_at <= INTERVAL '90 days' THEN '61-90 days'
        ELSE '90+ days'
    END AS bucket,
    COUNT(*)::bigint AS count,
    COALESCE(SUM(outstanding_amount), 0)::bigint AS total_outstanding
FROM reinsurance_recoveries
WHERE outstanding_amount > 0 AND status NOT IN ('PAID', 'WRITTEN_OFF')
GROUP BY bucket
ORDER BY bucket;

-- name: GetTotalRecoverableAmountAll :one
SELECT COALESCE(SUM(recoverable_amount), 0)::bigint as total_recoverable
FROM reinsurance_recoveries;

-- name: GetTotalRecoveredAmountAll :one
SELECT COALESCE(SUM(recovered_amount), 0)::bigint as total_recovered
FROM reinsurance_recoveries;
