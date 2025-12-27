package repo

import (
	"context"
	"errors"

	pg "Beanefits/internal/repository/postgres"
	pgdto "Beanefits/internal/repository/postgres/dto"
	"Beanefits/internal/repository/postgres/sqlc/gen"

	"github.com/jackc/pgx/v5"
)

type UsersRepo struct {
	q *gen.Queries
}

func NewUsersRepo(q *gen.Queries) *UsersRepo { return &UsersRepo{q: q} }

func (r *UsersRepo) Create(ctx context.Context, db pg.DBTX, phone, passwordHash string) (int64, error) {
	return r.q.CreateUser(ctx, db, gen.CreateUserParams{
		Phone:        phone,
		PasswordHash: passwordHash,
	})
}

func (r *UsersRepo) GetByPhone(ctx context.Context, db pg.DBTX, phone string) (pgdto.UserRow, bool, error) {
	row, err := r.q.GetUserByPhone(ctx, db, phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgdto.UserRow{}, false, nil
		}
		return pgdto.UserRow{}, false, err
	}
	return mapUserByPhone(row), true, nil
}

func (r *UsersRepo) GetByID(ctx context.Context, db pg.DBTX, id int64) (pgdto.UserRow, bool, error) {
	row, err := r.q.GetUserByID(ctx, db, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgdto.UserRow{}, false, nil
		}
		return pgdto.UserRow{}, false, err
	}
	return mapUserByID(row), true, nil
}

func (r *UsersRepo) List(ctx context.Context, db pg.DBTX, q string, limit, offset int) ([]pgdto.UserRow, error) {
	rows, err := r.q.ListUsers(ctx, db, gen.ListUsersParams{
		Column1: q,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, err
	}

	out := make([]pgdto.UserRow, 0, len(rows))
	for _, row := range rows {
		out = append(out, mapUserFromList(row))
	}
	return out, nil
}

func (r *UsersRepo) Deactivate(ctx context.Context, db pg.DBTX, id int64) error {
	return r.q.DeactivateUser(ctx, db, id)
}

// --- mapping ---

func mapUserBase(id int64, phone, passwordHash string, isActive bool, createdAt any) pgdto.UserRow {
	// createdAt in sqlc rows is pgtype.Timestamptz; we keep mapper signature flexible
	// only to share code between row types.
	switch t := createdAt.(type) {
	case gen.GetUserByIDRow:
		_ = t
	case gen.GetUserByPhoneRow:
		_ = t
	case gen.ListUsersRow:
		_ = t
	}
	return pgdto.UserRow{
		ID:           id,
		Phone:        phone,
		PasswordHash: passwordHash,
		IsActive:     isActive,
		// concrete mappers below set this field
	}
}

func mapUserByID(row gen.GetUserByIDRow) pgdto.UserRow {
	return pgdto.UserRow{
		ID:           row.ID,
		Phone:        row.Phone,
		PasswordHash: row.PasswordHash,
		IsActive:     row.IsActive,
		CreatedAt:    row.CreatedAt.Time,
	}
}

func mapUserByPhone(row gen.GetUserByPhoneRow) pgdto.UserRow {
	return pgdto.UserRow{
		ID:           row.ID,
		Phone:        row.Phone,
		PasswordHash: row.PasswordHash,
		IsActive:     row.IsActive,
		CreatedAt:    row.CreatedAt.Time,
	}
}

func mapUserFromList(row gen.ListUsersRow) pgdto.UserRow {
	return pgdto.UserRow{
		ID:           row.ID,
		Phone:        row.Phone,
		PasswordHash: row.PasswordHash,
		IsActive:     row.IsActive,
		CreatedAt:    row.CreatedAt.Time,
	}
}
