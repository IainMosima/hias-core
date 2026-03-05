-- name: CreateClaimDocument :one
INSERT INTO claim_documents (claim_id, file_name, file_type, file_size, s3_key, uploaded_by)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetClaimDocumentByID :one
SELECT * FROM claim_documents WHERE id = $1 AND is_deleted = false;

-- name: ListClaimDocumentsByClaim :many
SELECT * FROM claim_documents WHERE claim_id = $1 AND is_deleted = false ORDER BY created_at DESC;

-- name: SoftDeleteClaimDocument :one
UPDATE claim_documents SET is_deleted = true, updated_at = NOW() WHERE id = $1 RETURNING *;
