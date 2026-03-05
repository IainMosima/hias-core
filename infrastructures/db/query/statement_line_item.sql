-- name: CreateStatementLineItem :one
INSERT INTO statement_line_items (statement_id, claim_number, service_date, member_name, procedure_code, claimed_amount, match_status)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: ListStatementLineItemsByStatement :many
SELECT * FROM statement_line_items WHERE statement_id = $1 ORDER BY created_at;

-- name: MatchStatementLineItem :one
UPDATE statement_line_items SET
    matched_claim_id = $2,
    match_status = $3,
    discrepancy_amount = $4,
    notes = $5
WHERE id = $1 RETURNING *;
