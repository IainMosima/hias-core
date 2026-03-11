-- name: ListStandaloneDocuments :many
SELECT id, source_type, source_id, document_type, file_name, file_size, s3_key, created_by, created_at FROM (
    SELECT
        id,
        'policy' AS source_type,
        policy_id AS source_id,
        document_type,
        file_name,
        COALESCE(file_size, 0) AS file_size,
        s3_key,
        generated_by AS created_by,
        created_at
    FROM policy_documents
    UNION ALL
    SELECT
        id,
        'claim' AS source_type,
        claim_id AS source_id,
        file_type AS document_type,
        file_name,
        COALESCE(file_size, 0) AS file_size,
        s3_key,
        uploaded_by AS created_by,
        created_at
    FROM claim_documents WHERE is_deleted = false
    UNION ALL
    SELECT
        id,
        'quotation' AS source_type,
        quotation_id AS source_id,
        file_type AS document_type,
        file_name,
        COALESCE(file_size, 0) AS file_size,
        s3_key,
        COALESCE(uploaded_by, '00000000-0000-0000-0000-000000000000'::uuid) AS created_by,
        created_at
    FROM quotation_documents WHERE is_deleted = false
    UNION ALL
    SELECT
        id,
        entity_type AS source_type,
        entity_id AS source_id,
        document_type,
        file_name,
        COALESCE(file_size, 0) AS file_size,
        s3_key,
        uploaded_by AS created_by,
        created_at
    FROM documents WHERE status = 'ACTIVE' AND deleted_at IS NULL
) docs
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountStandaloneDocuments :one
SELECT (
    (SELECT COUNT(*) FROM policy_documents) +
    (SELECT COUNT(*) FROM claim_documents WHERE is_deleted = false) +
    (SELECT COUNT(*) FROM quotation_documents WHERE is_deleted = false) +
    (SELECT COUNT(*) FROM documents WHERE status = 'ACTIVE' AND deleted_at IS NULL)
)::bigint AS total;

-- name: FindDocumentS3Key :one
SELECT doc_id, s3_key, source_type FROM (
    SELECT pd.id AS doc_id, pd.s3_key, 'policy' AS source_type FROM policy_documents pd WHERE pd.id = $1
    UNION ALL
    SELECT cd.id AS doc_id, cd.s3_key, 'claim' AS source_type FROM claim_documents cd WHERE cd.id = $1 AND cd.is_deleted = false
    UNION ALL
    SELECT qd.id AS doc_id, qd.s3_key, 'quotation' AS source_type FROM quotation_documents qd WHERE qd.id = $1 AND qd.is_deleted = false
    UNION ALL
    SELECT d.id AS doc_id, d.s3_key, d.entity_type AS source_type FROM documents d WHERE d.id = $1 AND d.status = 'ACTIVE' AND d.deleted_at IS NULL
) docs
LIMIT 1;
