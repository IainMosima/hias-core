-- Scenario 7: Billing & Payment Gaps

CREATE TABLE IF NOT EXISTS premium_ledger_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    policy_id UUID NOT NULL REFERENCES policies(id),
    entry_type VARCHAR(10) NOT NULL,
    amount BIGINT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    reference_number VARCHAR(100) NOT NULL,
    effective_date TIMESTAMPTZ NOT NULL,
    balance_after BIGINT NOT NULL,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS commission_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id UUID NOT NULL REFERENCES plans(id),
    intermediary_id UUID NOT NULL REFERENCES users(id),
    rate_pct NUMERIC(5,2) NOT NULL DEFAULT 0,
    flat_amount BIGINT NOT NULL DEFAULT 0,
    effective_from TIMESTAMPTZ NOT NULL,
    effective_to TIMESTAMPTZ,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS commission_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    policy_id UUID NOT NULL REFERENCES policies(id),
    intermediary_id UUID NOT NULL REFERENCES users(id),
    commission_rule_id UUID NOT NULL REFERENCES commission_rules(id),
    amount BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'KES',
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    period_start TIMESTAMPTZ NOT NULL,
    period_end TIMESTAMPTZ NOT NULL,
    paid_at TIMESTAMPTZ,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS refunds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    policy_id UUID NOT NULL REFERENCES policies(id),
    credit_note_id UUID REFERENCES credit_notes(id),
    amount BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'KES',
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    reason TEXT NOT NULL DEFAULT '',
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMPTZ,
    processed_at TIMESTAMPTZ,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Withholding Tax columns on remittances
ALTER TABLE remittances ADD COLUMN IF NOT EXISTS wht_rate NUMERIC(5,4) NOT NULL DEFAULT 0.05;
ALTER TABLE remittances ADD COLUMN IF NOT EXISTS wht_amount BIGINT NOT NULL DEFAULT 0;
ALTER TABLE remittances ADD COLUMN IF NOT EXISTS net_amount BIGINT NOT NULL DEFAULT 0;
