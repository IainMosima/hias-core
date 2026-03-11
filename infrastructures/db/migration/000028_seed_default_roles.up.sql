-- Seed default roles so registration works out of the box
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

-- Grant Admin all permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT '10000000-0000-0000-0000-000000000001', id FROM permissions
ON CONFLICT (role_id, permission_id) DO NOTHING;