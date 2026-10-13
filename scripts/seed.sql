CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Reset database table
TRUNCATE TABLE menus CASCADE;

-- 1. Root Level (Depth 0)
INSERT INTO menus (id, name, parent_id, depth, sort_order) VALUES
  ('550e8400-e29b-41d4-a716-446655440001', 'System Management', NULL, 0, 0);

-- 2. Level 1 (Under System Management)
INSERT INTO menus (id, name, parent_id, depth, sort_order) VALUES
  ('550e8400-e29b-41d4-a716-446655440010', 'Systems', '550e8400-e29b-41d4-a716-446655440001', 1, 0),
  ('550e8400-e29b-41d4-a716-446655440020', 'Users & Groups', '550e8400-e29b-41d4-a716-446655440001', 1, 1);

-- 3. Level 2 (Under Systems)
INSERT INTO menus (id, name, parent_id, depth, sort_order) VALUES
  ('550e8400-e29b-41d4-a716-446655440011', 'System Code', '550e8400-e29b-41d4-a716-446655440010', 2, 0),
  ('550e8400-e29b-41d4-a716-446655440012', 'Properties', '550e8400-e29b-41d4-a716-446655440010', 2, 1),
  ('550e8400-e29b-41d4-a716-446655440013', 'Menus', '550e8400-e29b-41d4-a716-446655440010', 2, 2),
  ('550e8400-e29b-41d4-a716-446655440014', 'API List', '550e8400-e29b-41d4-a716-446655440010', 2, 3);

-- 4. Level 3 (Under System Code, Menus, API List)
INSERT INTO menus (id, name, parent_id, depth, sort_order) VALUES
  ('550e8400-e29b-41d4-a716-446655440031', 'Code Registration', '550e8400-e29b-41d4-a716-446655440011', 3, 0),
  ('550e8400-e29b-41d4-a716-446655440032', 'Code Registration - 2', '550e8400-e29b-41d4-a716-446655440011', 3, 1),
  ('550e8400-e29b-41d4-a716-446655440033', 'Menu Registration', '550e8400-e29b-41d4-a716-446655440013', 3, 0),
  ('550e8400-e29b-41d4-a716-446655440034', 'API Registration', '550e8400-e29b-41d4-a716-446655440014', 3, 0),
  ('550e8400-e29b-41d4-a716-446655440035', 'API Edit', '550e8400-e29b-41d4-a716-446655440014', 3, 1);

-- 5. Level 2 & 3 (Under Users & Groups)
INSERT INTO menus (id, name, parent_id, depth, sort_order) VALUES
  ('550e8400-e29b-41d4-a716-446655440021', 'Users', '550e8400-e29b-41d4-a716-446655440020', 2, 0),
  ('550e8400-e29b-41d4-a716-446655440022', 'Groups', '550e8400-e29b-41d4-a716-446655440020', 2, 1),
  ('550e8400-e29b-41d4-a716-446655440023', 'User Approval', '550e8400-e29b-41d4-a716-446655440020', 2, 2);

INSERT INTO menus (id, name, parent_id, depth, sort_order) VALUES
  ('550e8400-e29b-41d4-a716-446655440041', 'User Account Registration', '550e8400-e29b-41d4-a716-446655440021', 3, 0),
  ('550e8400-e29b-41d4-a716-446655440042', 'User Group Registration', '550e8400-e29b-41d4-a716-446655440022', 3, 0),
  ('550e8400-e29b-41d4-a716-446655440043', 'User Approval Detail', '550e8400-e29b-41d4-a716-446655440023', 3, 0);