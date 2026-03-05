-- Scenario 5: Reinsurance Processing

-- 1. Treaties
CREATE TABLE treaties (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    treaty_number VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(200) NOT NULL,
    treaty_type VARCHAR(30) NOT NULL DEFAULT 'QUOTA_SHARE',
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    effective_date DATE NOT NULL,
    expiry_date DATE NOT NULL,
    retention_limit BIGINT NOT NULL DEFAULT 0,
    currency VARCHAR(10) NOT NULL DEFAULT 'KES',
    notes TEXT,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2. Treaty Participants
CREATE TABLE treaty_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    treaty_id UUID NOT NULL REFERENCES treaties(id),
    reinsurer_name VARCHAR(200) NOT NULL,
    share_percentage NUMERIC(5,2) NOT NULL DEFAULT 0,
    commission_rate NUMERIC(5,2) NOT NULL DEFAULT 0,
    is_lead BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 3. Treaty Layers (for XOL treaties)
CREATE TABLE treaty_layers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    treaty_id UUID NOT NULL REFERENCES treaties(id),
    layer_number INT NOT NULL,
    attachment_point BIGINT NOT NULL DEFAULT 0,
    layer_limit BIGINT NOT NULL DEFAULT 0,
    deductible_amount BIGINT NOT NULL DEFAULT 0,
    premium_rate NUMERIC(8,4) NOT NULL DEFAULT 0,
    aggregate_limit BIGINT,
    aggregate_used BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(treaty_id, layer_number)
);

-- 4. Profit Commissions
CREATE TABLE profit_commissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    treaty_id UUID NOT NULL REFERENCES treaties(id),
    commission_type VARCHAR(30) NOT NULL DEFAULT 'SLIDING_SCALE',
    loss_ratio_from NUMERIC(8,4) NOT NULL DEFAULT 0,
    loss_ratio_to NUMERIC(8,4) NOT NULL DEFAULT 0,
    commission_rate NUMERIC(8,4) NOT NULL DEFAULT 0,
    carry_forward_years INT NOT NULL DEFAULT 0,
    carry_forward_balance BIGINT NOT NULL DEFAULT 0,
    period_start DATE,
    period_end DATE,
    calculated_amount BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 5. Cessions
CREATE TABLE cessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cession_number VARCHAR(50) NOT NULL UNIQUE,
    treaty_id UUID NOT NULL REFERENCES treaties(id),
    policy_id UUID NOT NULL REFERENCES policies(id),
    treaty_layer_id UUID REFERENCES treaty_layers(id),
    cession_type VARCHAR(30) NOT NULL DEFAULT 'PREMIUM',
    gross_amount BIGINT NOT NULL DEFAULT 0,
    ceded_amount BIGINT NOT NULL DEFAULT 0,
    retained_amount BIGINT NOT NULL DEFAULT 0,
    commission_amount BIGINT NOT NULL DEFAULT 0,
    share_percentage NUMERIC(5,2) NOT NULL DEFAULT 0,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 6. Reinsurance Recoveries
CREATE TABLE reinsurance_recoveries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recovery_number VARCHAR(50) NOT NULL UNIQUE,
    claim_id UUID NOT NULL REFERENCES claims(id),
    treaty_id UUID NOT NULL REFERENCES treaties(id),
    treaty_layer_id UUID REFERENCES treaty_layers(id),
    cession_id UUID REFERENCES cessions(id),
    gross_claim_amount BIGINT NOT NULL DEFAULT 0,
    recoverable_amount BIGINT NOT NULL DEFAULT 0,
    recovered_amount BIGINT NOT NULL DEFAULT 0,
    outstanding_amount BIGINT NOT NULL DEFAULT 0,
    status VARCHAR(30) NOT NULL DEFAULT 'NOTIFIED',
    workflow_status VARCHAR(30) NOT NULL DEFAULT 'NOTIFICATION',
    notes TEXT,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 7. Recovery Workflow Events
CREATE TABLE recovery_workflow_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recovery_id UUID NOT NULL REFERENCES reinsurance_recoveries(id),
    from_status VARCHAR(30) NOT NULL,
    to_status VARCHAR(30) NOT NULL,
    event_type VARCHAR(30) NOT NULL,
    notes TEXT,
    performed_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 8. Bordereaux
CREATE TABLE bordereaux (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bordereau_number VARCHAR(50) NOT NULL UNIQUE,
    treaty_id UUID NOT NULL REFERENCES treaties(id),
    bordereau_type VARCHAR(30) NOT NULL DEFAULT 'PREMIUM',
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    total_gross BIGINT NOT NULL DEFAULT 0,
    total_ceded BIGINT NOT NULL DEFAULT 0,
    total_commission BIGINT NOT NULL DEFAULT 0,
    item_count INT NOT NULL DEFAULT 0,
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 9. Bordereau Items
CREATE TABLE bordereau_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bordereau_id UUID NOT NULL REFERENCES bordereaux(id),
    cession_id UUID REFERENCES cessions(id),
    recovery_id UUID REFERENCES reinsurance_recoveries(id),
    policy_number VARCHAR(50),
    claim_number VARCHAR(50),
    gross_amount BIGINT NOT NULL DEFAULT 0,
    ceded_amount BIGINT NOT NULL DEFAULT 0,
    commission_amount BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 10. Reinsurer Statements
CREATE TABLE reinsurer_statements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    statement_number VARCHAR(50) NOT NULL UNIQUE,
    treaty_id UUID NOT NULL REFERENCES treaties(id),
    participant_id UUID NOT NULL REFERENCES treaty_participants(id),
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    premium_ceded BIGINT NOT NULL DEFAULT 0,
    claims_recovered BIGINT NOT NULL DEFAULT 0,
    commission_due BIGINT NOT NULL DEFAULT 0,
    profit_commission BIGINT NOT NULL DEFAULT 0,
    net_balance BIGINT NOT NULL DEFAULT 0,
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 11. Treaty Alerts
CREATE TABLE treaty_alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    treaty_id UUID NOT NULL REFERENCES treaties(id),
    treaty_layer_id UUID REFERENCES treaty_layers(id),
    alert_type VARCHAR(50) NOT NULL,
    severity VARCHAR(20) NOT NULL DEFAULT 'MEDIUM',
    message TEXT NOT NULL,
    threshold_value BIGINT NOT NULL DEFAULT 0,
    current_value BIGINT NOT NULL DEFAULT 0,
    is_acknowledged BOOLEAN NOT NULL DEFAULT false,
    acknowledged_by UUID,
    acknowledged_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_treaties_status ON treaties(status);
CREATE INDEX idx_treaties_type ON treaties(treaty_type);
CREATE INDEX idx_treaty_participants_treaty_id ON treaty_participants(treaty_id);
CREATE INDEX idx_treaty_layers_treaty_id ON treaty_layers(treaty_id);
CREATE INDEX idx_profit_commissions_treaty_id ON profit_commissions(treaty_id);
CREATE INDEX idx_cessions_treaty_id ON cessions(treaty_id);
CREATE INDEX idx_cessions_policy_id ON cessions(policy_id);
CREATE INDEX idx_cessions_status ON cessions(status);
CREATE INDEX idx_reinsurance_recoveries_claim_id ON reinsurance_recoveries(claim_id);
CREATE INDEX idx_reinsurance_recoveries_treaty_id ON reinsurance_recoveries(treaty_id);
CREATE INDEX idx_reinsurance_recoveries_status ON reinsurance_recoveries(status);
CREATE INDEX idx_recovery_workflow_events_recovery_id ON recovery_workflow_events(recovery_id);
CREATE INDEX idx_bordereaux_treaty_id ON bordereaux(treaty_id);
CREATE INDEX idx_bordereaux_status ON bordereaux(status);
CREATE INDEX idx_bordereau_items_bordereau_id ON bordereau_items(bordereau_id);
CREATE INDEX idx_reinsurer_statements_treaty_id ON reinsurer_statements(treaty_id);
CREATE INDEX idx_reinsurer_statements_participant_id ON reinsurer_statements(participant_id);
CREATE INDEX idx_treaty_alerts_treaty_id ON treaty_alerts(treaty_id);
CREATE INDEX idx_treaty_alerts_acknowledged ON treaty_alerts(is_acknowledged);
