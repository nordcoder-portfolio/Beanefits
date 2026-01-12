package repo

import (
	"context"
	"errors"

	pg "Beanefits/internal/repository/postgres"
	pgdto "Beanefits/internal/repository/postgres/dto"
	"Beanefits/internal/repository/postgres/sqlc/gen"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type AccountsRepo struct {
	q *gen.Queries
}

func NewAccountsRepo(q *gen.Queries) *AccountsRepo {
	return &AccountsRepo{q: q}
}

func (r *AccountsRepo) CreateForUser(
	ctx context.Context,
	db pg.DBTX,
	userID int64,
	publicCode string,
	initialLevelCode string,
) (pgdto.AccountRow, error) {
	row, err := r.q.CreateAccountForUser(ctx, db, gen.CreateAccountForUserParams{
		UserID:     userID,
		PublicCode: publicCode,
		LevelCode:  text(initialLevelCode),
	})
	if err != nil {
		return pgdto.AccountRow{}, err
	}
	return mapAccountRowFromCreate(row), nil
}

func (r *AccountsRepo) GetByUserID(ctx context.Context, db pg.DBTX, userID int64) (pgdto.AccountRow, bool, error) {
	row, err := r.q.GetAccountByUserID(ctx, db, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgdto.AccountRow{}, false, nil
		}
		return pgdto.AccountRow{}, false, err
	}
	return mapAccountRowFromGet(row), true, nil
}

func (r *AccountsRepo) GetByPublicCode(ctx context.Context, db pg.DBTX, publicCode string) (pgdto.AccountRow, bool, error) {
	row, err := r.q.GetAccountByPublicCode(ctx, db, publicCode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgdto.AccountRow{}, false, nil
		}
		return pgdto.AccountRow{}, false, err
	}
	return mapAccountRowFromGetPub(row), true, nil
}

func (r *AccountsRepo) LockByID(ctx context.Context, db pg.DBTX, accountID int64) (pgdto.AccountRow, error) {
	row, err := r.q.LockAccountByID(ctx, db, accountID)
	if err != nil {
		return pgdto.AccountRow{}, err
	}
	return mapAccountRowFromLock(row), nil
}

func (r *AccountsRepo) UpdateAfterEarn(
	ctx context.Context,
	db pg.DBTX,
	accountID int64,
	balancePoints int,
	totalSpend pgdto.Money,
	levelCode string,
) (pgdto.AccountRow, error) {
	row, err := r.q.UpdateAccountAfterEarn(ctx, db, gen.UpdateAccountAfterEarnParams{
		ID:              accountID,
		BalancePoints:   int32(balancePoints),
		TotalSpendMoney: totalSpend,
		LevelCode:       text(levelCode),
	})
	if err != nil {
		return pgdto.AccountRow{}, err
	}
	return mapAccountRowFromUpdateEarn(row), nil
}

func (r *AccountsRepo) UpdateAfterSpend(ctx context.Context, db pg.DBTX, accountID int64, balancePoints int) (pgdto.AccountRow, error) {
	row, err := r.q.UpdateAccountAfterSpend(ctx, db, gen.UpdateAccountAfterSpendParams{
		ID:            accountID,
		BalancePoints: int32(balancePoints),
	})
	if err != nil {
		return pgdto.AccountRow{}, err
	}
	return mapAccountRowFromUpdateSpend(row), nil
}

func text(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: true}
}

func mapAccountBase(
	id int64,
	userID int64,
	publicCode string,
	createdAt pgtype.Timestamptz,
	balancePoints int32,
	totalSpendMoney pgdto.Money,
	levelCode string,
) pgdto.AccountRow {
	return pgdto.AccountRow{
		ID:              id,
		UserID:          userID,
		PublicCode:      publicCode,
		BalancePoints:   int(balancePoints),
		TotalSpendMoney: totalSpendMoney,
		LevelCode:       levelCode,
		CreatedAt:       createdAt.Time,
	}
}

func mapAccountRowFromCreate(rw gen.CreateAccountForUserRow) pgdto.AccountRow {
	return mapAccountBase(
		rw.ID,
		rw.UserID,
		rw.PublicCode,
		rw.CreatedAt,
		rw.BalancePoints,
		rw.TotalSpendMoney,
		rw.LevelCode,
	)
}

func mapAccountRowFromGet(rw gen.GetAccountByUserIDRow) pgdto.AccountRow {
	return mapAccountBase(
		rw.ID,
		rw.UserID,
		rw.PublicCode,
		rw.CreatedAt,
		rw.BalancePoints,
		rw.TotalSpendMoney,
		rw.LevelCode,
	)
}

func mapAccountRowFromGetPub(rw gen.GetAccountByPublicCodeRow) pgdto.AccountRow {
	return mapAccountBase(
		rw.ID,
		rw.UserID,
		rw.PublicCode,
		rw.CreatedAt,
		rw.BalancePoints,
		rw.TotalSpendMoney,
		rw.LevelCode,
	)
}

func mapAccountRowFromLock(rw gen.LockAccountByIDRow) pgdto.AccountRow {
	return mapAccountBase(
		rw.ID,
		rw.UserID,
		rw.PublicCode,
		rw.CreatedAt,
		rw.BalancePoints,
		rw.TotalSpendMoney,
		rw.LevelCode,
	)
}

func mapAccountRowFromUpdateEarn(rw gen.UpdateAccountAfterEarnRow) pgdto.AccountRow {
	return mapAccountBase(
		rw.ID,
		rw.UserID,
		rw.PublicCode,
		rw.CreatedAt,
		rw.BalancePoints,
		rw.TotalSpendMoney,
		rw.LevelCode,
	)
}

func mapAccountRowFromUpdateSpend(rw gen.UpdateAccountAfterSpendRow) pgdto.AccountRow {
	return mapAccountBase(
		rw.ID,
		rw.UserID,
		rw.PublicCode,
		rw.CreatedAt,
		rw.BalancePoints,
		rw.TotalSpendMoney,
		rw.LevelCode,
	)
}
