-- name: CreateQuotation :one
INSERT INTO quotations (quotation_number, lead_id, plan_id, quotation_type, status, current_version, valid_from, valid_until, client_name, client_email, client_phone, currency, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING *;

-- name: GetQuotationByID :one
SELECT * FROM quotations WHERE id = $1;

-- name: GetQuotationByNumber :one
SELECT * FROM quotations WHERE quotation_number = $1;

-- name: ListQuotations :many
SELECT * FROM quotations ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: ListQuotationsByLead :many
SELECT * FROM quotations WHERE lead_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListQuotationsByStatus :many
SELECT * FROM quotations WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: UpdateQuotationStatus :one
UPDATE quotations SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: UpdateQuotationCurrentVersion :one
UPDATE quotations SET current_version = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: SetQuotationPolicyID :one
UPDATE quotations SET policy_id = $2, status = 'CONVERTED', updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: CountQuotations :one
SELECT COUNT(*) FROM quotations;

-- name: ListExpiredQuotations :many
SELECT * FROM quotations
WHERE status IN ('ISSUED', 'PENDING_DECISION')
  AND valid_until < NOW();
