-- name: CreateFraudFlag :one
INSERT INTO fraud_flags (claim_id, flag_type, severity, details, reference_claim_id)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetFraudFlagByID :one
SELECT * FROM fraud_flags WHERE id = $1;

-- name: ListFraudFlagsByClaim :many
SELECT * FROM fraud_flags WHERE claim_id = $1;

-- name: ListUnresolvedFraudFlags :many
SELECT * FROM fraud_flags WHERE resolved = FALSE ORDER BY severity DESC, created_at DESC LIMIT $1 OFFSET $2;

-- name: ResolveFraudFlag :one
UPDATE fraud_flags SET resolved = TRUE, resolved_by = $2, resolved_at = NOW(), updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: CheckDuplicateClaim :one
SELECT COUNT(*) FROM claims WHERE claim_number = $1 AND id != $2;

-- name: CheckFrequencyClaim :one
SELECT COUNT(*) FROM claims c
JOIN claim_line_items cli ON cli.claim_id = c.id
WHERE c.member_id = $1 AND c.provider_id = $2 AND cli.procedure_code = $3
  AND c.service_date > NOW() - INTERVAL '7 days' AND c.id != $4;
