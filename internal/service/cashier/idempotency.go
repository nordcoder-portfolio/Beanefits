package cashier

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"Beanefits/internal/domain/errs"
	pg "Beanefits/internal/repository/postgres"
	pgdto "Beanefits/internal/repository/postgres/dto"
	"Beanefits/internal/service/dto"
)

type opErrCache struct {
	Code errs.Code `json:"code"`
	Msg  string    `json:"msg"`
}

func (s *Service) replayOperation(ctx context.Context, tx pg.DBTX, accountID int64, opType pgdto.OperationType, operationID string, out *dto.OperationOut) error {
	rec, ok, err := s.operations.Get(ctx, tx, accountID, opType, operationID)
	if err != nil {
		return errs.Wrap(errs.CodeInternal, "operations.get", err)
	}
	if !ok || rec.HTTPStatus == nil || rec.ResponseJSON == nil {
		return errs.Wrap(errs.CodeInternal, "idempotency record incomplete", errors.New("missing cached response"))
	}

	if *rec.HTTPStatus != 200 {
		var cached opErrCache
		if err := json.Unmarshal(*rec.ResponseJSON, &cached); err != nil {
			return errs.Wrap(errs.CodeInternal, "unmarshal cached error", err)
		}
		if cached.Code == "" {
			return errs.Wrap(errs.CodeInternal, "cached error code missing", errors.New("invalid cache"))
		}
		return errs.New(cached.Code, cached.Msg)
	}

	if err := json.Unmarshal(*rec.ResponseJSON, out); err != nil {
		return errs.Wrap(errs.CodeInternal, "unmarshal cached success", err)
	}
	out.IdempotentReplay = true
	return nil
}

func (s *Service) finalizeOK(ctx context.Context, tx pg.DBTX, accountID int64, opType pgdto.OperationType, operationID string, out dto.OperationOut) error {
	// Canonical persisted payload should not include per-response replay flag.
	out.IdempotentReplay = false

	b, err := json.Marshal(out)
	if err != nil {
		return errs.Wrap(errs.CodeInternal, "marshal operation out", err)
	}

	return s.operations.Finalize(ctx, tx, pgdto.OperationFinalize{
		AccountID:    accountID,
		OpType:       opType,
		OperationID:  operationID,
		HTTPStatus:   200,
		ResponseJSON: json.RawMessage(b),
	})
}

func (s *Service) finalizeErr(ctx context.Context, tx pg.DBTX, accountID int64, opType pgdto.OperationType, operationID string, code errs.Code, msg string) error {
	b, err := json.Marshal(opErrCache{Code: code, Msg: msg})
	if err != nil {
		return errs.Wrap(errs.CodeInternal, "marshal cached error", err)
	}

	// For MVP: store “transport-ish” status approximation for replay logic.
	status := 409
	return s.operations.Finalize(ctx, tx, pgdto.OperationFinalize{
		AccountID:    accountID,
		OpType:       opType,
		OperationID:  operationID,
		HTTPStatus:   status,
		ResponseJSON: json.RawMessage(b),
	})
}

func marshalEarnRequest(in dto.EarnIn, ts time.Time) (pgdto.JSON, error) {
	type req struct {
		OperationID string    `json:"operationId"`
		PublicCode  string    `json:"publicCode"`
		AmountMoney string    `json:"amountMoney"`
		Ts          time.Time `json:"ts"`
	}
	b, err := json.Marshal(req{
		OperationID: in.OperationID,
		PublicCode:  in.PublicCode,
		AmountMoney: in.AmountMoney,
		Ts:          ts,
	})
	return pgdto.JSON(b), err
}

func marshalSpendRequest(in dto.SpendIn, ts time.Time) (pgdto.JSON, error) {
	type req struct {
		OperationID  string    `json:"operationId"`
		PublicCode   string    `json:"publicCode"`
		AmountPoints int       `json:"amountPoints"`
		Ts           time.Time `json:"ts"`
	}
	b, err := json.Marshal(req{
		OperationID:  in.OperationID,
		PublicCode:   in.PublicCode,
		AmountPoints: in.AmountPoints,
		Ts:           ts,
	})
	return pgdto.JSON(b), err
}
