package dto

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

type RoleCode string

const (
	RoleClient  RoleCode = "CLIENT"
	RoleCashier RoleCode = "CASHIER"
	RoleAdmin   RoleCode = "ADMIN"
)

type EventType string

const (
	EventEarn  EventType = "EARN"
	EventSpend EventType = "SPEND"
)

type OperationType string

const (
	OpEarn  OperationType = "EARN"
	OpSpend OperationType = "SPEND"
)

type JSON = json.RawMessage
type Money = decimal.Decimal
type Ts = time.Time
