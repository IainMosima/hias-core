-- name: CreateQuotationVersion :one
INSERT INTO quotation_versions (quotation_id, version_number, base_premium, discount_type, discount_value, discount_reason, loading_type, loading_value, loading_reason, final_premium, member_count, proposed_members, billing_frequency, requires_approval, approval_status, pricing_breakdown, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) RETURNING *;

-- name: GetQuotationVersionByID :one
SELECT * FROM quotation_versions WHERE id = $1;

-- name: GetQuotationVersionByNumber :one
SELECT * FROM quotation_versions WHERE quotation_id = $1 AND version_number = $2;

-- name: ListQuotationVersionsByQuotation :many
SELECT * FROM quotation_versions WHERE quotation_id = $1 ORDER BY version_number ASC;

-- name: GetLatestQuotationVersion :one
SELECT * FROM quotation_versions WHERE quotation_id = $1 ORDER BY version_number DESC LIMIT 1;

-- name: UpdateQuotationVersionApproval :one
UPDATE quotation_versions SET approval_status = $2, approved_by = $3, approved_at = NOW(), updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: RejectQuotationVersion :one
UPDATE quotation_versions SET approval_status = 'REJECTED', rejection_reason = $2, approved_by = $3, approved_at = NOW(), updated_at = NOW() WHERE id = $1 RETURNING *;
