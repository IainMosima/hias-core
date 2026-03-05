-- name: CreateReinsurerStatement :one
INSERT INTO reinsurer_statements (statement_number, treaty_id, participant_id, period_start, period_end, premium_ceded, claims_recovered, commission_due, profit_commission, net_balance, status, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING *;

-- name: GetReinsurerStatementByID :one
SELECT * FROM reinsurer_statements WHERE id = $1;

-- name: GetReinsurerStatementByNumber :one
SELECT * FROM reinsurer_statements WHERE statement_number = $1;

-- name: ListReinsurerStatementsByTreaty :many
SELECT * FROM reinsurer_statements WHERE treaty_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: UpdateReinsurerStatementStatus :one
UPDATE reinsurer_statements SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: UpdateReinsurerStatement :one
UPDATE reinsurer_statements SET
    premium_ceded = $2, claims_recovered = $3, commission_due = $4,
    profit_commission = $5, net_balance = $6, updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: CountReinsurerStatements :one
SELECT COUNT(*) FROM reinsurer_statements;
