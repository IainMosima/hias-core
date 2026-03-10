-- Add activation timestamp to track when policy was activated
ALTER TABLE policies ADD COLUMN activated_at TIMESTAMPTZ;

-- Change default member status from ACTIVE to PENDING
ALTER TABLE members ALTER COLUMN status SET DEFAULT 'PENDING';
