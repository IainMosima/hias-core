-- Revert to original status constraint
ALTER TABLE claims DROP CONSTRAINT IF EXISTS claims_status_check;
ALTER TABLE claims ADD CONSTRAINT claims_status_check CHECK (
    status IN ('RECEIVED', 'VALIDATED', 'ADJUDICATED', 'APPROVED', 'REJECTED', 'MANUAL_REVIEW', 'PAID')
);
