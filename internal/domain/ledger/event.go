package ledger

import "time"

type EventType string

const (
	EventEarn  EventType = "EARN"
	EventSpend EventType = "SPEND"
)

// EventDraft is a domain-level "event to be persisted" model.
// Repositories will store it into the events table and return the persisted EventOut.
type EventDraft struct {
	AccountID    int64
	Type         EventType
	DeltaPoints  Points
	BalanceAfter Points
	AmountMoney  *Money
	RulesetID    *int64
	ActorUserID  *int64
	Ts           time.Time
}

func NewEarnDraft(accountID int64, earned Points, balanceAfter Points, amount Money, rulesetID *int64, actorUserID *int64, ts time.Time) EventDraft {
	return EventDraft{
		AccountID:    accountID,
		Type:         EventEarn,
		DeltaPoints:  earned,
		BalanceAfter: balanceAfter,
		AmountMoney:  &amount,
		RulesetID:    rulesetID,
		ActorUserID:  actorUserID,
		Ts:           ts,
	}
}

func NewSpendDraft(accountID int64, spent Points, balanceAfter Points, actorUserID *int64, ts time.Time) EventDraft {
	neg := -spent
	return EventDraft{
		AccountID:    accountID,
		Type:         EventSpend,
		DeltaPoints:  neg,
		BalanceAfter: balanceAfter,
		AmountMoney:  nil,
		RulesetID:    nil,
		ActorUserID:  actorUserID,
		Ts:           ts,
	}
}
