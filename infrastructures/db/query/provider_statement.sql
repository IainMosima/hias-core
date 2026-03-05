-- name: CreateProviderStatement :one
INSERT INTO provider_statements (provider_id, statement_number, period_start, period_end, total_claimed, status, file_name, s3_key, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;

-- name: GetProviderStatementByID :one
SELECT * FROM provider_statements WHERE id = $1;

-- name: ListProviderStatementsByProvider :many
SELECT * FROM provider_statements WHERE provider_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ReconcileProviderStatement :one
UPDATE provider_statements SET
    total_matched = $2,
    total_discrepancy = $3,
    matched_count = $4,
    unmatched_count = $5,
    status = 'RECONCILED',
    reconciled_by = $6,
    reconciled_at = NOW(),
    updated_at = NOW()
WHERE id = $1 RETURNING *;
