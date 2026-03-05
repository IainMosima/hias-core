-- name: CreateCreditNote :one
INSERT INTO credit_notes (policy_id, member_id, credit_note_number, amount, currency, reason, status, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: GetCreditNoteByID :one
SELECT * FROM credit_notes WHERE id = $1;

-- name: GetCreditNoteByNumber :one
SELECT * FROM credit_notes WHERE credit_note_number = $1;

-- name: ListCreditNotesByPolicy :many
SELECT * FROM credit_notes WHERE policy_id = $1 ORDER BY created_at DESC;

-- name: ListCreditNotesByStatus :many
SELECT * FROM credit_notes WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ApproveCreditNote :one
UPDATE credit_notes SET
    status = 'APPROVED',
    approved_by = sqlc.arg('approved_by'),
    approved_at = NOW(),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;

-- name: ApplyCreditNote :one
UPDATE credit_notes SET
    status = 'APPLIED',
    applied_to_invoice_id = sqlc.arg('applied_to_invoice_id'),
    applied_at = NOW(),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;

-- name: UpdateCreditNoteStatus :one
UPDATE credit_notes SET
    status = sqlc.arg('status'),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;
