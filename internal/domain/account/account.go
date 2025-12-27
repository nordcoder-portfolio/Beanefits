package account

import (
	"Beanefits/internal/domain/errs"
	"github.com/google/uuid"
	"time"

	"Beanefits/internal/domain/ledger"
	"Beanefits/internal/domain/rules"
)

type PublicCode string

// ParsePublicCode validates a public code format (QR payload) as UUID.
func ParsePublicCode(s string) (PublicCode, error) {
	if _, err := uuid.Parse(s); err != nil {
		return "", ErrInvalidPublicCode
	}
	return PublicCode(s), nil
}

// Account is the main aggregate for loyalty balance.
type Account struct {
	ID         int64
	PublicCode PublicCode

	Balance    ledger.Points
	TotalSpend ledger.Money
	LevelCode  rules.LevelCode
}

// CanSpend checks spend invariants.
func (a Account) CanSpend(p ledger.Points) error {
	if err := p.ValidatePositive(); err != nil {
		return err
	}
	if a.Balance < p {
		return ErrNotEnoughBalance
	}
	return nil
}

// ApplySpend applies spending points and returns updated account and an event draft.
func (a Account) ApplySpend(p ledger.Points, actorUserID *int64, ts time.Time) (Account, ledger.EventDraft, error) {
	if err := a.CanSpend(p); err != nil {
		return Account{}, ledger.EventDraft{}, err
	}
	a.Balance -= p
	ev := ledger.NewSpendDraft(a.ID, p, a.Balance, actorUserID, ts)
	return a, ev, nil
}

// ApplyEarn applies earned points and purchase amount.
// Usecase is expected to compute earned points + new level based on ruleset.
func (a Account) ApplyEarn(earned ledger.Points, purchase ledger.Money, newLevel rules.LevelCode, rulesetID *int64, actorUserID *int64, ts time.Time) (Account, ledger.EventDraft, error) {
	if err := earned.ValidateNonNegative(); err != nil {
		return Account{}, ledger.EventDraft{}, err
	}
	if purchase.IsNegative() {
		return Account{}, ledger.EventDraft{}, ErrInvalidPurchaseAmount
	}
	a.Balance += earned
	a.TotalSpend = a.TotalSpend.Add(purchase)
	a.LevelCode = newLevel

	ev := ledger.NewEarnDraft(a.ID, earned, a.Balance, purchase, rulesetID, actorUserID, ts)
	return a, ev, nil
}

var (
	ErrNotEnoughBalance      = errs.New(errs.CodeNotEnoughBalance, "not enough balance")
	ErrInvalidPublicCode     = errs.New(errs.CodeInvalidPublicCode, "invalid public code format")
	ErrInvalidPurchaseAmount = errs.New(errs.CodeInvalidPurchaseAmount, "purchase amount must be >= 0")
)
