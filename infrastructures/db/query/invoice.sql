-- name: CreateInvoice :one
INSERT INTO invoices (policy_id, invoice_number, amount, currency, due_date, status, billing_period_start, billing_period_end, notes, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *;

-- name: GetInvoiceByID :one
SELECT * FROM invoices WHERE id = $1;

-- name: GetInvoiceByNumber :one
SELECT * FROM invoices WHERE invoice_number = $1;

-- name: ListInvoices :many
SELECT * FROM invoices ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: ListInvoicesByPolicy :many
SELECT * FROM invoices WHERE policy_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListInvoicesByStatus :many
SELECT * FROM invoices WHERE status = $1 ORDER BY due_date ASC LIMIT $2 OFFSET $3;

-- name: ListOverdueInvoices :many
SELECT * FROM invoices WHERE status = 'PENDING' AND due_date < NOW() ORDER BY due_date ASC;

-- name: CountInvoices :one
SELECT COUNT(*) FROM invoices;

-- name: UpdateInvoiceStatus :one
UPDATE invoices SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: ListInvoicesFiltered :many
SELECT * FROM invoices
WHERE (sqlc.narg('date_from')::timestamptz IS NULL OR created_at >= sqlc.narg('date_from'))
  AND (sqlc.narg('date_to')::timestamptz IS NULL OR created_at <= sqlc.narg('date_to'))
ORDER BY created_at DESC
LIMIT sqlc.arg('query_limit') OFFSET sqlc.arg('query_offset');

-- name: CountInvoicesFiltered :one
SELECT COUNT(*) FROM invoices
WHERE (sqlc.narg('date_from')::timestamptz IS NULL OR created_at >= sqlc.narg('date_from'))
  AND (sqlc.narg('date_to')::timestamptz IS NULL OR created_at <= sqlc.narg('date_to'));

-- name: GetInvoiceWithPolicy :one
SELECT i.*, p.policy_number, p.policyholder_name
FROM invoices i
JOIN policies p ON p.id = i.policy_id
WHERE i.id = $1;

-- name: ListInvoicesWithPolicy :many
SELECT i.*, p.policy_number, p.policyholder_name
FROM invoices i
JOIN policies p ON p.id = i.policy_id
ORDER BY i.created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListInvoicesFilteredWithPolicy :many
SELECT i.*, p.policy_number, p.policyholder_name
FROM invoices i
JOIN policies p ON p.id = i.policy_id
WHERE (sqlc.narg('date_from')::timestamptz IS NULL OR i.created_at >= sqlc.narg('date_from'))
  AND (sqlc.narg('date_to')::timestamptz IS NULL OR i.created_at <= sqlc.narg('date_to'))
ORDER BY i.created_at DESC
LIMIT sqlc.arg('query_limit') OFFSET sqlc.arg('query_offset');
