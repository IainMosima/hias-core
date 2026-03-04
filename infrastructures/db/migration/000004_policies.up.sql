CREATE TABLE policies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    plan_id UUID NOT NULL REFERENCES plans(id),
    policyholder_name VARCHAR(255) NOT NULL,
    policyholder_email VARCHAR(255) NOT NULL,
    policyholder_phone VARCHAR(20) NOT NULL,
    policy_number VARCHAR(20) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    start_date TIMESTAMPTZ,
    end_date TIMESTAMPTZ,
    premium_amount BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'KES',
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_policies_plan_id ON policies(plan_id);
CREATE INDEX idx_policies_policy_number ON policies(policy_number);
CREATE INDEX idx_policies_status ON policies(status);
CREATE INDEX idx_policies_policyholder_email ON policies(policyholder_email);
CREATE INDEX idx_policies_status_end_date ON policies(status, end_date);

CREATE TABLE members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    policy_id UUID NOT NULL REFERENCES policies(id) ON DELETE CASCADE,
    national_id VARCHAR(50),
    name VARCHAR(255) NOT NULL,
    date_of_birth DATE NOT NULL,
    gender VARCHAR(10) NOT NULL,
    relationship VARCHAR(20) NOT NULL,
    member_number VARCHAR(20) NOT NULL UNIQUE,
    phone VARCHAR(20),
    email VARCHAR(255),
    kra_pin VARCHAR(20),
    county VARCHAR(100),
    address TEXT,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_members_policy_id ON members(policy_id);
CREATE INDEX idx_members_national_id ON members(national_id);
CREATE INDEX idx_members_member_number ON members(member_number);
CREATE INDEX idx_members_verified ON members(verified);
