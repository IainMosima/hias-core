-- Underwriting Flags (auditable records of detected violations)
CREATE TABLE underwriting_flags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_id UUID REFERENCES underwriting_assessments(id),
    policy_id UUID NOT NULL REFERENCES policies(id),
    member_id UUID REFERENCES members(id),
    flag_type VARCHAR(50) NOT NULL,
    severity VARCHAR(20) NOT NULL,
    details TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN',
    resolved_by UUID,
    resolved_at TIMESTAMPTZ,
    resolution TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_uw_flags_policy ON underwriting_flags(policy_id);
CREATE INDEX idx_uw_flags_member ON underwriting_flags(member_id);
CREATE INDEX idx_uw_flags_status ON underwriting_flags(status);

-- Underwriting Rules (configurable per plan, like premium_rules)
CREATE TABLE underwriting_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id UUID NOT NULL REFERENCES plans(id),
    rule_type VARCHAR(50) NOT NULL,
    relationship VARCHAR(20),
    parameter_key VARCHAR(100) NOT NULL,
    parameter_value VARCHAR(500) NOT NULL,
    severity VARCHAR(20) NOT NULL DEFAULT 'MEDIUM',
    risk_score_weight INT NOT NULL DEFAULT 20,
    is_blocking BOOLEAN NOT NULL DEFAULT false,
    is_active BOOLEAN NOT NULL DEFAULT true,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_uw_rules_plan ON underwriting_rules(plan_id);

-- Credit Notes (financial refund records)
CREATE TABLE credit_notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    policy_id UUID NOT NULL REFERENCES policies(id),
    member_id UUID REFERENCES members(id),
    credit_note_number VARCHAR(50) NOT NULL UNIQUE,
    amount BIGINT NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'KES',
    reason TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    applied_to_invoice_id UUID REFERENCES invoices(id),
    approved_by UUID,
    approved_at TIMESTAMPTZ,
    applied_at TIMESTAMPTZ,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_credit_notes_policy ON credit_notes(policy_id);
CREATE INDEX idx_credit_notes_status ON credit_notes(status);
