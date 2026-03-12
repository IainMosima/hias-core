-- Drop API Partners
DROP INDEX IF EXISTS idx_api_partners_provider;
DROP INDEX IF EXISTS idx_api_partners_api_key;
DROP TABLE IF EXISTS api_partners;

-- Drop claims columns and indexes
DROP INDEX IF EXISTS idx_claims_source;
DROP INDEX IF EXISTS idx_claims_external_id;
DROP INDEX IF EXISTS idx_claims_idempotency;

ALTER TABLE claims DROP COLUMN IF EXISTS draft_completed_at;
ALTER TABLE claims DROP COLUMN IF EXISTS is_draft;
ALTER TABLE claims DROP COLUMN IF EXISTS source_metadata;
ALTER TABLE claims DROP COLUMN IF EXISTS external_claim_id;
ALTER TABLE claims DROP COLUMN IF EXISTS idempotency_key;
ALTER TABLE claims DROP COLUMN IF EXISTS claim_source;
