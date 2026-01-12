package validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
	"sort"
	"strings"

	"Beanefits/internal/domain/errs"
	pgdto "Beanefits/internal/repository/postgres/dto"
	sdto "Beanefits/internal/service/dto"

	"github.com/shopspring/decimal"
)

var rePhone = regexp.MustCompile(`^\+?[1-9]\d{10,14}$`)

func RegisterValidations(v *validator.Validate) error {
	if err := v.RegisterValidation("phone", validatePhone); err != nil {
		return err
	}
	if err := v.RegisterValidation("decimal2", validateDecimal2); err != nil {
		return err
	}
	if err := v.RegisterValidation("gtzero_decimal", validateGtZeroDecimal); err != nil {
		return err
	}
	return nil
}

func validatePhone(fl validator.FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	return ok && rePhone.MatchString(s)
}

func validateDecimal2(fl validator.FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	_, err := ParseDecimal2(s)
	return err == nil
}

func validateGtZeroDecimal(fl validator.FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	d, err := ParseDecimal2(s)
	if err != nil {
		return false
	}
	return d.GreaterThan(decimal.Zero)
}

func NormalizeListParams(limit, offset, defaultLimit, maxLimit int) (int, int) {
	if limit <= 0 {
		limit = defaultLimit
	}
	if maxLimit > 0 && limit > maxLimit {
		limit = maxLimit
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

func ParseDecimal2(s string) (decimal.Decimal, error) {
	v := strings.TrimSpace(s)
	d, err := decimal.NewFromString(v)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("decimal parse: %w", err)
	}
	if d.Exponent() < -2 {
		return decimal.Decimal{}, fmt.Errorf("too many fraction digits: %s", v)
	}
	return d, nil
}

func ValidateAndMapLevelRules(in []sdto.LevelRuleIn) ([]pgdto.LevelRuleRow, error) {
	if len(in) == 0 {
		return nil, errs.New(errs.CodeInvalidLevels, "levels must be non-empty")
	}

	const maxLevelCodeLen = 64

	seenCode := make(map[string]struct{}, len(in))
	seenThreshold := make(map[string]struct{}, len(in))

	out := make([]pgdto.LevelRuleRow, 0, len(in))

	for i, lv := range in {
		code := strings.TrimSpace(lv.LevelCode)
		if code == "" {
			return nil, errs.New(errs.CodeInvalidLevels, fmt.Sprintf("levels[%d].levelCode is required", i))
		}
		if len(code) > maxLevelCodeLen {
			return nil, errs.New(errs.CodeInvalidLevels, fmt.Sprintf("levels[%d].levelCode too long", i))
		}
		if _, ok := seenCode[code]; ok {
			return nil, errs.New(errs.CodeInvalidLevels, fmt.Sprintf("duplicate levelCode: %s", code))
		}
		seenCode[code] = struct{}{}

		thr, err := ParseDecimal2(lv.ThresholdTotalSpend)
		if err != nil {
			return nil, errs.Wrap(errs.CodeInvalidMoney, fmt.Sprintf("levels[%d].thresholdTotalSpend invalid", i), err)
		}
		if thr.Cmp(decimal.Zero) < 0 {
			return nil, errs.New(errs.CodeInvalidLevels, fmt.Sprintf("levels[%d].thresholdTotalSpend must be >= 0", i))
		}
		thrKey := thr.StringFixed(2)
		if _, ok := seenThreshold[thrKey]; ok {
			return nil, errs.New(errs.CodeInvalidLevels, fmt.Sprintf("duplicate thresholdTotalSpend: %s", thrKey))
		}
		seenThreshold[thrKey] = struct{}{}

		perc, err := ParseDecimal2(lv.PercentEarn)
		if err != nil {
			return nil, errs.Wrap(errs.CodeInvalidMoney, fmt.Sprintf("levels[%d].percentEarn invalid", i), err)
		}
		if perc.Cmp(decimal.Zero) <= 0 {
			return nil, errs.New(errs.CodeInvalidLevels, fmt.Sprintf("levels[%d].percentEarn must be > 0", i))
		}

		out = append(out, pgdto.LevelRuleRow{
			ID:                  0,
			RulesetID:           0,
			LevelCode:           code,
			ThresholdTotalSpend: pgdto.Money(thr),
			PercentEarn:         pgdto.Money(perc),
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].ThresholdTotalSpend.Cmp(out[j].ThresholdTotalSpend) < 0
	})

	if out[0].ThresholdTotalSpend.Cmp(decimal.Zero) != 0 {
		return nil, errs.New(errs.CodeInvalidLevels, "must include baseline level with thresholdTotalSpend=0.00")
	}

	return out, nil
}
