-- Set all existing providers to ACTIVE accreditation status
UPDATE providers SET accreditation_status = 'ACTIVE'
WHERE accreditation_status IS NULL OR accreditation_status = '' OR accreditation_status = 'NONE';
