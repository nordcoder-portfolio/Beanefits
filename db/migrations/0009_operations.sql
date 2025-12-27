-- +goose Up
CREATE TABLE operations
(
    id            BIGSERIAL PRIMARY KEY,
    operation_id  TEXT        NOT NULL,
    account_id    BIGINT      NOT NULL,
    op_type       event_type  NOT NULL,
    request_json  JSONB       NOT NULL,
    response_json JSONB,
    http_status   INT,
    event_id      BIGINT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_operations_account FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE RESTRICT,
    CONSTRAINT fk_operations_event FOREIGN KEY (event_id) REFERENCES events (id) ON DELETE SET NULL,

    CONSTRAINT uq_operations_idempotency UNIQUE (account_id, op_type, operation_id)
);

CREATE INDEX idx_operations_account_created_at ON operations (account_id, created_at);

-- +goose Down
DROP INDEX idx_operations_account_created_at;
DROP TABLE operations;
