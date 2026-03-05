-- Gap 8: Hierarchical sub-benefits
ALTER TABLE benefits ADD COLUMN IF NOT EXISTS parent_benefit_id UUID REFERENCES benefits(id);
CREATE INDEX IF NOT EXISTS idx_benefits_parent ON benefits(parent_benefit_id) WHERE parent_benefit_id IS NOT NULL;

-- Gap 9: Provider accreditation tracking
ALTER TABLE providers ADD COLUMN IF NOT EXISTS accreditation_status VARCHAR(50) DEFAULT 'NONE';
ALTER TABLE providers ADD COLUMN IF NOT EXISTS accreditation_expiry TIMESTAMPTZ;
ALTER TABLE providers ADD COLUMN IF NOT EXISTS accreditation_body VARCHAR(255);
CREATE INDEX IF NOT EXISTS idx_providers_accreditation ON providers(accreditation_status, accreditation_expiry);
