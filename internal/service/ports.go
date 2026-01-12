package service

import (
	"context"
	"time"

	"Beanefits/internal/service/dto"
)

type Auth interface {
	RegisterClient(ctx context.Context, in dto.RegisterIn) (dto.AuthOut, error)
	Login(ctx context.Context, in dto.LoginIn) (dto.AuthOut, error)
}

type Client interface {
	GetMe(ctx context.Context, userID int64) (dto.ClientProfileOut, error)
	GetBalance(ctx context.Context, userID int64) (dto.BalanceOut, error)
	GetEvents(ctx context.Context, userID int64, in dto.EventsIn) (dto.EventsOut, error)
}

type Cashier interface {
	LookupAccountByCode(ctx context.Context, publicCode string) (dto.AccountOut, error)
	GetAccountEventsByCode(ctx context.Context, publicCode string, in dto.EventsIn) (dto.EventsOut, error)

	Earn(ctx context.Context, actorUserID int64, in dto.EarnIn) (dto.OperationOut, error)
	Spend(ctx context.Context, actorUserID int64, in dto.SpendIn) (dto.OperationOut, error)
}

type Admin interface {
	ListUsers(ctx context.Context, in dto.ListUsersIn) (dto.UsersOut, error)
	DeactivateUser(ctx context.Context, userID int64) error

	CreateRuleset(ctx context.Context, actorUserID int64, in dto.CreateRulesetIn) (dto.RulesetOut, error)
	ListRulesets(ctx context.Context, in dto.ListRulesetsIn) (dto.RulesetsOut, error)
	GetCurrentRuleset(ctx context.Context, at time.Time) (dto.RulesetOut, error)
}
