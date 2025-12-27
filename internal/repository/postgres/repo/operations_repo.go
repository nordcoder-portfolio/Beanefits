package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	pg "Beanefits/internal/repository/postgres"
	pgdto "Beanefits/internal/repository/postgres/dto"
	"Beanefits/internal/repository/postgres/sqlc/gen"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type OperationsRepo struct {
	q *gen.Queries
}

func NewOperationsRepo(q *gen.Queries) *OperationsRepo {
	return &OperationsRepo{q: q}
}

func (r *OperationsRepo) Get(
	ctx context.Context,
	db pg.DBTX,
	accountID int64,
	opType pgdto.OperationType,
	operationID string,
) (pgdto.OperationRecord, bool, error) {
	row, err := r.q.GetOperation(ctx, db, gen.GetOperationParams{
		AccountID:   accountID,
		Column2:     gen.EventType(opType), // op_type in DB is event_type enum
		OperationID: operationID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgdto.OperationRecord{}, false, nil
		}
		return pgdto.OperationRecord{}, false, err
	}

	rec := pgdto.OperationRecord{
		AccountID:   row.AccountID,
		OpType:      pgdto.OperationType(row.OpType),
		OperationID: row.OperationID,

		RequestJSON:  json.RawMessage(bytesOrNullJSON(row.RequestJson)),
		ResponseJSON: nil,
		HTTPStatus:   nil,

		CreatedAt: row.CreatedAt.Time,
	}

	// response_json может быть NULL -> sqlc даёт []byte, но pgx возвращает nil для NULL
	if row.ResponseJson != nil {
		j := json.RawMessage(row.ResponseJson)
		rec.ResponseJSON = (*pgdto.JSON)(&j)
	}

	if row.HttpStatus.Valid {
		x := int(row.HttpStatus.Int32)
		rec.HTTPStatus = &x
	}

	return rec, true, nil
}

func (r *OperationsRepo) InsertPending(ctx context.Context, db pg.DBTX, in pgdto.OperationPendingInsert) (bool, error) {
	n, err := r.q.InsertOperationPending(ctx, db, gen.InsertOperationPendingParams{
		AccountID:   in.AccountID,
		Column2:     gen.EventType(in.OpType),
		OperationID: in.OperationID,
		RequestJson: []byte(in.RequestJSON),
	})
	if err != nil {
		return false, err
	}
	return n == 1, nil
}

func (r *OperationsRepo) Finalize(ctx context.Context, db pg.DBTX, in pgdto.OperationFinalize) error {
	status, err := int4(int32(in.HTTPStatus))
	if err != nil {
		return err
	}

	return r.q.FinalizeOperation(ctx, db, gen.FinalizeOperationParams{
		AccountID:    in.AccountID,
		Column2:      gen.EventType(in.OpType),
		OperationID:  in.OperationID,
		HttpStatus:   status,
		ResponseJson: []byte(in.ResponseJSON),
	})
}

func int4(v int32) (pgtype.Int4, error) {
	// pgtype.Int4 не валидирует диапазон, но оставим хук на случай будущих изменений.
	return pgtype.Int4{Int32: v, Valid: true}, nil
}

func bytesOrNullJSON(b []byte) []byte {
	if b == nil {
		// request_json в твоей схеме NOT NULL, но на всякий случай не паниковать.
		return []byte("null")
	}
	// защита от случайной передачи пустого массива байт как "невалидный json"
	if len(b) == 0 {
		return []byte("null")
	}
	// на всякий случай проверим, что это хотя бы валидный JSON (минимально)
	if !json.Valid(b) {
		// не падаем: кеш идемпотентности должен быть “best effort”, но лучше отдать ошибку наверх,
		// потому что это повреждение данных.
		// Если хочешь мягче — поменяй на return []byte("null").
		panic(fmt.Errorf("operations.request_json is not valid json: %q", string(b)))
	}
	return b
}
