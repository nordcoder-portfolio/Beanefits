package rules

import (
	"fmt"

	"github.com/shopspring/decimal"

	"Beanefits/internal/domain/ledger"
)

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

	basePointsDec := amountMoney.Decimal().Div(baseRubPerPoint.Decimal()).Floor()
	basePoints := basePointsDec.IntPart()

	earnedDec := decimal.NewFromInt(basePoints).
		Mul(percentEarn.Decimal()).
		Div(decimal.NewFromInt(100)).
		Floor()

	return ledger.Points(int(earnedDec.IntPart())), nil
}
