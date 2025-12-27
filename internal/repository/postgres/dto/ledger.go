package dto

type EventRow struct {
	ID           int64
	AccountID    int64
	Type         EventType
	DeltaPoints  int // signed: + for EARN, - for SPEND
	BalanceAfter int
	AmountMoney  *Money // present for EARN
	RulesetID    *int64
	ActorUserID  *int64
	Ts           Ts // event time (business timestamp)
	CreatedAt    Ts // insertion time (optional, if you store separately)
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
