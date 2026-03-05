-- name: CreatePolicyDocument :one
INSERT INTO policy_documents (policy_id, member_id, document_type, file_name, file_size, s3_key, generated_by)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetPolicyDocumentByID :one
SELECT * FROM policy_documents WHERE id = $1;

-- name: ListPolicyDocumentsByPolicy :many
SELECT * FROM policy_documents WHERE policy_id = $1 ORDER BY created_at DESC;

-- name: ListPolicyDocumentsByMember :many
SELECT * FROM policy_documents WHERE member_id = $1 ORDER BY created_at DESC;

-- name: DeletePolicyDocument :exec
DELETE FROM policy_documents WHERE id = $1;
