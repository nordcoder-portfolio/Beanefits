package admin

import (
	"context"
	"time"

	"Beanefits/internal/domain/errs"
	pg "Beanefits/internal/repository/postgres"
	pgdto "Beanefits/internal/repository/postgres/dto"
	"Beanefits/internal/service/dto"
	"Beanefits/internal/service/mapper"
	svcvalidation "Beanefits/internal/service/validation"

	"github.com/shopspring/decimal"
)

const (
	constraintRulesetEffectiveFrom = "ruleset_effective_from_key"

	constraintLevelRulesetLevelCode = "level_rules_ruleset_id_level_code_key"
	constraintLevelRulesetThreshold = "level_rules_ruleset_id_threshold_total_spend_key"
	defaultRulesetsLimit            = 20
	maxRulesetsLimit                = 100
)

func (s *Service) CreateRuleset(ctx context.Context, actorUserID int64, in dto.CreateRulesetIn) (dto.RulesetOut, error) {
	start := time.Now()
	s.log.InfoContext(ctx, "admin.create_ruleset start", "actorUserID", actorUserID, "effectiveFrom", in.EffectiveFrom)

	base, err := svcvalidation.ParseDecimal2(in.BaseRubPerPoint)
	if err != nil {
		e := errs.Wrap(errs.CodeInvalidMoney, "invalid baseRubPerPoint", err)
		s.log.ErrorContext(ctx, "admin.create_ruleset failed", "ms", time.Since(start).Milliseconds(), "err", e)
		return dto.RulesetOut{}, e
	}
	if base.Cmp(decimal.Zero) <= 0 {
		e := errs.New(errs.CodeInvalidMoney, "baseRubPerPoint must be > 0")
		s.log.ErrorContext(ctx, "admin.create_ruleset failed", "ms", time.Since(start).Milliseconds(), "err", e)
		return dto.RulesetOut{}, e
	}

	levelRows, err := svcvalidation.ValidateAndMapLevelRules(in.Levels)
	if err != nil {
		s.log.ErrorContext(ctx, "admin.create_ruleset failed", "ms", time.Since(start).Milliseconds(), "err", err)
		return dto.RulesetOut{}, err
	}

	var created pgdto.RulesetWithLevels

	err = s.txm.WithinTx(ctx, func(ctx context.Context, tx pg.DBTX) error {
		r, err := s.rules.CreateRuleset(ctx, tx, in.EffectiveFrom, pgdto.Money(base), levelRows)
		if err != nil {
			// human-friendly ошибки на уникальные ограничения.
			if pg.IsUniqueViolation(err, constraintRulesetEffectiveFrom) {
				return errs.New(errs.CodeInvalidRuleset, "ruleset.effectiveFrom already exists")
			}
			if pg.IsUniqueViolation(err, constraintLevelRulesetLevelCode) ||
				pg.IsUniqueViolation(err, constraintLevelRulesetThreshold) {
				return errs.New(errs.CodeInvalidLevels, "duplicate level definitions")
			}
			return errs.Wrap(errs.CodeInternal, "rules.create_ruleset", err)
		}
		created = r
		return nil
	})

	if err != nil {
		s.log.ErrorContext(ctx, "admin.create_ruleset failed", "ms", time.Since(start).Milliseconds(), "err", err)
		return dto.RulesetOut{}, err
	}

	out := mapper.RulesetOut(created)
	s.log.InfoContext(ctx, "admin.create_ruleset ok",
		"ms", time.Since(start).Milliseconds(),
		"rulesetID", out.ID,
	)

	return out, nil
}

func (s *Service) ListRulesets(ctx context.Context, in dto.ListRulesetsIn) (dto.RulesetsOut, error) {
	limit, offset := svcvalidation.NormalizeListParams(in.Limit, in.Offset, defaultRulesetsLimit, maxRulesetsLimit)

	start := time.Now()
	s.log.InfoContext(ctx, "admin.list_rulesets start", "limit", limit, "offset", offset)

	items, err := s.rules.ListRulesets(ctx, s.db, limit, offset)
	if err != nil {
		wrapped := errs.Wrap(errs.CodeInternal, "rules.list_rulesets", err)
		s.log.ErrorContext(ctx, "admin.list_rulesets failed", "ms", time.Since(start).Milliseconds(), "err", wrapped)
		return dto.RulesetsOut{}, wrapped
	}

	out := dto.RulesetsOut{
		Items: make([]dto.RulesetOut, 0, len(items)),
		Total: nil,
	}
	for _, it := range items {
		out.Items = append(out.Items, mapper.RulesetOut(it))
	}

	s.log.InfoContext(ctx, "admin.list_rulesets ok",
		"ms", time.Since(start).Milliseconds(),
		"items", len(out.Items),
	)

	return out, nil
}

func (s *Service) GetCurrentRuleset(ctx context.Context, at time.Time) (dto.RulesetOut, error) {
	start := time.Now()
	s.log.InfoContext(ctx, "admin.get_current_ruleset start", "at", at)

	r, ok, err := s.rules.GetEffectiveAt(ctx, s.db, at)
	if err != nil {
		wrapped := errs.Wrap(errs.CodeInternal, "rules.get_effective_at", err)
		s.log.ErrorContext(ctx, "admin.get_current_ruleset failed", "ms", time.Since(start).Milliseconds(), "err", wrapped)
		return dto.RulesetOut{}, wrapped
	}
	if !ok {
		e := errs.New(errs.CodeInvalidRuleset, "ruleset not found")
		s.log.ErrorContext(ctx, "admin.get_current_ruleset failed", "ms", time.Since(start).Milliseconds(), "err", e)
		return dto.RulesetOut{}, e
	}

	out := mapper.RulesetOut(r)
	s.log.InfoContext(ctx, "admin.get_current_ruleset ok",
		"ms", time.Since(start).Milliseconds(),
		"rulesetID", out.ID,
	)
	return out, nil
}
