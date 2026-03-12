-- Source tracking on claims
ALTER TABLE claims ADD COLUMN claim_source VARCHAR(30) NOT NULL DEFAULT 'INTERNAL';
ALTER TABLE claims ADD COLUMN idempotency_key VARCHAR(100);
ALTER TABLE claims ADD COLUMN external_claim_id VARCHAR(100);
ALTER TABLE claims ADD COLUMN source_metadata JSONB NOT NULL DEFAULT '{}';
ALTER TABLE claims ADD COLUMN is_draft BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE claims ADD COLUMN draft_completed_at TIMESTAMPTZ;

CREATE UNIQUE INDEX idx_claims_idempotency ON claims(idempotency_key) WHERE idempotency_key IS NOT NULL;
CREATE INDEX idx_claims_external_id ON claims(external_claim_id) WHERE external_claim_id IS NOT NULL;
CREATE INDEX idx_claims_source ON claims(claim_source);

-- API Partners table
CREATE TABLE api_partners (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    partner_type VARCHAR(30) NOT NULL,
    api_key VARCHAR(64) NOT NULL UNIQUE,
    api_secret_hash VARCHAR(128) NOT NULL,
    provider_id UUID REFERENCES providers(id),
    is_active BOOLEAN NOT NULL DEFAULT true,
    rate_limit_per_minute INT NOT NULL DEFAULT 60,
    allowed_claim_types TEXT[] DEFAULT '{"DIRECT"}',
    webhook_url VARCHAR(500),
    contact_email VARCHAR(200),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_api_partners_api_key ON api_partners(api_key);
CREATE INDEX idx_api_partners_provider ON api_partners(provider_id) WHERE provider_id IS NOT NULL;
