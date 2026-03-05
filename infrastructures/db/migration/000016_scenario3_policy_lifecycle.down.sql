DROP TABLE IF EXISTS policy_documents;
DROP TABLE IF EXISTS underwriting_assessments;
DROP TABLE IF EXISTS policy_renewals;
DROP TABLE IF EXISTS endorsements;

ALTER TABLE policies DROP COLUMN IF EXISTS renewed_from_id;

DROP INDEX IF EXISTS idx_members_status;
ALTER TABLE members DROP COLUMN IF EXISTS status;
