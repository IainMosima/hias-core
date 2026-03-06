ALTER TABLE remittances DROP COLUMN IF EXISTS wht_rate;
ALTER TABLE remittances DROP COLUMN IF EXISTS wht_amount;
ALTER TABLE remittances DROP COLUMN IF EXISTS net_amount;
DROP TABLE IF EXISTS refunds;
DROP TABLE IF EXISTS commission_payments;
DROP TABLE IF EXISTS commission_rules;
DROP TABLE IF EXISTS premium_ledger_entries;
