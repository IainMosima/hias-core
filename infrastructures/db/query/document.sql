-- name: CreateDocument :one
INSERT INTO documents (entity_type, entity_id, document_type, status, file_name, file_size, mime_type, s3_key, uploaded_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;

-- name: GetDocumentByID :one
SELECT * FROM documents WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateDocumentStatus :exec
UPDATE documents SET status = $2, confirmed_at = $3, updated_at = NOW() WHERE id = $1;

-- name: SoftDeleteDocument :exec
UPDATE documents SET status = 'DELETED', deleted_at = $2, updated_at = NOW() WHERE id = $1;

-- name: ListDocumentsByEntity :many
SELECT * FROM documents WHERE entity_type = $1 AND entity_id = $2 AND deleted_at IS NULL ORDER BY created_at DESC;
