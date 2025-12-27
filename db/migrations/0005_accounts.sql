-- +goose Up
CREATE TABLE accounts
(
    id                BIGSERIAL PRIMARY KEY,
    user_id           BIGINT         NOT NULL,
    public_code       TEXT           NOT NULL,
    created_at        TIMESTAMPTZ    NOT NULL DEFAULT now(),
    balance_points    INT            NOT NULL DEFAULT 0,
    total_spend_money NUMERIC(12, 2) NOT NULL DEFAULT 0,
    level_code        TEXT,

    CONSTRAINT uq_accounts_user_id UNIQUE (user_id),
    CONSTRAINT uq_accounts_public_code UNIQUE (public_code),

    CONSTRAINT fk_accounts_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT,

    CONSTRAINT chk_accounts_balance_nonnegative CHECK (balance_points >= 0),
    CONSTRAINT chk_accounts_total_spend_nonnegative CHECK (total_spend_money >= 0)
);

-- +goose Down
DROP TABLE accounts;
