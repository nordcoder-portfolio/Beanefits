package auth

import (
	"context"
	"errors"
	"time"

	"Beanefits/internal/domain/errs"
	"Beanefits/internal/domain/user"
	"Beanefits/internal/repository/postgres"
	pg "Beanefits/internal/repository/postgres"
	pgdto "Beanefits/internal/repository/postgres/dto"
	sdto "Beanefits/internal/service/dto"
	"Beanefits/internal/service/mapper"
)

func (s *Service) RegisterClient(ctx context.Context, in sdto.RegisterIn) (sdto.AuthOut, error) {
	start := time.Now()
	s.log.InfoContext(ctx, "auth.register_client start", "phone", in.Phone)

	if _, err := user.ParsePhone(in.Phone); err != nil {
		s.log.ErrorContext(ctx, "auth.register_client failed", "ms", time.Since(start).Milliseconds(), "err", err)
		return sdto.AuthOut{}, err
	}

	var (
		created pgdto.UserWithRoles
		acc     pgdto.AccountRow
	)

	if err := s.txm.WithinTx(ctx, func(ctx context.Context, tx pg.DBTX) error {
		_, exists, err := s.users.GetByPhone(ctx, tx, in.Phone)
		if err != nil {
			return errs.Wrap(errs.CodeInternal, "users.get_by_phone", err)
		}
		if exists {
			return errs.New(errs.CodePhoneAlreadyExists, "phone already exists")
		}

		hash, err := s.hasher.Hash(ctx, in.Password)
		if err != nil {
			return errs.Wrap(errs.CodeInternal, "hasher.hash", err)
		}

		userID, err := s.users.Create(ctx, tx, in.Phone, hash)
		if err != nil {
			if postgres.IsUniqueViolation(err, "users_phone_key") {
				return errs.New(errs.CodePhoneAlreadyExists, "phone already exists")
			}
			return errs.Wrap(errs.CodeInternal, "users.create", err)
		}

		if err := s.roles.AddRole(ctx, tx, userID, pgdto.RoleClient); err != nil {
			return errs.Wrap(errs.CodeInternal, "roles.add_role", err)
		}

		var lastCollision error
		for i := 0; i < s.publicCodeRetries; i++ {
			code, err := s.codeGen.New(ctx)
			if err != nil {
				return errs.Wrap(errs.CodeInternal, "public_code.new", err)
			}

			a, err := s.accounts.CreateForUser(ctx, tx, userID, code, s.initialLevelCode)
			if err == nil {
				acc = a
				lastCollision = nil
				break
			}

			if postgres.IsUniqueViolation(err, "accounts_public_code_key") {
				lastCollision = err
				continue
			}

			return errs.Wrap(errs.CodeInternal, "accounts.create_for_user", err)
		}

		if lastCollision != nil {
			return errs.Wrap(errs.CodePublicCodeCollision, "public code collision retries exhausted", lastCollision)
		}

		u, ok, err := s.users.GetByID(ctx, tx, userID)
		if err != nil {
			return errs.Wrap(errs.CodeInternal, "users.get_by_id", err)
		}
		if !ok {
			return errs.Wrap(errs.CodeInternal, "user not found after create", errors.New("invariant broken"))
		}

		rs, err := s.roles.GetRoles(ctx, tx, userID)
		if err != nil {
			return errs.Wrap(errs.CodeInternal, "roles.get_roles", err)
		}
		if len(rs) == 0 {
			return errs.New(errs.CodeRolesNotFound, "roles not found")
		}

		created = pgdto.UserWithRoles{UserRow: u, Roles: rs}
		return nil
	}); err != nil {
		s.log.ErrorContext(ctx, "auth.register_client failed", "ms", time.Since(start).Milliseconds(), "err", err)
		return sdto.AuthOut{}, err
	}

	token, err := s.issuer.Issue(ctx, created)
	if err != nil {
		wrapped := errs.Wrap(errs.CodeInternal, "issuer.issue", err)
		s.log.ErrorContext(ctx, "auth.register_client failed", "ms", time.Since(start).Milliseconds(), "err", wrapped)
		return sdto.AuthOut{}, wrapped
	}

	out := mapper.AuthOut(token, created, acc)
	s.log.InfoContext(ctx, "auth.register_client ok", "ms", time.Since(start).Milliseconds(), "userID", created.ID)
	return out, nil
}
