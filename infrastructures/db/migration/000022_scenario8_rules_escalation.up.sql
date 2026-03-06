-- Scenario 8: Configurable Adjudication Rules & Escalation

CREATE TABLE IF NOT EXISTS adjudication_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    rule_type VARCHAR(50) NOT NULL,
    parameters JSONB NOT NULL DEFAULT '{}',
    priority INT NOT NULL DEFAULT 100,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS escalation_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    condition_type VARCHAR(50) NOT NULL,
    threshold_amount BIGINT NOT NULL DEFAULT 0,
    escalation_role VARCHAR(50) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE claims ADD COLUMN IF NOT EXISTS escalated_to VARCHAR(50);

-- Approval hierarchy for claims (Step 7)
ALTER TABLE approval_limits ADD COLUMN IF NOT EXISTS entity_type VARCHAR(50) NOT NULL DEFAULT 'QUOTATION';
ALTER TABLE approval_limits ADD COLUMN IF NOT EXISTS max_claim_amount BIGINT NOT NULL DEFAULT 0;

-- Seed default adjudication rules
INSERT INTO adjudication_rules (name, rule_type, parameters, priority) VALUES
('Low-value auto-approve', 'AMOUNT_THRESHOLD', '{"max_amount": 5000000}', 10),
('Frequency limit', 'FREQUENCY_LIMIT', '{"max_per_month": 4}', 20),
('Benefit coverage check', 'BENEFIT_CHECK', '{"require_active_benefit": true}', 30),
('Auto-approve pre-authorized', 'AUTO_APPROVE', '{"require_preauth": true}', 5);

-- Seed default escalation rules
INSERT INTO escalation_rules (name, condition_type, threshold_amount, escalation_role) VALUES
('High-value claim', 'AMOUNT_EXCEEDS', 50000000, 'Manager'),
('Fraud-flagged', 'FRAUD_FLAG', 0, 'Manager');
