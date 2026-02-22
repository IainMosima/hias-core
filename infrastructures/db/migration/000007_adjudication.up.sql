CREATE TABLE adjudication_decisions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL REFERENCES claims(id) ON DELETE CASCADE,
    decision VARCHAR(20) NOT NULL CHECK (decision IN ('APPROVE', 'REJECT', 'MANUAL_REVIEW')),
    payable_amount BIGINT NOT NULL DEFAULT 0,
    member_responsibility BIGINT NOT NULL DEFAULT 0,
    reasons JSONB NOT NULL DEFAULT '[]',
    rule_results JSONB NOT NULL DEFAULT '[]',
    adjudicated_by UUID,
    adjudicated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_adjudication_decisions_claim_id ON adjudication_decisions(claim_id);
CREATE INDEX idx_adjudication_decisions_decision ON adjudication_decisions(decision);

CREATE TABLE fraud_flags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL REFERENCES claims(id) ON DELETE CASCADE,
    flag_type VARCHAR(30) NOT NULL CHECK (flag_type IN ('DUPLICATE', 'FREQUENCY', 'AMOUNT_THRESHOLD')),
    severity VARCHAR(20) NOT NULL DEFAULT 'MEDIUM' CHECK (severity IN ('LOW', 'MEDIUM', 'HIGH', 'CRITICAL')),
    details TEXT NOT NULL DEFAULT '',
    resolved BOOLEAN NOT NULL DEFAULT FALSE,
    resolved_by UUID,
    resolved_at TIMESTAMPTZ,
    reference_claim_id UUID REFERENCES claims(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_fraud_flags_claim_id ON fraud_flags(claim_id);
CREATE INDEX idx_fraud_flags_flag_type ON fraud_flags(flag_type);
CREATE INDEX idx_fraud_flags_resolved ON fraud_flags(resolved);
CREATE INDEX idx_fraud_flags_severity ON fraud_flags(severity);
