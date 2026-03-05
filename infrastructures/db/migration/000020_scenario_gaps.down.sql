DROP INDEX IF EXISTS idx_providers_accreditation;
ALTER TABLE providers DROP COLUMN IF EXISTS accreditation_body;
ALTER TABLE providers DROP COLUMN IF EXISTS accreditation_expiry;
ALTER TABLE providers DROP COLUMN IF EXISTS accreditation_status;

DROP INDEX IF EXISTS idx_benefits_parent;
ALTER TABLE benefits DROP COLUMN IF EXISTS parent_benefit_id;
