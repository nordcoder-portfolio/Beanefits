-- name: CreateAccountForUser :one
INSERT INTO accounts (user_id, public_code, level_code)
VALUES ($1, $2, $3)
RETURNING
    id, user_id, public_code, created_at,
    balance_points, total_spend_money,
    COALESCE(level_code, '')::text AS level_code;

-- name: GetAccountByUserID :one
SELECT
    id, user_id, public_code, created_at,
    balance_points, total_spend_money,
    COALESCE(level_code, '')::text AS level_code
FROM accounts
WHERE user_id = $1;

-- name: GetAccountByPublicCode :one
SELECT
    id, user_id, public_code, created_at,
    balance_points, total_spend_money,
    COALESCE(level_code, '')::text AS level_code
FROM accounts
WHERE public_code = $1;

-- name: LockAccountByID :one
SELECT
    id, user_id, public_code, created_at,
    balance_points, total_spend_money,
    COALESCE(level_code, '')::text AS level_code
FROM accounts
WHERE id = $1
    FOR UPDATE;

-- name: UpdateAccountAfterEarn :one
UPDATE accounts
SET
    balance_points = $2,
    total_spend_money = $3,
    level_code = $4
WHERE id = $1
RETURNING
    id, user_id, public_code, created_at,
    balance_points, total_spend_money,
    COALESCE(level_code, '')::text AS level_code;

-- name: UpdateAccountAfterSpend :one
UPDATE accounts
SET balance_points = $2
WHERE id = $1
RETURNING
    id, user_id, public_code, created_at,
    balance_points, total_spend_money,
    COALESCE(level_code, '')::text AS level_code;
