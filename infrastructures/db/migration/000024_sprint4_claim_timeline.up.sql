CREATE TABLE claim_status_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    claim_id UUID NOT NULL REFERENCES claims(id),
    from_status VARCHAR(50) NOT NULL DEFAULT '',
    to_status VARCHAR(50) NOT NULL,
    action VARCHAR(100) NOT NULL,
    notes TEXT NOT NULL DEFAULT '',
    performed_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_claim_status_history_claim_id ON claim_status_history(claim_id);
