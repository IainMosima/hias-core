-- name: CreatePolicyDocument :one
INSERT INTO policy_documents (policy_id, member_id, document_type, file_name, file_size, s3_key, generated_by)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetPolicyDocumentByID :one
SELECT pd.*, COALESCE(u.name, '') AS generated_by_name
FROM policy_documents pd
LEFT JOIN users u ON u.id = pd.generated_by
WHERE pd.id = $1;

-- name: ListPolicyDocumentsByPolicy :many
SELECT pd.*, COALESCE(u.name, '') AS generated_by_name
FROM policy_documents pd
LEFT JOIN users u ON u.id = pd.generated_by
WHERE pd.policy_id = $1 ORDER BY pd.created_at DESC;

-- name: ListPolicyDocumentsByMember :many
SELECT pd.*, COALESCE(u.name, '') AS generated_by_name
FROM policy_documents pd
LEFT JOIN users u ON u.id = pd.generated_by
WHERE pd.member_id = $1 ORDER BY pd.created_at DESC;

-- name: DeletePolicyDocument :exec
DELETE FROM policy_documents WHERE id = $1;

-- name: CreatePolicyDocumentV2 :one
INSERT INTO policy_documents (
    policy_id, member_id, document_type, file_name, file_size, s3_key, generated_by,
    version, status, generation_mode, entity_type, entity_id, mime_type
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;

-- name: GetNextDocumentVersion :one
SELECT COALESCE(MAX(version), 0) + 1 AS next_version
FROM policy_documents
WHERE entity_type = $1 AND entity_id = $2 AND document_type = $3;

-- name: GetLatestDocumentByEntity :one
SELECT pd.*, COALESCE(u.name, '') AS generated_by_name
FROM policy_documents pd
LEFT JOIN users u ON u.id = pd.generated_by
WHERE pd.entity_type = $1 AND pd.entity_id = $2 AND pd.document_type = $3 AND pd.status = 'GENERATED'
ORDER BY pd.version DESC
LIMIT 1;

-- name: UpdatePolicyDocumentStatus :one
UPDATE policy_documents
SET status = $2, error_message = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdatePolicyDocumentGenerated :one
UPDATE policy_documents
SET file_name = $2, file_size = $3, s3_key = $4, status = $5, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: SupersedePolicyDocument :exec
UPDATE policy_documents
SET superseded_by = $2, updated_at = NOW()
WHERE id = $1;

-- name: ConfirmPolicyDocumentUpload :one
UPDATE policy_documents
SET status = 'GENERATED', file_size = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ListPolicyDocumentsByEntity :many
SELECT pd.*, COALESCE(u.name, '') AS generated_by_name
FROM policy_documents pd
LEFT JOIN users u ON u.id = pd.generated_by
WHERE pd.entity_type = $1 AND pd.entity_id = $2
ORDER BY pd.document_type, pd.version DESC;
