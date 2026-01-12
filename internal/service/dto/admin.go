package dto

import "time"

type ListUsersIn struct {
	Limit  int    `validate:"omitempty,gte=1,lte=100"`
	Offset int    `validate:"omitempty,gte=0"`
	Q      string `validate:"omitempty,max=64"`
}

type UsersOut struct {
	Items []UserWithRoles `validate:"required"`
	Total *int            `validate:"omitempty,gte=0"`
}

type CreateRulesetIn struct {
	EffectiveFrom   time.Time     `validate:"required"`
	BaseRubPerPoint string        `validate:"required,decimal2,gtzero_decimal"`
	Levels          []LevelRuleIn `validate:"required,min=1,dive"`
}

type LevelRuleIn struct {
	LevelCode           string `validate:"required,min=1,max=64"`
	ThresholdTotalSpend string `validate:"required,decimal2"`
	PercentEarn         string `validate:"required,decimal2,gtzero_decimal"`
}

type ListRulesetsIn struct {
	Limit  int `validate:"omitempty,gte=1,lte=100"`
	Offset int `validate:"omitempty,gte=0"`
}

type RulesetsOut struct {
	Items []RulesetOut `validate:"required"`
	Total *int         `validate:"omitempty,gte=0"`
}

type RulesetOut struct {
	ID              int64          `validate:"required,gt=0"`
	EffectiveFrom   time.Time      `validate:"required"`
	BaseRubPerPoint string         `validate:"required,decimal2"`
	Levels          []LevelRuleOut `validate:"required"`
	CreatedAt       time.Time      `validate:"required"`
}

type LevelRuleOut struct {
	ID                  int64  `validate:"required,gt=0"`
	LevelCode           string `validate:"required,min=1,max=64"`
	ThresholdTotalSpend string `validate:"required,decimal2"`
	PercentEarn         string `validate:"required,decimal2"`
}
