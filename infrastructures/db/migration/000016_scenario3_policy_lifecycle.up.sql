-- Add status to members (default ACTIVE for existing rows)
ALTER TABLE members ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE';
CREATE INDEX idx_members_status ON members(status);

-- Add renewed_from_id to policies (self-referencing for renewal chain)
ALTER TABLE policies ADD COLUMN renewed_from_id UUID REFERENCES policies(id);

-- Endorsements table
CREATE TABLE endorsements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    policy_id UUID NOT NULL REFERENCES policies(id),
    endorsement_type VARCHAR(30) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    effective_date TIMESTAMPTZ NOT NULL,
    changes JSONB NOT NULL DEFAULT '{}',
    reason TEXT,
    premium_adjustment BIGINT DEFAULT 0,
    requested_by UUID NOT NULL,
    approved_by UUID,
    approved_at TIMESTAMPTZ,
    applied_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_endorsements_policy_id ON endorsements(policy_id);
CREATE INDEX idx_endorsements_status ON endorsements(status);

-- Policy renewals table
CREATE TABLE policy_renewals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    policy_id UUID NOT NULL REFERENCES policies(id),
    renewed_policy_id UUID REFERENCES policies(id),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    renewal_date TIMESTAMPTZ NOT NULL,
    new_premium BIGINT,
    premium_change_reason TEXT,
    new_plan_id UUID REFERENCES plans(id),
    approved_by UUID,
    approved_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_policy_renewals_policy_id ON policy_renewals(policy_id);
CREATE INDEX idx_policy_renewals_status ON policy_renewals(status);

-- Underwriting assessments table
CREATE TABLE underwriting_assessments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    policy_id UUID NOT NULL REFERENCES policies(id),
    member_id UUID REFERENCES members(id),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    questionnaire JSONB NOT NULL DEFAULT '{}',
    medical_declarations JSONB DEFAULT '{}',
    risk_score INT DEFAULT 0,
    risk_flags JSONB DEFAULT '[]',
    decision_reason TEXT,
    assessed_by UUID,
    assessed_at TIMESTAMPTZ,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_underwriting_policy_id ON underwriting_assessments(policy_id);
CREATE INDEX idx_underwriting_member_id ON underwriting_assessments(member_id);
CREATE INDEX idx_underwriting_status ON underwriting_assessments(status);

-- Policy documents table
CREATE TABLE policy_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    policy_id UUID NOT NULL REFERENCES policies(id),
    member_id UUID REFERENCES members(id),
    document_type VARCHAR(30) NOT NULL,
    file_name VARCHAR(500) NOT NULL,
    file_size BIGINT DEFAULT 0,
    s3_key VARCHAR(1000) NOT NULL,
    generated_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_policy_documents_policy_id ON policy_documents(policy_id);
CREATE INDEX idx_policy_documents_document_type ON policy_documents(document_type);
