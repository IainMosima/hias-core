-- Add coverage date columns to members table
ALTER TABLE members ADD COLUMN coverage_start_date TIMESTAMPTZ;
ALTER TABLE members ADD COLUMN coverage_end_date TIMESTAMPTZ;

-- Backfill existing data
UPDATE members SET coverage_start_date = created_at WHERE status IN ('ACTIVE', 'PENDING');
UPDATE members SET coverage_end_date = updated_at WHERE status = 'REMOVED';
