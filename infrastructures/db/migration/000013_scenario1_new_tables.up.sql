CREATE TABLE premium_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    plan_id UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    calculation_type VARCHAR(20) NOT NULL,
    relationship VARCHAR(20),
    rate_amount BIGINT NOT NULL,
    discount_type VARCHAR(20),
    discount_value BIGINT NOT NULL DEFAULT 0,
    min_members INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_premium_rules_plan_id ON premium_rules(plan_id);

CREATE TABLE provider_networks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    plan_id UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    provider_id UUID NOT NULL REFERENCES providers(id) ON DELETE CASCADE,
    benefit_category VARCHAR(30),
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(plan_id, provider_id, benefit_category)
);

CREATE INDEX idx_provider_networks_plan_id ON provider_networks(plan_id);
CREATE INDEX idx_provider_networks_provider_id ON provider_networks(provider_id);
CREATE INDEX idx_provider_networks_status ON provider_networks(status);

CREATE TABLE installment_schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    policy_id UUID NOT NULL REFERENCES policies(id) ON DELETE CASCADE,
    frequency VARCHAR(20) NOT NULL,
    total_installments INT NOT NULL,
    amount_per_installment BIGINT NOT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_installment_schedules_policy_id ON installment_schedules(policy_id);

CREATE TABLE installments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    schedule_id UUID NOT NULL REFERENCES installment_schedules(id) ON DELETE CASCADE,
    installment_number INT NOT NULL,
    due_date TIMESTAMPTZ NOT NULL,
    amount BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    paid_at TIMESTAMPTZ,
    invoice_id UUID REFERENCES invoices(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_installments_schedule_id ON installments(schedule_id);
CREATE INDEX idx_installments_status ON installments(status);
CREATE INDEX idx_installments_due_date ON installments(due_date);
