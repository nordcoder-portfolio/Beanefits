package cashier

import (
	"context"
	"strings"
	"time"

	"Beanefits/internal/domain/account"
	"Beanefits/internal/domain/errs"
	"Beanefits/internal/service/dto"
	"Beanefits/internal/service/mapper"
	svcvalidation "Beanefits/internal/service/validation"
)

func (s *Service) LookupAccountByCode(ctx context.Context, publicCode string) (dto.AccountOut, error) {
	start := time.Now()
	s.log.InfoContext(ctx, "cashier.lookup_account_by_code start", "publicCode", publicCode)

	if _, err := account.ParsePublicCode(publicCode); err != nil {
		s.log.ErrorContext(ctx, "cashier.lookup_account_by_code failed", "ms", time.Since(start).Milliseconds(), "err", err)
		return dto.AccountOut{}, err
	}

	row, ok, err := s.accounts.GetByPublicCode(ctx, s.db, publicCode)
	if err != nil {
		wrapped := errs.Wrap(errs.CodeInternal, "accounts.get_by_public_code", err)
		s.log.ErrorContext(ctx, "cashier.lookup_account_by_code failed", "ms", time.Since(start).Milliseconds(), "err", wrapped)
		return dto.AccountOut{}, wrapped
	}
	if !ok {
		e := errs.New(errs.CodeAccountNotFound, "account not found")
		s.log.ErrorContext(ctx, "cashier.lookup_account_by_code failed", "ms", time.Since(start).Milliseconds(), "err", e)
		return dto.AccountOut{}, e
	}

	out := dto.AccountOut{AccountBase: mapper.AccountBase(row)}
	s.log.InfoContext(ctx, "cashier.lookup_account_by_code ok", "ms", time.Since(start).Milliseconds(), "accountID", row.ID)
	return out, nil
}

func (s *Service) GetAccountEventsByCode(ctx context.Context, publicCode string, in dto.EventsIn) (dto.EventsOut, error) {
	const (
		defaultLimit = 20
		maxLimit     = 100
	)

	start := time.Now()
	s.log.InfoContext(ctx, "cashier.get_account_events_by_code start", "publicCode", publicCode)

	if _, err := account.ParsePublicCode(publicCode); err != nil {
		s.log.ErrorContext(ctx, "cashier.get_account_events_by_code failed", "ms", time.Since(start).Milliseconds(), "err", err)
		return dto.EventsOut{}, err
	}

	limit, _ := svcvalidation.NormalizeListParams(in.Limit, 0, defaultLimit, maxLimit)

	acc, ok, err := s.accounts.GetByPublicCode(ctx, s.db, strings.TrimSpace(publicCode))
	if err != nil {
		wrapped := errs.Wrap(errs.CodeInternal, "accounts.get_by_public_code", err)
		s.log.ErrorContext(ctx, "cashier.get_account_events_by_code failed", "ms", time.Since(start).Milliseconds(), "err", wrapped)
		return dto.EventsOut{}, wrapped
	}
	if !ok {
		e := errs.New(errs.CodeAccountNotFound, "account not found")
		s.log.ErrorContext(ctx, "cashier.get_account_events_by_code failed", "ms", time.Since(start).Milliseconds(), "err", e)
		return dto.EventsOut{}, e
	}

	rows, err := s.events.ListByAccount(ctx, s.db, acc.ID, limit, in.BeforeTs)
	if err != nil {
		wrapped := errs.Wrap(errs.CodeInternal, "events.list_by_account", err)
		s.log.ErrorContext(ctx, "cashier.get_account_events_by_code failed", "ms", time.Since(start).Milliseconds(), "err", wrapped)
		return dto.EventsOut{}, wrapped
	}

	out := dto.EventsOut{
		Items:        make([]dto.EventOut, 0, len(rows)),
		NextBeforeTs: nil,
	}
	for _, r := range rows {
		out.Items = append(out.Items, mapper.EventOut(r))
	}
	if len(rows) == limit {
		ts := rows[len(rows)-1].Ts
		out.NextBeforeTs = &ts
	}

	s.log.InfoContext(ctx, "cashier.get_account_events_by_code ok",
		"ms", time.Since(start).Milliseconds(),
		"accountID", acc.ID,
		"items", len(out.Items),
	)

	return out, nil
}
