-- Reverse Scenario 4: Claims Processing & Case Management

-- Drop added columns
ALTER TABLE benefits DROP COLUMN IF EXISTS deductible_amount;
ALTER TABLE providers DROP COLUMN IF EXISTS tier;
ALTER TABLE claims DROP COLUMN IF EXISTS sla_breach_at;
ALTER TABLE claims DROP COLUMN IF EXISTS vetted_at;
ALTER TABLE claims DROP COLUMN IF EXISTS vetted_by;
ALTER TABLE claims DROP COLUMN IF EXISTS vetted_amount;
ALTER TABLE claims DROP COLUMN IF EXISTS claim_type;

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS statement_line_items;
DROP TABLE IF EXISTS provider_statements;
DROP TABLE IF EXISTS claim_documents;
DROP TABLE IF EXISTS case_records;
