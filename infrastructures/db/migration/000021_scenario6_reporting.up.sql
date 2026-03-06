-- Scenario 6: Operational & Management Reporting

-- Table 1: Report Definitions - Pre-built report templates + ad-hoc saved reports
CREATE TABLE report_definitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT DEFAULT '',
    category VARCHAR(50) NOT NULL,
    report_type VARCHAR(20) NOT NULL,
    query_template TEXT DEFAULT '',
    default_parameters JSONB DEFAULT '{}',
    allowed_roles TEXT[] NOT NULL,
    columns JSONB NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Table 2: Report Schedules - Recurring report schedules
CREATE TABLE report_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_definition_id UUID NOT NULL REFERENCES report_definitions(id),
    name VARCHAR(200) NOT NULL,
    cron_expression VARCHAR(100) NOT NULL,
    parameters JSONB DEFAULT '{}',
    export_format VARCHAR(10) NOT NULL,
    recipients UUID[] NOT NULL,
    is_active BOOLEAN DEFAULT true,
    last_run_at TIMESTAMPTZ,
    next_run_at TIMESTAMPTZ,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Table 3: Generated Reports - Report execution history + stored output
CREATE TABLE generated_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_definition_id UUID NOT NULL REFERENCES report_definitions(id),
    schedule_id UUID REFERENCES report_schedules(id),
    report_number VARCHAR(30) UNIQUE NOT NULL,
    name VARCHAR(200) NOT NULL,
    parameters JSONB DEFAULT '{}',
    format VARCHAR(10) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'GENERATING',
    row_count INT DEFAULT 0,
    file_data BYTEA,
    file_size BIGINT DEFAULT 0,
    error_message TEXT DEFAULT '',
    generated_by UUID REFERENCES users(id),
    generated_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_generated_reports_definition ON generated_reports(report_definition_id);
CREATE INDEX idx_generated_reports_status ON generated_reports(status);
CREATE INDEX idx_generated_reports_generated_by ON generated_reports(generated_by);

-- Seed Data: 8 Pre-Built Report Definitions
INSERT INTO report_definitions (code, name, description, category, report_type, allowed_roles, columns, default_parameters) VALUES
('CLAIMS_EXPERIENCE', 'Claims Experience Report', 'Loss ratio analysis by policy/corporate client with claims breakdown', 'CLAIMS', 'PRE_BUILT',
 ARRAY['Admin','Manager','Finance','Underwriter'],
 '[{"name":"policy_number","label":"Policy No","type":"string"},{"name":"policyholder_name","label":"Policyholder","type":"string"},{"name":"total_premium","label":"Total Premium","type":"money"},{"name":"total_claims","label":"Total Claims","type":"money"},{"name":"approved_claims","label":"Approved","type":"money"},{"name":"rejected_claims","label":"Rejected","type":"money"},{"name":"loss_ratio","label":"Loss Ratio %","type":"percentage"},{"name":"claim_count","label":"# Claims","type":"number"},{"name":"avg_tat_hours","label":"Avg TAT (hrs)","type":"decimal"}]',
 '{"period":"year"}'),

('CLAIMS_REGISTER', 'Claims Register', 'Complete register of all claims with status, amounts, and provider details', 'CLAIMS', 'PRE_BUILT',
 ARRAY['Admin','Manager','ClaimsOfficer','Finance'],
 '[{"name":"claim_number","label":"Claim No","type":"string"},{"name":"policy_number","label":"Policy No","type":"string"},{"name":"member_name","label":"Member","type":"string"},{"name":"provider_name","label":"Provider","type":"string"},{"name":"claim_type","label":"Type","type":"string"},{"name":"service_date","label":"Service Date","type":"date"},{"name":"total_amount","label":"Claimed","type":"money"},{"name":"approved_amount","label":"Approved","type":"money"},{"name":"status","label":"Status","type":"string"},{"name":"created_at","label":"Submitted","type":"datetime"}]',
 '{"period":"month"}'),

('PREMIUM_DEBTORS_AGEING', 'Premium Debtors Ageing', 'Outstanding premium balances with ageing buckets (current, 30, 60, 90+ days)', 'PREMIUM', 'PRE_BUILT',
 ARRAY['Admin','Manager','Finance'],
 '[{"name":"policy_number","label":"Policy No","type":"string"},{"name":"policyholder_name","label":"Policyholder","type":"string"},{"name":"total_premium","label":"Total Premium","type":"money"},{"name":"total_paid","label":"Total Paid","type":"money"},{"name":"outstanding","label":"Outstanding","type":"money"},{"name":"current_bucket","label":"Current","type":"money"},{"name":"days_30","label":"30 Days","type":"money"},{"name":"days_60","label":"60 Days","type":"money"},{"name":"days_90_plus","label":"90+ Days","type":"money"}]',
 '{}'),

('PREMIUM_REGISTER', 'Premium Register', 'Complete premium collection register with payment details', 'PREMIUM', 'PRE_BUILT',
 ARRAY['Admin','Manager','Finance'],
 '[{"name":"policy_number","label":"Policy No","type":"string"},{"name":"policyholder_name","label":"Policyholder","type":"string"},{"name":"plan_name","label":"Plan","type":"string"},{"name":"premium_amount","label":"Premium","type":"money"},{"name":"payment_amount","label":"Paid","type":"money"},{"name":"payment_date","label":"Payment Date","type":"date"},{"name":"payment_status","label":"Status","type":"string"},{"name":"payment_method","label":"Method","type":"string"}]',
 '{"period":"month"}'),

('MEMBERSHIP', 'Membership Report', 'Member listing with KYC data (DOB, phone, gender, relationship)', 'MEMBERSHIP', 'PRE_BUILT',
 ARRAY['Admin','Manager','Underwriter','SalesAgent'],
 '[{"name":"member_number","label":"Member No","type":"string"},{"name":"full_name","label":"Full Name","type":"string"},{"name":"date_of_birth","label":"DOB","type":"date"},{"name":"gender","label":"Gender","type":"string"},{"name":"phone","label":"Phone","type":"string"},{"name":"email","label":"Email","type":"string"},{"name":"relationship","label":"Relationship","type":"string"},{"name":"policy_number","label":"Policy No","type":"string"},{"name":"plan_name","label":"Plan","type":"string"},{"name":"status","label":"Status","type":"string"},{"name":"enrollment_date","label":"Enrolled","type":"date"}]',
 '{}'),

('PROVIDER_PERFORMANCE', 'Provider Performance Report', 'Provider utilization, claim patterns, and turnaround analysis', 'PROVIDER', 'PRE_BUILT',
 ARRAY['Admin','Manager','ClaimsOfficer'],
 '[{"name":"provider_name","label":"Provider","type":"string"},{"name":"tier","label":"Tier","type":"string"},{"name":"status","label":"Status","type":"string"},{"name":"total_claims","label":"# Claims","type":"number"},{"name":"total_claimed","label":"Total Claimed","type":"money"},{"name":"total_approved","label":"Total Approved","type":"money"},{"name":"rejection_rate","label":"Rejection %","type":"percentage"},{"name":"avg_claim_amount","label":"Avg Claim","type":"money"},{"name":"fraud_flag_count","label":"Fraud Flags","type":"number"}]',
 '{"period":"year"}'),

('LOSS_RATIO', 'Loss Ratio Report', 'Loss ratio breakdown by plan, policy type, and time period', 'MANAGEMENT', 'PRE_BUILT',
 ARRAY['Admin','Manager','Finance','Underwriter'],
 '[{"name":"plan_name","label":"Plan","type":"string"},{"name":"active_policies","label":"Active Policies","type":"number"},{"name":"total_members","label":"Members","type":"number"},{"name":"earned_premium","label":"Earned Premium","type":"money"},{"name":"incurred_claims","label":"Incurred Claims","type":"money"},{"name":"loss_ratio","label":"Loss Ratio %","type":"percentage"},{"name":"expense_ratio","label":"Expense Ratio %","type":"percentage"},{"name":"combined_ratio","label":"Combined Ratio %","type":"percentage"}]',
 '{"period":"year"}'),

('RENEWAL', 'Renewal Report', 'Policy renewal tracking with retention analysis', 'MANAGEMENT', 'PRE_BUILT',
 ARRAY['Admin','Manager','SalesAgent','Underwriter'],
 '[{"name":"policy_number","label":"Policy No","type":"string"},{"name":"policyholder_name","label":"Policyholder","type":"string"},{"name":"plan_name","label":"Plan","type":"string"},{"name":"expiry_date","label":"Expiry Date","type":"date"},{"name":"renewal_status","label":"Renewal Status","type":"string"},{"name":"current_premium","label":"Current Premium","type":"money"},{"name":"proposed_premium","label":"Proposed Premium","type":"money"},{"name":"premium_change_pct","label":"Change %","type":"percentage"},{"name":"member_count","label":"Members","type":"number"}]',
 '{"period":"quarter"}');
