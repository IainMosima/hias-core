-- Remove restrictive CHECK constraint on claims status
-- Status validation is handled at the application level in shared/types.go
ALTER TABLE claims DROP CONSTRAINT IF EXISTS claims_status_check;
