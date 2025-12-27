// internal/repository/postgres/repo/events_repo.go
package repo

import (
	"context"
	"database/sql/driver"
	"fmt"
	"time"

	pg "Beanefits/internal/repository/postgres"
	pgdto "Beanefits/internal/repository/postgres/dto"
	"Beanefits/internal/repository/postgres/sqlc/gen"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

type EventsRepo struct {
	q *gen.Queries
}

func NewEventsRepo(q *gen.Queries) *EventsRepo { return &EventsRepo{q: q} }

func (r *EventsRepo) Insert(ctx context.Context, db pg.DBTX, in pgdto.EventInsert) (pgdto.EventRow, error) {
	row, err := r.q.InsertEvent(ctx, db, gen.InsertEventParams{
		AccountID:    in.AccountID,
		Column2:      gen.EventType(in.Type), // важно: это sqlc enum
		DeltaPoints:  int32(in.DeltaPoints),
		BalanceAfter: int32(in.BalanceAfter),
		AmountMoney:  numericFromMoneyPtr(in.AmountMoney),
		RulesetID:    int8FromPtr(in.RulesetID),
		ActorUserID:  int8FromPtr(in.ActorUserID),
		Ts:           timestamptz(in.Ts),
	})
	if err != nil {
		return pgdto.EventRow{}, err
	}

	return mapEventInsertRow(row)
}

func (r *EventsRepo) ListByAccount(ctx context.Context, db pg.DBTX, accountID int64, limit int, beforeTs *time.Time) ([]pgdto.EventRow, error) {
	if beforeTs == nil {
		rows, err := r.q.ListEventsByAccount(ctx, db, gen.ListEventsByAccountParams{
			AccountID: accountID,
			Limit:     int32(limit),
		})
		if err != nil {
			return nil, err
		}
		return mapEventListRows(rows)
	}

	rows, err := r.q.ListEventsByAccountBefore(ctx, db, gen.ListEventsByAccountBeforeParams{
		AccountID: accountID,
		Limit:     int32(limit),
		Ts:        timestamptz(*beforeTs),
	})
	if err != nil {
		return nil, err
	}

	out := make([]pgdto.EventRow, 0, len(rows))
	for _, rw := range rows {
		ev, err := mapEventListBeforeRow(rw)
		if err != nil {
			return nil, err
		}
		out = append(out, ev)
	}
	return out, nil
}

// ---------- mapping ----------

func mapEventInsertRow(rw gen.InsertEventRow) (pgdto.EventRow, error) {
	return mapEventBase(
		rw.ID,
		rw.AccountID,
		rw.Type,
		rw.DeltaPoints,
		rw.BalanceAfter,
		rw.AmountMoney,
		rw.RulesetID,
		rw.ActorUserID,
		rw.Ts,
		rw.CreatedAt,
	)
}

func mapEventListRows(rows []gen.ListEventsByAccountRow) ([]pgdto.EventRow, error) {
	out := make([]pgdto.EventRow, 0, len(rows))
	for _, rw := range rows {
		ev, err := mapEventBase(
			rw.ID,
			rw.AccountID,
			rw.Type,
			rw.DeltaPoints,
			rw.BalanceAfter,
			rw.AmountMoney,
			rw.RulesetID,
			rw.ActorUserID,
			rw.Ts,
			rw.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		out = append(out, ev)
	}
	return out, nil
}

func mapEventListBeforeRow(rw gen.ListEventsByAccountBeforeRow) (pgdto.EventRow, error) {
	return mapEventBase(
		rw.ID,
		rw.AccountID,
		rw.Type,
		rw.DeltaPoints,
		rw.BalanceAfter,
		rw.AmountMoney,
		rw.RulesetID,
		rw.ActorUserID,
		rw.Ts,
		rw.CreatedAt,
	)
}

func mapEventBase(
	id int64,
	accountID int64,
	typ string,
	delta int32,
	balanceAfter int32,
	amountMoney pgtype.Numeric,
	rulesetID pgtype.Int8,
	actorUserID pgtype.Int8,
	ts pgtype.Timestamptz,
	createdAt pgtype.Timestamptz,
) (pgdto.EventRow, error) {
	amt, err := moneyPtrFromNumeric(amountMoney)
	if err != nil {
		return pgdto.EventRow{}, err
	}

	return pgdto.EventRow{
		ID:           id,
		AccountID:    accountID,
		Type:         pgdto.EventType(typ),
		DeltaPoints:  int(delta),
		BalanceAfter: int(balanceAfter),
		AmountMoney:  amt,
		RulesetID:    ptrFromInt8(rulesetID),
		ActorUserID:  ptrFromInt8(actorUserID),
		Ts:           ts.Time,
		CreatedAt:    createdAt.Time,
	}, nil
}

// ---------- pgtype helpers ----------

func timestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func int8FromPtr(v *int64) pgtype.Int8 {
	if v == nil {
		return pgtype.Int8{Valid: false}
	}
	return pgtype.Int8{Int64: *v, Valid: true}
}

func ptrFromInt8(v pgtype.Int8) *int64 {
	if !v.Valid {
		return nil
	}
	x := v.Int64
	return &x
}

func numericFromMoneyPtr(m *pgdto.Money) pgtype.Numeric {
	if m == nil {
		return pgtype.Numeric{Valid: false}
	}
	var n pgtype.Numeric
	// Надёжно: сканим из string (а не из decimal напрямую).
	if err := n.Scan(decimal.Decimal(*m).String()); err != nil {
		// это означает несостыковку типов/кодека, лучше увидеть сразу
		panic(fmt.Errorf("pgtype.Numeric.Scan(string): %w", err))
	}
	n.Valid = true
	return n
}

func moneyPtrFromNumeric(n pgtype.Numeric) (*pgdto.Money, error) {
	if !n.Valid {
		return nil, nil
	}

	// pgtype.Numeric implements driver.Valuer: Value() -> string/[]byte обычно.
	var (
		val driver.Value
		err error
	)
	val, err = n.Value()
	if err != nil {
		return nil, fmt.Errorf("pgtype.Numeric.Value(): %w", err)
	}

	var s string
	switch v := val.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		s = fmt.Sprint(v)
	}

	d, err := decimal.NewFromString(s)
	if err != nil {
		return nil, fmt.Errorf("decimal parse from numeric %q: %w", s, err)
	}

	m := pgdto.Money(d)
	return &m, nil
}
