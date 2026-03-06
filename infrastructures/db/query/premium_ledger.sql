-- name: CreatePremiumLedgerEntry :one
INSERT INTO premium_ledger_entries (policy_id, entry_type, amount, description, reference_number, effective_date, balance_after, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: ListPremiumLedgerByPolicy :many
SELECT * FROM premium_ledger_entries WHERE policy_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: GetPremiumBalanceByPolicy :one
SELECT COALESCE(
    SUM(CASE WHEN entry_type = 'CREDIT' THEN amount ELSE -amount END), 0
)::bigint AS balance
FROM premium_ledger_entries WHERE policy_id = $1;
