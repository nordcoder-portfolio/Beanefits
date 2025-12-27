package dto

import "time"

// ListUsersIn is the usecase input for listing users with optional search and pagination.
type ListUsersIn struct {
	Limit  int    `validate:"omitempty,gte=1,lte=100"`
	Offset int    `validate:"omitempty,gte=0"`
	Q      string `validate:"omitempty,max=64"`
}

// UsersOut is the usecase output for listing users.
type UsersOut struct {
	Items []UserWithRoles `validate:"required"`
	Total *int            `validate:"omitempty,gte=0"` // optional total count
}

// CreateRulesetIn is the usecase input for creating a new ruleset (not retroactive).
type CreateRulesetIn struct {
	EffectiveFrom   time.Time     `validate:"required"`
	BaseRubPerPoint string        `validate:"required,decimal2,gtzero_decimal"`
	Levels          []LevelRuleIn `validate:"required,min=1,dive"`
}

// LevelRuleIn is a single level definition inside a ruleset.
type LevelRuleIn struct {
	LevelCode           string `validate:"required,min=1,max=64"`
	ThresholdTotalSpend string `validate:"required,decimal2"`
	PercentEarn         string `validate:"required,decimal2,gtzero_decimal"`
}

// ListRulesetsIn is the usecase input for listing rulesets with pagination.
type ListRulesetsIn struct {
	Limit  int `validate:"omitempty,gte=1,lte=100"`
	Offset int `validate:"omitempty,gte=0"`
}

// RulesetsOut is the usecase output for listing rulesets.
type RulesetsOut struct {
	Items []RulesetOut `validate:"required"`
	Total *int         `validate:"omitempty,gte=0"`
}

// RulesetOut is the ruleset representation returned to admin.
type RulesetOut struct {
	ID              int64          `validate:"required,gt=0"`
	EffectiveFrom   time.Time      `validate:"required"`
	BaseRubPerPoint string         `validate:"required,decimal2"`
	Levels          []LevelRuleOut `validate:"required"`
	CreatedAt       time.Time      `validate:"required"`
}

// LevelRuleOut is the stored level rule representation returned to admin.
type LevelRuleOut struct {
	ID                  int64  `validate:"required,gt=0"`
	LevelCode           string `validate:"required,min=1,max=64"`
	ThresholdTotalSpend string `validate:"required,decimal2"`
	PercentEarn         string `validate:"required,decimal2"`
}
