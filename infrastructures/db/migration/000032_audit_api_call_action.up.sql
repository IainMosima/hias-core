ALTER TABLE audit_events DROP CONSTRAINT audit_events_action_check;
ALTER TABLE audit_events ADD CONSTRAINT audit_events_action_check
    CHECK (action IN ('CREATE', 'UPDATE', 'DELETE', 'STATE_CHANGE', 'API_CALL'));
