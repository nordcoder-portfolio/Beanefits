package dto

type RulesetRow struct {
	ID              int64
	EffectiveFrom   Ts
	BaseRubPerPoint Money
	CreatedAt       Ts
}

type LevelRuleRow struct {
	ID                  int64
	RulesetID           int64
	LevelCode           string
	ThresholdTotalSpend Money
	PercentEarn         Money
}

type RulesetWithLevels struct {
	Ruleset RulesetRow
	Levels  []LevelRuleRow
}
