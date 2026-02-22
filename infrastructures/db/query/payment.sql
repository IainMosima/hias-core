-- name: CreatePayment :one
INSERT INTO payments (invoice_id, claim_id, type, amount, currency, method, reference_number, status, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;

-- name: GetPaymentByID :one
SELECT * FROM payments WHERE id = $1;

-- name: GetPaymentByReference :one
SELECT * FROM payments WHERE reference_number = $1;

-- name: ListPayments :many
SELECT * FROM payments ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: ListPaymentsByStatus :many
SELECT * FROM payments WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListPaymentsByInvoice :many
SELECT * FROM payments WHERE invoice_id = $1 ORDER BY created_at DESC;

-- name: CountPayments :one
SELECT COUNT(*) FROM payments;

-- name: UpdatePaymentStatus :one
UPDATE payments SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: ConfirmPayment :one
UPDATE payments SET status = 'CONFIRMED', paid_at = NOW(), gateway_response = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: ReconcilePayment :one
UPDATE payments SET status = 'RECONCILED', reconciled_at = NOW(), updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: IncrementPaymentRetry :one
UPDATE payments SET retry_count = retry_count + 1, status = 'PROCESSING', updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: GetFailedPaymentsForRetry :many
SELECT * FROM payments WHERE status = 'FAILED' AND retry_count < max_retries ORDER BY created_at ASC LIMIT $1;
