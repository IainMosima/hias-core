-- Revert accreditation status back to NONE
UPDATE providers SET accreditation_status = 'NONE'
WHERE accreditation_status = 'ACTIVE';
