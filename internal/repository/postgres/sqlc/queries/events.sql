-- internal/repository/postgres/sqlc/queries/events.sql

-- name: InsertEvent :one
INSERT INTO events (
    account_id, type, delta_points, balance_after,
    amount_money, ruleset_id, actor_user_id, ts
)
VALUES (
           $1,
           $2::event_type,
           $3,
           $4,
           $5,
           $6,
           $7,
           $8
       )
RETURNING
    id,
    account_id,
    type::text AS type,
    delta_points,
    balance_after,
    amount_money,
    ruleset_id,
    actor_user_id,
    ts,
    created_at;

-- name: ListEventsByAccount :many
SELECT
    id,
    account_id,
    type::text AS type,
    delta_points,
    balance_after,
    amount_money,
    ruleset_id,
    actor_user_id,
    ts,
    created_at
FROM events
WHERE account_id = $1
ORDER BY ts DESC, id DESC
LIMIT $2;

-- name: ListEventsByAccountBefore :many
SELECT
    id,
    account_id,
    type::text AS type,
    delta_points,
    balance_after,
    amount_money,
    ruleset_id,
    actor_user_id,
    ts,
    created_at
FROM events
WHERE account_id = $1
  AND ts < $3
ORDER BY ts DESC, id DESC
LIMIT $2;
