-- Scenario 4: Claims Processing & Case Management

-- Case Management (inpatient case tracking)
CREATE TABLE case_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    case_number VARCHAR(50) NOT NULL UNIQUE,
    preauth_id UUID NOT NULL REFERENCES pre_authorizations(id),
    policy_id UUID NOT NULL REFERENCES policies(id),
    member_id UUID NOT NULL REFERENCES members(id),
    provider_id UUID NOT NULL REFERENCES providers(id),
    status VARCHAR(30) NOT NULL DEFAULT 'SCHEDULED',
    admission_date TIMESTAMPTZ,
    expected_discharge TIMESTAMPTZ,
    actual_discharge TIMESTAMPTZ,
    diagnosis TEXT,
    treating_doctor VARCHAR(200),
    room_type VARCHAR(50),
    total_estimated_cost BIGINT NOT NULL DEFAULT 0,
    total_actual_cost BIGINT NOT NULL DEFAULT 0,
    notes TEXT,
    closed_at TIMESTAMPTZ,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_case_records_policy ON case_records(policy_id);
CREATE INDEX idx_case_records_member ON case_records(member_id);
CREATE INDEX idx_case_records_status ON case_records(status);

-- Claim Documents (attachments on claims)
CREATE TABLE claim_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    claim_id UUID NOT NULL REFERENCES claims(id),
    file_name VARCHAR(500) NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    file_size BIGINT NOT NULL DEFAULT 0,
    s3_key VARCHAR(1000) NOT NULL,
    uploaded_by UUID NOT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_claim_documents_claim ON claim_documents(claim_id);

-- Provider Statements (uploaded by provider for reconciliation)
CREATE TABLE provider_statements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_id UUID NOT NULL REFERENCES providers(id),
    statement_number VARCHAR(50) NOT NULL UNIQUE,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    total_claimed BIGINT NOT NULL DEFAULT 0,
    total_matched BIGINT NOT NULL DEFAULT 0,
    total_discrepancy BIGINT NOT NULL DEFAULT 0,
    matched_count INT NOT NULL DEFAULT 0,
    unmatched_count INT NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'UPLOADED',
    file_name VARCHAR(500),
    s3_key VARCHAR(1000),
    reconciled_by UUID,
    reconciled_at TIMESTAMPTZ,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_provider_statements_provider ON provider_statements(provider_id);

-- Statement Line Items (individual lines from provider statement)
CREATE TABLE statement_line_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    statement_id UUID NOT NULL REFERENCES provider_statements(id),
    claim_number VARCHAR(50),
    service_date DATE,
    member_name VARCHAR(200),
    procedure_code VARCHAR(50),
    claimed_amount BIGINT NOT NULL DEFAULT 0,
    matched_claim_id UUID REFERENCES claims(id),
    match_status VARCHAR(20) NOT NULL DEFAULT 'UNMATCHED',
    discrepancy_amount BIGINT NOT NULL DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_statement_items_statement ON statement_line_items(statement_id);

-- Claims: add claim_type, vetted_amount, vetted_by, vetted_at, sla_breach_at
ALTER TABLE claims ADD COLUMN claim_type VARCHAR(30) NOT NULL DEFAULT 'DIRECT';
ALTER TABLE claims ADD COLUMN vetted_amount BIGINT;
ALTER TABLE claims ADD COLUMN vetted_by UUID;
ALTER TABLE claims ADD COLUMN vetted_at TIMESTAMPTZ;
ALTER TABLE claims ADD COLUMN sla_breach_at TIMESTAMPTZ;

-- Providers: add tier
ALTER TABLE providers ADD COLUMN tier VARCHAR(20) NOT NULL DEFAULT 'TIER_2';

-- Benefits: add deductible_amount
ALTER TABLE benefits ADD COLUMN deductible_amount BIGINT NOT NULL DEFAULT 0;
