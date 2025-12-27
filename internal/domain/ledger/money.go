package ledger

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// Money is a decimal amount with 2 fraction digits semantics (rubles).
// It is stored/processed as decimal to avoid float errors.
type Money struct {
	d decimal.Decimal
}

// ParseMoney parses decimal string (e.g. "10.00", "450", "12500.5").
// Callers may enforce scale separately if needed.
func ParseMoney(s string) (Money, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return Money{}, fmt.Errorf("money parse: %w", err)
	}
	return Money{d: d}, nil
}

func MustMoney(s string) Money {
	m, err := ParseMoney(s)
	if err != nil {
		panic(err)
	}
	return m
}

func ZeroMoney() Money {
	return Money{d: decimal.Zero}
}

func (m Money) Decimal() decimal.Decimal { return m.d }

func (m Money) String() string {
	// Keep scale as-is. If you want fixed 2 digits for API, format in usecase/mapper.
	return m.d.String()
}

func (m Money) IsNegative() bool { return m.d.IsNegative() }
func (m Money) IsZero() bool     { return m.d.IsZero() }

func (m Money) Add(x Money) Money { return Money{d: m.d.Add(x.d)} }
func (m Money) Sub(x Money) Money { return Money{d: m.d.Sub(x.d)} }

func (m Money) Cmp(x Money) int  { return m.d.Cmp(x.d) }
func (m Money) GT(x Money) bool  { return m.Cmp(x) > 0 }
func (m Money) GTE(x Money) bool { return m.Cmp(x) >= 0 }
func (m Money) LT(x Money) bool  { return m.Cmp(x) < 0 }
func (m Money) LTE(x Money) bool { return m.Cmp(x) <= 0 }
