-- name: CreateClaim :one
INSERT INTO claims (claim_number, policy_id, member_id, provider_id, preauth_id, status, total_amount, diagnosis_codes, service_date, admission_date, discharge_date, notes, created_by, claim_type, sla_breach_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) RETURNING *;

-- name: GetMaxClaimCounterForYear :one
SELECT COALESCE(MAX(CAST(SPLIT_PART(claim_number, '-', 3) AS BIGINT)), 0)::bigint AS max_counter
FROM claims
WHERE claim_number LIKE $1;

-- name: GetClaimByID :one
SELECT * FROM claims WHERE id = $1;

-- name: GetClaimByNumber :one
SELECT * FROM claims WHERE claim_number = $1;

-- name: ListClaims :many
SELECT * FROM claims ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: ListClaimsByStatus :many
SELECT * FROM claims WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListClaimsByProvider :many
SELECT * FROM claims WHERE provider_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListClaimsByMember :many
SELECT * FROM claims WHERE member_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListClaimsByPolicy :many
SELECT * FROM claims WHERE policy_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: CountClaims :one
SELECT COUNT(*) FROM claims;

-- name: CountClaimsByStatus :one
SELECT COUNT(*) FROM claims WHERE status = $1;

-- name: GetClaimsForAdjudication :many
SELECT * FROM claims WHERE status = 'VALIDATED' ORDER BY created_at ASC LIMIT $1;

-- name: GetApprovedClaimsForRemittance :many
SELECT * FROM claims WHERE status = 'APPROVED' AND provider_id = $1;

-- name: UpdateClaimStatus :one
UPDATE claims SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: UpdateClaimAmounts :one
UPDATE claims SET
    approved_amount = $2,
    co_pay_amount = $3,
    member_responsibility = $4,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: UpdateClaimRejection :one
UPDATE claims SET status = 'REJECTED', rejection_reason = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: VetClaim :one
UPDATE claims SET
    vetted_amount = $2,
    vetted_by = $3,
    vetted_at = NOW(),
    status = $4,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: MarkClaimReadyForPayment :one
UPDATE claims SET status = 'READY_FOR_PAYMENT', updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: ListSLABreachedClaims :many
SELECT * FROM claims
WHERE sla_breach_at IS NOT NULL
  AND sla_breach_at < NOW()
  AND status NOT IN ('PAID', 'REJECTED')
ORDER BY sla_breach_at ASC
LIMIT $1 OFFSET $2;

-- name: FindClaimByProviderAndDate :one
SELECT * FROM claims
WHERE provider_id = $1
AND service_date = $2
AND ABS(total_amount - $3) < 100
AND status NOT IN ('REJECTED', 'CANCELLED')
ORDER BY created_at DESC
LIMIT 1;

-- name: ListApproachingSLAClaims :many
SELECT * FROM claims
WHERE sla_breach_at IS NOT NULL
  AND sla_breach_at > NOW()
  AND sla_breach_at <= NOW() + INTERVAL '24 hours'
  AND status NOT IN ('PAID', 'REJECTED', 'VETTED', 'READY_FOR_PAYMENT')
ORDER BY sla_breach_at ASC
LIMIT $1 OFFSET $2;

-- name: CountClaimsByMemberThisMonth :one
SELECT COUNT(*) FROM claims WHERE member_id = $1 AND created_at >= date_trunc('month', CURRENT_DATE) AND deleted_at IS NULL;

-- name: SetClaimEscalatedTo :exec
UPDATE claims SET escalated_to = $2, updated_at = now() WHERE id = $1;

-- name: GetApprovedAmountForBenefitThisYear :one
SELECT COALESCE(SUM(c.approved_amount), 0)::bigint as total_approved
FROM claims c
JOIN claim_line_items cli ON cli.claim_id = c.id
JOIN benefits b ON b.plan_id = (SELECT plan_id FROM policies WHERE id = c.policy_id)
WHERE c.member_id = $1
  AND c.status IN ('APPROVED', 'PAID')
  AND b.category = $2
  AND EXTRACT(YEAR FROM c.service_date) = EXTRACT(YEAR FROM NOW());

-- name: ListClaimsFiltered :many
SELECT * FROM claims
WHERE (sqlc.arg('status_filter')::text = '' OR status = sqlc.arg('status_filter'))
  AND (sqlc.narg('date_from')::timestamptz IS NULL OR created_at >= sqlc.narg('date_from'))
  AND (sqlc.narg('date_to')::timestamptz IS NULL OR created_at <= sqlc.narg('date_to'))
  AND (sqlc.arg('search')::text = '' OR (claim_number ILIKE '%' || sqlc.arg('search') || '%' OR status ILIKE '%' || sqlc.arg('search') || '%'))
ORDER BY created_at DESC
LIMIT sqlc.arg('query_limit') OFFSET sqlc.arg('query_offset');

-- name: CountClaimsFiltered :one
SELECT COUNT(*) FROM claims
WHERE (sqlc.arg('status_filter')::text = '' OR status = sqlc.arg('status_filter'))
  AND (sqlc.narg('date_from')::timestamptz IS NULL OR created_at >= sqlc.narg('date_from'))
  AND (sqlc.narg('date_to')::timestamptz IS NULL OR created_at <= sqlc.narg('date_to'))
  AND (sqlc.arg('search')::text = '' OR (claim_number ILIKE '%' || sqlc.arg('search') || '%' OR status ILIKE '%' || sqlc.arg('search') || '%'));

-- name: GetClaimByIdempotencyKey :one
SELECT * FROM claims WHERE idempotency_key = $1;

-- name: GetClaimByExternalClaimID :one
SELECT * FROM claims WHERE external_claim_id = $1;

-- name: CreateDraftClaim :one
INSERT INTO claims (claim_number, policy_id, member_id, provider_id, preauth_id, status, total_amount, diagnosis_codes, service_date, admission_date, discharge_date, notes, created_by, claim_type, sla_breach_at, claim_source, is_draft)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, 'INTERNAL', true) RETURNING *;

-- name: UpdateDraftClaim :one
UPDATE claims SET
    policy_id = $2,
    member_id = $3,
    provider_id = $4,
    preauth_id = $5,
    diagnosis_codes = $6,
    service_date = $7,
    notes = $8,
    claim_type = $9,
    total_amount = $10,
    updated_at = NOW()
WHERE id = $1 AND is_draft = true RETURNING *;

-- name: ListDraftClaims :many
SELECT * FROM claims WHERE is_draft = true AND created_by = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: CompleteDraftClaim :one
UPDATE claims SET is_draft = false, draft_completed_at = NOW(), updated_at = NOW() WHERE id = $1 AND is_draft = true RETURNING *;

-- name: DeleteDraftClaim :exec
DELETE FROM claims WHERE id = $1 AND is_draft = true;

-- name: UpdateClaimSource :exec
UPDATE claims SET
    claim_source = $2,
    idempotency_key = $3,
    external_claim_id = $4,
    source_metadata = $5,
    updated_at = NOW()
WHERE id = $1;
