-- Scenario 2: Sales Domain — Leads, Quotations, Approvals, Documents, Activities

-- 1. Leads
CREATE TABLE leads (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    lead_number VARCHAR(20) NOT NULL UNIQUE,
    contact_name VARCHAR(255) NOT NULL,
    contact_email VARCHAR(255),
    contact_phone VARCHAR(50),
    company_name VARCHAR(255),
    source VARCHAR(50) NOT NULL DEFAULT 'direct',
    segment VARCHAR(50) NOT NULL DEFAULT 'retail',
    plan_type VARCHAR(50) NOT NULL DEFAULT 'individual',
    estimated_members INT NOT NULL DEFAULT 1,
    expected_premium BIGINT NOT NULL DEFAULT 0,
    closure_probability INT NOT NULL DEFAULT 0,
    currency VARCHAR(10) NOT NULL DEFAULT 'KES',
    status VARCHAR(50) NOT NULL DEFAULT 'NEW',
    assigned_to UUID REFERENCES users(id),
    next_follow_up_date TIMESTAMPTZ,
    notes TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_leads_status ON leads(status);
CREATE INDEX idx_leads_assigned_to ON leads(assigned_to);
CREATE INDEX idx_leads_next_follow_up ON leads(next_follow_up_date);
CREATE INDEX idx_leads_source ON leads(source);

-- 2. Quotations
CREATE TABLE quotations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    quotation_number VARCHAR(20) NOT NULL UNIQUE,
    lead_id UUID NOT NULL REFERENCES leads(id),
    plan_id UUID NOT NULL REFERENCES plans(id),
    quotation_type VARCHAR(50) NOT NULL DEFAULT 'standard',
    status VARCHAR(50) NOT NULL DEFAULT 'DRAFT',
    current_version INT NOT NULL DEFAULT 1,
    policy_id UUID REFERENCES policies(id),
    valid_from TIMESTAMPTZ,
    valid_until TIMESTAMPTZ,
    client_name VARCHAR(255) NOT NULL,
    client_email VARCHAR(255),
    client_phone VARCHAR(50),
    currency VARCHAR(10) NOT NULL DEFAULT 'KES',
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_quotations_status ON quotations(status);
CREATE INDEX idx_quotations_lead_id ON quotations(lead_id);
CREATE INDEX idx_quotations_valid_until ON quotations(valid_until);

-- 3. Quotation Versions
CREATE TABLE quotation_versions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    quotation_id UUID NOT NULL REFERENCES quotations(id),
    version_number INT NOT NULL,
    base_premium BIGINT NOT NULL DEFAULT 0,
    discount_type VARCHAR(50) NOT NULL DEFAULT 'percentage',
    discount_value BIGINT NOT NULL DEFAULT 0,
    discount_reason TEXT,
    loading_type VARCHAR(50) NOT NULL DEFAULT 'percentage',
    loading_value BIGINT NOT NULL DEFAULT 0,
    loading_reason TEXT,
    final_premium BIGINT NOT NULL DEFAULT 0,
    member_count INT NOT NULL DEFAULT 1,
    proposed_members JSONB NOT NULL DEFAULT '[]'::jsonb,
    billing_frequency VARCHAR(50) NOT NULL DEFAULT 'monthly',
    requires_approval BOOLEAN NOT NULL DEFAULT false,
    approval_status VARCHAR(50) NOT NULL DEFAULT 'NONE',
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMPTZ,
    rejection_reason TEXT,
    pricing_breakdown JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(quotation_id, version_number)
);

CREATE INDEX idx_quotation_versions_quotation_id ON quotation_versions(quotation_id);
CREATE INDEX idx_quotation_versions_approval_status ON quotation_versions(approval_status);

-- 4. Approval Limits
CREATE TABLE approval_limits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role_name VARCHAR(100) NOT NULL UNIQUE,
    max_discount_percentage BIGINT NOT NULL DEFAULT 0,
    max_discount_amount BIGINT NOT NULL DEFAULT 0,
    max_loading_percentage BIGINT NOT NULL DEFAULT 0,
    max_loading_amount BIGINT NOT NULL DEFAULT 0,
    escalation_role VARCHAR(100),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Seed default approval limits
INSERT INTO approval_limits (role_name, max_discount_percentage, max_discount_amount, max_loading_percentage, max_loading_amount, escalation_role)
VALUES
    ('Underwriter', 1000, 10000000, 1000, 10000000, 'Admin'),
    ('Admin', 10000, 0, 10000, 0, NULL);

-- 5. Quotation Documents
CREATE TABLE quotation_documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    quotation_id UUID NOT NULL REFERENCES quotations(id),
    version_number INT,
    file_name VARCHAR(500) NOT NULL,
    file_type VARCHAR(100) NOT NULL,
    file_size BIGINT NOT NULL DEFAULT 0,
    s3_key VARCHAR(1000) NOT NULL,
    uploaded_by UUID REFERENCES users(id),
    can_edit_roles JSONB NOT NULL DEFAULT '["Admin"]'::jsonb,
    can_delete_roles JSONB NOT NULL DEFAULT '["Admin"]'::jsonb,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_quotation_documents_quotation_id ON quotation_documents(quotation_id);

-- 6. Lead Activities
CREATE TABLE lead_activities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    lead_id UUID NOT NULL REFERENCES leads(id),
    activity_type VARCHAR(50) NOT NULL,
    description TEXT,
    scheduled_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_lead_activities_lead_id ON lead_activities(lead_id);
CREATE INDEX idx_lead_activities_type ON lead_activities(activity_type);
