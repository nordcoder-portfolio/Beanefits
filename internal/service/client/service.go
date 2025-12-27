package client

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"Beanefits/internal/domain/errs"
	pg "Beanefits/internal/repository/postgres"
	pgdto "Beanefits/internal/repository/postgres/dto"
)

type Clock func() time.Time

type Service struct {
	db pg.DBTX

	users    pg.UsersRepo
	roles    pg.RolesRepo
	accounts pg.AccountsRepo
	events   pg.EventsRepo

	now Clock
	log *slog.Logger
}

type Deps struct {
	DB pg.DBTX

	Users    pg.UsersRepo
	Roles    pg.RolesRepo
	Accounts pg.AccountsRepo
	Events   pg.EventsRepo

	Now Clock
	Log *slog.Logger
}

func New(deps Deps) *Service {
	n := deps.Now
	if n == nil {
		n = time.Now
	}

	l := deps.Log
	if l == nil {
		l = slog.Default()
	}
	l = l.With("layer", "service", "svc", "client")

	if deps.DB == nil {
		panic("client.New: deps.DB is nil")
	}
	if deps.Users == nil || deps.Roles == nil || deps.Accounts == nil || deps.Events == nil {
		panic("client.New: repos are nil")
	}

	return &Service{
		db:       deps.DB,
		users:    deps.Users,
		roles:    deps.Roles,
		accounts: deps.Accounts,
		events:   deps.Events,
		now:      n,
		log:      l,
	}
}

func (s *Service) requireActiveUser(ctx context.Context, userID int64) (pgdto.UserRow, error) {
	u, ok, err := s.users.GetByID(ctx, s.db, userID)
	if err != nil {
		return pgdto.UserRow{}, errs.Wrap(errs.CodeInternal, "users.get_by_id", err)
	}
	if !ok {
		return pgdto.UserRow{}, errs.Wrap(errs.CodeUserNotFound, "user not found", errors.New("invariant broken"))
	}
	if !u.IsActive {
		return pgdto.UserRow{}, errs.New(errs.CodeUserInactive, "user is inactive")
	}
	return u, nil
}

func (s *Service) requireRoles(ctx context.Context, userID int64) ([]pgdto.RoleCode, error) {
	rs, err := s.roles.GetRoles(ctx, s.db, userID)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternal, "roles.get_roles", err)
	}
	if len(rs) == 0 {
		return nil, errs.New(errs.CodeRolesNotFound, "roles not found")
	}
	return rs, nil
}

func hasRole(roles []pgdto.RoleCode, want pgdto.RoleCode) bool {
	for _, r := range roles {
		if r == want {
			return true
		}
	}
	return false
}

func (s *Service) requireAccount(ctx context.Context, userID int64) (pgdto.AccountRow, error) {
	acc, ok, err := s.accounts.GetByUserID(ctx, s.db, userID)
	if err != nil {
		return pgdto.AccountRow{}, errs.Wrap(errs.CodeInternal, "accounts.get_by_user_id", err)
	}
	if !ok {
		return pgdto.AccountRow{}, errs.New(errs.CodeAccountNotFound, "account not found")
	}
	return acc, nil
}
