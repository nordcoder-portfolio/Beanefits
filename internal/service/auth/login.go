package auth

import (
	"context"
	"time"

	"Beanefits/internal/domain/errs"
	"Beanefits/internal/domain/user"
	pgdto "Beanefits/internal/repository/postgres/dto"
	sdto "Beanefits/internal/service/dto"
	"Beanefits/internal/service/mapper"
)

func (s *Service) Login(ctx context.Context, in sdto.LoginIn) (sdto.AuthOut, error) {
	start := time.Now()
	s.log.InfoContext(ctx, "auth.login start", "phone", in.Phone)

	if _, err := user.ParsePhone(in.Phone); err != nil {
		s.log.ErrorContext(ctx, "auth.login failed", "ms", time.Since(start).Milliseconds(), "err", err)
		return sdto.AuthOut{}, err
	}

	u, ok, err := s.users.GetByPhone(ctx, s.db, in.Phone)
	if err != nil {
		wrapped := errs.Wrap(errs.CodeInternal, "users.get_by_phone", err)
		s.log.ErrorContext(ctx, "auth.login failed", "ms", time.Since(start).Milliseconds(), "err", wrapped)
		return sdto.AuthOut{}, wrapped
	}
	if !ok {
		e := s.invalidCredentials()
		s.log.ErrorContext(ctx, "auth.login failed", "ms", time.Since(start).Milliseconds(), "err", e)
		return sdto.AuthOut{}, e
	}
	if !u.IsActive {
		e := errs.New(errs.CodeUserInactive, "user is inactive")
		s.log.ErrorContext(ctx, "auth.login failed", "ms", time.Since(start).Milliseconds(), "err", e)
		return sdto.AuthOut{}, e
	}

	match, err := s.hasher.Compare(ctx, in.Password, u.PasswordHash)
	if err != nil {
		wrapped := errs.Wrap(errs.CodeInternal, "hasher.compare", err)
		s.log.ErrorContext(ctx, "auth.login failed", "ms", time.Since(start).Milliseconds(), "err", wrapped)
		return sdto.AuthOut{}, wrapped
	}
	if !match {
		e := s.invalidCredentials()
		s.log.ErrorContext(ctx, "auth.login failed", "ms", time.Since(start).Milliseconds(), "err", e)
		return sdto.AuthOut{}, e
	}

	rs, err := s.roles.GetRoles(ctx, s.db, u.ID)
	if err != nil {
		wrapped := errs.Wrap(errs.CodeInternal, "roles.get_roles", err)
		s.log.ErrorContext(ctx, "auth.login failed", "ms", time.Since(start).Milliseconds(), "err", wrapped)
		return sdto.AuthOut{}, wrapped
	}
	if len(rs) == 0 {
		e := errs.New(errs.CodeRolesNotFound, "roles not found")
		s.log.ErrorContext(ctx, "auth.login failed", "ms", time.Since(start).Milliseconds(), "err", e)
		return sdto.AuthOut{}, e
	}

	acc, ok, err := s.accounts.GetByUserID(ctx, s.db, u.ID)
	if err != nil {
		wrapped := errs.Wrap(errs.CodeInternal, "accounts.get_by_user_id", err)
		s.log.ErrorContext(ctx, "auth.login failed", "ms", time.Since(start).Milliseconds(), "err", wrapped)
		return sdto.AuthOut{}, wrapped
	}
	if !ok {
		e := errs.New(errs.CodeAccountNotFound, "account not found")
		s.log.ErrorContext(ctx, "auth.login failed", "ms", time.Since(start).Milliseconds(), "err", e)
		return sdto.AuthOut{}, e
	}

	userWithRoles := pgdto.UserWithRoles{UserRow: u, Roles: rs}

	token, err := s.issuer.Issue(ctx, userWithRoles)
	if err != nil {
		wrapped := errs.Wrap(errs.CodeInternal, "issuer.issue", err)
		s.log.ErrorContext(ctx, "auth.login failed", "ms", time.Since(start).Milliseconds(), "err", wrapped)
		return sdto.AuthOut{}, wrapped
	}

	out := mapper.AuthOut(token, userWithRoles, acc)
	s.log.InfoContext(ctx, "auth.login ok", "ms", time.Since(start).Milliseconds(), "userID", u.ID)
	return out, nil
}
