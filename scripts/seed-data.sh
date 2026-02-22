#!/bin/bash
set -e
DB_URL="${DB_URL:-postgresql://root:supersecret@localhost:5432/hias_db?sslmode=disable}"

echo "Seeding default data..."

psql "$DB_URL" << 'EOF'
-- Seed default roles
INSERT INTO roles (id, name, description) VALUES
  (uuid_generate_v4(), 'Admin', 'Full system access'),
  (uuid_generate_v4(), 'Underwriter', 'Policy management and member enrollment'),
  (uuid_generate_v4(), 'ClaimsOfficer', 'Claims review and adjudication'),
  (uuid_generate_v4(), 'Finance', 'Financial operations and reporting'),
  (uuid_generate_v4(), 'Provider', 'Healthcare provider portal access'),
  (uuid_generate_v4(), 'Member', 'Member self-service portal access')
ON CONFLICT (name) DO NOTHING;

-- Seed default permissions
INSERT INTO permissions (id, resource, action, description) VALUES
  (uuid_generate_v4(), 'users', 'create', 'Create users'),
  (uuid_generate_v4(), 'users', 'read', 'View users'),
  (uuid_generate_v4(), 'users', 'update', 'Update users'),
  (uuid_generate_v4(), 'users', 'delete', 'Delete users'),
  (uuid_generate_v4(), 'plans', 'create', 'Create insurance plans'),
  (uuid_generate_v4(), 'plans', 'read', 'View insurance plans'),
  (uuid_generate_v4(), 'plans', 'update', 'Update insurance plans'),
  (uuid_generate_v4(), 'policies', 'create', 'Create policies'),
  (uuid_generate_v4(), 'policies', 'read', 'View policies'),
  (uuid_generate_v4(), 'policies', 'update', 'Update policies'),
  (uuid_generate_v4(), 'policies', 'activate', 'Activate policies'),
  (uuid_generate_v4(), 'providers', 'create', 'Register providers'),
  (uuid_generate_v4(), 'providers', 'read', 'View providers'),
  (uuid_generate_v4(), 'providers', 'update', 'Update providers'),
  (uuid_generate_v4(), 'providers', 'credential', 'Credential providers'),
  (uuid_generate_v4(), 'claims', 'create', 'Submit claims'),
  (uuid_generate_v4(), 'claims', 'read', 'View claims'),
  (uuid_generate_v4(), 'claims', 'review', 'Review claims'),
  (uuid_generate_v4(), 'claims', 'approve', 'Approve claims'),
  (uuid_generate_v4(), 'claims', 'reject', 'Reject claims'),
  (uuid_generate_v4(), 'preauth', 'create', 'Submit pre-authorizations'),
  (uuid_generate_v4(), 'preauth', 'read', 'View pre-authorizations'),
  (uuid_generate_v4(), 'preauth', 'approve', 'Approve pre-authorizations'),
  (uuid_generate_v4(), 'billing', 'read', 'View billing information'),
  (uuid_generate_v4(), 'billing', 'create', 'Create invoices'),
  (uuid_generate_v4(), 'payments', 'read', 'View payments'),
  (uuid_generate_v4(), 'payments', 'process', 'Process payments'),
  (uuid_generate_v4(), 'analytics', 'read', 'View analytics'),
  (uuid_generate_v4(), 'audit', 'read', 'View audit trail'),
  (uuid_generate_v4(), 'notifications', 'read', 'View notifications')
ON CONFLICT DO NOTHING;

-- Seed a sample plan with benefits
INSERT INTO plans (id, name, type, base_premium, currency, status, created_at, updated_at) VALUES
  ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Afya Basic', 'individual', 250000, 'KES', 'ACTIVE', NOW(), NOW()),
  ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'Afya Family', 'group', 750000, 'KES', 'ACTIVE', NOW(), NOW()),
  ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Corporate Gold', 'group', 1500000, 'KES', 'ACTIVE', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Benefits for Afya Basic
INSERT INTO benefits (id, plan_id, name, category, annual_limit, co_pay_type, co_pay_value, waiting_period_days, created_at, updated_at) VALUES
  (uuid_generate_v4(), 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Outpatient Care', 'outpatient', 10000000, 'percentage', 1000, 0, NOW(), NOW()),
  (uuid_generate_v4(), 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Inpatient Care', 'inpatient', 50000000, 'percentage', 1000, 30, NOW(), NOW()),
  (uuid_generate_v4(), 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Dental Care', 'dental', 5000000, 'fixed', 50000, 90, NOW(), NOW()),
  (uuid_generate_v4(), 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Optical Care', 'optical', 3000000, 'fixed', 30000, 90, NOW(), NOW()),
  (uuid_generate_v4(), 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Maternity', 'maternity', 30000000, 'percentage', 500, 270, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Exclusions for Afya Basic
INSERT INTO exclusions (id, plan_id, description, type, icd_codes, created_at, updated_at) VALUES
  (uuid_generate_v4(), 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Pre-existing conditions (first 12 months)', 'pre_existing', '[]', NOW(), NOW()),
  (uuid_generate_v4(), 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Cosmetic procedures', 'cosmetic', '[]', NOW(), NOW()),
  (uuid_generate_v4(), 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Experimental treatments', 'experimental', '[]', NOW(), NOW())
ON CONFLICT DO NOTHING;

SELECT 'Seed data inserted successfully' as result;
EOF

echo "Seeding complete."
