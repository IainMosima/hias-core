CREATE TABLE claims (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_number VARCHAR(20) NOT NULL UNIQUE,
    policy_id UUID NOT NULL REFERENCES policies(id),
    member_id UUID NOT NULL REFERENCES members(id),
    provider_id UUID NOT NULL REFERENCES providers(id),
    preauth_id UUID,
    status VARCHAR(20) NOT NULL DEFAULT 'RECEIVED' CHECK (status IN ('RECEIVED', 'VALIDATED', 'ADJUDICATED', 'APPROVED', 'REJECTED', 'MANUAL_REVIEW', 'PAID')),
    total_amount BIGINT NOT NULL,
    approved_amount BIGINT NOT NULL DEFAULT 0,
    co_pay_amount BIGINT NOT NULL DEFAULT 0,
    member_responsibility BIGINT NOT NULL DEFAULT 0,
    diagnosis_codes JSONB NOT NULL DEFAULT '[]',
    service_date TIMESTAMPTZ NOT NULL,
    admission_date TIMESTAMPTZ,
    discharge_date TIMESTAMPTZ,
    notes TEXT NOT NULL DEFAULT '',
    rejection_reason TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_claims_claim_number ON claims(claim_number);
CREATE INDEX idx_claims_policy_id ON claims(policy_id);
CREATE INDEX idx_claims_member_id ON claims(member_id);
CREATE INDEX idx_claims_provider_id ON claims(provider_id);
CREATE INDEX idx_claims_preauth_id ON claims(preauth_id);
CREATE INDEX idx_claims_status ON claims(status);
CREATE INDEX idx_claims_service_date ON claims(service_date);
CREATE INDEX idx_claims_status_provider ON claims(status, provider_id);
CREATE INDEX idx_claims_status_member ON claims(status, member_id);
CREATE INDEX idx_claims_created_at ON claims(created_at);

CREATE TABLE claim_line_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL REFERENCES claims(id) ON DELETE CASCADE,
    procedure_code VARCHAR(20) NOT NULL,
    procedure_name VARCHAR(255) NOT NULL,
    diagnosis_code VARCHAR(20),
    quantity INT NOT NULL DEFAULT 1,
    unit_price BIGINT NOT NULL,
    total_price BIGINT NOT NULL,
    approved_amount BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_claim_line_items_claim_id ON claim_line_items(claim_id);
CREATE INDEX idx_claim_line_items_procedure_code ON claim_line_items(procedure_code);
