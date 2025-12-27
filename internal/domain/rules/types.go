package rules

import (
	"Beanefits/internal/domain/errs"
	"fmt"
	"sort"

	"github.com/shopspring/decimal"

	"Beanefits/internal/domain/ledger"
)

type LevelCode string

// Percent is stored as decimal percent, e.g. 100.00, 110.00.
type Percent struct {
	d decimal.Decimal
}

func ParsePercent(s string) (Percent, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return Percent{}, fmt.Errorf("percent parse: %w", err)
	}
	return Percent{d: d}, nil
}

func (p Percent) Decimal() decimal.Decimal { return p.d }
func (p Percent) IsPositive() bool         { return p.d.GreaterThan(decimal.Zero) }

// LevelRule defines a threshold and percent for a single customer level.
type LevelRule struct {
	ID                  int64
	LevelCode           LevelCode
	ThresholdTotalSpend ledger.Money
	PercentEarn         Percent
}

// Ruleset is a versioned rules configuration (no retroactive changes).
type Ruleset struct {
	ID              int64
	BaseRubPerPoint ledger.Money
	Levels          []LevelRule // must be validated & sorted by ThresholdTotalSpend
}

// ValidateLevels checks basic invariants for levels.
func ValidateLevels(levels []LevelRule) error {
	if len(levels) == 0 {
		return ErrInvalidLevels
	}

	seenCode := make(map[LevelCode]struct{}, len(levels))
	seenThreshold := make(map[string]struct{}, len(levels))

	for _, lr := range levels {
		if lr.LevelCode == "" {
			return fmt.Errorf("%w: empty levelCode", ErrInvalidLevels)
		}
		if _, ok := seenCode[lr.LevelCode]; ok {
			return fmt.Errorf("%w: duplicate levelCode %q", ErrInvalidLevels, lr.LevelCode)
		}
		seenCode[lr.LevelCode] = struct{}{}

		if lr.ThresholdTotalSpend.IsNegative() {
			return fmt.Errorf("%w: negative threshold for %q", ErrInvalidLevels, lr.LevelCode)
		}
		// Use normalized string as a uniqueness key.
		thKey := lr.ThresholdTotalSpend.Decimal().String()
		if _, ok := seenThreshold[thKey]; ok {
			return fmt.Errorf("%w: duplicate threshold %s", ErrInvalidLevels, thKey)
		}
		seenThreshold[thKey] = struct{}{}

		if !lr.PercentEarn.IsPositive() {
			return fmt.Errorf("%w: percent must be > 0 for %q", ErrInvalidLevels, lr.LevelCode)
		}
	}

	// Ensure thresholds are non-decreasing once sorted.
	sorted := append([]LevelRule(nil), levels...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].ThresholdTotalSpend.LT(sorted[j].ThresholdTotalSpend)
	})
	for i := 1; i < len(sorted); i++ {
		if sorted[i].ThresholdTotalSpend.LT(sorted[i-1].ThresholdTotalSpend) {
			return ErrInvalidLevels
		}
	}

	return nil
}

// SortLevels sorts levels by threshold ascending (recommended to call after validation).
func SortLevels(levels []LevelRule) {
	sort.Slice(levels, func(i, j int) bool {
		return levels[i].ThresholdTotalSpend.LT(levels[j].ThresholdTotalSpend)
	})
}

var (
	ErrInvalidRuleset = errs.New(errs.CodeInvalidRuleset, "invalid ruleset")
	ErrInvalidLevels  = errs.New(errs.CodeInvalidLevels, "invalid levels")
)
