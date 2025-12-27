package dto

import "time"

// ClientProfileOut — ответ для Client.GetMe(userID)
type ClientProfileOut struct {
	User    UserWithRoles `validate:"required"`
	Account AccountBase   `validate:"required"`
}

// BalanceOut — ответ для Client.GetBalance(userID)
type BalanceOut struct {
	AccountID       int64     `validate:"required,gt=0"`
	BalancePoints   int       `validate:"required,gte=0"`
	TotalSpendMoney string    `validate:"required"`
	LevelCode       string    `validate:"required,min=1,max=64"`
	AsOf            time.Time `validate:"required"`
}

// EventsIn — вход для Client.GetEvents(userID, in)
type EventsIn struct {
	Limit    int        `validate:"omitempty,gte=1,lte=100"`
	BeforeTs *time.Time `validate:"omitempty"`
}

// EventsOut — ответ для Client.GetEvents(...)
type EventsOut struct {
	Items        []EventOut `validate:"required"`
	NextBeforeTs *time.Time `validate:"omitempty"`
}

// EventType — типы событий для истории
type EventType string

const (
	EventEarn  EventType = "EARN"
	EventSpend EventType = "SPEND"
)

// EventOut — строка истории для клиента
type EventOut struct {
	ID           int64     `validate:"required,gt=0"`
	AccountID    int64     `validate:"required,gt=0"`
	Type         EventType `validate:"required,oneof=EARN SPEND"`
	DeltaPoints  int       `validate:"required"`
	BalanceAfter int       `validate:"required,gte=0"`
	AmountMoney  *string   `validate:"omitempty"`
	RulesetID    *int64    `validate:"omitempty"`
	ActorUserID  *int64    `validate:"omitempty"`
	Ts           time.Time `validate:"required"`
}
