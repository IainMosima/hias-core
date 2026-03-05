-- name: CreateQuotationDocument :one
INSERT INTO quotation_documents (quotation_id, version_number, file_name, file_type, file_size, s3_key, uploaded_by, can_edit_roles, can_delete_roles)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;

-- name: GetQuotationDocumentByID :one
SELECT * FROM quotation_documents WHERE id = $1 AND is_deleted = false;

-- name: ListQuotationDocumentsByQuotation :many
SELECT * FROM quotation_documents WHERE quotation_id = $1 AND is_deleted = false ORDER BY created_at DESC;

-- name: SoftDeleteQuotationDocument :exec
UPDATE quotation_documents SET is_deleted = true, updated_at = NOW() WHERE id = $1;

-- name: UpdateQuotationDocument :one
UPDATE quotation_documents SET
    file_name = COALESCE(NULLIF(sqlc.arg('file_name')::text, ''), file_name),
    can_edit_roles = COALESCE(sqlc.narg('can_edit_roles'), can_edit_roles),
    can_delete_roles = COALESCE(sqlc.narg('can_delete_roles'), can_delete_roles),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;
