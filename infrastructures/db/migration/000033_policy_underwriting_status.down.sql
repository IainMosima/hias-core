ALTER TABLE policies DROP COLUMN underwriting_status;

ALTER TABLE underwriting_assessments ALTER COLUMN status TYPE VARCHAR(20);
