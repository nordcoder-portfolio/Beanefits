-- +goose Up
CREATE TABLE level_rules
(
    id                    BIGSERIAL PRIMARY KEY,
    ruleset_id            BIGINT         NOT NULL,
    level_code            TEXT           NOT NULL,
    threshold_total_spend NUMERIC(12, 2) NOT NULL,
    percent_earn          NUMERIC(5, 2)  NOT NULL,

    CONSTRAINT fk_level_rules_ruleset FOREIGN KEY (ruleset_id) REFERENCES ruleset (id) ON DELETE CASCADE,

    CONSTRAINT uq_level_rules_ruleset_level UNIQUE (ruleset_id, level_code),
    CONSTRAINT uq_level_rules_ruleset_threshold UNIQUE (ruleset_id, threshold_total_spend),

    CONSTRAINT chk_level_rules_threshold_nonnegative CHECK (threshold_total_spend >= 0),
    CONSTRAINT chk_level_rules_percent_positive CHECK (percent_earn > 0)
);

-- +goose Down
DROP TABLE level_rules;
