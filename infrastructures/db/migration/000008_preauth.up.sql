CREATE TABLE preauthorizations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    policy_id UUID NOT NULL REFERENCES policies(id),
    member_id UUID NOT NULL REFERENCES members(id),
    provider_id UUID NOT NULL REFERENCES providers(id),
    auth_code VARCHAR(50) UNIQUE,
    procedure_codes JSONB NOT NULL DEFAULT '[]',
    diagnosis_codes JSONB NOT NULL DEFAULT '[]',
    estimated_cost BIGINT NOT NULL,
    approved_amount BIGINT NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'SUBMITTED' CHECK (status IN ('SUBMITTED', 'UNDER_REVIEW', 'APPROVED', 'DENIED', 'INFO_REQUESTED', 'EXPIRED', 'CLAIMED')),
    validity_start TIMESTAMPTZ,
    validity_end TIMESTAMPTZ,
    notes TEXT NOT NULL DEFAULT '',
    denial_reason TEXT,
    reviewed_by UUID,
    reviewed_at TIMESTAMPTZ,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_preauthorizations_policy_id ON preauthorizations(policy_id);
CREATE INDEX idx_preauthorizations_member_id ON preauthorizations(member_id);
CREATE INDEX idx_preauthorizations_provider_id ON preauthorizations(provider_id);
CREATE INDEX idx_preauthorizations_auth_code ON preauthorizations(auth_code);
CREATE INDEX idx_preauthorizations_status ON preauthorizations(status);
CREATE INDEX idx_preauthorizations_validity_end ON preauthorizations(validity_end);
CREATE INDEX idx_preauthorizations_status_validity ON preauthorizations(status, validity_end);
