-- Add document generation V1 columns to policy_documents
ALTER TABLE policy_documents
    ADD COLUMN version INT NOT NULL DEFAULT 1,
    ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'GENERATED',
    ADD COLUMN generation_mode VARCHAR(10) NOT NULL DEFAULT 'MANUAL',
    ADD COLUMN entity_type VARCHAR(30) NOT NULL DEFAULT 'policy',
    ADD COLUMN entity_id UUID,
    ADD COLUMN superseded_by UUID REFERENCES policy_documents(id),
    ADD COLUMN error_message TEXT,
    ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- Backfill entity_id from policy_id for existing rows
UPDATE policy_documents SET entity_id = policy_id WHERE entity_id IS NULL;

-- Now make entity_id NOT NULL
ALTER TABLE policy_documents ALTER COLUMN entity_id SET NOT NULL;

-- Update entity_type for member cards (they have member_id set)
UPDATE policy_documents SET entity_type = 'member', entity_id = member_id
    WHERE document_type = 'MEMBER_CARD' AND member_id IS NOT NULL;

-- Indexes
CREATE INDEX idx_policy_documents_entity ON policy_documents(entity_type, entity_id);
CREATE INDEX idx_policy_documents_status ON policy_documents(status);
CREATE INDEX idx_policy_documents_entity_version ON policy_documents(entity_type, entity_id, document_type, version DESC);
