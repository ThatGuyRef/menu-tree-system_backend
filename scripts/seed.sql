CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Inisialisasi tabel jika belum ada
CREATE TABLE IF NOT EXISTS menus (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    parent_id UUID REFERENCES menus(id) ON DELETE CASCADE,
    depth INT NOT NULL DEFAULT 0,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Seed Data: Root Menus
INSERT INTO menus (id, name, parent_id, depth, sort_order) VALUES
  ('550e8400-e29b-41d4-a716-446655440001', 'System Management', NULL, 0, 0),
  ('550e8400-e29b-41d4-a716-446655440002', 'User Management', NULL, 0, 1)
ON CONFLICT (id) DO NOTHING;

-- Seed Data: Sub Menus level 1
INSERT INTO menus (id, name, parent_id, depth, sort_order) VALUES
  ('550e8400-e29b-41d4-a716-446655440003', 'System Code', '550e8400-e29b-41d4-a716-446655440001', 1, 0),
  ('550e8400-e29b-41d4-a716-446655440004', 'Properties', '550e8400-e29b-41d4-a716-446655440001', 1, 1),
  ('550e8400-e29b-41d4-a716-446655440005', 'Menus', '550e8400-e29b-41d4-a716-446655440001', 1, 2),
  ('550e8400-e29b-41d4-a716-446655440006', 'API List', '550e8400-e29b-41d4-a716-446655440001', 1, 3)
ON CONFLICT (id) DO NOTHING;

-- Seed Data: Sub Menus level 2
INSERT INTO menus (id, name, parent_id, depth, sort_order) VALUES
  ('550e8400-e29b-41d4-a716-446655440007', 'Code Management', '550e8400-e29b-41d4-a716-446655440003', 2, 0),
  ('550e8400-e29b-41d4-a716-446655440008', 'Users', '550e8400-e29b-41d4-a716-446655440002', 1, 0),
  ('550e8400-e29b-41d4-a716-446655440009', 'Roles & Permissions', '550e8400-e29b-41d4-a716-446655440002', 1, 1)
ON CONFLICT (id) DO NOTHING;