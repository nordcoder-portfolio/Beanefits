-- +goose Up
-- Seed baseline roles + test users for local/dev.
-- Password for all seeded users: "password"

-- 1) Roles dictionary
INSERT INTO roles (code) VALUES
                             ('ADMIN'),
                             ('CASHIER'),
                             ('CLIENT')
ON CONFLICT (code) DO NOTHING;

-- 2) Users
-- NOTE: bcrypt hash below corresponds to password "password"
-- Generated once with bcrypt default cost.
INSERT INTO users (phone, password_hash, is_active)
VALUES
    ('+79990000001', '$2a$10$RiPgB5moANhOIqbI5eSZkOREyEnR.ktJmq3dABX1hm8yNGvBOZwJ.', true), -- admin
    ('+79990000002', '$2a$10$RiPgB5moANhOIqbI5eSZkOREyEnR.ktJmq3dABX1hm8yNGvBOZwJ.', true), -- cashier
    ('+79990000003', '$2a$10$RiPgB5moANhOIqbI5eSZkOREyEnR.ktJmq3dABX1hm8yNGvBOZwJ.', true)  -- client
ON CONFLICT (phone) DO NOTHING;

-- 3) User roles
INSERT INTO user_roles (user_id, role_code)
SELECT u.id, 'ADMIN'
FROM users u
WHERE u.phone = '+79990000001'
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_code)
SELECT u.id, 'CASHIER'
FROM users u
WHERE u.phone = '+79990000002'
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_code)
SELECT u.id, 'CLIENT'
FROM users u
WHERE u.phone = '+79990000003'
ON CONFLICT DO NOTHING;

-- 4) Accounts for ALL seeded users (so login/me/etc works uniformly)
-- Stable UUIDs to simplify Postman.
-- IMPORTANT: accounts.user_id is UNIQUE, so ON CONFLICT (user_id) works.
INSERT INTO accounts (user_id, public_code, level_code)
SELECT u.id, '550e8400-e29b-41d4-a716-446655440001', 'Green Bean'
FROM users u
WHERE u.phone = '+79990000001'
ON CONFLICT (user_id) DO NOTHING;

INSERT INTO accounts (user_id, public_code, level_code)
SELECT u.id, '550e8400-e29b-41d4-a716-446655440002', 'Green Bean'
FROM users u
WHERE u.phone = '+79990000002'
ON CONFLICT (user_id) DO NOTHING;

INSERT INTO accounts (user_id, public_code, level_code)
SELECT u.id, '550e8400-e29b-41d4-a716-446655440000', 'Green Bean'
FROM users u
WHERE u.phone = '+79990000003'
ON CONFLICT (user_id) DO NOTHING;

-- +goose Down
-- Remove seeded data (best-effort).
-- Order matters because of FKs.

-- Delete operations/events by account public_code (stable anchor)
DELETE FROM operations
WHERE account_id IN (
    SELECT id FROM accounts
    WHERE public_code IN (
                          '550e8400-e29b-41d4-a716-446655440001',
                          '550e8400-e29b-41d4-a716-446655440002',
                          '550e8400-e29b-41d4-a716-446655440000'
        )
);

DELETE FROM events
WHERE account_id IN (
    SELECT id FROM accounts
    WHERE public_code IN (
                          '550e8400-e29b-41d4-a716-446655440001',
                          '550e8400-e29b-41d4-a716-446655440002',
                          '550e8400-e29b-41d4-a716-446655440000'
        )
);

DELETE FROM user_roles
WHERE user_id IN (
    SELECT id FROM users WHERE phone IN ('+79990000001', '+79990000002', '+79990000003')
);

DELETE FROM accounts
WHERE public_code IN (
                      '550e8400-e29b-41d4-a716-446655440001',
                      '550e8400-e29b-41d4-a716-446655440002',
                      '550e8400-e29b-41d4-a716-446655440000'
    );

DELETE FROM users
WHERE phone IN ('+79990000001', '+79990000002', '+79990000003');

-- Roles can be shared by other data; usually do not delete them.
