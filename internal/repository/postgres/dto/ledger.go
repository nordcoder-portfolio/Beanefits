package dto

type EventRow struct {
	ID           int64
	AccountID    int64
	Type         EventType
	DeltaPoints  int
	BalanceAfter int
	AmountMoney  *Money
	RulesetID    *int64
	ActorUserID  *int64
	Ts           Ts
	CreatedAt    Ts
}

type EventInsert struct {
	AccountID    int64
	Type         EventType
	DeltaPoints  int
	BalanceAfter int
	AmountMoney  *Money
	RulesetID    *int64
	ActorUserID  *int64
	Ts           Ts
}
