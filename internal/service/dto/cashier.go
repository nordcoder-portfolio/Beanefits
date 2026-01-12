package dto

import "time"

type AccountOut struct {
	AccountBase
}

type EarnIn struct {
	OperationID string     `validate:"required,uuid"`
	PublicCode  string     `validate:"required,min=6,max=64"`
	AmountMoney string     `validate:"required,decimal2"`
	Ts          *time.Time `validate:"omitempty"`
}

type SpendIn struct {
	OperationID  string     `validate:"required,uuid"`
	PublicCode   string     `validate:"required,min=6,max=64"`
	AmountPoints int        `validate:"required,gt=0"`
	Ts           *time.Time `validate:"omitempty"`
}

type OperationType string

const (
	OpEarn  OperationType = "EARN"
	OpSpend OperationType = "SPEND"
)

type OperationOut struct {
	OperationID      string        `validate:"required,uuid"`
	OpType           OperationType `validate:"required,oneof=EARN SPEND"`
	Event            EventOut      `validate:"required"`
	Balance          BalanceOut    `validate:"required"`
	IdempotentReplay bool          `validate:"-"`
}
