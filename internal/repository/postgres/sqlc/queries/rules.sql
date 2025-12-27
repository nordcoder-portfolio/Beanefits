-- internal/repository/postgres/sqlc/queries/rules.sql

-- name: GetRulesetEffectiveAt :one
SELECT id, effective_from, base_rub_per_point, created_at
FROM ruleset
WHERE effective_from <= $1
ORDER BY effective_from DESC
LIMIT 1;

-- name: InsertRuleset :one
INSERT INTO ruleset (effective_from, base_rub_per_point)
VALUES ($1, $2)
RETURNING id, effective_from, base_rub_per_point, created_at;

-- name: InsertLevelRule :one
INSERT INTO level_rules (ruleset_id, level_code, threshold_total_spend, percent_earn)
VALUES ($1, $2, $3, $4)
RETURNING id, ruleset_id, level_code, threshold_total_spend, percent_earn;

-- name: GetRulesetByID :one
SELECT id, effective_from, base_rub_per_point, created_at
FROM ruleset
WHERE id = $1;

-- name: ListLevelRulesByRulesetID :many
SELECT id, ruleset_id, level_code, threshold_total_spend, percent_earn
FROM level_rules
WHERE ruleset_id = $1
ORDER BY threshold_total_spend ASC;

-- name: ListRulesetsBase :many
SELECT id, effective_from, base_rub_per_point, created_at
FROM ruleset
ORDER BY effective_from DESC
LIMIT $1 OFFSET $2;
