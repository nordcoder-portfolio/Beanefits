package ledger

import "fmt"

// Points are integer loyalty points.
type Points int

func (p Points) Int() int { return int(p) }

func (p Points) IsNegative() bool { return p < 0 }

func (p Points) ValidateNonNegative() error {
	if p < 0 {
		return fmt.Errorf("%w: %d", ErrInvalidPoints, p)
	}
	return nil
}

func (p Points) ValidatePositive() error {
	if p <= 0 {
		return fmt.Errorf("%w: %d", ErrInvalidPoints, p)
	}
	return nil
}
