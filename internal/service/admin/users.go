package admin

import (
	"context"
	"fmt"
	"strings"
	"time"

	"Beanefits/internal/domain/errs"
	pg "Beanefits/internal/repository/postgres"
	"Beanefits/internal/service/dto"
	"Beanefits/internal/service/mapper"
	svcvalidation "Beanefits/internal/service/validation"
)

func (s *Service) ListUsers(ctx context.Context, in dto.ListUsersIn) (dto.UsersOut, error) {
	const (
		defaultLimit = 20
		maxLimit     = 100
	)

	limit, offset := svcvalidation.NormalizeListParams(in.Limit, in.Offset, defaultLimit, maxLimit)
	q := strings.TrimSpace(in.Q)

	start := time.Now()
	s.log.InfoContext(ctx, "admin.list_users start", "q", q, "limit", limit, "offset", offset)

	rows, err := s.users.List(ctx, s.db, q, limit, offset)
	if err != nil {
		wrapped := errs.Wrap(errs.CodeInternal, "users.list", err)
		s.log.ErrorContext(ctx, "admin.list_users failed", "ms", time.Since(start).Milliseconds(), "err", wrapped)
		return dto.UsersOut{}, wrapped
	}

	out := dto.UsersOut{
		Items: make([]dto.UserWithRoles, 0, len(rows)),
		Total: nil, // optional для MVP
	}

	for _, u := range rows {
		rs, err := s.roles.GetRoles(ctx, s.db, u.ID)
		if err != nil {
			wrapped := errs.Wrap(errs.CodeInternal, "roles.get_roles", err)
			s.log.ErrorContext(ctx, "admin.list_users failed", "ms", time.Since(start).Milliseconds(), "userID", u.ID, "err", wrapped)
			return dto.UsersOut{}, wrapped
		}
		if len(rs) == 0 {
			// Это скорее повреждение данных: user без ролей.
			e := errs.New(errs.CodeRolesNotFound, fmt.Sprintf("roles not found for user_id=%d", u.ID))
			s.log.ErrorContext(ctx, "admin.list_users failed", "ms", time.Since(start).Milliseconds(), "userID", u.ID, "err", e)
			return dto.UsersOut{}, e
		}

		out.Items = append(out.Items, mapper.UserWithRoles(u, rs))
	}

	s.log.InfoContext(ctx, "admin.list_users ok",
		"ms", time.Since(start).Milliseconds(),
		"items", len(out.Items),
	)

	return out, nil
}

func (s *Service) DeactivateUser(ctx context.Context, userID int64) error {
	start := time.Now()
	s.log.InfoContext(ctx, "admin.deactivate_user start", "userID", userID)

	if userID <= 0 {
		e := errs.New(errs.CodeInternal, "invalid userID")
		s.log.ErrorContext(ctx, "admin.deactivate_user failed", "ms", time.Since(start).Milliseconds(), "userID", userID, "err", e)
		return e
	}

	err := s.txm.WithinTx(ctx, func(ctx context.Context, tx pg.DBTX) error {
		_, ok, err := s.users.GetByID(ctx, tx, userID)
		if err != nil {
			return errs.Wrap(errs.CodeInternal, "users.get_by_id", err)
		}
		if !ok {
			// В идеале: завести errs.CodeUserNotFound и вернуть его.
			return errs.New(errs.CodeInternal, fmt.Sprintf("user not found: id=%d", userID))
		}

		if err := s.users.Deactivate(ctx, tx, userID); err != nil {
			return errs.Wrap(errs.CodeInternal, "users.deactivate", err)
		}
		return nil
	})

	if err != nil {
		s.log.ErrorContext(ctx, "admin.deactivate_user failed", "ms", time.Since(start).Milliseconds(), "userID", userID, "err", err)
		return err
	}

	s.log.InfoContext(ctx, "admin.deactivate_user ok", "ms", time.Since(start).Milliseconds(), "userID", userID)
	return nil
}
