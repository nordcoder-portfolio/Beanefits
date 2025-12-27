package cashier

import (
	"context"
	"fmt"
	"time"

	"Beanefits/internal/domain/account"
	"Beanefits/internal/domain/errs"
	"Beanefits/internal/domain/ledger"
	"Beanefits/internal/domain/rules"
	pg "Beanefits/internal/repository/postgres"
	pgdto "Beanefits/internal/repository/postgres/dto"
	"Beanefits/internal/service/dto"
	"Beanefits/internal/service/mapper"
)

func (s *Service) Earn(ctx context.Context, actorUserID int64, in dto.EarnIn) (dto.OperationOut, error) {
	start := time.Now()
	s.log.InfoContext(ctx, "cashier.earn start", "actorUserID", actorUserID, "operationID", in.OperationID, "publicCode", in.PublicCode)

	if _, err := account.ParsePublicCode(in.PublicCode); err != nil {
		s.log.ErrorContext(ctx, "cashier.earn failed", "ms", time.Since(start).Milliseconds(), "err", err)
		return dto.OperationOut{}, err
	}

	opTs := s.now()
	if in.Ts != nil {
		opTs = *in.Ts
	}

	purchase, err := parseMoney2(in.AmountMoney)
	if err != nil {
		wrapped := errs.Wrap(errs.CodeInvalidMoney, "invalid amountMoney", err)
		s.log.ErrorContext(ctx, "cashier.earn failed", "ms", time.Since(start).Milliseconds(), "err", wrapped)
		return dto.OperationOut{}, wrapped
	}
	if purchase.IsNegative() {
		e := errs.New(errs.CodeInvalidPurchaseAmount, "purchase amount must be >= 0")
		s.log.ErrorContext(ctx, "cashier.earn failed", "ms", time.Since(start).Milliseconds(), "err", e)
		return dto.OperationOut{}, e
	}

	var result dto.OperationOut

	err = s.txm.WithinTx(ctx, func(ctx context.Context, tx pg.DBTX) error {
		accRow, ok, err := s.accounts.GetByPublicCode(ctx, tx, in.PublicCode)
		if err != nil {
			return errs.Wrap(errs.CodeInternal, "accounts.get_by_public_code", err)
		}
		if !ok {
			return errs.New(errs.CodeAccountNotFound, "account not found")
		}

		reqJSON, err := marshalEarnRequest(in, opTs)
		if err != nil {
			return errs.Wrap(errs.CodeInternal, "marshal earn request", err)
		}

		inserted, err := s.operations.InsertPending(ctx, tx, pgdto.OperationPendingInsert{
			AccountID:   accRow.ID,
			OpType:      pgdto.OpEarn,
			OperationID: in.OperationID,
			RequestJSON: reqJSON,
		})
		if err != nil {
			return errs.Wrap(errs.CodeInternal, "operations.insert_pending", err)
		}
		if !inserted {
			// idempotent replay
			return s.replayOperation(ctx, tx, accRow.ID, pgdto.OpEarn, in.OperationID, &result)
		}

		// concurrency gate
		lockedRow, err := s.accounts.LockByID(ctx, tx, accRow.ID)
		if err != nil {
			return errs.Wrap(errs.CodeInternal, "accounts.lock_by_id", err)
		}

		rs, ok, err := s.rules.GetEffectiveAt(ctx, tx, opTs)
		if err != nil {
			return errs.Wrap(errs.CodeInternal, "rules.get_effective_at", err)
		}
		if !ok {
			return errs.New(errs.CodeInvalidRuleset, "no ruleset effective at provided ts")
		}

		agg, err := accountFromRow(lockedRow)
		if err != nil {
			return errs.Wrap(errs.CodeInternal, "account aggregate build failed", err)
		}

		earned, levelBefore, levelAfter, baseRubPerPoint, err := computeEarnDomain(rs, agg.TotalSpend, purchase)
		if err != nil {
			return err
		}

		actor := actorUserID
		rulesetID := rs.Ruleset.ID

		updatedAgg, evDraft, err := agg.ApplyEarn(
			earned,
			purchase,
			levelAfter,
			&rulesetID,
			&actor,
			opTs,
		)
		if err != nil {
			return err
		}

		evRow, err := s.events.Insert(ctx, tx, eventInsertFromDraft(evDraft))
		if err != nil {
			return errs.Wrap(errs.CodeInternal, "events.insert", err)
		}

		updatedRow, err := s.accounts.UpdateAfterEarn(
			ctx, tx,
			updatedAgg.ID,
			updatedAgg.Balance.Int(),
			pgdto.Money(updatedAgg.TotalSpend.Decimal()),
			string(updatedAgg.LevelCode),
		)
		if err != nil {
			return errs.Wrap(errs.CodeInternal, "accounts.update_after_earn", err)
		}

		result = dto.OperationOut{
			OperationID:      in.OperationID,
			OpType:           dto.OpEarn,
			Event:            mapper.EventOut(evRow),
			Balance:          mapper.BalanceOut(updatedRow, s.now()),
			IdempotentReplay: false,
		}

		if err := s.finalizeOK(ctx, tx, updatedRow.ID, pgdto.OpEarn, in.OperationID, result); err != nil {
			return err
		}

		_ = levelBefore
		_ = baseRubPerPoint
		return nil
	})

	if err != nil {
		s.log.ErrorContext(ctx, "cashier.earn failed", "ms", time.Since(start).Milliseconds(), "err", err)
		return dto.OperationOut{}, err
	}

	s.log.InfoContext(ctx, "cashier.earn ok",
		"ms", time.Since(start).Milliseconds(),
		"operationID", in.OperationID,
		"publicCode", in.PublicCode,
		"replay", result.IdempotentReplay,
	)

	return result, nil
}

func (s *Service) Spend(ctx context.Context, actorUserID int64, in dto.SpendIn) (dto.OperationOut, error) {
	start := time.Now()
	s.log.InfoContext(ctx, "cashier.spend start", "actorUserID", actorUserID, "operationID", in.OperationID, "publicCode", in.PublicCode)

	if _, err := account.ParsePublicCode(in.PublicCode); err != nil {
		s.log.ErrorContext(ctx, "cashier.spend failed", "ms", time.Since(start).Milliseconds(), "err", err)
		return dto.OperationOut{}, err
	}

	opTs := s.now()
	if in.Ts != nil {
		opTs = *in.Ts
	}

	if in.AmountPoints <= 0 {
		e := errs.New(errs.CodeInvalidPoints, "amountPoints must be > 0")
		s.log.ErrorContext(ctx, "cashier.spend failed", "ms", time.Since(start).Milliseconds(), "err", e)
		return dto.OperationOut{}, e
	}

	var result dto.OperationOut

	err := s.txm.WithinTx(ctx, func(ctx context.Context, tx pg.DBTX) error {
		accRow, ok, err := s.accounts.GetByPublicCode(ctx, tx, in.PublicCode)
		if err != nil {
			return errs.Wrap(errs.CodeInternal, "accounts.get_by_public_code", err)
		}
		if !ok {
			return errs.New(errs.CodeAccountNotFound, "account not found")
		}

		reqJSON, err := marshalSpendRequest(in, opTs)
		if err != nil {
			return errs.Wrap(errs.CodeInternal, "marshal spend request", err)
		}

		inserted, err := s.operations.InsertPending(ctx, tx, pgdto.OperationPendingInsert{
			AccountID:   accRow.ID,
			OpType:      pgdto.OpSpend,
			OperationID: in.OperationID,
			RequestJSON: reqJSON,
		})
		if err != nil {
			return errs.Wrap(errs.CodeInternal, "operations.insert_pending", err)
		}
		if !inserted {
			return s.replayOperation(ctx, tx, accRow.ID, pgdto.OpSpend, in.OperationID, &result)
		}

		lockedRow, err := s.accounts.LockByID(ctx, tx, accRow.ID)
		if err != nil {
			return errs.Wrap(errs.CodeInternal, "accounts.lock_by_id", err)
		}

		agg, err := accountFromRow(lockedRow)
		if err != nil {
			return errs.Wrap(errs.CodeInternal, "account aggregate build failed", err)
		}

		actor := actorUserID
		updatedAgg, evDraft, err := agg.ApplySpend(ledger.Points(in.AmountPoints), &actor, opTs)
		if err != nil {
			// важно: кешируем “бизнес-ошибку” для replay
			if code, ok := errs.CodeOf(err); ok && code == errs.CodeNotEnoughBalance {
				if ferr := s.finalizeErr(ctx, tx, agg.ID, pgdto.OpSpend, in.OperationID, errs.CodeNotEnoughBalance, "not enough balance"); ferr != nil {
					return ferr
				}
			}
			return err
		}

		evRow, err := s.events.Insert(ctx, tx, eventInsertFromDraft(evDraft))
		if err != nil {
			return errs.Wrap(errs.CodeInternal, "events.insert", err)
		}

		updatedRow, err := s.accounts.UpdateAfterSpend(ctx, tx, updatedAgg.ID, updatedAgg.Balance.Int())
		if err != nil {
			return errs.Wrap(errs.CodeInternal, "accounts.update_after_spend", err)
		}

		result = dto.OperationOut{
			OperationID:      in.OperationID,
			OpType:           dto.OpSpend,
			Event:            mapper.EventOut(evRow),
			Balance:          mapper.BalanceOut(updatedRow, s.now()),
			IdempotentReplay: false,
		}

		if err := s.finalizeOK(ctx, tx, updatedRow.ID, pgdto.OpSpend, in.OperationID, result); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		s.log.ErrorContext(ctx, "cashier.spend failed", "ms", time.Since(start).Milliseconds(), "err", err)
		return dto.OperationOut{}, err
	}

	s.log.InfoContext(ctx, "cashier.spend ok",
		"ms", time.Since(start).Milliseconds(),
		"operationID", in.OperationID,
		"publicCode", in.PublicCode,
		"replay", result.IdempotentReplay,
	)

	return result, nil
}

// computeEarnDomain uses domain rules to compute points and resolve levels.
func computeEarnDomain(
	rs pgdto.RulesetWithLevels,
	totalSpendBefore ledger.Money,
	purchase ledger.Money,
) (earned ledger.Points, levelBefore rules.LevelCode, levelAfter rules.LevelCode, base ledger.Money, _ error) {
	base, err := ledger.ParseMoney(rs.Ruleset.BaseRubPerPoint.String())
	if err != nil {
		return 0, "", "", ledger.Money{}, errs.Wrap(errs.CodeInvalidRuleset, "baseRubPerPoint parse failed", err)
	}

	levels, err := levelRulesFromRepo(rs.Levels)
	if err != nil {
		return 0, "", "", ledger.Money{}, err
	}
	if err := rules.ValidateLevels(levels); err != nil {
		return 0, "", "", ledger.Money{}, err
	}
	rules.SortLevels(levels)

	beforeRule, err := rules.ResolveLevel(totalSpendBefore, levels)
	if err != nil {
		return 0, "", "", ledger.Money{}, err
	}

	afterSpend := totalSpendBefore.Add(purchase)
	afterRule, err := rules.ResolveLevel(afterSpend, levels)
	if err != nil {
		return 0, "", "", ledger.Money{}, err
	}

	pts, err := rules.ComputeEarnPoints(purchase, base, beforeRule.PercentEarn)
	if err != nil {
		return 0, "", "", ledger.Money{}, err
	}

	return pts, beforeRule.LevelCode, afterRule.LevelCode, base, nil
}

func levelRulesFromRepo(in []pgdto.LevelRuleRow) ([]rules.LevelRule, error) {
	out := make([]rules.LevelRule, 0, len(in))
	for _, lv := range in {
		thr, err := ledger.ParseMoney(lv.ThresholdTotalSpend.String())
		if err != nil {
			return nil, errs.Wrap(errs.CodeInvalidLevels, fmt.Sprintf("thresholdTotalSpend parse failed: %s", lv.ThresholdTotalSpend.String()), err)
		}
		perc, err := rules.ParsePercent(lv.PercentEarn.String())
		if err != nil {
			return nil, errs.Wrap(errs.CodeInvalidLevels, fmt.Sprintf("percentEarn parse failed: %s", lv.PercentEarn.String()), err)
		}
		out = append(out, rules.LevelRule{
			ID:                  lv.ID,
			LevelCode:           rules.LevelCode(lv.LevelCode),
			ThresholdTotalSpend: thr,
			PercentEarn:         perc,
		})
	}
	return out, nil
}

func accountFromRow(r pgdto.AccountRow) (account.Account, error) {
	ts, err := ledger.ParseMoney(r.TotalSpendMoney.String())
	if err != nil {
		return account.Account{}, err
	}
	return account.Account{
		ID:         r.ID,
		PublicCode: account.PublicCode(r.PublicCode),
		Balance:    ledger.Points(r.BalancePoints),
		TotalSpend: ts,
		LevelCode:  rules.LevelCode(r.LevelCode),
	}, nil
}

func eventInsertFromDraft(d ledger.EventDraft) pgdto.EventInsert {
	var typ pgdto.EventType
	switch d.Type {
	case ledger.EventEarn:
		typ = pgdto.EventEarn
	case ledger.EventSpend:
		typ = pgdto.EventSpend
	default:
		typ = pgdto.EventType(d.Type)
	}

	var amt *pgdto.Money
	if d.AmountMoney != nil {
		m := pgdto.Money(d.AmountMoney.Decimal())
		amt = &m
	}

	return pgdto.EventInsert{
		AccountID:    d.AccountID,
		Type:         typ,
		DeltaPoints:  int(d.DeltaPoints),
		BalanceAfter: int(d.BalanceAfter),
		AmountMoney:  amt,
		RulesetID:    d.RulesetID,
		ActorUserID:  d.ActorUserID,
		Ts:           d.Ts,
	}
}

// Money parsing helper: enforces <=2 fraction digits.
func parseMoney2(s string) (ledger.Money, error) {
	m, err := ledger.ParseMoney(s)
	if err != nil {
		return ledger.Money{}, err
	}
	if -m.Decimal().Exponent() > 2 {
		return ledger.Money{}, fmt.Errorf("too many fraction digits: %s", s)
	}
	return m, nil
}
