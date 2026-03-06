ALTER TABLE claims DROP COLUMN IF EXISTS escalated_to;
ALTER TABLE approval_limits DROP COLUMN IF EXISTS entity_type;
ALTER TABLE approval_limits DROP COLUMN IF EXISTS max_claim_amount;
DROP TABLE IF EXISTS escalation_rules;
DROP TABLE IF EXISTS adjudication_rules;
