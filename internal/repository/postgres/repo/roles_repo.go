package repo

import (
	"context"

	pg "Beanefits/internal/repository/postgres"
	pgdto "Beanefits/internal/repository/postgres/dto"
	"Beanefits/internal/repository/postgres/sqlc/gen"
)

type RolesRepo struct {
	q *gen.Queries
}

func NewRolesRepo(q *gen.Queries) *RolesRepo {
	return &RolesRepo{q: q}
}

func (r *RolesRepo) GetRoles(ctx context.Context, db pg.DBTX, userID int64) ([]pgdto.RoleCode, error) {
	rows, err := r.q.GetRolesByUserID(ctx, db, userID)
	if err != nil {
		return nil, err
	}

	out := make([]pgdto.RoleCode, 0, len(rows))
	for _, code := range rows {
		out = append(out, pgdto.RoleCode(code))
	}
	return out, nil
}

func (r *RolesRepo) AddRole(ctx context.Context, db pg.DBTX, userID int64, role pgdto.RoleCode) error {
	return r.q.AddUserRole(ctx, db, gen.AddUserRoleParams{
		UserID:   userID,
		RoleCode: string(role),
	})
}

func (r *RolesRepo) RemoveRole(ctx context.Context, db pg.DBTX, userID int64, role pgdto.RoleCode) error {
	return r.q.RemoveUserRole(ctx, db, gen.RemoveUserRoleParams{
		UserID:   userID,
		RoleCode: string(role),
	})
}
