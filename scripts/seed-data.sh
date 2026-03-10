#!/bin/bash
set -e
DB_URL="${DB_URL:-postgresql://root:supersecret@localhost:5432/hias_db?sslmode=disable}"

echo "Seeding comprehensive demo data..."

psql "$DB_URL" << 'EOF'
BEGIN;

-- ============================================================
-- LAYER 1: ROLES & PERMISSIONS (foundation)
-- ============================================================
INSERT INTO roles (id, name, description) VALUES
  ('10000000-0000-0000-0000-000000000001', 'Admin',         'Full system access'),
  ('10000000-0000-0000-0000-000000000002', 'Underwriter',   'Policy management and member enrollment'),
  ('10000000-0000-0000-0000-000000000003', 'ClaimsOfficer', 'Claims review and adjudication'),
  ('10000000-0000-0000-0000-000000000004', 'Finance',       'Financial operations and reporting'),
  ('10000000-0000-0000-0000-000000000005', 'Provider',      'Healthcare provider portal access'),
  ('10000000-0000-0000-0000-000000000006', 'Member',        'Member self-service portal access'),
  ('10000000-0000-0000-0000-000000000007', 'SalesAgent',    'Sales pipeline and quotation management'),
  ('10000000-0000-0000-0000-000000000008', 'Manager',       'Team management and approvals')
ON CONFLICT (name) DO NOTHING;

INSERT INTO permissions (id, resource, action, description) VALUES
  (uuid_generate_v4(), 'users',         'create',     'Create users'),
  (uuid_generate_v4(), 'users',         'read',       'View users'),
  (uuid_generate_v4(), 'users',         'update',     'Update users'),
  (uuid_generate_v4(), 'users',         'delete',     'Delete users'),
  (uuid_generate_v4(), 'plans',         'create',     'Create insurance plans'),
  (uuid_generate_v4(), 'plans',         'read',       'View insurance plans'),
  (uuid_generate_v4(), 'plans',         'update',     'Update insurance plans'),
  (uuid_generate_v4(), 'policies',      'create',     'Create policies'),
  (uuid_generate_v4(), 'policies',      'read',       'View policies'),
  (uuid_generate_v4(), 'policies',      'update',     'Update policies'),
  (uuid_generate_v4(), 'policies',      'activate',   'Activate policies'),
  (uuid_generate_v4(), 'providers',     'create',     'Register providers'),
  (uuid_generate_v4(), 'providers',     'read',       'View providers'),
  (uuid_generate_v4(), 'providers',     'update',     'Update providers'),
  (uuid_generate_v4(), 'providers',     'credential', 'Credential providers'),
  (uuid_generate_v4(), 'claims',        'create',     'Submit claims'),
  (uuid_generate_v4(), 'claims',        'read',       'View claims'),
  (uuid_generate_v4(), 'claims',        'review',     'Review claims'),
  (uuid_generate_v4(), 'claims',        'approve',    'Approve claims'),
  (uuid_generate_v4(), 'claims',        'reject',     'Reject claims'),
  (uuid_generate_v4(), 'preauth',       'create',     'Submit pre-authorizations'),
  (uuid_generate_v4(), 'preauth',       'read',       'View pre-authorizations'),
  (uuid_generate_v4(), 'preauth',       'approve',    'Approve pre-authorizations'),
  (uuid_generate_v4(), 'billing',       'read',       'View billing information'),
  (uuid_generate_v4(), 'billing',       'create',     'Create invoices'),
  (uuid_generate_v4(), 'payments',      'read',       'View payments'),
  (uuid_generate_v4(), 'payments',      'process',    'Process payments'),
  (uuid_generate_v4(), 'analytics',     'read',       'View analytics'),
  (uuid_generate_v4(), 'audit',         'read',       'View audit trail'),
  (uuid_generate_v4(), 'notifications', 'read',       'View notifications')
ON CONFLICT (resource, action) DO NOTHING;

-- Grant Admin all permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT '10000000-0000-0000-0000-000000000001', id FROM permissions
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- ============================================================
-- LAYER 2: USERS
-- ============================================================
-- Password: Password123!  bcrypt hash at cost 10
INSERT INTO users (id, email, name, phone, role_id, status, password_hash) VALUES
  ('20000000-0000-0000-0000-000000000001', 'admin@hias.co.ke',  'Admin User',     '+254700000001', '10000000-0000-0000-0000-000000000001', 'ACTIVE', '$2a$10$Zzk5j8vZSCG27ODa9ctJoOTy.6vfudm6pzWhKV4Um69Tay/TKF8OK'),
  ('20000000-0000-0000-0000-000000000002', 'jane@hias.co.ke',   'Jane Muthoni',   '+254700000002', '10000000-0000-0000-0000-000000000003', 'ACTIVE', '$2a$10$Zzk5j8vZSCG27ODa9ctJoOTy.6vfudm6pzWhKV4Um69Tay/TKF8OK'),
  ('20000000-0000-0000-0000-000000000003', 'peter@hias.co.ke',  'Peter Kamau',    '+254700000003', '10000000-0000-0000-0000-000000000002', 'ACTIVE', '$2a$10$Zzk5j8vZSCG27ODa9ctJoOTy.6vfudm6pzWhKV4Um69Tay/TKF8OK'),
  ('20000000-0000-0000-0000-000000000004', 'sarah@hias.co.ke',  'Sarah Wanjiku',  '+254700000004', '10000000-0000-0000-0000-000000000004', 'ACTIVE', '$2a$10$Zzk5j8vZSCG27ODa9ctJoOTy.6vfudm6pzWhKV4Um69Tay/TKF8OK')
ON CONFLICT (email) DO NOTHING;

-- ============================================================
-- LAYER 3: PRODUCTS (Plans + Benefits + Exclusions)
-- ============================================================
INSERT INTO plans (id, name, type, segment, base_premium, currency, status, description) VALUES
  ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Afya Basic',      'individual', 'retail',    250000,  'KES', 'ACTIVE', 'Affordable individual cover with essential benefits'),
  ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'Afya Family',     'group',      'retail',    750000,  'KES', 'ACTIVE', 'Comprehensive family cover for up to 6 members'),
  ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Corporate Gold',  'group',      'corporate', 1500000, 'KES', 'ACTIVE', 'Premium corporate cover with enhanced limits')
ON CONFLICT (id) DO NOTHING;

-- Benefits: Afya Basic
INSERT INTO benefits (id, plan_id, name, category, annual_limit, co_pay_type, co_pay_value, waiting_period_days, deductible_amount) VALUES
  ('30000000-0000-0000-0001-000000000001', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Outpatient Care', 'outpatient', 10000000, 'percentage', 1000, 0,  0),
  ('30000000-0000-0000-0001-000000000002', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Inpatient Care',  'inpatient',  50000000, 'percentage', 1000, 30, 0),
  ('30000000-0000-0000-0001-000000000003', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Dental Care',     'dental',      5000000, 'fixed',       50000, 90, 0),
  ('30000000-0000-0000-0001-000000000004', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Optical Care',    'optical',     3000000, 'fixed',       30000, 90, 0),
  ('30000000-0000-0000-0001-000000000005', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Maternity',       'maternity',  30000000, 'percentage',    500, 270, 0),
  ('30000000-0000-0000-0001-000000000006', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Pharmacy',        'pharmacy',   2000000,  'fixed',         20000, 0, 0)
ON CONFLICT DO NOTHING;

-- Benefits: Afya Family
INSERT INTO benefits (id, plan_id, name, category, annual_limit, co_pay_type, co_pay_value, waiting_period_days, deductible_amount) VALUES
  ('30000000-0000-0000-0002-000000000001', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'Outpatient Care', 'outpatient', 20000000, 'percentage', 1000, 0,  0),
  ('30000000-0000-0000-0002-000000000002', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'Inpatient Care',  'inpatient',  80000000, 'percentage',  500, 14, 0),
  ('30000000-0000-0000-0002-000000000003', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'Dental Care',     'dental',     10000000, 'fixed',       30000, 60, 0),
  ('30000000-0000-0000-0002-000000000004', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'Optical Care',    'optical',     5000000, 'fixed',       20000, 60, 0),
  ('30000000-0000-0000-0002-000000000005', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'Maternity',       'maternity',  50000000, 'percentage',    500, 270, 0),
  ('30000000-0000-0000-0002-000000000006', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'Pharmacy',        'pharmacy',   4000000,  'fixed',         15000, 0, 0)
ON CONFLICT DO NOTHING;

-- Benefits: Corporate Gold
INSERT INTO benefits (id, plan_id, name, category, annual_limit, co_pay_type, co_pay_value, waiting_period_days, deductible_amount) VALUES
  ('30000000-0000-0000-0003-000000000001', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Outpatient Care', 'outpatient',  50000000, 'percentage',  500, 0,  0),
  ('30000000-0000-0000-0003-000000000002', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Inpatient Care',  'inpatient',  200000000, 'percentage',  500, 0,  0),
  ('30000000-0000-0000-0003-000000000003', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Dental Care',     'dental',      20000000, 'fixed',       10000, 30, 0),
  ('30000000-0000-0000-0003-000000000004', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Optical Care',    'optical',     10000000, 'fixed',       10000, 30, 0),
  ('30000000-0000-0000-0003-000000000005', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Maternity',       'maternity',  100000000, 'percentage',    0, 180, 0),
  ('30000000-0000-0000-0003-000000000006', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Pharmacy',        'pharmacy',   10000000, 'fixed',         5000, 0, 0),
  ('30000000-0000-0000-0003-000000000007', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Specialist Care', 'specialist',  30000000, 'percentage',  500, 30, 0),
  ('30000000-0000-0000-0003-000000000008', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Emergency',       'emergency',  50000000, 'percentage',  200, 0, 0),
  ('30000000-0000-0000-0003-000000000009', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Chronic Care',    'chronic',   40000000, 'percentage',  300, 90, 0),
  ('30000000-0000-0000-0003-000000000010', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Wellness',        'wellness',    5000000, 'fixed',       10000, 0, 0)
ON CONFLICT DO NOTHING;

-- Exclusions for all plans
INSERT INTO exclusions (plan_id, description, type, icd_codes) VALUES
  ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Pre-existing conditions (first 12 months)', 'pre_existing',  '[]'),
  ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Cosmetic procedures',                       'cosmetic',      '[]'),
  ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Experimental treatments',                   'experimental',  '[]'),
  ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'Pre-existing conditions (first 12 months)', 'pre_existing',  '[]'),
  ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'Cosmetic procedures',                       'cosmetic',      '[]'),
  ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'Experimental treatments',                   'experimental',  '[]'),
  ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Pre-existing conditions (first 6 months)',  'pre_existing',  '[]'),
  ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Cosmetic procedures',                       'cosmetic',      '[]'),
  ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Experimental treatments',                   'experimental',  '[]')
ON CONFLICT DO NOTHING;

-- ============================================================
-- LAYER 4: PROVIDERS + Contracts + Rate Cards
-- ============================================================
INSERT INTO providers (id, name, type, license_number, status, county, address, phone, email, contact_person, tier) VALUES
  ('40000000-0000-0000-0000-000000000001', 'Nairobi General Hospital', 'hospital', 'LIC-NGH-001', 'ACTIVE', 'Nairobi', 'Kenyatta Avenue, Nairobi CBD',     '+254201234567', 'admin@nairobigeneral.co.ke', 'Dr. Otieno',      'TIER_1'),
  ('40000000-0000-0000-0000-000000000002', 'City Medical Centre',      'clinic',   'LIC-CMC-002', 'ACTIVE', 'Nairobi', 'Westlands Road, Westlands',        '+254202345678', 'info@citymedical.co.ke',    'Dr. Njeri',       'TIER_2'),
  ('40000000-0000-0000-0000-000000000003', 'Wellness Pharmacy',        'pharmacy', 'LIC-WPH-003', 'ACTIVE', 'Nairobi', 'Moi Avenue, Nairobi CBD',          '+254203456789', 'orders@wellnesspharm.co.ke','James Ochieng',   'TIER_3')
ON CONFLICT (license_number) DO NOTHING;

-- Contracts (active, current year)
INSERT INTO contracts (id, provider_id, start_date, end_date, terms, status) VALUES
  ('41000000-0000-0000-0000-000000000001', '40000000-0000-0000-0000-000000000001', '2026-01-01', '2026-12-31', 'Standard hospital contract — capitation model', 'ACTIVE'),
  ('41000000-0000-0000-0000-000000000002', '40000000-0000-0000-0000-000000000002', '2026-01-01', '2026-12-31', 'Standard clinic contract — fee for service',    'ACTIVE'),
  ('41000000-0000-0000-0000-000000000003', '40000000-0000-0000-0000-000000000003', '2026-01-01', '2026-12-31', 'Standard pharmacy contract — cost plus margin', 'ACTIVE')
ON CONFLICT DO NOTHING;

-- Rate Cards: Nairobi General Hospital
INSERT INTO rate_cards (provider_id, procedure_code, procedure_name, rate_amount, effective_date) VALUES
  ('40000000-0000-0000-0000-000000000001', 'CONS-001', 'General Consultation',  250000,  '2026-01-01'),
  ('40000000-0000-0000-0000-000000000001', 'LAB-001',  'Complete Blood Count',  150000,  '2026-01-01'),
  ('40000000-0000-0000-0000-000000000001', 'PHARM-001','Prescription Medication',80000,  '2026-01-01'),
  ('40000000-0000-0000-0000-000000000001', 'PROC-001', 'Minor Surgery',         1500000, '2026-01-01')
ON CONFLICT DO NOTHING;

-- Rate Cards: City Medical Centre
INSERT INTO rate_cards (provider_id, procedure_code, procedure_name, rate_amount, effective_date) VALUES
  ('40000000-0000-0000-0000-000000000002', 'CONS-001', 'General Consultation',  200000,  '2026-01-01'),
  ('40000000-0000-0000-0000-000000000002', 'LAB-001',  'Complete Blood Count',  120000,  '2026-01-01'),
  ('40000000-0000-0000-0000-000000000002', 'LAB-002',  'Urinalysis',             50000,  '2026-01-01'),
  ('40000000-0000-0000-0000-000000000002', 'PHARM-001','Prescription Medication',60000,  '2026-01-01')
ON CONFLICT DO NOTHING;

-- Rate Cards: Wellness Pharmacy
INSERT INTO rate_cards (provider_id, procedure_code, procedure_name, rate_amount, effective_date) VALUES
  ('40000000-0000-0000-0000-000000000003', 'PHARM-001','Prescription Medication',50000,  '2026-01-01'),
  ('40000000-0000-0000-0000-000000000003', 'PHARM-002','Over-the-counter Meds',  30000,  '2026-01-01'),
  ('40000000-0000-0000-0000-000000000003', 'PHARM-003','Medical Supplies',        80000, '2026-01-01')
ON CONFLICT DO NOTHING;

-- Provider Networks (link plans to providers)
INSERT INTO provider_networks (plan_id, provider_id, benefit_category, status) VALUES
  ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '40000000-0000-0000-0000-000000000001', 'inpatient',  'ACTIVE'),
  ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '40000000-0000-0000-0000-000000000001', 'emergency', 'ACTIVE'),
  ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '40000000-0000-0000-0000-000000000002', 'outpatient', 'ACTIVE'),
  ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '40000000-0000-0000-0000-000000000002', 'specialist', 'ACTIVE'),
  ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '40000000-0000-0000-0000-000000000002', 'chronic', 'ACTIVE'),
  ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '40000000-0000-0000-0000-000000000002', 'wellness', 'ACTIVE'),
  ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '40000000-0000-0000-0000-000000000003', 'outpatient', 'ACTIVE'),
  ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '40000000-0000-0000-0000-000000000003', 'pharmacy', 'ACTIVE'),
  ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', '40000000-0000-0000-0000-000000000001', 'inpatient',  'ACTIVE'),
  ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', '40000000-0000-0000-0000-000000000001', 'emergency', 'ACTIVE'),
  ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', '40000000-0000-0000-0000-000000000002', 'outpatient', 'ACTIVE'),
  ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', '40000000-0000-0000-0000-000000000002', 'specialist', 'ACTIVE'),
  ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', '40000000-0000-0000-0000-000000000003', 'outpatient', 'ACTIVE'),
  ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', '40000000-0000-0000-0000-000000000003', 'pharmacy', 'ACTIVE'),
  ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', '40000000-0000-0000-0000-000000000001', 'inpatient',  'ACTIVE'),
  ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', '40000000-0000-0000-0000-000000000001', 'emergency', 'ACTIVE'),
  ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', '40000000-0000-0000-0000-000000000002', 'outpatient', 'ACTIVE'),
  ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', '40000000-0000-0000-0000-000000000002', 'specialist', 'ACTIVE'),
  ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', '40000000-0000-0000-0000-000000000002', 'chronic', 'ACTIVE'),
  ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', '40000000-0000-0000-0000-000000000002', 'wellness', 'ACTIVE'),
  ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', '40000000-0000-0000-0000-000000000003', 'outpatient', 'ACTIVE'),
  ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', '40000000-0000-0000-0000-000000000003', 'pharmacy', 'ACTIVE')
ON CONFLICT (plan_id, provider_id, benefit_category) DO NOTHING;

-- ============================================================
-- LAYER 5: POLICIES + MEMBERS
-- ============================================================

-- Policy 1: Afya Basic — John Doe (individual)
INSERT INTO policies (id, plan_id, policyholder_name, policyholder_email, policyholder_phone, policy_number, status, start_date, end_date, premium_amount, currency, created_by) VALUES
  ('50000000-0000-0000-0000-000000000001', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'John Doe',       'john.doe@email.com',       '+254711000001', 'POL-2026-000001', 'ACTIVE', '2026-01-01', '2026-12-31', 250000,  'KES', '20000000-0000-0000-0000-000000000003')
ON CONFLICT (policy_number) DO NOTHING;

-- Policy 2: Afya Family — James Mwangi + family
INSERT INTO policies (id, plan_id, policyholder_name, policyholder_email, policyholder_phone, policy_number, status, start_date, end_date, premium_amount, currency, created_by) VALUES
  ('50000000-0000-0000-0000-000000000002', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'James Mwangi',   'james.mwangi@email.com',   '+254711000002', 'POL-2026-000002', 'ACTIVE', '2026-01-01', '2026-12-31', 750000,  'KES', '20000000-0000-0000-0000-000000000003')
ON CONFLICT (policy_number) DO NOTHING;

-- Policy 3: Corporate Gold — Acme Corp (Mary + David)
INSERT INTO policies (id, plan_id, policyholder_name, policyholder_email, policyholder_phone, policy_number, status, start_date, end_date, premium_amount, currency, created_by) VALUES
  ('50000000-0000-0000-0000-000000000003', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Acme Corporation','hr@acmecorp.co.ke',        '+254711000003', 'POL-2026-000003', 'ACTIVE', '2026-01-01', '2026-12-31', 1500000, 'KES', '20000000-0000-0000-0000-000000000003')
ON CONFLICT (policy_number) DO NOTHING;

-- Members
INSERT INTO members (id, policy_id, national_id, name, date_of_birth, gender, relationship, member_number, phone, email, status) VALUES
  -- Policy 1: John Doe (principal)
  ('60000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000001', '12345678', 'John Doe',        '1990-05-15', 'MALE',   'PRINCIPAL', 'MBR-2026-000001', '+254711000001', 'john.doe@email.com',       'ACTIVE'),
  -- Policy 2: James Mwangi (principal) + spouse + child
  ('60000000-0000-0000-0000-000000000002', '50000000-0000-0000-0000-000000000002', '23456789', 'James Mwangi',    '1985-03-20', 'MALE',   'PRINCIPAL', 'MBR-2026-000002', '+254711000002', 'james.mwangi@email.com',   'ACTIVE'),
  ('60000000-0000-0000-0000-000000000003', '50000000-0000-0000-0000-000000000002', '23456790', 'Grace Mwangi',    '1987-08-10', 'FEMALE', 'SPOUSE',    'MBR-2026-000003', '+254711000012', 'grace.mwangi@email.com',   'ACTIVE'),
  ('60000000-0000-0000-0000-000000000004', '50000000-0000-0000-0000-000000000002', NULL,       'Brian Mwangi',    '2015-11-25', 'MALE',   'CHILD',     'MBR-2026-000004', NULL,            NULL,                       'ACTIVE'),
  -- Policy 3: Mary Akinyi + David Ouma
  ('60000000-0000-0000-0000-000000000005', '50000000-0000-0000-0000-000000000003', '34567891', 'Mary Akinyi',     '1992-07-12', 'FEMALE', 'PRINCIPAL', 'MBR-2026-000005', '+254711000003', 'mary.akinyi@acmecorp.co.ke','ACTIVE'),
  ('60000000-0000-0000-0000-000000000006', '50000000-0000-0000-0000-000000000003', '34567892', 'David Ouma',      '1988-01-30', 'MALE',   'PRINCIPAL', 'MBR-2026-000006', '+254711000004', 'david.ouma@acmecorp.co.ke', 'ACTIVE')
ON CONFLICT (member_number) DO NOTHING;

-- ============================================================
-- LAYER 6: CLAIMS + Line Items + Adjudication Decisions
-- ============================================================

-- CLM-1: John Doe @ Nairobi General — ADJUDICATED (15,000 KES)
INSERT INTO claims (id, claim_number, policy_id, member_id, provider_id, status, total_amount, approved_amount, co_pay_amount, member_responsibility, diagnosis_codes, service_date, notes, claim_type, created_by) VALUES
  ('70000000-0000-0000-0000-000000000001', 'CLM-2026-000001', '50000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000001', '40000000-0000-0000-0000-000000000001', 'ADJUDICATED', 1500000, 1350000, 150000, 150000, '["J06.9"]', '2026-02-10', 'Upper respiratory infection treatment', 'DIRECT', '20000000-0000-0000-0000-000000000002')
ON CONFLICT (claim_number) DO NOTHING;

INSERT INTO claim_line_items (claim_id, procedure_code, procedure_name, diagnosis_code, quantity, unit_price, total_price, approved_amount) VALUES
  ('70000000-0000-0000-0000-000000000001', 'CONS-001', 'General Consultation',  'J06.9', 1, 250000,  250000,  250000),
  ('70000000-0000-0000-0000-000000000001', 'LAB-001',  'Complete Blood Count',  'J06.9', 1, 150000,  150000,  150000),
  ('70000000-0000-0000-0000-000000000001', 'PHARM-001','Prescription Medication','J06.9', 2, 550000, 1100000,  950000);

INSERT INTO adjudication_decisions (claim_id, decision, payable_amount, member_responsibility, reasons, rule_results, adjudicated_by, adjudicated_at) VALUES
  ('70000000-0000-0000-0000-000000000001', 'APPROVE', 1350000, 150000, '["Within benefit limits","Active policy"]', '[{"rule":"benefit_check","pass":true},{"rule":"policy_active","pass":true}]', '20000000-0000-0000-0000-000000000002', '2026-02-11');

-- CLM-2: James Mwangi @ City Medical — APPROVED (8,000 KES)
INSERT INTO claims (id, claim_number, policy_id, member_id, provider_id, status, total_amount, approved_amount, co_pay_amount, member_responsibility, diagnosis_codes, service_date, notes, claim_type, created_by) VALUES
  ('70000000-0000-0000-0000-000000000002', 'CLM-2026-000002', '50000000-0000-0000-0000-000000000002', '60000000-0000-0000-0000-000000000002', '40000000-0000-0000-0000-000000000002', 'APPROVED', 800000, 760000, 40000, 40000, '["K29.7"]', '2026-02-15', 'Gastritis consultation and medication', 'DIRECT', '20000000-0000-0000-0000-000000000002')
ON CONFLICT (claim_number) DO NOTHING;

INSERT INTO claim_line_items (claim_id, procedure_code, procedure_name, diagnosis_code, quantity, unit_price, total_price, approved_amount) VALUES
  ('70000000-0000-0000-0000-000000000002', 'CONS-001', 'General Consultation',  'K29.7', 1, 200000, 200000, 200000),
  ('70000000-0000-0000-0000-000000000002', 'PHARM-001','Prescription Medication','K29.7', 1, 600000, 600000, 560000);

INSERT INTO adjudication_decisions (claim_id, decision, payable_amount, member_responsibility, reasons, rule_results, adjudicated_by, adjudicated_at) VALUES
  ('70000000-0000-0000-0000-000000000002', 'APPROVE', 760000, 40000, '["Within benefit limits","Provider in network"]', '[{"rule":"benefit_check","pass":true},{"rule":"network_check","pass":true}]', '20000000-0000-0000-0000-000000000002', '2026-02-16');

-- CLM-3: Grace Mwangi (spouse) @ Wellness Pharmacy — PAID (3,500 KES)
INSERT INTO claims (id, claim_number, policy_id, member_id, provider_id, status, total_amount, approved_amount, co_pay_amount, member_responsibility, diagnosis_codes, service_date, notes, claim_type, created_by) VALUES
  ('70000000-0000-0000-0000-000000000003', 'CLM-2026-000003', '50000000-0000-0000-0000-000000000002', '60000000-0000-0000-0000-000000000003', '40000000-0000-0000-0000-000000000003', 'PAID', 350000, 350000, 0, 0, '["R51"]', '2026-02-20', 'Migraine medication refill', 'DIRECT', '20000000-0000-0000-0000-000000000002')
ON CONFLICT (claim_number) DO NOTHING;

INSERT INTO claim_line_items (claim_id, procedure_code, procedure_name, diagnosis_code, quantity, unit_price, total_price, approved_amount) VALUES
  ('70000000-0000-0000-0000-000000000003', 'PHARM-001','Prescription Medication','R51', 1, 350000, 350000, 350000);

INSERT INTO adjudication_decisions (claim_id, decision, payable_amount, member_responsibility, reasons, rule_results, adjudicated_by, adjudicated_at) VALUES
  ('70000000-0000-0000-0000-000000000003', 'APPROVE', 350000, 0, '["Below co-pay threshold","Pharmacy claim auto-approved"]', '[{"rule":"benefit_check","pass":true}]', '20000000-0000-0000-0000-000000000002', '2026-02-20');

-- CLM-4: Mary Akinyi @ Nairobi General — REJECTED (25,000 KES)
INSERT INTO claims (id, claim_number, policy_id, member_id, provider_id, status, total_amount, approved_amount, co_pay_amount, member_responsibility, diagnosis_codes, service_date, notes, rejection_reason, claim_type, created_by) VALUES
  ('70000000-0000-0000-0000-000000000004', 'CLM-2026-000004', '50000000-0000-0000-0000-000000000003', '60000000-0000-0000-0000-000000000005', '40000000-0000-0000-0000-000000000001', 'REJECTED', 2500000, 0, 0, 2500000, '["Z41.1"]', '2026-02-25', 'Cosmetic procedure — rhinoplasty consultation', 'Excluded benefit: cosmetic procedures are not covered under this plan', 'DIRECT', '20000000-0000-0000-0000-000000000002')
ON CONFLICT (claim_number) DO NOTHING;

INSERT INTO claim_line_items (claim_id, procedure_code, procedure_name, diagnosis_code, quantity, unit_price, total_price, approved_amount) VALUES
  ('70000000-0000-0000-0000-000000000004', 'CONS-001', 'Specialist Consultation','Z41.1', 1, 500000,  500000,  0),
  ('70000000-0000-0000-0000-000000000004', 'PROC-001', 'Procedure Assessment',   'Z41.1', 1, 2000000, 2000000, 0);

INSERT INTO adjudication_decisions (claim_id, decision, payable_amount, member_responsibility, reasons, rule_results, adjudicated_by, adjudicated_at) VALUES
  ('70000000-0000-0000-0000-000000000004', 'REJECT', 0, 2500000, '["Cosmetic procedure excluded"]', '[{"rule":"exclusion_check","pass":false,"reason":"Cosmetic exclusion applies"}]', '20000000-0000-0000-0000-000000000002', '2026-02-26');

-- CLM-5: John Doe @ City Medical — MANUAL_REVIEW (120,000 KES — high amount triggers review)
INSERT INTO claims (id, claim_number, policy_id, member_id, provider_id, status, total_amount, approved_amount, co_pay_amount, member_responsibility, diagnosis_codes, service_date, notes, claim_type, created_by) VALUES
  ('70000000-0000-0000-0000-000000000005', 'CLM-2026-000005', '50000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000001', '40000000-0000-0000-0000-000000000002', 'MANUAL_REVIEW', 12000000, 0, 0, 0, '["M54.5"]', '2026-03-01', 'Lower back pain — MRI and physiotherapy series', 'DIRECT', '20000000-0000-0000-0000-000000000002')
ON CONFLICT (claim_number) DO NOTHING;

INSERT INTO claim_line_items (claim_id, procedure_code, procedure_name, diagnosis_code, quantity, unit_price, total_price, approved_amount) VALUES
  ('70000000-0000-0000-0000-000000000005', 'PROC-001', 'MRI Lumbar Spine',       'M54.5', 1, 8000000, 8000000, 0),
  ('70000000-0000-0000-0000-000000000005', 'CONS-001', 'Physiotherapy Sessions',  'M54.5', 8, 500000,  4000000, 0);

INSERT INTO adjudication_decisions (claim_id, decision, payable_amount, member_responsibility, reasons, rule_results, adjudicated_by, adjudicated_at) VALUES
  ('70000000-0000-0000-0000-000000000005', 'MANUAL_REVIEW', 0, 0, '["Amount exceeds auto-approval threshold of KES 100,000"]', '[{"rule":"amount_threshold","pass":false,"reason":"Claim amount 120,000 KES exceeds 100,000 KES limit"}]', NULL, '2026-03-01');

-- CLM-6: David Ouma @ Nairobi General — RECEIVED (5,000 KES — just submitted)
INSERT INTO claims (id, claim_number, policy_id, member_id, provider_id, status, total_amount, approved_amount, co_pay_amount, member_responsibility, diagnosis_codes, service_date, notes, claim_type, created_by) VALUES
  ('70000000-0000-0000-0000-000000000006', 'CLM-2026-000006', '50000000-0000-0000-0000-000000000003', '60000000-0000-0000-0000-000000000006', '40000000-0000-0000-0000-000000000001', 'RECEIVED', 500000, 0, 0, 0, '["J03.9"]', '2026-03-05', 'Acute tonsillitis — consultation and antibiotics', 'DIRECT', '20000000-0000-0000-0000-000000000002')
ON CONFLICT (claim_number) DO NOTHING;

INSERT INTO claim_line_items (claim_id, procedure_code, procedure_name, diagnosis_code, quantity, unit_price, total_price, approved_amount) VALUES
  ('70000000-0000-0000-0000-000000000006', 'CONS-001', 'General Consultation',  'J03.9', 1, 250000, 250000, 0),
  ('70000000-0000-0000-0000-000000000006', 'PHARM-001','Antibiotics',           'J03.9', 1, 250000, 250000, 0);

-- ============================================================
-- LAYER 7: SALES PIPELINE (Leads + Quotations)
-- ============================================================

-- Lead 1: NEW
INSERT INTO leads (id, lead_number, contact_name, contact_email, contact_phone, company_name, source, segment, plan_type, estimated_members, expected_premium, closure_probability, status, assigned_to, next_follow_up_date, notes, created_by) VALUES
  ('80000000-0000-0000-0000-000000000001', 'LEAD-000001', 'Michael Njoroge', 'michael@techstartup.co.ke', '+254722111222', 'Tech Startup Ltd', 'referral', 'corporate', 'group', 15, 22500000, 30, 'NEW', '20000000-0000-0000-0000-000000000003', '2026-03-10', 'Referred by existing client. Interested in Corporate Gold plan for 15 employees.', '20000000-0000-0000-0000-000000000003')
ON CONFLICT (lead_number) DO NOTHING;

-- Lead 2: QUALIFIED
INSERT INTO leads (id, lead_number, contact_name, contact_email, contact_phone, company_name, source, segment, plan_type, estimated_members, expected_premium, closure_probability, status, assigned_to, next_follow_up_date, notes, created_by) VALUES
  ('80000000-0000-0000-0000-000000000002', 'LEAD-000002', 'Fatima Hassan', 'fatima@email.com', '+254733222333', NULL, 'website', 'retail', 'individual', 1, 250000, 60, 'QUALIFIED', '20000000-0000-0000-0000-000000000003', '2026-03-08', 'Young professional looking for individual health cover. Prefers outpatient-heavy plan.', '20000000-0000-0000-0000-000000000003')
ON CONFLICT (lead_number) DO NOTHING;

-- Quotation 1: DRAFT (for Lead 1)
INSERT INTO quotations (id, quotation_number, lead_id, plan_id, quotation_type, status, current_version, valid_from, valid_until, client_name, client_email, client_phone, created_by) VALUES
  ('81000000-0000-0000-0000-000000000001', 'QUO-000001', '80000000-0000-0000-0000-000000000001', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'standard', 'DRAFT', 1, '2026-03-06', '2026-04-06', 'Michael Njoroge', 'michael@techstartup.co.ke', '+254722111222', '20000000-0000-0000-0000-000000000003')
ON CONFLICT (quotation_number) DO NOTHING;

INSERT INTO quotation_versions (quotation_id, version_number, base_premium, discount_type, discount_value, discount_reason, loading_type, loading_value, final_premium, member_count, billing_frequency, pricing_breakdown, created_by) VALUES
  ('81000000-0000-0000-0000-000000000001', 1, 22500000, 'percentage', 500, 'Volume discount for 15+ members', 'percentage', 0, 21375000, 15, 'monthly', '{"per_member": 1500000, "discount_pct": 5, "monthly_total": 21375000}', '20000000-0000-0000-0000-000000000003')
ON CONFLICT (quotation_id, version_number) DO NOTHING;

-- Quotation 2: SENT (for Lead 2)
INSERT INTO quotations (id, quotation_number, lead_id, plan_id, quotation_type, status, current_version, valid_from, valid_until, client_name, client_email, client_phone, created_by) VALUES
  ('81000000-0000-0000-0000-000000000002', 'QUO-000002', '80000000-0000-0000-0000-000000000002', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'standard', 'SENT', 1, '2026-03-01', '2026-04-01', 'Fatima Hassan', 'fatima@email.com', '+254733222333', '20000000-0000-0000-0000-000000000003')
ON CONFLICT (quotation_number) DO NOTHING;

INSERT INTO quotation_versions (quotation_id, version_number, base_premium, discount_type, discount_value, loading_type, loading_value, final_premium, member_count, billing_frequency, pricing_breakdown, created_by) VALUES
  ('81000000-0000-0000-0000-000000000002', 1, 250000, 'percentage', 0, 'percentage', 0, 250000, 1, 'monthly', '{"per_member": 250000, "monthly_total": 250000}', '20000000-0000-0000-0000-000000000003')
ON CONFLICT (quotation_id, version_number) DO NOTHING;

-- ============================================================
-- LAYER 8: BILLING (Invoices + Payments)
-- ============================================================

-- Invoice 1: POL-1, PAID
INSERT INTO invoices (id, policy_id, invoice_number, amount, currency, due_date, status, billing_period_start, billing_period_end, notes, created_by) VALUES
  ('90000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000001', 'INV-2026-000001', 250000, 'KES', '2026-01-31', 'PAID', '2026-01-01', '2026-01-31', 'January 2026 premium', '20000000-0000-0000-0000-000000000004')
ON CONFLICT (invoice_number) DO NOTHING;

-- Invoice 2: POL-2, PENDING
INSERT INTO invoices (id, policy_id, invoice_number, amount, currency, due_date, status, billing_period_start, billing_period_end, notes, created_by) VALUES
  ('90000000-0000-0000-0000-000000000002', '50000000-0000-0000-0000-000000000002', 'INV-2026-000002', 750000, 'KES', '2026-03-15', 'PENDING', '2026-03-01', '2026-03-31', 'March 2026 premium', '20000000-0000-0000-0000-000000000004')
ON CONFLICT (invoice_number) DO NOTHING;

-- Invoice 3: POL-3, OVERDUE
INSERT INTO invoices (id, policy_id, invoice_number, amount, currency, due_date, status, billing_period_start, billing_period_end, notes, created_by) VALUES
  ('90000000-0000-0000-0000-000000000003', '50000000-0000-0000-0000-000000000003', 'INV-2026-000003', 1500000, 'KES', '2026-02-28', 'OVERDUE', '2026-02-01', '2026-02-28', 'February 2026 premium — overdue', '20000000-0000-0000-0000-000000000004')
ON CONFLICT (invoice_number) DO NOTHING;

-- Payment 1: Full payment for Invoice 1
INSERT INTO payments (id, invoice_id, type, amount, currency, method, reference_number, status, paid_at, created_by) VALUES
  ('91000000-0000-0000-0000-000000000001', '90000000-0000-0000-0000-000000000001', 'PREMIUM', 250000, 'KES', 'MPESA', 'MPESA-REF-20260125-001', 'CONFIRMED', '2026-01-25', '20000000-0000-0000-0000-000000000004')
ON CONFLICT (reference_number) DO NOTHING;

-- Payment 2: Partial payment for Invoice 3 (overdue)
INSERT INTO payments (id, invoice_id, type, amount, currency, method, reference_number, status, paid_at, created_by) VALUES
  ('91000000-0000-0000-0000-000000000002', '90000000-0000-0000-0000-000000000003', 'PREMIUM', 500000, 'KES', 'BANK_TRANSFER', 'BNK-REF-20260301-001', 'CONFIRMED', '2026-03-01', '20000000-0000-0000-0000-000000000004')
ON CONFLICT (reference_number) DO NOTHING;

-- Payment for CLM-3 (paid claim — remittance type for provider payment)
INSERT INTO payments (id, claim_id, type, amount, currency, method, reference_number, status, paid_at, created_by) VALUES
  ('91000000-0000-0000-0000-000000000003', '70000000-0000-0000-0000-000000000003', 'REMITTANCE', 350000, 'KES', 'BANK_TRANSFER', 'BNK-CLM-20260222-001', 'CONFIRMED', '2026-02-22', '20000000-0000-0000-0000-000000000004')
ON CONFLICT (reference_number) DO NOTHING;

-- ============================================================
-- LAYER 9: PRE-AUTHORIZATION
-- ============================================================
INSERT INTO preauthorizations (id, policy_id, member_id, provider_id, auth_code, procedure_codes, diagnosis_codes, estimated_cost, approved_amount, status, validity_start, validity_end, notes, reviewed_by, reviewed_at, created_by) VALUES
  ('A0000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000002', '60000000-0000-0000-0000-000000000002', '40000000-0000-0000-0000-000000000001', 'PA-2026-000001', '["PROC-001"]', '["K35.8"]', 5000000, 5000000, 'APPROVED', '2026-03-01', '2026-03-31', 'Appendectomy — pre-authorized for inpatient surgery', '20000000-0000-0000-0000-000000000002', '2026-02-28', '20000000-0000-0000-0000-000000000002')
ON CONFLICT (auth_code) DO NOTHING;

-- ============================================================
-- LAYER 10: NOTIFICATIONS
-- ============================================================
INSERT INTO notifications (user_id, channel, type, subject, body, metadata, status, sent_at) VALUES
  ('20000000-0000-0000-0000-000000000001', 'IN_APP', 'CLAIM_APPROVED',    'Claim CLM-2026-000002 Approved',     'Claim for James Mwangi has been approved for KES 7,600. Ready for payment processing.', '{"claim_id": "70000000-0000-0000-0000-000000000002"}', 'DELIVERED', '2026-02-16'),
  ('20000000-0000-0000-0000-000000000001', 'IN_APP', 'PAYMENT_RECEIVED',  'Payment Received — INV-2026-000001', 'Payment of KES 2,500 received via M-Pesa for policy POL-2026-000001.',                    '{"invoice_id": "90000000-0000-0000-0000-000000000001"}', 'DELIVERED', '2026-01-25'),
  ('20000000-0000-0000-0000-000000000001', 'IN_APP', 'SLA_WARNING',       'SLA Breach Warning — CLM-2026-000005','Claim CLM-2026-000005 has been in MANUAL_REVIEW for 48+ hours. SLA breach imminent.',    '{"claim_id": "70000000-0000-0000-0000-000000000005"}', 'SENT', NOW());

-- ============================================================
-- LAYER 11: SUPPORTING DATA
-- ============================================================

-- Endorsement: member addition to POL-2
INSERT INTO endorsements (policy_id, endorsement_type, status, effective_date, changes, reason, premium_adjustment, requested_by, approved_by, approved_at, applied_at) VALUES
  ('50000000-0000-0000-0000-000000000002', 'MEMBER_ADDITION', 'APPROVED', '2026-02-01', '{"added_member": "Brian Mwangi", "relationship": "CHILD"}', 'Adding dependent child to family policy', 0, '20000000-0000-0000-0000-000000000003', '20000000-0000-0000-0000-000000000001', '2026-01-28', '2026-02-01');

-- Claim status history (one entry per claim showing latest transition)
INSERT INTO claim_status_history (claim_id, from_status, to_status, action, notes, performed_by) VALUES
  ('70000000-0000-0000-0000-000000000001', 'RECEIVED',      'ADJUDICATED',   'auto_adjudicate', 'Auto-adjudication completed',                          '20000000-0000-0000-0000-000000000002'),
  ('70000000-0000-0000-0000-000000000002', 'ADJUDICATED',   'APPROVED',      'approve',         'Approved by claims officer',                            '20000000-0000-0000-0000-000000000002'),
  ('70000000-0000-0000-0000-000000000003', 'APPROVED',      'PAID',          'process_payment', 'Payment processed via bank transfer',                   '20000000-0000-0000-0000-000000000004'),
  ('70000000-0000-0000-0000-000000000004', 'RECEIVED',      'REJECTED',      'reject',          'Excluded benefit: cosmetic procedures',                 '20000000-0000-0000-0000-000000000002'),
  ('70000000-0000-0000-0000-000000000005', 'RECEIVED',      'MANUAL_REVIEW', 'escalate',        'Amount exceeds auto-approval threshold',                '20000000-0000-0000-0000-000000000002'),
  ('70000000-0000-0000-0000-000000000006', '',              'RECEIVED',      'submit',          'Claim submitted by provider',                           '20000000-0000-0000-0000-000000000002');

-- ============================================================
-- LAYER 12: PREMIUM MODEL ENHANCEMENTS
-- ============================================================

-- Set premium frequency on plans
UPDATE plans SET premium_frequency = 'annual';

-- Classify existing premium rules
UPDATE premium_rules SET rule_type = 'age_band', sort_order = 10 WHERE min_age > 0 OR max_age < 150;
UPDATE premium_rules SET rule_type = 'base_rate', sort_order = 0 WHERE rule_type = 'base_rate' AND min_age = 0;

-- Mark maternity as optional add-on
UPDATE benefits SET is_optional = true, addon_premium = 1500000 WHERE category = 'maternity';

-- Set members of non-active policies to PENDING
UPDATE members SET status = 'PENDING'
WHERE policy_id IN (SELECT id FROM policies WHERE status = 'DRAFT');

SELECT 'Seed data inserted successfully' as result;

COMMIT;
EOF

echo "Seeding complete."
