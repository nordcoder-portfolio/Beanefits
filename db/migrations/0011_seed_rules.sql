-- +goose Up
-- Seed baseline ruleset + levels for local/dev.
-- Levels: Green Bean, Light Roast, Medium Roast, Dark Roast, Premium Roast

-- We use a fixed effective_from so Postman/env can rely on deterministic rules.
-- Any operation with ts >= this effective_from will pick this ruleset (until a newer one appears).
WITH ins AS (
    INSERT INTO ruleset (effective_from, base_rub_per_point, created_by)
        VALUES (
                   '2025-01-01T00:00:00Z'::timestamptz,
                   10.00::numeric(10,2),
                   (SELECT id FROM users WHERE phone = '+79990000001' LIMIT 1) -- admin
               )
        ON CONFLICT (effective_from) DO NOTHING
        RETURNING id
),
     rid AS (
         SELECT id FROM ins
         UNION ALL
         SELECT id FROM ruleset WHERE effective_from = '2025-01-01T00:00:00Z'::timestamptz
         LIMIT 1
     )
INSERT INTO level_rules (ruleset_id, level_code, threshold_total_spend, percent_earn)
SELECT rid.id, x.level_code, x.threshold_total_spend, x.percent_earn
FROM rid
         JOIN (
    VALUES
        ('Green Bean',   0.00::numeric(12,2),   100.00::numeric(5,2)),
        ('Light Roast',  5000.00::numeric(12,2),110.00::numeric(5,2)),
        ('Medium Roast', 15000.00::numeric(12,2),120.00::numeric(5,2)),
        ('Dark Roast',   30000.00::numeric(12,2),130.00::numeric(5,2)),
        ('Premium Roast',60000.00::numeric(12,2),150.00::numeric(5,2))
) AS x(level_code, threshold_total_spend, percent_earn) ON TRUE
ON CONFLICT DO NOTHING;

-- +goose Down
-- Best-effort rollback for dev: also removes events/operations that reference this ruleset.
WITH rid AS (
    SELECT id
    FROM ruleset
    WHERE effective_from = '2025-01-01T00:00:00Z'::timestamptz
    LIMIT 1
),
     ev AS (
         SELECT e.id
         FROM events e
                  JOIN rid ON rid.id = e.ruleset_id
     )
DELETE FROM operations o
WHERE o.event_id IN (SELECT id FROM ev);

WITH rid AS (
    SELECT id
    FROM ruleset
    WHERE effective_from = '2025-01-01T00:00:00Z'::timestamptz
    LIMIT 1
)
DELETE FROM events e
WHERE e.ruleset_id IN (SELECT id FROM rid);

WITH rid AS (
    SELECT id
    FROM ruleset
    WHERE effective_from = '2025-01-01T00:00:00Z'::timestamptz
    LIMIT 1
)
DELETE FROM level_rules lr
WHERE lr.ruleset_id IN (SELECT id FROM rid);

DELETE FROM ruleset
WHERE effective_from = '2025-01-01T00:00:00Z'::timestamptz;
