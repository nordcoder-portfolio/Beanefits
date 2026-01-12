-- +goose Up
CREATE TABLE events
(
    id            BIGSERIAL PRIMARY KEY,
    account_id    BIGINT      NOT NULL,
    type          event_type  NOT NULL,
    delta_points  INT         NOT NULL,
    balance_after INT         NOT NULL,
    amount_money  NUMERIC(12, 2),
    ts            TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    ruleset_id    BIGINT,
    actor_user_id BIGINT,

    CONSTRAINT fk_events_account FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE RESTRICT,
    CONSTRAINT fk_events_ruleset FOREIGN KEY (ruleset_id) REFERENCES ruleset (id) ON DELETE SET NULL,
    CONSTRAINT fk_events_actor FOREIGN KEY (actor_user_id) REFERENCES users (id) ON DELETE SET NULL,

    CONSTRAINT chk_events_balance_after_nonnegative CHECK (balance_after >= 0),
    CONSTRAINT chk_events_amount_money_nonnegative CHECK (amount_money IS NULL OR amount_money >= 0)
);

CREATE INDEX idx_events_account_ts ON events (account_id, ts);
CREATE INDEX idx_events_actor_ts ON events (actor_user_id, ts);

-- +goose Down
DROP INDEX idx_events_actor_ts;
DROP INDEX idx_events_account_ts;
DROP TABLE events;
