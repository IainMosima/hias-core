-- 1a. Plans: add premium frequency
ALTER TABLE plans ADD COLUMN premium_frequency VARCHAR(20) NOT NULL DEFAULT 'annual';

-- 1b. Premium rules: add rule classification, effective dates, ordering
ALTER TABLE premium_rules ADD COLUMN rule_type VARCHAR(30) NOT NULL DEFAULT 'base_rate';
ALTER TABLE premium_rules ADD COLUMN effective_from DATE NOT NULL DEFAULT '2020-01-01';
ALTER TABLE premium_rules ADD COLUMN effective_to DATE;
ALTER TABLE premium_rules ADD COLUMN sort_order INT NOT NULL DEFAULT 0;

CREATE INDEX idx_premium_rules_effective ON premium_rules(plan_id, effective_from, effective_to);

-- 1c. Benefits: add optional/add-on support
ALTER TABLE benefits ADD COLUMN is_optional BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE benefits ADD COLUMN addon_premium BIGINT NOT NULL DEFAULT 0;
