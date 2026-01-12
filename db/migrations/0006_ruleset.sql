-- +goose Up
CREATE TABLE ruleset
(
    id                 BIGSERIAL PRIMARY KEY,
    effective_from     TIMESTAMPTZ    NOT NULL,
    base_rub_per_point NUMERIC(10, 2) NOT NULL,
    created_by         BIGINT,
    created_at         TIMESTAMPTZ    NOT NULL DEFAULT now(),

    CONSTRAINT uq_ruleset_effective_from UNIQUE (effective_from),

    CONSTRAINT fk_ruleset_created_by FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,

    CONSTRAINT chk_ruleset_base_rub_per_point_positive CHECK (base_rub_per_point > 0)
);

-- +goose Down
DROP TABLE ruleset;
