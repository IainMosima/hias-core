-- Reverse document generation V1 changes
DROP INDEX IF EXISTS idx_policy_documents_entity_version;
DROP INDEX IF EXISTS idx_policy_documents_status;
DROP INDEX IF EXISTS idx_policy_documents_entity;

ALTER TABLE policy_documents
    DROP COLUMN IF EXISTS updated_at,
    DROP COLUMN IF EXISTS error_message,
    DROP COLUMN IF EXISTS superseded_by,
    DROP COLUMN IF EXISTS entity_id,
    DROP COLUMN IF EXISTS entity_type,
    DROP COLUMN IF EXISTS generation_mode,
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS version;
