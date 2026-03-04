-- name: CreateInstallment :one
INSERT INTO installments (schedule_id, installment_number, due_date, amount, status)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetInstallmentByID :one
SELECT * FROM installments WHERE id = $1;

-- name: ListInstallmentsBySchedule :many
SELECT * FROM installments WHERE schedule_id = $1 ORDER BY installment_number;

-- name: ListOverdueInstallments :many
SELECT * FROM installments WHERE status = 'PENDING' AND due_date < NOW() ORDER BY due_date;

-- name: UpdateInstallmentStatus :one
UPDATE installments SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: MarkInstallmentPaid :one
UPDATE installments SET status = 'PAID', paid_at = NOW(), invoice_id = $2, updated_at = NOW() WHERE id = $1 RETURNING *;
