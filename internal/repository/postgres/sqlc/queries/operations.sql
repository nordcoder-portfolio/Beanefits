-- name: GetOperation :one
SELECT
    account_id,
    op_type::text AS op_type,
    operation_id,
    request_json,
    response_json,
    http_status,
    created_at
FROM operations
WHERE account_id = $1
  AND op_type = $2::event_type
  AND operation_id = $3;

-- name: InsertOperationPending :execrows
INSERT INTO operations (account_id, op_type, operation_id, request_json)
VALUES ($1, $2::event_type, $3, $4)
ON CONFLICT (account_id, op_type, operation_id) DO NOTHING;

-- name: FinalizeOperation :exec
UPDATE operations
SET http_status = $4,
    response_json = $5
WHERE account_id = $1
  AND op_type = $2::event_type
  AND operation_id = $3;
