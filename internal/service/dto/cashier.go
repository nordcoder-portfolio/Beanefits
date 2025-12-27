package dto

import "time"

// AccountOut is returned to a cashier after resolving a customer by publicCode (QR payload).
type AccountOut struct {
	AccountBase
}

// EarnIn is the usecase input for earning points (idempotent).
type EarnIn struct {
	OperationID string     `validate:"required,uuid"`
	PublicCode  string     `validate:"required,min=6,max=64"`
	AmountMoney string     `validate:"required,decimal2"` // decimal-as-string, up to 2 fractional digits
	Ts          *time.Time `validate:"omitempty"`
}

// SpendIn is the usecase input for spending points (idempotent, concurrency-safe).
type SpendIn struct {
	OperationID  string     `validate:"required,uuid"`
	PublicCode   string     `validate:"required,min=6,max=64"`
	AmountPoints int        `validate:"required,gt=0"`
	Ts           *time.Time `validate:"omitempty"`
}

// OperationType is a stable operation kind for idempotency.
type OperationType string

const (
	OpEarn  OperationType = "EARN"
	OpSpend OperationType = "SPEND"
)

// OperationOut is returned by Earn/Spend usecases.
type OperationOut struct {
	OperationID      string        `validate:"required,uuid"`
	OpType           OperationType `validate:"required,oneof=EARN SPEND"`
	Event            EventOut      `validate:"required"`
	Balance          BalanceOut    `validate:"required"`
	IdempotentReplay bool          `validate:"-"`
}
