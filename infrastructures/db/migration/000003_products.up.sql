CREATE TABLE plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(20) NOT NULL,
    segment VARCHAR(20) NOT NULL DEFAULT 'retail',
    base_premium BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'KES',
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    description TEXT NOT NULL DEFAULT '',
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_plans_status ON plans(status);
CREATE INDEX idx_plans_type ON plans(type);

CREATE TABLE benefits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    plan_id UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(30) NOT NULL,
    annual_limit BIGINT NOT NULL,
    co_pay_type VARCHAR(20) NOT NULL DEFAULT 'percentage',
    co_pay_value BIGINT NOT NULL DEFAULT 0,
    waiting_period_days INT NOT NULL DEFAULT 0,
    sub_limit_type VARCHAR(20) NOT NULL DEFAULT 'none',
    sub_limit_value BIGINT NOT NULL DEFAULT 0,
    min_age INT NOT NULL DEFAULT 0,
    max_age INT NOT NULL DEFAULT 150,
    waiting_period_type VARCHAR(30) NOT NULL DEFAULT 'general',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_benefits_plan_id ON benefits(plan_id);
CREATE INDEX idx_benefits_category ON benefits(category);

CREATE TABLE exclusions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    plan_id UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    type VARCHAR(30) NOT NULL,
    icd_codes JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_exclusions_plan_id ON exclusions(plan_id);
CREATE INDEX idx_exclusions_type ON exclusions(type);
