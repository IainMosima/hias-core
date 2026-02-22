-- name: CreateAdjudicationDecision :one
INSERT INTO adjudication_decisions (claim_id, decision, payable_amount, member_responsibility, reasons, rule_results, adjudicated_by, adjudicated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: GetAdjudicationByClaimID :one
SELECT * FROM adjudication_decisions WHERE claim_id = $1 ORDER BY created_at DESC LIMIT 1;

-- name: ListAdjudicationsByDecision :many
SELECT * FROM adjudication_decisions WHERE decision = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;
