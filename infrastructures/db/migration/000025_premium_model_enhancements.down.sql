ALTER TABLE benefits DROP COLUMN IF EXISTS addon_premium;
ALTER TABLE benefits DROP COLUMN IF EXISTS is_optional;

DROP INDEX IF EXISTS idx_premium_rules_effective;
ALTER TABLE premium_rules DROP COLUMN IF EXISTS sort_order;
ALTER TABLE premium_rules DROP COLUMN IF EXISTS effective_to;
ALTER TABLE premium_rules DROP COLUMN IF EXISTS effective_from;
ALTER TABLE premium_rules DROP COLUMN IF EXISTS rule_type;

ALTER TABLE plans DROP COLUMN IF EXISTS premium_frequency;
