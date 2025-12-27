package rules

import (
	"fmt"

	"github.com/shopspring/decimal"

	"Beanefits/internal/domain/ledger"
)

// ResolveLevel picks the best level rule for the given totalSpend.
// Levels must be validated and sorted ascending by threshold.
func ResolveLevel(totalSpend ledger.Money, levels []LevelRule) (LevelRule, error) {
	if len(levels) == 0 {
		return LevelRule{}, fmt.Errorf("%w: empty levels", ErrInvalidLevels)
	}

	best := levels[0]
	for _, lr := range levels {
		if totalSpend.GTE(lr.ThresholdTotalSpend) {
			best = lr
		} else {
			break
		}
	}
	return best, nil
}

// ComputeEarnPoints computes earned points using:
// basePoints = floor(amountMoney / baseRubPerPoint)
// earned = floor(basePoints * (percentEarn / 100))
//
// Note: percent is determined by current level (usually based on totalSpend BEFORE this purchase).
func ComputeEarnPoints(amountMoney ledger.Money, baseRubPerPoint ledger.Money, percentEarn Percent) (ledger.Points, error) {
	if amountMoney.IsNegative() {
		return 0, fmt.Errorf("%w: amountMoney", ErrInvalidRuleset)
	}
	if !baseRubPerPoint.Decimal().GreaterThan(decimal.Zero) {
		return 0, fmt.Errorf("%w: baseRubPerPoint must be > 0", ErrInvalidRuleset)
	}
	if !percentEarn.IsPositive() {
		return 0, fmt.Errorf("%w: percentEarn must be > 0", ErrInvalidRuleset)
	}

	// basePoints = floor(amount / base)
	basePointsDec := amountMoney.Decimal().Div(baseRubPerPoint.Decimal()).Floor()
	basePoints := basePointsDec.IntPart() // safe: floor => integer

	// earned = floor(basePoints * percent / 100)
	earnedDec := decimal.NewFromInt(basePoints).
		Mul(percentEarn.Decimal()).
		Div(decimal.NewFromInt(100)).
		Floor()

	return ledger.Points(int(earnedDec.IntPart())), nil
}
