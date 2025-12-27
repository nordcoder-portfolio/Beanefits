package client

import (
	"context"
	"time"

	"Beanefits/internal/domain/errs"
	sdto "Beanefits/internal/service/dto"
	"Beanefits/internal/service/mapper"
)

func (s *Service) GetEvents(ctx context.Context, userID int64, in sdto.EventsIn) (sdto.EventsOut, error) {
	start := time.Now()
	s.log.InfoContext(ctx, "client.get_events start", "userID", userID, "limit", in.Limit, "beforeTs", in.BeforeTs)

	limit := in.Limit
	if limit == 0 {
		limit = 20
	}
	if limit < 0 {
		e := errs.New(errs.CodeInternal, "invalid limit")
		s.log.ErrorContext(ctx, "client.get_events failed", "ms", time.Since(start).Milliseconds(), "err", e)
		return sdto.EventsOut{}, e
	}

	if _, err := s.requireActiveUser(ctx, userID); err != nil {
		s.log.ErrorContext(ctx, "client.get_events failed", "ms", time.Since(start).Milliseconds(), "err", err)
		return sdto.EventsOut{}, err
	}

	acc, err := s.requireAccount(ctx, userID)
	if err != nil {
		s.log.ErrorContext(ctx, "client.get_events failed", "ms", time.Since(start).Milliseconds(), "err", err)
		return sdto.EventsOut{}, err
	}

	rows, err := s.events.ListByAccount(ctx, s.db, acc.ID, limit, in.BeforeTs)
	if err != nil {
		wrapped := errs.Wrap(errs.CodeInternal, "events.list_by_account", err)
		s.log.ErrorContext(ctx, "client.get_events failed", "ms", time.Since(start).Milliseconds(), "err", wrapped)
		return sdto.EventsOut{}, wrapped
	}

	out := sdto.EventsOut{
		Items:        make([]sdto.EventOut, 0, len(rows)),
		NextBeforeTs: nil,
	}
	for _, r := range rows {
		out.Items = append(out.Items, mapper.EventOut(r))
	}

	// newest-first -> next cursor is last element ts
	if len(rows) == limit {
		ts := rows[len(rows)-1].Ts
		out.NextBeforeTs = &ts
	}

	s.log.InfoContext(ctx, "client.get_events ok", "ms", time.Since(start).Milliseconds(), "userID", userID, "accountID", acc.ID, "items", len(out.Items))
	return out, nil
}
