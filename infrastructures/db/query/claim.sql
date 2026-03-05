-- name: CreateClaim :one
INSERT INTO claims (claim_number, policy_id, member_id, provider_id, preauth_id, status, total_amount, diagnosis_codes, service_date, admission_date, discharge_date, notes, created_by, claim_type, sla_breach_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) RETURNING *;

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

-- name: GetApprovedAmountForBenefitThisYear :one
SELECT COALESCE(SUM(c.approved_amount), 0)::bigint as total_approved
FROM claims c
JOIN claim_line_items cli ON cli.claim_id = c.id
JOIN benefits b ON b.plan_id = (SELECT plan_id FROM policies WHERE id = c.policy_id)
WHERE c.member_id = $1
  AND c.status IN ('APPROVED', 'PAID')
  AND b.category = $2
  AND EXTRACT(YEAR FROM c.service_date) = EXTRACT(YEAR FROM NOW());
