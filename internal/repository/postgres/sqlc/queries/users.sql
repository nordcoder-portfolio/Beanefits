-- name: CreateUser :one
INSERT INTO users (phone, password_hash)
VALUES ($1, $2)
RETURNING id;

-- name: GetUserByPhone :one
SELECT id, phone, password_hash, is_active, created_at
FROM users
WHERE phone = $1;

-- name: GetUserByID :one
SELECT id, phone, password_hash, is_active, created_at
FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT id, phone, password_hash, is_active, created_at
FROM users
WHERE ($1::text = '' OR phone ILIKE '%' || $1 || '%')
ORDER BY id DESC
LIMIT $2 OFFSET $3;

-- name: DeactivateUser :exec
UPDATE users
SET is_active = false
WHERE id = $1;
