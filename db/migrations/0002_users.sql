-- +goose Up
CREATE TABLE users
(
    id            BIGSERIAL PRIMARY KEY,
    phone         TEXT        NOT NULL,
    password_hash TEXT        NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_active     BOOLEAN     NOT NULL DEFAULT true,
    CONSTRAINT uq_users_phone UNIQUE (phone)
);

-- +goose Down
DROP TABLE users;
